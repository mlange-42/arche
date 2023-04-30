package ecs

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/mlange-42/arche/ecs/stats"
)

// ComponentID returns the [ID] for a component type via generics.
// Registers the type if it is not already registered.
//
// The number of unique component types per [World] is limited to [MaskTotalBits].
func ComponentID[T any](w *World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.componentID(tp)
}

// TypeID returns the [ID] for a component type.
// Registers the type if it is not already registered.
//
// The number of unique component types per [World] is limited to [MaskTotalBits].
func TypeID(w *World, tp reflect.Type) ID {
	return w.componentID(tp)
}

// ResourceID returns the [ResID] for a resource type via generics.
// Registers the type if it is not already registered.
//
// The number of resources per [World] is limited to [MaskTotalBits].
func ResourceID[T any](w *World) ResID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.resourceID(tp)
}

// GetResource returns a pointer to the given resource type in world.
//
// Returns nil if there is no such resource.
//
// Uses reflection. For more efficient access, see [World.Resources],
// and [github.com/mlange-42/arche/generic.Resource.Get] for a generic variant.
// These methods are more than 20 times faster than the GetResource function.
func GetResource[T any](w *World) *T {
	return w.resources.Get(ResourceID[T](w)).(*T)
}

// AddResource adds a resource to the world.
// Returns the ID for the added resource.
//
// Panics if there is already such a resource.
//
// Uses reflection. For more efficient access, see [World.AddResource],
// and [github.com/mlange-42/arche/generic.Resource.Add] for a generic variant.
//
// The number of resources per [World] is limited to [MaskTotalBits].
func AddResource[T any](w *World, res *T) ResID {
	id := ResourceID[T](w)
	w.resources.Add(id, res)
	return id
}

// World is the central type holding [Entity] and component data, as well as resources.
type World struct {
	config         Config                    // World configuration.
	listener       func(e *EntityEvent)      // Component change listener.
	resources      Resources                 // World resources.
	entities       []entityIndex             // Mapping from entities to archetype and index.
	entityPool     entityPool                // Pool for entities.
	archetypes     pagedSlice[archetype]     // The archetypes.
	freeArchetypes []int32                   // Indices of free archetypes.
	graph          pagedSlice[archetypeNode] // The archetype graph.
	relationNodes  []*archetypeNode          // Archetype nodes that have an entity relation.
	locks          lockMask                  // World locks.
	registry       componentRegistry[ID]     // Component registry.
	filterCache    Cache                     // Cache for registered filters.
	stats          stats.WorldStats          // Cached world statistics
}

// NewWorld creates a new [World] from an optional [Config].
//
// Uses the default [Config] if called without an argument.
// Accepts zero or one arguments.
func NewWorld(config ...Config) World {
	if len(config) > 1 {
		panic("can't use more than one Config")
	}
	if len(config) == 1 {
		return fromConfig(config[0])
	}
	return fromConfig(NewConfig())
}

// fromConfig creates a new [World] from a [Config].
func fromConfig(conf Config) World {
	if conf.CapacityIncrement < 1 {
		panic("invalid CapacityIncrement in config, must be > 0")
	}
	entities := make([]entityIndex, 1, conf.CapacityIncrement)
	entities[0] = entityIndex{arch: nil, index: 0}
	w := World{
		config:        conf,
		entities:      entities,
		entityPool:    newEntityPool(uint32(conf.CapacityIncrement)),
		registry:      newComponentRegistry(),
		archetypes:    newPagedSlice[archetype](),
		graph:         newPagedSlice[archetypeNode](),
		relationNodes: []*archetypeNode{},
		locks:         lockMask{},
		listener:      nil,
		resources:     newResources(),
		filterCache:   newCache(),
	}
	node := w.createArchetypeNode(Mask{}, -1)
	w.createArchetype(node, Entity{}, 0, false)
	return w
}

// NewEntity returns a new or recycled [Entity].
// The given component types are added to the entity.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// Note that calling a method with varargs in Go causes a slice allocation.
// For maximum performance, pre-allocate a slice of component IDs and pass it using ellipsis:
//
//	// fast
//	world.NewEntity(idA, idB, idC)
//	// even faster
//	world.NewEntity(ids...)
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
func (w *World) NewEntity(comps ...ID) Entity {
	w.checkLocked()

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil, Entity{}, -1)
	}

	entity := w.createEntity(arch)

	if w.listener != nil {
		w.listener(&EntityEvent{entity, Mask{}, arch.Mask, comps, nil, arch.graphNode.Ids, 1})
	}
	return entity
}

// NewEntityWith returns a new or recycled [Entity].
// The given component values are assigned to the entity.
//
// The components in the Comp field of [Component] must be pointers.
// The passed pointers are no valid references to the assigned memory!
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
func (w *World) NewEntityWith(comps ...Component) Entity {
	w.checkLocked()

	if len(comps) == 0 {
		return w.NewEntity()
	}

	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}

	arch := w.archetypes.Get(0)
	arch = w.findOrCreateArchetype(arch, ids, nil, Entity{}, -1)

	entity := w.createEntity(arch)

	for _, c := range comps {
		w.copyTo(entity, c.ID, c.Comp)
	}

	if w.listener != nil {
		w.listener(&EntityEvent{entity, Mask{}, arch.Mask, ids, nil, arch.graphNode.Ids, 1})
	}
	return entity
}

func (w *World) newEntityTarget(targetID ID, target Entity, comps ...ID) Entity {
	w.checkLocked()

	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	arch := w.archetypes.Get(0)

	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil, target, int8(targetID))
	}
	w.checkRelation(arch, targetID)

	entity := w.createEntity(arch)

	if !target.IsZero() {
		w.entities[target.id].isTarget = true
	}

	if w.listener != nil {
		w.listener(&EntityEvent{entity, Mask{}, arch.Mask, comps, nil, arch.graphNode.Ids, 1})
	}
	return entity
}

func (w *World) newEntityTargetWith(targetID ID, target Entity, comps ...Component) Entity {
	w.checkLocked()

	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}

	arch := w.archetypes.Get(0)
	arch = w.findOrCreateArchetype(arch, ids, nil, target, int8(targetID))
	w.checkRelation(arch, targetID)

	entity := w.createEntity(arch)

	if !target.IsZero() {
		w.entities[target.id].isTarget = true
	}

	for _, c := range comps {
		w.copyTo(entity, c.ID, c.Comp)
	}

	if w.listener != nil {
		w.listener(&EntityEvent{entity, Mask{}, arch.Mask, ids, nil, arch.graphNode.Ids, 1})
	}
	return entity
}

// Creates new entities without returning a query over them.
// Used via [World.Batch].
func (w *World) newEntities(count int, targetID int8, target Entity, comps ...ID) (*archetype, uint32) {
	arch, startIdx := w.newEntitiesNoNotify(count, targetID, target, comps...)

	if w.listener != nil {
		cnt := uint32(count)
		var i uint32
		for i = 0; i < cnt; i++ {
			idx := startIdx + i
			entity := arch.GetEntity(uintptr(idx))
			w.listener(&EntityEvent{entity, Mask{}, arch.Mask, comps, nil, arch.graphNode.Ids, 1})
		}
	}

	return arch, startIdx
}

// Creates new entities and returns a query over them.
// Used via [World.Batch].
func (w *World) newEntitiesQuery(count int, targetID int8, target Entity, comps ...ID) Query {
	arch, startIdx := w.newEntitiesNoNotify(count, targetID, target, comps...)
	lock := w.lock()
	return newArchQuery(w, lock, arch, startIdx)
}

// Creates new entities with component values without returning a query over them.
// Used via [World.Batch].
func (w *World) newEntitiesWith(count int, targetID int8, target Entity, comps ...Component) (*archetype, uint32) {
	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}

	arch, startIdx := w.newEntitiesWithNoNotify(count, targetID, target, ids, comps...)

	if w.listener != nil {
		var i uint32
		cnt := uint32(count)
		for i = 0; i < cnt; i++ {
			idx := startIdx + i
			entity := arch.GetEntity(uintptr(idx))
			w.listener(&EntityEvent{entity, Mask{}, arch.Mask, ids, nil, arch.graphNode.Ids, 1})
		}
	}

	return arch, startIdx
}

// Creates new entities with component values and returns a query over them.
// Used via [World.Batch].
func (w *World) newEntitiesWithQuery(count int, targetID int8, target Entity, comps ...Component) Query {
	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}

	arch, startIdx := w.newEntitiesWithNoNotify(count, targetID, target, ids, comps...)
	lock := w.lock()
	return newArchQuery(w, lock, arch, startIdx)
}

// RemoveEntity removes and recycles an [Entity].
//
// Panics when called for a removed (and potentially recycled) entity.
//
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
func (w *World) RemoveEntity(entity Entity) {
	w.checkLocked()

	if !w.entityPool.Alive(entity) {
		panic("can't remove a dead entity")
	}

	index := &w.entities[entity.id]
	oldArch := index.arch

	if w.listener != nil {
		lock := w.lock()
		w.listener(&EntityEvent{entity, oldArch.Mask, Mask{}, nil, oldArch.graphNode.Ids, nil, -1})
		w.unlock(lock)
	}

	swapped := oldArch.Remove(index.index)

	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := oldArch.GetEntity(index.index)
		w.entities[swapEntity.id].index = index.index
	}
	index.arch = nil

	if index.isTarget {
		w.cleanupArchetypes(entity)
		index.isTarget = false
	}

	w.cleanupArchetype(oldArch)
}

// RemoveEntities removes and recycles all entities matching a filter.
//
// Returns the number of removed entities.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
func (w *World) RemoveEntities(filter Filter) int {
	w.checkLocked()

	lock := w.lock()
	numArches := w.archetypes.Len()
	var count uintptr
	var i int32
	for i = 0; i < numArches; i++ {
		arch := w.archetypes.Get(i)

		if !arch.IsActive() {
			continue
		}

		if !arch.Matches(filter) {
			continue
		}

		len := uintptr(arch.Len())
		count += len

		var j uintptr
		for j = 0; j < len; j++ {
			entity := arch.GetEntity(j)
			if w.listener != nil {
				w.listener(&EntityEvent{entity, arch.Mask, Mask{}, nil, arch.graphNode.Ids, nil, -1})
			}
			index := &w.entities[entity.id]
			index.arch = nil

			if index.isTarget {
				w.cleanupArchetypes(entity)
				index.isTarget = false
			}

			w.entityPool.Recycle(entity)
		}

		arch.Reset()
		w.cleanupArchetype(arch)
	}
	w.unlock(lock)

	return int(count)
}

// Alive reports whether an entity is still alive.
func (w *World) Alive(entity Entity) bool {
	return w.entityPool.Alive(entity)
}

// Get returns a pointer to the given component of an [Entity].
// Returns nil if the entity has no such component.
//
// Panics when called for a removed (and potentially recycled) entity.
//
// See [World.GetUnchecked] for an optimized version for static entities.
// See also [github.com/mlange-42/arche/generic.Map.Get] for a generic variant.
func (w *World) Get(entity Entity, comp ID) unsafe.Pointer {
	if !w.entityPool.Alive(entity) {
		panic("can't get component of a dead entity")
	}
	index := &w.entities[entity.id]
	return index.arch.Get(index.index, comp)
}

// GetUnchecked returns a pointer to the given component of an [Entity].
// Returns nil if the entity has no such component.
//
// GetUnchecked is an optimized version of [World.Get],
// for cases where entities are static or checked with [World.Alive] in user code.
// It can also be used after getting another component of the same entity with [World.Get].
//
// Panics when called for a removed entity, but not for a recycled entity.
//
// See also [github.com/mlange-42/arche/generic.Map.Get] for a generic variant.
func (w *World) GetUnchecked(entity Entity, comp ID) unsafe.Pointer {
	index := &w.entities[entity.id]
	return index.arch.Get(index.index, comp)
}

// Has returns whether an [Entity] has a given component.
//
// Panics when called for a removed (and potentially recycled) entity.
//
// See [World.HasUnchecked] for an optimized version for static entities.
// See also [github.com/mlange-42/arche/generic.Map.Has] for a generic variant.
func (w *World) Has(entity Entity, comp ID) bool {
	if !w.entityPool.Alive(entity) {
		panic("can't check for component of a dead entity")
	}
	return w.entities[entity.id].arch.HasComponent(comp)
}

// HasUnchecked returns whether an [Entity] has a given component.
//
// HasUnchecked is an optimized version of [World.Has],
// for cases where entities are static or checked with [World.Alive] in user code.
//
// Panics when called for a removed entity, but not for a recycled entity.
//
// See also [github.com/mlange-42/arche/generic.Map.Has] for a generic variant.
func (w *World) HasUnchecked(entity Entity, comp ID) bool {
	return w.entities[entity.id].arch.HasComponent(comp)
}

// Add adds components to an [Entity].
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - when called with components that can't be added because they are already present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// Note that calling a method with varargs in Go causes a slice allocation.
// For maximum performance, pre-allocate a slice of component IDs and pass it using ellipsis:
//
//	// fast
//	world.Add(entity, idA, idB, idC)
//	// even faster
//	world.Add(entity, ids...)
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
func (w *World) Add(entity Entity, comps ...ID) {
	w.Exchange(entity, comps, nil)
}

// Assign assigns multiple components to an [Entity], using pointers for the content.
//
// The components in the Comp field of [Component] must be pointers.
// The passed pointers are no valid references to the assigned memory!
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - when called with components that can't be added because they are already present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
func (w *World) Assign(entity Entity, comps ...Component) {
	len := len(comps)
	if len == 0 {
		panic("no components given to assign")
	}
	if len == 1 {
		c := comps[0]
		w.Exchange(entity, []ID{c.ID}, nil)
		w.copyTo(entity, c.ID, c.Comp)
		return
	}
	ids := make([]ID, len)
	for i, c := range comps {
		ids[i] = c.ID
	}
	w.Exchange(entity, ids, nil)
	for _, c := range comps {
		w.copyTo(entity, c.ID, c.Comp)
	}
}

// Set overwrites a component for an [Entity], using a given pointer for the content.
//
// The passed component must be a pointer.
// Returns a pointer to the assigned memory.
// The passed in pointer is not a valid reference to that memory!
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - if the entity does not have a component of that type.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [github.com/mlange-42/arche/generic.Map.Set] for a generic variant.
func (w *World) Set(entity Entity, id ID, comp interface{}) unsafe.Pointer {
	return w.copyTo(entity, id, comp)
}

// Remove removes components from an entity.
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - when called with components that can't be removed because they are not present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
func (w *World) Remove(entity Entity, comps ...ID) {
	w.Exchange(entity, nil, comps)
}

// Exchange adds and removes components in one pass.
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Exchange].
func (w *World) Exchange(entity Entity, add []ID, rem []ID) {
	w.checkLocked()

	if !w.entityPool.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}

	if len(add) == 0 && len(rem) == 0 {
		return
	}
	index := &w.entities[entity.id]
	oldArch := index.arch
	mask := oldArch.Mask
	oldMask := mask
	for _, comp := range add {
		if oldArch.HasComponent(comp) {
			panic("entity already has this component, can't add")
		}
		mask.Set(comp, true)
	}
	for _, comp := range rem {
		if !oldArch.HasComponent(comp) {
			panic("entity does not have this component, can't remove")
		}
		mask.Set(comp, false)
	}

	oldIDs := oldArch.Components()

	arch := w.findOrCreateArchetype(oldArch, add, rem, Entity{}, -1)
	newIndex := arch.Alloc(entity)

	for _, id := range oldIDs {
		if mask.Get(id) {
			comp := oldArch.Get(index.index, id)
			arch.SetPointer(newIndex, id, comp)
		}
	}

	swapped := oldArch.Remove(index.index)

	if swapped {
		swapEntity := oldArch.GetEntity(index.index)
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{arch: arch, index: newIndex, isTarget: index.isTarget}

	w.cleanupArchetype(oldArch)

	if w.listener != nil {
		w.listener(&EntityEvent{entity, oldMask, arch.Mask, add, rem, arch.graphNode.Ids, 0})
	}
}

// GetRelation returns the target entity for an entity relation.
func (w *World) GetRelation(entity Entity, comp ID) Entity {
	if !w.entityPool.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}

	index := &w.entities[entity.id]
	w.checkRelation(index.arch, comp)

	return index.arch.Relation
}

// SetRelation sets the target entity for an entity relation.
func (w *World) SetRelation(entity Entity, comp ID, target Entity) {
	w.checkLocked()

	if !w.entityPool.Alive(entity) {
		panic("can't set relation for a dead entity")
	}
	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	index := &w.entities[entity.id]
	w.checkRelation(index.arch, comp)

	oldArch := index.arch

	if index.arch.Relation == target {
		return
	}

	arch := oldArch.graphNode.GetArchetype(target)
	if arch == nil {
		arch = w.createArchetype(oldArch.graphNode, target, int8(oldArch.graphNode.relation), true)
	}

	newIndex := arch.Alloc(entity)
	for _, id := range oldArch.graphNode.Ids {
		comp := oldArch.Get(index.index, id)
		arch.SetPointer(newIndex, id, comp)
	}

	swapped := oldArch.Remove(index.index)

	if swapped {
		swapEntity := oldArch.GetEntity(index.index)
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{arch: arch, index: newIndex, isTarget: index.isTarget}
	w.entities[target.id].isTarget = true

	w.cleanupArchetype(oldArch)
}

func (w *World) checkRelation(arch *archetype, comp ID) {
	if arch.graphNode.relation != int8(comp) {
		if !arch.HasComponent(comp) {
			panic(fmt.Sprintf("entity does not have relation component %v", w.registry.Types[comp]))
		}
		panic(fmt.Sprintf("not a relation component: %v", w.registry.Types[comp]))
	}
}

// Reset removes all entities and resources from the world.
//
// Does NOT free reserved memory, remove archetypes, clear the registry, clear cached filters, etc.
//
// Can be used to run systematic simulations without the need to re-allocate memory for each run.
// Accelerates re-populating the world by a factor of 2-3.
func (w *World) Reset() {
	w.checkLocked()

	w.entities = w.entities[:1]
	w.entityPool.Reset()
	w.locks.Reset()
	w.resources.reset()

	len := w.archetypes.Len()
	var i int32
	for i = 0; i < len; i++ {
		w.archetypes.Get(i).Reset()
	}
}

// Query creates a [Query] iterator.
//
// The [ecs] core package provides only the filter [All] for querying the given components.
// Further, it can be chained with [Mask.Without] (see the examples) to exclude components.
//
// Locks the world to prevent changes to component compositions.
// The lock is released automatically when the query finishes iteration, or when [Query.Close] is called.
// The number of simultaneous locks (and thus open queries) at a given time is limited to [MaskTotalBits].
//
// For type-safe generics queries, see package [github.com/mlange-42/arche/generic].
// For advanced filtering, see package [github.com/mlange-42/arche/filter].
func (w *World) Query(filter Filter) Query {
	l := w.lock()
	if cached, ok := filter.(*CachedFilter); ok {
		archetypes := &w.filterCache.get(cached).Archetypes
		return newQuery(w, cached, l, archetypes, true)
	}

	return newQuery(w, filter, l, &w.archetypes, false)
}

// Resources of the world.
//
// Resources are component-like data that is not associated to an entity, but unique to the world.
func (w *World) Resources() *Resources {
	return &w.resources
}

// Cache returns the [Cache] of the world, for registering filters.
//
// See [Cache] for details on filter caching.
func (w *World) Cache() *Cache {
	if w.filterCache.getArchetypes == nil {
		w.filterCache.getArchetypes = w.getArchetypes
	}
	return &w.filterCache
}

// IsLocked returns whether the world is locked by any queries.
func (w *World) IsLocked() bool {
	return w.locks.IsLocked()
}

// Mask returns the archetype [Mask] for the given [Entity].
//
// Can be used for fast checks of the entity composition, e.g. using a [Filter].
func (w *World) Mask(entity Entity) Mask {
	if !w.entityPool.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}
	return w.entities[entity.id].arch.Mask
}

// ComponentType returns the reflect.Type for a given component ID, as well as whether the ID is in use.
func (w *World) ComponentType(id ID) (reflect.Type, bool) {
	return w.registry.ComponentType(id)
}

// SetListener sets a listener callback func(e EntityEvent) for the world.
// The listener is immediately called on every [ecs.Entity] change.
// Replaces the current listener. Call with nil to remove a listener.
//
// For details, see [EntityEvent].
func (w *World) SetListener(listener func(e *EntityEvent)) {
	w.listener = listener
}

// Stats reports statistics for inspecting the World.
func (w *World) Stats() *stats.WorldStats {
	w.stats.Entities = stats.EntityStats{
		Used:     w.entityPool.Len(),
		Total:    w.entityPool.Cap(),
		Recycled: w.entityPool.Available(),
		Capacity: w.entityPool.TotalCap(),
	}

	compCount := len(w.registry.Components)
	types := append([]reflect.Type{}, w.registry.Types[:compCount]...)

	memory := cap(w.entities)*int(entityIndexSize) + w.entityPool.TotalCap()*int(entitySize)

	cntOld := int32(len(w.stats.Archetypes))
	cntNew := int32(w.archetypes.Len())
	var i int32
	for i = 0; i < cntOld; i++ {
		arch := &w.stats.Archetypes[i]
		w.archetypes.Get(i).UpdateStats(arch, &w.registry)
		memory += arch.Memory
	}
	for i = cntOld; i < cntNew; i++ {
		w.stats.Archetypes = append(w.stats.Archetypes, w.archetypes.Get(i).Stats(&w.registry))
		memory += w.stats.Archetypes[i].Memory
	}

	w.stats.ComponentCount = compCount
	w.stats.ComponentTypes = types
	w.stats.Locked = w.IsLocked()
	w.stats.Memory = memory
	w.stats.CachedFilters = len(w.filterCache.filters)

	return &w.stats
}

// lock the world and get the lock bit for later unlocking.
func (w *World) lock() uint8 {
	return w.locks.Lock()
}

// unlock unlocks the given lock bit.
func (w *World) unlock(l uint8) {
	w.locks.Unlock(l)
}

// checkLocked checks if the world is locked, and panics if so.
func (w *World) checkLocked() {
	if w.IsLocked() {
		panic("attempt to modify a locked world")
	}
}

// Internal method to create new entities.
func (w *World) newEntitiesNoNotify(count int, targetID int8, target Entity, comps ...ID) (*archetype, uint32) {
	w.checkLocked()

	if count < 1 {
		panic("can only create a positive number of entities")
	}

	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil, target, targetID)
	}
	if targetID >= 0 {
		w.checkRelation(arch, uint8(targetID))
		if !target.IsZero() {
			w.entities[target.id].isTarget = true
		}
	}

	startIdx := arch.Len()
	w.createEntities(arch, uint32(count))

	return arch, startIdx
}

// Internal method to create new entities with component values.
func (w *World) newEntitiesWithNoNotify(count int, targetID int8, target Entity, ids []ID, comps ...Component) (*archetype, uint32) {
	w.checkLocked()

	if count < 1 {
		panic("can only create a positive number of entities")
	}

	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	if len(comps) == 0 {
		return w.newEntitiesNoNotify(count, targetID, target)
	}

	cnt := uint32(count)

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, ids, nil, target, targetID)
	}
	if targetID >= 0 {
		w.checkRelation(arch, uint8(targetID))
		if !target.IsZero() {
			w.entities[target.id].isTarget = true
		}
	}

	startIdx := arch.Len()
	w.createEntities(arch, uint32(count))

	var i uint32
	for i = 0; i < cnt; i++ {
		idx := startIdx + i
		entity := arch.GetEntity(uintptr(idx))
		for _, c := range comps {
			w.copyTo(entity, c.ID, c.Comp)
		}
	}

	return arch, startIdx
}

// createEntity creates an Entity and adds it to the given archetype.
func (w *World) createEntity(arch *archetype) Entity {
	entity := w.entityPool.Get()
	idx := arch.Alloc(entity)
	len := len(w.entities)
	if int(entity.id) == len {
		if len == cap(w.entities) {
			old := w.entities
			w.entities = make([]entityIndex, len, len+w.config.CapacityIncrement)
			copy(w.entities, old)
		}
		w.entities = append(w.entities, entityIndex{arch: arch, index: idx})
	} else {
		w.entities[entity.id] = entityIndex{arch: arch, index: idx}
	}
	return entity
}

// createEntity creates multiple Entities and adds them to the given archetype.
func (w *World) createEntities(arch *archetype, count uint32) {
	startIdx := arch.Len()
	arch.AllocN(uint32(count))

	len := len(w.entities)
	required := len + int(count) - w.entityPool.Available()
	if required > cap(w.entities) {
		cap := capacity(required, w.config.CapacityIncrement)
		old := w.entities
		w.entities = make([]entityIndex, required, cap)
		copy(w.entities, old)
	} else if required > len {
		w.entities = w.entities[:required]
	}

	var i uint32
	for i = 0; i < count; i++ {
		idx := startIdx + i
		entity := w.entityPool.Get()
		arch.SetEntity(uintptr(idx), entity)
		w.entities[entity.id] = entityIndex{arch: arch, index: uintptr(idx)}
	}
}

// Copies a component to an entity
func (w *World) copyTo(entity Entity, id ID, comp interface{}) unsafe.Pointer {
	if !w.Has(entity, id) {
		panic("can't copy component into entity that has no such component type")
	}
	index := &w.entities[entity.id]
	arch := index.arch

	return arch.Set(index.index, id, comp)
}

// Tries to find an archetype by traversing the archetype graph,
// searching by mask and extending the graph if necessary.
// A new archetype is created for the final graph node if not already present.
func (w *World) findOrCreateArchetype(start *archetype, add []ID, rem []ID, target Entity, targetComponent int8) *archetype {
	curr := start.graphNode
	mask := start.Mask
	relation := start.graphNode.relation
	for _, id := range rem {
		mask.Set(id, false)
		if w.registry.IsRelation.Get(id) {
			relation = -1
		}
		if next, ok := curr.TransitionRemove.Get(id); ok {
			curr = next
		} else {
			next, _ := w.findOrCreateArchetypeSlow(mask, relation)
			next.TransitionAdd.Set(id, curr)
			curr.TransitionRemove.Set(id, next)
			curr = next
		}
	}
	for _, id := range add {
		mask.Set(id, true)
		if w.registry.IsRelation.Get(id) {
			if relation >= 0 {
				panic("entity already has a relation component")
			}
			relation = int8(id)
		}
		if next, ok := curr.TransitionAdd.Get(id); ok {
			curr = next
		} else {
			next, _ := w.findOrCreateArchetypeSlow(mask, relation)
			next.TransitionRemove.Set(id, curr)
			curr.TransitionAdd.Set(id, next)
			curr = next
		}
	}
	arch := curr.GetArchetype(target)
	if arch == nil {
		arch = w.createArchetype(curr, target, targetComponent, true)
	}
	return arch
}

// Tries to find an archetype for a mask, when it can't be reached through the archetype graph.
// Creates an archetype graph node.
func (w *World) findOrCreateArchetypeSlow(mask Mask, relation int8) (*archetypeNode, bool) {
	if arch, ok := w.findArchetypeSlow(mask); ok {
		return arch, false
	}
	return w.createArchetypeNode(mask, relation), true
}

// Searches for an archetype by a mask.
func (w *World) findArchetypeSlow(mask Mask) (*archetypeNode, bool) {
	length := w.graph.Len()
	var i int32
	for i = 0; i < length; i++ {
		arch := w.graph.Get(i)
		if arch.mask == mask {
			return arch, true
		}
	}
	return nil, false
}

// Creates a node in the archetype graph.
func (w *World) createArchetypeNode(mask Mask, relation int8) *archetypeNode {
	w.graph.Add(newArchetypeNode(mask, relation, w.config.CapacityIncrement))
	node := w.graph.Get(w.graph.Len() - 1)
	w.relationNodes = append(w.relationNodes, node)
	return node
}

// Creates an archetype for the given archetype graph node.
// Initializes the archetype with a capacity according to CapacityIncrement if forStorage is true,
// and with a capacity of 1 otherwise.
func (w *World) createArchetype(node *archetypeNode, target Entity, targetComponent int8, forStorage bool) *archetype {
	mask := node.mask
	count := int(mask.TotalBitsSet())
	types := make([]componentType, count)

	start := 0
	end := MaskTotalBits
	if mask.Lo == 0 {
		start = wordSize
	}
	if mask.Hi == 0 {
		end = wordSize
	}

	idx := 0
	for i := start; i < end; i++ {
		id := ID(i)
		if mask.Get(id) {
			types[idx] = componentType{ID: id, Type: w.registry.Types[id]}
			idx++
		}
	}

	var arch *archetype
	var archIndex int32
	lenFree := len(w.freeArchetypes)
	if lenFree > 0 {
		archIndex = w.freeArchetypes[lenFree-1]
		arch = w.archetypes.Get(archIndex)
		w.freeArchetypes = w.freeArchetypes[:lenFree-1]

		if int(archIndex) < len(w.stats.Archetypes) {
			w.stats.Archetypes[archIndex].Dirty = true
		}
	} else {
		w.archetypes.Add(archetype{})
		archIndex = w.archetypes.Len() - 1
		arch = w.archetypes.Get(archIndex)
	}
	arch.Init(node, archIndex, forStorage, target, targetComponent, types...)

	node.SetArchetype(target, arch)

	w.filterCache.addArchetype(arch)
	return arch
}

// Returns all archetypes that match the given filter. Used by [Cache].
func (w *World) getArchetypes(filter Filter) archetypePointers {
	arches := []*archetype{}
	ln := int32(w.archetypes.Len())
	var i int32
	for i = 0; i < ln; i++ {
		a := w.archetypes.Get(i)
		if a.IsActive() && a.Matches(filter) {
			arches = append(arches, a)
		}
	}
	return archetypePointers{arches}
}

// Removes the archetype if it is empty, and has a relation to a dead target.
func (w *World) cleanupArchetype(arch *archetype) {
	if arch.Len() > 0 || arch.graphNode.relation < 0 {
		return
	}
	target := arch.Relation
	if target.IsZero() || w.Alive(target) {
		return
	}

	w.deleteArchetype(arch, target)
}

// Removes empty archetypes that have a target relation to the given entity.
func (w *World) cleanupArchetypes(target Entity) {
	for _, node := range w.relationNodes {
		if arch, ok := node.archetypes[target]; ok && arch.Len() == 0 {
			w.deleteArchetype(arch, target)
		}
	}
}

func (w *World) deleteArchetype(arch *archetype, target Entity) {
	delete(arch.graphNode.archetypes, target)
	idx := arch.index
	w.freeArchetypes = append(w.freeArchetypes, idx)
	w.archetypes.Get(idx).Deactivate()

	w.filterCache.removeArchetype(arch)

	if int(idx) < len(w.stats.Archetypes) {
		w.stats.Archetypes[idx].Dirty = true
	}
}

// componentID returns the ID for a component type, and registers it if not already registered.
func (w *World) componentID(tp reflect.Type) ID {
	return w.registry.ComponentID(tp)
}

// resourceID returns the ID for a resource type, and registers it if not already registered.
func (w *World) resourceID(tp reflect.Type) ResID {
	return w.resources.registry.ComponentID(tp)
}

// closeQuery closes a query and unlocks the world.
func (w *World) closeQuery(query *Query) {
	query.archIndex = -2
	w.unlock(query.lockBit)

	if w.listener != nil {
		if arch, ok := query.archetypes.(batchArchetype); ok {
			w.notifyQuery(&arch)
		}
	}
}

// notifies the listener for all entities on a batch query.
func (w *World) notifyQuery(batchArch *batchArchetype) {
	arch := batchArch.Archetype
	var i uintptr
	len := uintptr(arch.Len())
	event := EntityEvent{Entity{}, Mask{}, arch.Mask, arch.graphNode.Ids, nil, arch.graphNode.Ids, 1}

	for i = uintptr(batchArch.StartIndex); i < len; i++ {
		entity := arch.GetEntity(i)
		event.Entity = entity
		w.listener(&event)
	}
}
