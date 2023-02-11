package ecs

import (
	"internal/base"
	"reflect"
	"unsafe"

	f "github.com/mlange-42/arche/filter"
)

// NewWorld creates a new [World]
func NewWorld() World {
	return FromConfig(NewConfig())
}

// FromConfig creates a new [World] from a [Config]
func FromConfig(conf Config) World {
	arch := base.Archetype{}
	arch.Init(conf.CapacityIncrement)
	arches := base.PagedArr32[base.Archetype]{}
	arches.Add(arch)
	return World{
		config:     conf,
		entities:   []base.EntityIndex{{Arch: nil, Index: 0}},
		entityPool: base.NewEntityPool(conf.CapacityIncrement),
		bitPool:    base.NewBitPool(),
		registry:   base.NewComponentRegistry(),
		archetypes: arches,
		locks:      bitMask(0),
	}
}

// World is the central type holding [Entity] and component data.
type World struct {
	config     Config
	entities   []base.EntityIndex
	archetypes base.PagedArr32[base.Archetype]
	entityPool base.EntityPool
	bitPool    base.BitPool
	registry   base.ComponentRegistry
	locks      bitMask
}

// NewEntity returns a new or recycled [Entity].
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
func (w *World) NewEntity() Entity {
	w.checkLocked()

	entity := w.entityPool.Get()
	idx := w.archetypes.Get(0).Add(entity)
	len := len(w.entities)
	if int(entity.ID) == len {
		if len == cap(w.entities) {
			old := w.entities
			w.entities = make([]base.EntityIndex, len, len+w.config.CapacityIncrement)
			copy(w.entities, old)
		}
		w.entities = append(w.entities, base.EntityIndex{Arch: w.archetypes.Get(0), Index: idx})
	} else {
		w.entities[entity.ID] = base.EntityIndex{Arch: w.archetypes.Get(0), Index: idx}
	}
	return entity
}

// RemEntity removes and recycles an [Entity].
//
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
func (w *World) RemEntity(entity Entity) {
	w.checkLocked()

	index := w.entities[entity.ID]
	oldArch := index.Arch
	swapped := oldArch.Remove(int(index.Index))

	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.Index))
		w.entities[swapEntity.ID].Index = index.Index
	}

	w.entities[entity.ID].Arch = nil
}

// Query creates a [Query] iterator for the given components.
//
// Locks the world to prevent changes to component compositions.
//
// See also the generic alternatives [Query1], [Query2], [Query3], ...
func (w *World) Query(comps ...ID) Query {
	mask := base.NewBitMask(comps...)
	lock := w.bitPool.Get()
	w.locks.Set(ID(lock), true)
	return newQuery(w, mask, 0, lock)
}

func (w *World) query(mask, exclude bitMask) Query {
	lock := w.bitPool.Get()
	w.locks.Set(ID(lock), true)
	return newQuery(w, mask, exclude, lock)
}

// Filter creates an advanced [Filter] iterator.
//
// Locks the world to prevent changes to component compositions.
//
// There is no generic alternative for filters.
func (w *World) Filter(filter f.MaskFilter) Filter {
	lock := w.bitPool.Get()
	w.locks.Set(ID(lock), true)
	return newFilter(w, filter, lock)
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
// See also [Map.Get] for a generic variant.
func (w *World) Get(entity Entity, comp ID) unsafe.Pointer {
	index := w.entities[entity.ID]
	return index.Arch.Get(int(index.Index), comp)
}

// Has returns whether an [Entity] has a given component.
//
// Panics when called for an already removed entity.
//
// See also [Map.Has] for a generic variant.
func (w *World) Has(entity Entity, comp ID) bool {
	return w.entities[entity.ID].Arch.HasComponent(comp)
}

// Add adds components to an [Entity].
//
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
//
// See also the generic variants [Add], [Add2], [Add3], ...
func (w *World) Add(entity Entity, comps ...ID) {
	w.Exchange(entity, comps, []ID{})
}

// Assign assigns a component to an [Entity], using a given pointer for the content.
// See also [World.AssignN].
//
// The passed component must be a pointer.
// Returns a pointer to the assigned memory.
// The passed in pointer is not a valid reference to that memory!
//
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
//
// See also the generic variants [Assign], [Assign2], [Assign3], ...
func (w *World) Assign(entity Entity, id ID, comp interface{}) unsafe.Pointer {
	w.Exchange(entity, []ID{id}, []ID{})
	return w.copyTo(entity, id, comp)
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
// See also [Map.Set] for a generic variant.
func (w *World) Set(entity Entity, id ID, comp interface{}) unsafe.Pointer {
	return w.copyTo(entity, id, comp)
}

// AssignN assigns multiple components to an [Entity], using pointers for the content.
// See also [World.Assign].
//
// The passed components must be pointers.
// The passed in pointers are no valid references to the assigned memory!
//
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
//
// See also the generic variants [Assign], [Assign2], [Assign3], ...
func (w *World) AssignN(entity Entity, comps ...Component) {
	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}
	w.Exchange(entity, ids, []ID{})
	for _, c := range comps {
		w.copyTo(entity, c.ID, c.Component)
	}
}

// Remove removes components from an entity.
//
// Panics when called on a locked world or for an already removed entity.
//
// Do not use during [Query] iteration!
//
// See also the generic variants [Remove], [Remove2], [Remove3], ...
func (w *World) Remove(entity Entity, comps ...ID) {
	w.Exchange(entity, []ID{}, comps)
}

// Exchange adds and removes components in one pass
//
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
func (w *World) Exchange(entity Entity, add []ID, rem []ID) {
	w.checkLocked()

	if len(add) == 0 && len(rem) == 0 {
		return
	}
	index := w.entities[entity.ID]
	oldArch := index.Arch
	mask := oldArch.Mask
	for _, comp := range add {
		if mask.Get(comp) {
			panic("entity already has this component, can't add")
		}
		mask.Set(comp, true)
	}
	for _, comp := range rem {
		if !mask.Get(comp) {
			panic("entity does not have this component, can't remove")
		}
		mask.Set(comp, false)
	}

	oldIDs := oldArch.Components()
	keepIDs := make([]ID, 0, len(oldIDs))
	for _, id := range oldIDs {
		if mask.Get(id) {
			keepIDs = append(keepIDs, id)
		}
	}
	addIDs := make([]ID, 0, len(add))
	for _, id := range add {
		if mask.Get(id) {
			addIDs = append(addIDs, id)
		}
	}

	arch := w.findOrCreateArchetype(oldArch, addIDs, rem)

	allComps := make([]base.ComponentPointer, 0, len(keepIDs)+len(addIDs))
	for _, id := range keepIDs {
		comp := oldArch.Get(int(index.Index), id)
		allComps = append(allComps, base.ComponentPointer{ID: id, Pointer: comp})
	}
	for _, id := range addIDs {
		allComps = append(allComps, base.ComponentPointer{ID: id, Pointer: nil})
	}

	newIndex := arch.AddPointer(entity, allComps...)
	swapped := oldArch.Remove(int(index.Index))

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.Index))
		w.entities[swapEntity.ID].Index = index.Index
	}
	w.entities[entity.ID] = base.EntityIndex{Arch: arch, Index: newIndex}
}

// IsLocked returns whether the world is locked by any queries.
func (w *World) IsLocked() bool {
	return w.locks != 0
}

func (w *World) copyTo(entity Entity, id ID, comp interface{}) unsafe.Pointer {
	if !w.Has(entity, id) {
		panic("can't copy component into entity that has no such component type")
	}
	index := w.entities[entity.ID]
	arch := index.Arch
	return arch.Set(index.Index, id, comp)
}

func (w *World) findOrCreateArchetype(start *base.Archetype, add []ID, rem []ID) *base.Archetype {
	curr := start
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
	return curr
}

func (w *World) findOrCreateArchetypeSlow(mask bitMask) (*base.Archetype, bool) {
	if arch, ok := w.findArchetype(mask); ok {
		return arch, false
	}
	return w.createArchetype(mask), true
}

func (w *World) findArchetype(mask bitMask) (*base.Archetype, bool) {
	length := w.archetypes.Len()
	for i := 0; i < length; i++ {
		arch := w.archetypes.Get(i)
		if arch.Mask == mask {
			return arch, true
		}
	}
	return nil, false
}

func (w *World) createArchetype(mask bitMask) *base.Archetype {
	count := int(mask.TotalBitsSet())
	types := make([]base.ComponentType, count)

	idx := 0
	for i := 0; i < MaskTotalBits; i++ {
		id := ID(i)
		if mask.Get(id) {
			types[idx] = base.ComponentType{ID: id, Type: w.registry.Types[id]}
			idx++
		}
	}
	w.archetypes.Add(base.Archetype{})
	arch := w.archetypes.Get(w.archetypes.Len() - 1)
	arch.Init(w.config.CapacityIncrement, types...)
	return arch
}

// componentID returns the ID for a component type, and registers it if not already registered.
func (w *World) componentID(tp reflect.Type) ID {
	return w.registry.ComponentID(tp)
}

func (w *World) nextArchetype(mask, exclude bitMask, index int) (int, archetypeIter, bool) {
	len := w.archetypes.Len()
	if index >= len {
		panic("exceeded end of query")
	}
	for i := index + 1; i < len; i++ {
		a := w.archetypes.Get(i)
		if a.Len() > 0 && a.Mask.Contains(mask) && !a.Mask.ContainsAny(exclude) {
			return i, newArchetypeIter(a), true
		}
	}
	return len, archetypeIter{}, false
}

func (w *World) nextArchetypeFilter(filter f.MaskFilter, index int) (int, archetypeIter, bool) {
	len := w.archetypes.Len()
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
	if w.locks != 0 {
		panic("attempt to modify a locked world")
	}
}
