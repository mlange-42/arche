package ecs

import (
	"reflect"
	"unsafe"

	"github.com/mlange-42/arche/ecs/stats"
)

// ComponentID returns the ID for a component type via generics. Registers the type if it is not already registered.
func ComponentID[T any](w *World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.componentID(tp)
}

// TypeID returns the ID for a component type. Registers the type if it is not already registered.
func TypeID(w *World, tp reflect.Type) ID {
	return w.componentID(tp)
}

// World is the central type holding [Entity] and component data.
type World struct {
	config     Config
	entities   []entityIndex
	archetypes pagedArr32[archetype]
	entityPool entityPool
	bitPool    bitPool
	registry   componentRegistry
	locks      BitMask
	listener   func(e EntityEvent)
}

// NewWorld creates a new [World] from an optional [Config].
//
// Uses the default [Config] if called without an argument.
// Accepts maximum one argument.
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
		archetypes: pagedArr32[archetype]{},
		locks:      BitMask{},
		listener:   nil,
	}
	w.createArchetype(BitMask{}, false)
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

	entity := w.entityPool.Get()

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil)
	}

	idx := arch.Alloc(entity, true)
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

	if w.listener != nil {
		w.listener(EntityEvent{entity, BitMask{}, arch.Mask, comps, nil, arch.Ids, 1})
	}
	return entity
}

// NewEntityWith returns a new or recycled [Entity].
// The given component values are assigned to the entity.
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

	entity := w.entityPool.Get()
	arch := w.archetypes.Get(0)
	arch = w.findOrCreateArchetype(arch, ids, nil)

	idx := arch.Alloc(entity, false)
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

	for _, c := range comps {
		w.copyTo(entity, c.ID, c.Comp)
	}

	if w.listener != nil {
		w.listener(EntityEvent{entity, BitMask{}, arch.Mask, ids, nil, arch.Ids, 1})
	}
	return entity
}

// RemoveEntity removes and recycles an [Entity].
//
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
func (w *World) RemoveEntity(entity Entity) {
	w.checkLocked()

	index := w.entities[entity.id]
	oldArch := index.arch

	if w.listener != nil {
		w.listener(EntityEvent{entity, oldArch.Mask, oldArch.Mask, nil, nil, oldArch.Ids, -1})
	}

	swapped := oldArch.Remove(index.index)

	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := oldArch.GetEntity(index.index)
		w.entities[swapEntity.id].index = index.index
	}

	w.entities[entity.id].arch = nil
}

// Query creates a [Query] iterator.
//
// The [ecs] core package provides only the filter [All] for querying the given components.
// Further, it can be chained with [Mask.Without] (see the examples) to exclude components.
//
// Example:
//
//	query := world.Query(All(idA, idB).Without(idC))
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
	lock := w.bitPool.Get()
	w.locks.Set(ID(lock), true)
	return newQuery(w, filter, lock)
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

// Mask returns the archetype [BitMask] for the given [Entity].
//
// Can be used for fast checks of the entity composition, e.g. using a [Filter].
func (w *World) Mask(entity Entity) BitMask {
	return w.entities[entity.id].arch.Mask
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
// The passed components must be pointers.
// The passed in pointers are no valid references to the assigned memory!
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

// Exchange adds and removes components in one pass
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
	index := w.entities[entity.id]
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

// IsLocked returns whether the world is locked by any queries.
func (w *World) IsLocked() bool {
	return !w.locks.IsZero()
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

func (w *World) copyTo(entity Entity, id ID, comp interface{}) unsafe.Pointer {
	if !w.Has(entity, id) {
		panic("can't copy component into entity that has no such component type")
	}
	index := w.entities[entity.id]
	arch := index.arch

	return arch.Set(index.index, id, comp)
}

func (w *World) findOrCreateArchetype(start *archetype, add []ID, rem []ID) *archetype {
	curr := start
	mask := start.Mask
	maxRem := len(rem) - 1
	maxAdd := len(add) - 1
	for i, id := range rem {
		mask.Set(id, false)
		if next, ok := curr.GetTransitionRemove(id); ok {
			curr = next
		} else {
			next, _ := w.findOrCreateArchetypeSlow(mask, i == maxRem && maxAdd == 0)
			next.SetTransitionAdd(id, curr)
			curr.SetTransitionRemove(id, next)
			curr = next
		}
	}
	for i, id := range add {
		mask.Set(id, true)
		if next, ok := curr.GetTransitionAdd(id); ok {
			curr = next
		} else {
			next, _ := w.findOrCreateArchetypeSlow(mask, i == maxAdd)
			next.SetTransitionRemove(id, curr)
			curr.SetTransitionAdd(id, next)
			curr = next
		}
	}
	return curr
}

func (w *World) findOrCreateArchetypeSlow(mask BitMask, forStorage bool) (*archetype, bool) {
	if arch, ok := w.findArchetype(mask); ok {
		return arch, false
	}
	return w.createArchetype(mask, forStorage), true
}

func (w *World) findArchetype(mask BitMask) (*archetype, bool) {
	length := w.archetypes.Len()
	for i := 0; i < length; i++ {
		arch := w.archetypes.Get(i)
		if arch.Mask == mask {
			return arch, true
		}
	}
	return nil, false
}

func (w *World) createArchetype(mask BitMask, forStorage bool) *archetype {
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
	arch.Init(w.config.CapacityIncrement, forStorage, types...)
	return arch
}

// componentID returns the ID for a component type, and registers it if not already registered.
func (w *World) componentID(tp reflect.Type) ID {
	return w.registry.ComponentID(tp)
}

func (w *World) nextArchetype(filter Filter, index int) (int, archetypeIter, bool) {
	len := int(w.archetypes.Len())
	if index >= len {
		panic("exceeded end of query")
	}
	for i := index + 1; i < len; i++ {
		a := w.archetypes.Get(i)
		if a.Len() > 0 && filter.Matches(a.Mask) {
			return i, newArchetypeIter(a), true
		}
	}
	return len, archetypeIter{}, false
}

// closeQuery closes a query and unlocks the world
func (w *World) closeQuery(query *queryIter) {
	l := query.lockBit
	if !w.locks.Get(ID(l)) {
		panic("unbalanced query unlock")
	}
	w.locks.Set(ID(l), false)
	w.bitPool.Recycle(l)
}

func (w *World) checkLocked() {
	if !w.locks.IsZero() {
		panic("attempt to modify a locked world")
	}
}
