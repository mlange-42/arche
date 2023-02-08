package ecs

import (
	"reflect"
	"sort"
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
	arches := PagedArr32[archetype]{}
	arches.Add(arch)
	return World{
		config:     conf,
		entities:   []entityIndex{{nil, 0}},
		entityPool: newEntityPool(conf.CapacityIncrement),
		bitPool:    newBitPool(),
		registry:   newComponentRegistry(),
		archetypes: arches,
		locks:      Mask(0),
	}
}

// World is the central type holding [Entity] and component data.
type World struct {
	config     Config
	entities   []entityIndex
	archetypes PagedArr32[archetype]
	entityPool entityPool
	bitPool    bitPool
	registry   componentRegistry
	locks      Mask
}

// NewEntity returns a new or recycled [Entity].
//
// Panics when called on a locked world.
//
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
//
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
//
// Do not use during [Query] iteration!
func (w *World) Add(entity Entity, comps ...ID) {
	w.checkLocked()

	if len(comps) == 0 {
		return
	}
	index := w.entities[entity.id]
	oldArch := index.arch
	mask := oldArch.mask
	for _, comp := range comps {
		mask.Set(comp, true)
	}

	oldIDs := oldArch.Components()

	arch, ok := w.findArchetype(mask)
	if !ok {
		ids := append(oldIDs, comps...)
		arch = w.createArchetype(ids...)
	}

	allComps := make([]componentPointer, 0, len(oldIDs)+len(comps))
	for _, id := range oldIDs {
		comp := oldArch.Get(int(index.index), id)
		allComps = append(allComps, componentPointer{id, comp})
	}
	for _, id := range comps {
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

// Remove removes components from an entity.
//
// Panics when called on a locked world or for an already removed entity.
//
// Do not use during [Query] iteration!
func (w *World) Remove(entity Entity, comps ...ID) {
	w.checkLocked()

	if len(comps) == 0 {
		return
	}
	index := w.entities[entity.id]
	oldArch := index.arch
	mask := oldArch.mask
	for _, comp := range comps {
		mask.Set(comp, false)
	}

	oldIDs := oldArch.Components()
	newIDs := make([]ID, 0, len(oldIDs))
	for _, id := range oldIDs {
		if mask.Get(id) {
			newIDs = append(newIDs, id)
		}
	}

	arch, ok := w.findArchetype(mask)
	if !ok {
		arch = w.createArchetype(newIDs...)
	}

	allComps := make([]componentPointer, 0, len(newIDs))
	for _, id := range newIDs {
		comp := oldArch.Get(int(index.index), id)
		allComps = append(allComps, componentPointer{id, comp})
	}

	newIndex := arch.AddPointer(entity, allComps...)

	swapped := oldArch.Remove(int(index.index))

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{arch, newIndex}
}

func (w *World) findArchetype(mask Mask) (*archetype, bool) {
	length := w.archetypes.Len()
	for i := 0; i < length; i++ {
		arch := w.archetypes.Get(i)
		if arch.mask == mask {
			return arch, true
		}
	}
	return nil, false
}

func (w *World) createArchetype(comps ...ID) *archetype {
	sort.Slice(comps, func(i, j int) bool { return comps[i] < comps[j] })
	types := make([]componentType, len(comps))
	for i, id := range comps {
		types[i] = componentType{id, w.registry.types[id]}
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

func (w *World) nextArchetype(mask Mask, index int) (int, archetypeIter, bool) {
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
	mask := NewMask(comps...)
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
