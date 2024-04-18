package ecs

import (
	"reflect"
	"unsafe"

	"github.com/mlange-42/arche/ecs/stats"
)

// World is the central type holding entity and component data, as well as resources.
//
// The World provides all the basic ECS functionality of Arche,
// like [World.Query], [World.NewEntity], [World.Add], [World.Remove] or [World.RemoveEntity].
//
// For more advanced functionality, see [World.Relations], [World.Resources],
// [World.Batch], [World.Cache] and [Builder].
type World struct {
	listener       Listener                  // EntityEvent listener.
	nodePointers   []*archNode               // Helper list of all node pointers for queries.
	entities       []entityIndex             // Mapping from entities to archetype and index.
	targetEntities bitSet                    // Whether entities are potential relation targets. Used for archetype cleanup.
	relationNodes  []*archNode               // Archetype nodes that have an entity relation.
	filterCache    Cache                     // Cache for registered filters.
	nodes          pagedSlice[archNode]      // The archetype graph.
	archetypeData  pagedSlice[archetypeData] // Storage for the actual archetype data (components).
	nodeData       pagedSlice[nodeData]      // The archetype graph's data.
	archetypes     pagedSlice[archetype]     // Archetypes that have no relations components.
	entityPool     entityPool                // Pool for entities.
	stats          stats.World               // Cached world statistics.
	resources      Resources                 // World resources.
	registry       componentRegistry         // Component registry.
	locks          lockMask                  // World locks.
	config         Config                    // World configuration.
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

// NewEntity returns a new or recycled [Entity].
// The given component types are added to the entity.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// ⚠️ Important:
// Entities are intended to be stored and passed around via copy, not via pointers! See [Entity].
//
// Note that calling a method with varargs in Go causes a slice allocation.
// For maximum performance, pre-allocate a slice of component IDs and pass it using ellipsis:
//
//	// fast
//	world.NewEntity(idA, idB, idC)
//	// even faster
//	world.NewEntity(ids...)
//
// For more advanced and batched entity creation, see [Builder].
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
func (w *World) NewEntity(comps ...ID) Entity {
	w.checkLocked()

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil, Entity{})
	}

	entity := w.createEntity(arch)

	if w.listener != nil {
		var newRel *ID
		if arch.HasRelationComponent {
			newRel = &arch.RelationComponent
		}
		bits := subscription(true, false, len(comps) > 0, false, newRel != nil, newRel != nil)
		trigger := w.listener.Subscriptions() & bits
		if trigger != 0 && subscribes(trigger, &arch.Mask, nil, w.listener.Components(), nil, newRel) {
			w.listener.Notify(w, EntityEvent{Entity: entity, Added: arch.Mask, AddedIDs: comps, NewRelation: newRel, EventTypes: bits})
		}
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
// ⚠️ Important:
// Entities are intended to be stored and passed around via copy, not via pointers! See [Entity].
//
// For more advanced and batched entity creation, see [Builder].
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
	arch = w.findOrCreateArchetype(arch, ids, nil, Entity{})

	entity := w.createEntity(arch)

	for _, c := range comps {
		w.copyTo(entity, c.ID, c.Comp)
	}

	if w.listener != nil {
		var newRel *ID
		if arch.HasRelationComponent {
			newRel = &arch.RelationComponent
		}
		bits := subscription(true, false, len(comps) > 0, false, newRel != nil, newRel != nil)
		trigger := w.listener.Subscriptions() & bits
		if trigger != 0 && subscribes(trigger, &arch.Mask, nil, w.listener.Components(), nil, newRel) {
			w.listener.Notify(w, EntityEvent{Entity: entity, Added: arch.Mask, AddedIDs: ids, NewRelation: newRel, EventTypes: bits})
		}
	}
	return entity
}

// RemoveEntity removes an [Entity], making it eligible for recycling.
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
		var oldRel *ID
		if oldArch.HasRelationComponent {
			oldRel = &oldArch.RelationComponent
		}
		var oldIds []ID
		if len(oldArch.node.Ids) > 0 {
			oldIds = oldArch.node.Ids
		}

		bits := subscription(false, true, false, len(oldIds) > 0, oldRel != nil, oldRel != nil)
		trigger := w.listener.Subscriptions() & bits
		if trigger != 0 && subscribes(trigger, nil, &oldArch.Mask, w.listener.Components(), oldRel, nil) {
			lock := w.lock()
			w.listener.Notify(w, EntityEvent{Entity: entity, Removed: oldArch.Mask, RemovedIDs: oldIds, OldRelation: oldRel, OldTarget: oldArch.RelationTarget, EventTypes: bits})
			w.unlock(lock)
		}
	}

	swapped := oldArch.Remove(index.index)

	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := oldArch.GetEntity(index.index)
		w.entities[swapEntity.id].index = index.index
	}
	index.arch = nil

	if w.targetEntities.Get(entity.id) {
		w.cleanupArchetypes(entity)
		w.targetEntities.Set(entity.id, false)
	}

	w.cleanupArchetype(oldArch)
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
// See also [World.Exchange].
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
	w.assign(entity, ID{}, false, Entity{}, comps...)
}

// Set overwrites a component for an [Entity], using the given pointer for the content.
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
// See also [World.Exchange].
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
func (w *World) Remove(entity Entity, comps ...ID) {
	w.Exchange(entity, nil, comps)
}

// Exchange adds and removes components in one pass.
// This is more efficient than subsequent use of [World.Add] and [World.Remove].
//
// When a [Relation] component is removed and another one is added,
// the target entity of the relation is reset to zero.
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Relations.Exchange] and the generic variants under [github.com/mlange-42/arche/generic.Exchange].
func (w *World) Exchange(entity Entity, add []ID, rem []ID) {
	w.exchange(entity, add, rem, ID{}, false, Entity{})
}

// Reset removes all entities and resources from the world.
//
// Does NOT free reserved memory, remove archetypes, clear the registry, clear cached filters, etc.
// However, it removes archetypes with a relation component that is not zero.
//
// Can be used to run systematic simulations without the need to re-allocate memory for each run.
// Accelerates re-populating the world by a factor of 2-3.
func (w *World) Reset() {
	w.checkLocked()

	w.entities = w.entities[:1]
	w.targetEntities.Reset()
	w.entityPool.Reset()
	w.locks.Reset()
	w.resources.reset()

	len := w.nodes.Len()
	var i int32
	for i = 0; i < len; i++ {
		w.nodes.Get(i).Reset(w.Cache())
	}
}

// Query creates a [Query] iterator.
//
// Locks the world to prevent changes to component compositions.
// The lock is released automatically when the query finishes iteration, or when [Query.Close] is called.
// The number of simultaneous locks (and thus open queries) at a given time is limited to [MaskTotalBits] (256).
//
// A query can iterate through its entities only once, and can't be used anymore afterwards.
//
// To create a [Filter] for querying, see [All], [Mask.Without], [Mask.Exclusive] and [RelationFilter].
//
// For type-safe generics queries, see package [github.com/mlange-42/arche/generic].
// For advanced filtering, see package [github.com/mlange-42/arche/filter].
func (w *World) Query(filter Filter) Query {
	l := w.lock()
	if cached, ok := filter.(*CachedFilter); ok {
		return newCachedQuery(w, cached.filter, l, w.filterCache.get(cached).Archetypes.pointers)
	}

	return newQuery(w, filter, l, w.nodePointers)
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

// Batch creates a [Batch] processing helper.
// It provides the functionality to manipulate large numbers of entities in batches,
// which is more efficient than handling them one by one.
func (w *World) Batch() *Batch {
	return &Batch{w}
}

// Relations returns the [Relations] of the world, for accessing entity [Relation] targets.
//
// See [Relations] for details.
func (w *World) Relations() *Relations {
	return &Relations{world: w}
}

// IsLocked returns whether the world is locked by any queries.
func (w *World) IsLocked() bool {
	return w.locks.IsLocked()
}

// Mask returns the archetype [Mask] for the given [Entity].
func (w *World) Mask(entity Entity) Mask {
	if !w.entityPool.Alive(entity) {
		panic("can't get mask for a dead entity")
	}
	return w.entities[entity.id].arch.Mask
}

// Ids returns the component IDs for the archetype of the given [Entity].
//
// Returns a copy of the archetype's component IDs slice, for safety.
// This means that the result can be manipulated safely,
// but also that calling the method may incur some significant cost.
func (w *World) Ids(entity Entity) []ID {
	if !w.entityPool.Alive(entity) {
		panic("can't get component IDs for a dead entity")
	}
	return append([]ID{}, w.entities[entity.id].arch.node.Ids...)
}

// SetListener sets a [Listener] for the world.
// The listener is immediately called on every [ecs.Entity] change.
// Replaces the current listener. Call with nil to remove a listener.
//
// For details, see [EntityEvent], [Listener] and sub-package [event].
func (w *World) SetListener(listener Listener) {
	w.listener = listener
}

// Stats reports statistics for inspecting the World.
//
// The underlying [stats.World] object is re-used and updated between calls.
// The returned pointer should thus not be stored for later analysis.
// Rather, the required data should be extracted immediately.
func (w *World) Stats() *stats.World {
	w.stats.Entities = stats.Entities{
		Used:     w.entityPool.Len(),
		Total:    w.entityPool.Cap(),
		Recycled: w.entityPool.Available(),
		Capacity: w.entityPool.TotalCap(),
	}

	compCount := len(w.registry.Components)
	types := append([]reflect.Type{}, w.registry.Types[:compCount]...)

	memory := cap(w.entities)*int(entityIndexSize) + w.entityPool.TotalCap()*int(entitySize)

	cntOld := int32(len(w.stats.Nodes))
	cntNew := int32(w.nodes.Len())
	cntActive := 0
	var i int32
	for i = 0; i < cntOld; i++ {
		node := w.nodes.Get(i)
		nodeStats := &w.stats.Nodes[i]
		node.UpdateStats(nodeStats, &w.registry)
		if node.IsActive {
			memory += nodeStats.Memory
			cntActive++
		}
	}
	for i = cntOld; i < cntNew; i++ {
		node := w.nodes.Get(i)
		w.stats.Nodes = append(w.stats.Nodes, node.Stats(&w.registry))
		if node.IsActive {
			memory += w.stats.Nodes[i].Memory
			cntActive++
		}
	}

	w.stats.ComponentCount = compCount
	w.stats.ComponentTypes = types
	w.stats.Locked = w.IsLocked()
	w.stats.Memory = memory
	w.stats.CachedFilters = len(w.filterCache.filters)
	w.stats.ActiveNodeCount = cntActive

	return &w.stats
}

// DumpEntities dumps entity information into an [EntityDump] object.
// This dump can be used with [World.LoadEntities] to set the World's entity state.
//
// For world serialization with components and resources, see module [github.com/mlange-42/arche-serde].
func (w *World) DumpEntities() EntityDump {
	alive := []uint32{}

	query := w.Query(All())
	for query.Next() {
		alive = append(alive, uint32(query.Entity().id))
	}

	data := EntityDump{
		Entities:  append([]Entity{}, w.entityPool.entities...),
		Alive:     alive,
		Next:      uint32(w.entityPool.next),
		Available: w.entityPool.available,
	}

	return data
}

// LoadEntities resets all entities to the state saved with [World.DumpEntities].
//
// Use this only on an empty world! Can be used after [World.Reset].
//
// The resulting world will have the same entities (in terms of ID, generation and alive state)
// as the original world. This is necessary for proper serialization of entity relations.
// However, the entities will not have any components.
//
// Panics if the world has any dead or alive entities.
//
// For world serialization with components and resources, see module [github.com/mlange-42/arche-serde].
func (w *World) LoadEntities(data *EntityDump) {
	w.checkLocked()

	if len(w.entityPool.entities) > 1 || w.entityPool.available > 0 {
		panic("can set entity data only on a fresh or reset world")
	}

	capacity := capacity(len(data.Entities), w.config.CapacityIncrement)

	entities := make([]Entity, 0, capacity)
	entities = append(entities, data.Entities...)

	w.entityPool.entities = entities
	w.entityPool.next = eid(data.Next)
	w.entityPool.available = data.Available

	w.entities = make([]entityIndex, len(data.Entities), capacity)
	w.targetEntities = bitSet{}
	w.targetEntities.ExtendTo(capacity)

	arch := w.archetypes.Get(0)
	for _, idx := range data.Alive {
		entity := w.entityPool.entities[idx]
		archIdx := arch.Alloc(entity)
		w.entities[entity.id] = entityIndex{arch: arch, index: archIdx}
	}
}
