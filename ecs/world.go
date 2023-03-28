package ecs

import (
	"reflect"
	"unsafe"

	"github.com/mlange-42/arche/ecs/stats"
)

// ComponentID returns the [ID] for a component type via generics. Registers the type if it is not already registered.
func ComponentID[T any](w *World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.componentID(tp)
}

// TypeID returns the [ID] for a component type. Registers the type if it is not already registered.
func TypeID(w *World, tp reflect.Type) ID {
	return w.componentID(tp)
}

// ResourceID returns the [ResID] for a resource type via generics. Registers the type if it is not already registered.
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
//
// Panics if there is already such a resource.
//
// Uses reflection. For more efficient access, see [World.AddResource],
// and [github.com/mlange-42/arche/generic.Resource.Add] for a generic variant.
func AddResource[T any](w *World, res *T) {
	w.resources.Add(ResourceID[T](w), res)
}

// World is the central type holding [Entity] and component data, as well as resources.
type World struct {
	config      Config                    // World configuration.
	listener    func(e *EntityEvent)      // Component change listener.
	resources   Resources                 // World resources.
	entities    []entityIndex             // Mapping from entities to archetype and index.
	entityPool  entityPool                // Pool for entities.
	archetypes  pagedArr32[archetype]     // The archetypes.
	graph       pagedArr32[archetypeNode] // The archetype graph.
	locks       lockMask                  // World locks.
	registry    componentRegistry[ID]     // Component registry.
	filterCache Cache                     // Cache for registered filters.
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
		config:      conf,
		entities:    entities,
		entityPool:  newEntityPool(conf.CapacityIncrement),
		registry:    newComponentRegistry(),
		archetypes:  pagedArr32[archetype]{},
		graph:       pagedArr32[archetypeNode]{},
		locks:       lockMask{},
		listener:    nil,
		resources:   newResources(),
		filterCache: newCache(),
	}
	node := w.createArchetypeNode(Mask{})
	w.createArchetype(node, false)
	return w
}

// NewEntity returns a new or recycled [Entity].
// The given component types are added to the entity.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
//
// Note that calling a method with varargs in Go causes a slice allocation.
// For maximum performance, pre-allocate a slice of component IDs and pass it using ellipsis:
//
//	// fast
//	world.NewEntity(idA, idB, idC)
//	// even faster
//	world.NewEntity(ids...)
func (w *World) NewEntity(comps ...ID) Entity {
	w.checkLocked()

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil)
	}

	entity := w.createEntity(arch)

	if w.listener != nil {
		w.listener(&EntityEvent{entity, Mask{}, arch.Mask, comps, nil, arch.Ids, 1})
	}
	return entity
}

// NewEntityWith returns a new or recycled [Entity].
// The given component values are assigned to the entity.
//
// The components in the `Comp` field of [Component] must be pointers.
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
	arch = w.findOrCreateArchetype(arch, ids, nil)

	entity := w.createEntity(arch)

	for _, c := range comps {
		w.copyTo(entity, c.ID, c.Comp)
	}

	if w.listener != nil {
		w.listener(&EntityEvent{entity, Mask{}, arch.Mask, ids, nil, arch.Ids, 1})
	}
	return entity
}

// Creates new entities without returning a query over them.
// Used via [World.Batch].
func (w *World) newEntities(count int, comps ...ID) (*archetype, uint32) {
	arch, startIdx := w.newEntitiesNoNotify(count, comps...)

	if w.listener != nil {
		cnt := uint32(count)
		var i uint32
		for i = 0; i < cnt; i++ {
			idx := startIdx + i
			entity := arch.GetEntity(uintptr(idx))
			w.listener(&EntityEvent{entity, Mask{}, arch.Mask, comps, nil, arch.Ids, 1})
		}
	}

	return arch, startIdx
}

// Creates new entities and returns a query over them.
// Used via [World.Batch].
func (w *World) newEntitiesQuery(count int, comps ...ID) Query {
	arch, startIdx := w.newEntitiesNoNotify(count, comps...)
	lock := w.lock()
	return newArchQuery(w, lock, arch, startIdx)
}

// Creates new entities with component values without returning a query over them.
// Used via [World.Batch].
func (w *World) newEntitiesWith(count int, comps ...Component) (*archetype, uint32) {
	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}

	arch, startIdx := w.newEntitiesWithNoNotify(count, ids, comps...)

	if w.listener != nil {
		var i uint32
		cnt := uint32(count)
		for i = 0; i < cnt; i++ {
			idx := startIdx + i
			entity := arch.GetEntity(uintptr(idx))
			w.listener(&EntityEvent{entity, Mask{}, arch.Mask, ids, nil, arch.Ids, 1})
		}
	}

	return arch, startIdx
}

// Creates new entities with component values and returns a query over them.
// Used via [World.Batch].
func (w *World) newEntitiesWithQuery(count int, comps ...Component) Query {
	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}

	arch, startIdx := w.newEntitiesWithNoNotify(count, ids, comps...)
	lock := w.lock()
	return newArchQuery(w, lock, arch, startIdx)
}

// RemoveEntity removes and recycles an [Entity].
//
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
func (w *World) RemoveEntity(entity Entity) {
	w.checkLocked()

	index := &w.entities[entity.id]
	oldArch := index.arch

	if w.listener != nil {
		lock := w.lock()
		w.listener(&EntityEvent{entity, oldArch.Mask, Mask{}, nil, oldArch.Ids, nil, -1})
		w.unlock(lock)
	}

	swapped := oldArch.Remove(index.index)

	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := oldArch.GetEntity(index.index)
		w.entities[swapEntity.id].index = index.index
	}

	index.arch = nil
}

// removeEntities removes and recycles all entities matching a filter.
//
// Returns the number of removed entities.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
func (w *World) removeEntities(filter Filter) int {
	w.checkLocked()

	lock := w.lock()
	numArches := w.archetypes.Len()
	var count uintptr
	for i := 0; i < numArches; i++ {
		arch := w.archetypes.Get(i)

		if !filter.Matches(arch.Mask) {
			continue
		}

		len := uintptr(arch.Len())
		count += len

		var j uintptr
		for j = 0; j < len; j++ {
			entity := arch.GetEntity(j)
			if w.listener != nil {
				w.listener(&EntityEvent{entity, arch.Mask, Mask{}, nil, arch.Ids, nil, -1})
			}
			w.entities[entity.id].arch = nil
			w.entityPool.Recycle(entity)
		}

		arch.Reset()
	}
	w.unlock(lock)

	return int(count)
}

// Alive reports whether an entity is still alive.
func (w *World) Alive(entity Entity) bool {
	return w.entityPool.Alive(entity)
}

// Get returns a pointer th the given component of an [Entity].
//
// Returns `nil` if the entity has no such component.
// Panics when called for an already removed entity.
//
// See also [github.com/mlange-42/arche/generic.Map.Get] for a generic variant.
func (w *World) Get(entity Entity, comp ID) unsafe.Pointer {
	index := &w.entities[entity.id]
	return index.arch.Get(index.index, comp)
}

// Has returns whether an [Entity] has a given component.
//
// Panics when called for an already removed entity.
//
// See also [github.com/mlange-42/arche/generic.Map.Has] for a generic variant.
func (w *World) Has(entity Entity, comp ID) bool {
	return w.entities[entity.id].arch.HasComponent(comp)
}

// Add adds components to an [Entity].
//
// Panics when called with component that can't be added because they are already present.
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
//
// Note that calling a method with varargs in Go causes a slice allocation.
// For maximum performance, pre-allocate a slice of component IDs and pass it using ellipsis:
//
//	// fast
//	world.Add(entity, idA, idB, idC)
//	// even faster
//	world.Add(entity, ids...)
func (w *World) Add(entity Entity, comps ...ID) {
	w.Exchange(entity, comps, nil)
}

// Assign assigns multiple components to an [Entity], using pointers for the content.
//
// The components in the `Comp` field of [Component] must be pointers.
// The passed pointers are no valid references to the assigned memory!
//
// Panics when called with components that can't be added because they are already present.
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
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
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
//
// Panics if the entity does not have a component of that type.
//
// See also [github.com/mlange-42/arche/generic.Map.Set] for a generic variant.
func (w *World) Set(entity Entity, id ID, comp interface{}) unsafe.Pointer {
	return w.copyTo(entity, id, comp)
}

// Remove removes components from an entity.
//
// Panics when called with components that can't be removed because they are not present.
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
func (w *World) Remove(entity Entity, comps ...ID) {
	w.Exchange(entity, nil, comps)
}

// Exchange adds and removes components in one pass.
//
// Panics when called with components that can't be added or removed because
// they are already present/not present, respectively.
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Exchange].
func (w *World) Exchange(entity Entity, add []ID, rem []ID) {
	w.checkLocked()

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

	arch := w.findOrCreateArchetype(oldArch, add, rem)
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
	w.entities[entity.id] = entityIndex{arch: arch, index: newIndex}

	if w.listener != nil {
		w.listener(&EntityEvent{entity, oldMask, arch.Mask, add, rem, arch.Ids, 0})
	}
}

// Reset removes all entities and resources from the world.
//
// Does NOT free reserved memory, remove archetypes, clear the registry etc.
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
	for i := 0; i < len; i++ {
		w.archetypes.Get(i).Reset()
	}
}

// Query creates a [Query] iterator.
//
// The [ecs] core package provides only the filter [All] for querying the given components.
// Further, it can be chained with [Mask.Without] (see the examples) to exclude components.
//
// Example:
//
//	filter := ecs.All(idA, idB).Without(idC)
//	query := world.Query(filter)
//	for query.Next() {
//	    pos := (*position)(query.Get(posID))
//	    pos.X += 1.0
//	}
//
// For type-safe generics queries, see package [github.com/mlange-42/arche/generic].
// For advanced filtering, see package [github.com/mlange-42/arche/filter].
//
// Locks the world to prevent changes to component compositions.
func (w *World) Query(filter Filter) Query {
	l := w.lock()
	if cached, ok := filter.(*CachedFilter); ok {
		archetypes := &w.filterCache.get(cached).Archetypes
		return newQuery(w, cached, l, archetypes)
	}

	return newQuery(w, filter, l, &w.archetypes)
}

// Resources of the world.
//
// Resources are component-like data that is not associated to an entity, but unique to the world.
func (w *World) Resources() *Resources {
	return &w.resources
}

// Batch creates a [Batch] processing helper.
//
// It provides the functionality to create and remove large numbers of entities in batches,
// in a more efficient way.
func (w *World) Batch() *Batch {
	return &Batch{w}
}

// Cache returns the [Cache] of the world, for registering filters.
//
// See [Cache] for details.
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

// Mask returns the archetype [BitMask] for the given [Entity].
//
// Can be used for fast checks of the entity composition, e.g. using a [Filter].
func (w *World) Mask(entity Entity) Mask {
	return w.entities[entity.id].arch.Mask
}

// SetListener sets a listener callback func(e EntityEvent) for the world.
// The listener is immediately called on every [ecs.Entity] change.
// Replaces the current listener. Call with `nil` to remove a listener.
//
// For details, see [EntityEvent].
func (w *World) SetListener(listener func(e *EntityEvent)) {
	w.listener = listener
}

// Stats reports statistics for inspecting the World.
func (w *World) Stats() *stats.WorldStats {
	entities := stats.EntityStats{
		Used:     w.entityPool.Len(),
		Total:    w.entityPool.Cap(),
		Recycled: w.entityPool.Available(),
		Capacity: w.entityPool.TotalCap(),
	}

	compCount := len(w.registry.Components)
	types := append([]reflect.Type{}, w.registry.Types[:compCount]...)

	memory := cap(w.entities)*int(entityIndexSize) + w.entityPool.TotalCap()*int(entitySize)
	archetypes := make([]stats.ArchetypeStats, w.archetypes.Len())
	for i := 0; i < w.archetypes.Len(); i++ {
		archetypes[i] = w.archetypes.Get(i).Stats(&w.registry)
		memory += archetypes[i].Memory
	}

	return &stats.WorldStats{
		Entities:       entities,
		ComponentCount: compCount,
		ComponentTypes: types,
		Locked:         w.IsLocked(),
		Archetypes:     archetypes,
		Memory:         memory,
	}
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
func (w *World) newEntitiesNoNotify(count int, comps ...ID) (*archetype, uint32) {
	w.checkLocked()

	if count < 1 {
		panic("can only create a positive number of entities")
	}

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil)
	}
	startIdx := arch.Len()
	w.createEntities(arch, uint32(count))

	return arch, startIdx
}

// Internal method to create new entities with component values.
func (w *World) newEntitiesWithNoNotify(count int, ids []ID, comps ...Component) (*archetype, uint32) {
	w.checkLocked()

	if count < 1 {
		panic("can only create a positive number of entities")
	}
	if len(comps) == 0 {
		return w.newEntitiesNoNotify(count)
	}

	cnt := uint32(count)

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, ids, nil)
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
func (w *World) findOrCreateArchetype(start *archetype, add []ID, rem []ID) *archetype {
	curr := start.graphNode
	mask := start.Mask
	for _, id := range rem {
		mask.Set(id, false)
		if next, ok := curr.GetTransitionRemove(id); ok {
			curr = next
		} else {
			next, _ := w.findOrCreateArchetypeSlow(mask)
			next.SetTransitionAdd(id, curr)
			curr.SetTransitionRemove(id, next)
			curr = next
		}
	}
	for _, id := range add {
		mask.Set(id, true)
		if next, ok := curr.GetTransitionAdd(id); ok {
			curr = next
		} else {
			next, _ := w.findOrCreateArchetypeSlow(mask)
			next.SetTransitionRemove(id, curr)
			curr.SetTransitionAdd(id, next)
			curr = next
		}
	}
	if curr.archetype == nil {
		w.createArchetype(curr, true)
	}
	return curr.archetype
}

// Tries to find an archetype for a mask, when it can't be reached through the archetype graph.
// Creates an archetype graph node.
func (w *World) findOrCreateArchetypeSlow(mask Mask) (*archetypeNode, bool) {
	if arch, ok := w.findArchetypeSlow(mask); ok {
		return arch, false
	}
	return w.createArchetypeNode(mask), true
}

// Searches for an archetype by a mask.
func (w *World) findArchetypeSlow(mask Mask) (*archetypeNode, bool) {
	length := w.graph.Len()
	for i := 0; i < length; i++ {
		arch := w.graph.Get(i)
		if arch.mask == mask {
			return arch, true
		}
	}
	return nil, false
}

// Creates a node in the archetype graph.
func (w *World) createArchetypeNode(mask Mask) *archetypeNode {
	w.graph.Add(newArchetypeNode(mask))
	node := w.graph.Get(w.graph.Len() - 1)
	return node
}

// Creates an archetype for the given archetype graph node.
// Initializes the archetype with a capacity according to CapacityIncrement if forStorage is true,
// and with a capacity of 1 otherwise.
func (w *World) createArchetype(node *archetypeNode, forStorage bool) *archetype {
	mask := node.mask
	count := int(mask.TotalBitsSet())
	types := make([]componentType, count)

	idx := 0
	for i := 0; i < MaskTotalBits; i++ {
		id := ID(i)
		if mask.Get(id) {
			types[idx] = componentType{ID: id, Type: w.registry.Types[id]}
			idx++
		}
	}

	w.archetypes.Add(archetype{})
	arch := w.archetypes.Get(w.archetypes.Len() - 1)
	arch.Init(node, w.config.CapacityIncrement, forStorage, types...)
	node.archetype = arch

	w.filterCache.addArchetype(arch)
	return arch
}

// Returns all archetypes that match the given filter. Used by [Cache].
func (w *World) getArchetypes(filter Filter) pagedPointerArr32[archetype] {
	arches := pagedPointerArr32[archetype]{}
	len := int(w.archetypes.Len())
	for i := 0; i < len; i++ {
		a := w.archetypes.Get(i)
		if filter.Matches(a.Mask) {
			arches.Add(a)
		}
	}
	return arches
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
	query.index = -2
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
	event := EntityEvent{Entity{}, Mask{}, arch.Mask, arch.Ids, nil, arch.Ids, 1}

	for i = uintptr(batchArch.StartIndex); i < len; i++ {
		entity := arch.GetEntity(i)
		event.Entity = entity
		w.listener(&event)
	}
}
