package ecs

import (
	"reflect"
	"unsafe"
)

// ComponentID returns the ID for a component type. Registers the type if it is not already registered.
func ComponentID[T any](w *World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.componentID(tp)
}

// NewWorld creates a new [World]
func NewWorld() World {
	return FromConfig(NewConfig())
}

// FromConfig creates a new [World] from a [Config]
func FromConfig(conf Config) World {
	arch := archetype{}
	arch.init(conf.CapacityIncrement)
	arches := pagedArr32[archetype]{}
	arches.Add(arch)
	return World{
		config:     conf,
		entities:   []entityIndex{{nil, 0}},
		entityPool: newEntityPool(conf.CapacityIncrement),
		bitPool:    newBitPool(),
		registry:   newComponentRegistry(),
		archetypes: arches,
		locks:      bitMask(0),
	}
}

// World is the central type holding [Entity] and component data.
type World struct {
	config     Config
	entities   []entityIndex
	archetypes pagedArr32[archetype]
	entityPool entityPool
	bitPool    bitPool
	registry   componentRegistry
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
	if int(entity.id) == len {
		if len == cap(w.entities) {
			old := w.entities
			w.entities = make([]entityIndex, len, len+w.config.CapacityIncrement)
			copy(w.entities, old)
		}
		w.entities = append(w.entities, entityIndex{w.archetypes.Get(0), idx})
	} else {
		w.entities[entity.id] = entityIndex{w.archetypes.Get(0), idx}
	}
	return entity
}

// RemEntity removes and recycles an [Entity].
//
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
func (w *World) RemEntity(entity Entity) {
	w.checkLocked()

	index := w.entities[entity.id]
	oldArch := index.arch
	swapped := oldArch.Remove(int(index.index))

	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}

	w.entities[entity.id].arch = nil
}

// Get returns a pointer th the given component of an [Entity].
//
// Returns `nil` if the entity has no such component.
// Panics when called for an already removed entity.
func (w *World) Get(entity Entity, comp ID) unsafe.Pointer {
	index := w.entities[entity.id]
	arch := index.arch

	if !arch.HasComponent(comp) {
		return nil
	}

	return arch.Get(int(index.index), comp)
}

// Has returns whether an [Entity] has a given component.
//
// Panics when called for an already removed entity.
func (w *World) Has(entity Entity, comp ID) bool {
	index := w.entities[entity.id]
	return index.arch.HasComponent(comp)
}

// Add adds components to an [Entity].
//
// Panics when called on a locked world or for an already removed entity.
// Do not use during [Query] iteration!
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
func (w *World) Assign(entity Entity, id ID, comp interface{}) unsafe.Pointer {
	w.Exchange(entity, []ID{id}, []ID{})
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
	index := w.entities[entity.id]
	oldArch := index.arch
	mask := oldArch.mask
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

	allComps := make([]componentPointer, 0, len(keepIDs)+len(addIDs))
	for _, id := range keepIDs {
		comp := oldArch.Get(int(index.index), id)
		allComps = append(allComps, componentPointer{id, comp})
	}
	for _, id := range addIDs {
		allComps = append(allComps, componentPointer{id, nil})
	}

	newIndex := arch.AddPointer(entity, allComps...)
	swapped := oldArch.Remove(int(index.index))

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{arch, newIndex}
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
	mask := start.mask
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

func (w *World) findOrCreateArchetypeSlow(mask bitMask) (*archetype, bool) {
	if arch, ok := w.findArchetype(mask); ok {
		return arch, false
	}
	return w.createArchetype(mask), true
}

func (w *World) findArchetype(mask bitMask) (*archetype, bool) {
	length := w.archetypes.Len()
	for i := 0; i < length; i++ {
		arch := w.archetypes.Get(i)
		if arch.mask == mask {
			return arch, true
		}
	}
	return nil, false
}

func (w *World) createArchetype(mask bitMask) *archetype {
	count := int(mask.TotalBitsSet())
	types := make([]componentType, count)

	idx := 0
	for i := 0; i < MaskTotalBits; i++ {
		id := ID(i)
		if mask.Get(id) {
			types[idx] = componentType{id, w.registry.types[id]}
			idx++
		}
	}
	w.archetypes.Add(archetype{})
	arch := w.archetypes.Get(w.archetypes.Len() - 1)
	arch.init(w.config.CapacityIncrement, types...)
	return arch
}

// Alive reports whether an entity is still alive.
func (w *World) Alive(entity Entity) bool {
	return w.entityPool.Alive(entity)
}

// componentID returns the ID for a component type, and registers it if not already registered.
func (w *World) componentID(tp reflect.Type) ID {
	return w.registry.ComponentID(tp)
}

func (w *World) nextArchetype(mask bitMask, index int) (int, archetypeIter, bool) {
	len := w.archetypes.Len()
	if index >= len {
		panic("exceeded end of query")
	}
	for i := index + 1; i < len; i++ {
		a := w.archetypes.Get(i)
		if a.Len() > 0 && a.mask.Contains(mask) {
			return i, newArchetypeIter(a), true
		}
	}
	return len, archetypeIter{}, false
}

// Query creates a [Query] iterator for the given components.
//
// Locks the world to prevent changes to component compositions.
func (w *World) Query(comps ...ID) Query {
	mask := newMask(comps...)
	lock := w.bitPool.Get()
	w.locks.Set(ID(lock), true)
	return newQuery(w, mask, lock)
}

// closeQuery closes a query and unlocks the world
func (w *World) closeQuery(query *Query) {
	l := query.lockBit
	if !w.locks.Get(ID(l)) {
		panic("unbalanced query unlock")
	}
	w.locks.Set(ID(l), false)
	w.bitPool.Recycle(l)
}

// IsLocked returns whether the world is locked by any queries.
func (w *World) IsLocked() bool {
	return w.locks != 0
}

func (w *World) checkLocked() {
	if w.locks != 0 {
		panic("attempt to modify a locked world")
	}
}
