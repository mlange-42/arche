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

// GetResource returns a pointer to the given resource type.
//
// Returns nil if there is no such resource.
//
// Uses reflection. For more efficient access, see [World.GetResource],
// and [github.com/mlange-42/arche/generic.Resource.Get] for a generic variant.
// These methods are more than 20 times faster than the GetResource function.
func GetResource[T any](w *World) *T {
	return w.GetResource(ResourceID[T](w)).(*T)
}

// AddResource adds a resource to the world.
//
// Panics if there is already such a resource.
//
// Uses reflection. For more efficient access, see [World.AddResource],
// and [github.com/mlange-42/arche/generic.Resource.Add] for a generic variant.
func AddResource[T any](w *World, res *T) {
	w.AddResource(ResourceID[T](w), res)
}

// World is the central type holding [Entity] and component data, as well as resources.
type World struct {
	config     Config
	entities   []entityIndex
	archetypes archetypeArr
	graph      pagedArr32[archetypeNode]
	entityPool entityPool
	bitPool    bitPool
	registry   componentRegistry[ID]
	locks      Mask
	listener   func(e EntityEvent)
	resources  resources
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
		config:     conf,
		entities:   entities,
		entityPool: newEntityPool(conf.CapacityIncrement),
		bitPool:    newBitPool(),
		registry:   newComponentRegistry(),
		archetypes: archetypeArr{},
		graph:      pagedArr32[archetypeNode]{},
		locks:      Mask{},
		listener:   nil,
		resources:  newResources(),
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
func (w *World) NewEntity(comps ...ID) Entity {
	w.checkLocked()

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil)
	}

	entity := w.createEntity(arch, true)

	if w.listener != nil {
		w.listener(EntityEvent{entity, Mask{}, arch.Mask, comps, nil, arch.Ids, 1})
	}
	return entity
}

func (w *World) newEntities(count int, comps ...ID) Query {
	w.checkLocked()

	if count < 1 {
		panic("can only create a positive number of entities")
	}
	cnt := uint32(count)

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil)
	}
	startIdx := arch.Len()
	w.createEntities(arch, uint32(count), true)

	if w.listener != nil {
		lock := w.lock()
		var i uint32
		for i = 0; i < cnt; i++ {
			idx := startIdx + i
			entity := arch.GetEntity(uintptr(idx))
			w.listener(EntityEvent{entity, Mask{}, arch.Mask, comps, nil, arch.Ids, 1})
		}
		w.unlock(lock)
	}

	lock := w.lock()
	return newArchQuery(w, lock, arch, startIdx)
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

	entity := w.createEntity(arch, false)

	for _, c := range comps {
		w.copyTo(entity, c.ID, c.Comp)
	}

	if w.listener != nil {
		w.listener(EntityEvent{entity, Mask{}, arch.Mask, ids, nil, arch.Ids, 1})
	}
	return entity
}

func (w *World) newEntitiesWith(count int, comps ...Component) Query {
	w.checkLocked()

	if count < 1 {
		panic("can only create a positive number of entities")
	}
	if len(comps) == 0 {
		return w.newEntities(count)
	}

	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}

	cnt := uint32(count)

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, ids, nil)
	}
	startIdx := arch.Len()
	w.createEntities(arch, uint32(count), true)

	var i uint32
	for i = 0; i < cnt; i++ {
		idx := startIdx + i
		entity := arch.GetEntity(uintptr(idx))
		for _, c := range comps {
			w.copyTo(entity, c.ID, c.Comp)
		}
	}

	if w.listener != nil {
		lock := w.lock()
		var i uint32
		for i = 0; i < cnt; i++ {
			idx := startIdx + i
			entity := arch.GetEntity(uintptr(idx))
			w.listener(EntityEvent{entity, Mask{}, arch.Mask, ids, nil, arch.Ids, 1})
		}
		w.unlock(lock)
	}

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
		w.listener(EntityEvent{entity, oldArch.Mask, Mask{}, nil, oldArch.Ids, nil, -1})
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
// Panics when called on a locked world.
// Do not use during [Query] iteration!
func (w *World) removeEntities(filter Filter) {
	w.checkLocked()

	lock := w.lock()
	for i := 0; i < w.archetypes.Len(); i++ {
		arch := w.archetypes.Get(i)

		if !filter.Matches(arch.Mask) {
			continue
		}

		len := uintptr(arch.Len())
		var j uintptr
		for j = 0; j < len; j++ {
			entity := arch.GetEntity(j)
			if w.listener != nil {
				w.listener(EntityEvent{entity, arch.Mask, Mask{}, nil, arch.Ids, nil, -1})
			}
			w.entities[entity.id].arch = nil
			w.entityPool.Recycle(entity)
		}

		arch.Reset()
	}
	w.unlock(lock)
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
	newIndex := arch.Alloc(entity, false)

	for _, id := range oldIDs {
		if mask.Get(id) {
			comp := oldArch.Get(index.index, id)
			arch.SetPointer(newIndex, id, comp)
		}
	}
	for _, id := range add {
		arch.Zero(newIndex, id)
	}

	swapped := oldArch.Remove(index.index)

	if swapped {
		swapEntity := oldArch.GetEntity(index.index)
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{arch: arch, index: newIndex}

	if w.listener != nil {
		w.listener(EntityEvent{entity, oldMask, arch.Mask, add, rem, arch.Ids, 0})
	}
}

// AddResource adds a resource to the world.
// The resource should always be a pointer.
//
// Panics if there is already a resource of the given type.
//
// See also [github.com/mlange-42/arche/generic.Resource.Add] for a generic variant.
func (w *World) AddResource(id ResID, res any) {
	w.resources.Add(id, res)
}

// RemoveResource removes a resource from the world.
// The resource should always be a pointer.
//
// Panics if there is no resource of the given type.
//
// See also [github.com/mlange-42/arche/generic.Resource.Remove] for a generic variant.
func (w *World) RemoveResource(id ResID) {
	w.resources.Remove(id)
}

// GetResource returns a pointer to the given resource type.
//
// Returns nil if there is no such resource.
//
// See also [github.com/mlange-42/arche/generic.Resource.Get] for a generic variant.
func (w *World) GetResource(id ResID) interface{} {
	return w.resources.Get(id)
}

// HasResource returns whether the world has the given resource type.
//
// See also [github.com/mlange-42/arche/generic.Resource.Has] for a generic variant.
func (w *World) HasResource(id ResID) bool {
	return w.resources.Has(id)
}

// Reset removes all entities and resources as well as the listener from the world.
//
// Does NOT free reserved memory, remove archetypes, clear the registry etc.
//
// Can be used to run systematic simulations without the need to re-allocate memory for each run.
// Accelerates re-populating the world by a factor of 2-3.
func (w *World) Reset() {
	w.checkLocked()

	w.entities = w.entities[:1]
	w.entityPool.Reset()
	w.bitPool.Reset()
	w.resources.Reset()

	w.listener = nil

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
//	filter := All(idA, idB).Without(idC)
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
	return newQuery(w, filter, l, &w.archetypes)
}

// Batch creates a [Batch] processing helper.
//
// It provides the functionality to create and remove large numbers of entities in batches,
// in a more efficient way.
func (w *World) Batch() Batch {
	return Batch{w}
}

// lock the world and get the lock bit for later unlocking.
func (w *World) lock() uint8 {
	lock := w.bitPool.Get()
	w.locks.Set(ID(lock), true)
	return lock
}

// unlock unlocks the given lock bit.
func (w *World) unlock(l uint8) {
	if !w.locks.Get(ID(l)) {
		panic("unbalanced query unlock")
	}
	w.locks.Set(ID(l), false)
	w.bitPool.Recycle(l)
}

// IsLocked returns whether the world is locked by any queries.
func (w *World) IsLocked() bool {
	return !w.locks.IsZero()
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
// Events notified are entity creation, removal and changes to the component composition.
// Events are emitted immediately after the change is applied.
// Except for removal of an entity, where the event is emitted before removal.
// This allows for inspection of the to-be-removed Entity.
func (w *World) SetListener(listener func(e EntityEvent)) {
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

// createEntity creates an Entity and adds it to the given archetype.
func (w *World) createEntity(arch *archetype, zero bool) Entity {
	entity := w.entityPool.Get()
	idx := arch.Alloc(entity, zero)
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

func (w *World) createEntities(arch *archetype, count uint32, zero bool) {
	startIdx := arch.Len()
	arch.AllocN(uint32(count), zero)

	var i uint32
	for i = 0; i < count; i++ {
		idx := startIdx + i
		entity := w.entityPool.Get()
		arch.SetEntity(uintptr(idx), entity)
		len := len(w.entities)
		if int(entity.id) == len {
			if len == cap(w.entities) {
				old := w.entities
				w.entities = make([]entityIndex, len, len+w.config.CapacityIncrement)
				copy(w.entities, old)
			}
			w.entities = append(w.entities, entityIndex{arch: arch, index: uintptr(idx)})
		} else {
			w.entities[entity.id] = entityIndex{arch: arch, index: uintptr(idx)}
		}
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
	return arch
}

// componentID returns the ID for a component type, and registers it if not already registered.
func (w *World) componentID(tp reflect.Type) ID {
	return w.registry.ComponentID(tp)
}

// resourceID returns the ID for a resource type, and registers it if not already registered.
func (w *World) resourceID(tp reflect.Type) ResID {
	return w.resources.registry.ComponentID(tp)
}

// closeQuery closes a query and unlocks the world
func (w *World) closeQuery(query *Query) {
	query.index = -2
	w.unlock(query.lockBit)
}

// checkLocked checks if the world is locked, and panics if so.
func (w *World) checkLocked() {
	if !w.locks.IsZero() {
		panic("attempt to modify a locked world")
	}
}
