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
	return World{
		config:     conf,
		entities:   []entityIndex{{-1, 0}},
		entityPool: newEntityPool(conf.CapacityIncrement),
		bitPool:    newBitPool(),
		registry:   newComponentRegistry(),
		archetypes: []archetype{newArchetype(conf.CapacityIncrement)},
		locks:      Mask(0),
	}
}

// World is the central type holding [Entity] and component data.
type World struct {
	config     Config
	entities   []entityIndex
	archetypes []archetype
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
	idx := w.archetypes[0].Add(entity)
	if int(entity.id) == len(w.entities) {
		if len(w.entities) == cap(w.entities) {
			old := w.entities
			w.entities = make([]entityIndex, len(w.entities), len(w.entities)+w.config.CapacityIncrement)
			copy(w.entities, old)
		}
		w.entities = append(w.entities, entityIndex{0, idx})
	} else {
		w.entities[entity.id] = entityIndex{0, idx}
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
	oldArch := &w.archetypes[index.arch]
	swapped := oldArch.Remove(int(index.index))

	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}

	w.entities[entity.id].arch = -1
}

// Get returns a pointer th the given component of an [Entity].
//
// Returns `nil` if the entity has no such component.
// Panics when called for an already removed entity.
func (w *World) Get(entity Entity, comp ID) unsafe.Pointer {
	index := w.entities[entity.id]
	arch := &w.archetypes[index.arch]

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
	arch := &w.archetypes[index.arch]
	return arch.HasComponent(comp)
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
	oldArch := &w.archetypes[index.arch]
	mask := oldArch.mask
	for _, comp := range comps {
		mask.Set(comp, true)
	}

	oldIDs := oldArch.Components()

	archIdx, ok := w.findArchetype(mask)
	if !ok {
		ids := append(oldIDs, comps...)
		archIdx = w.createArchetype(ids...)
	}
	arch := &w.archetypes[archIdx]
	oldArch = &w.archetypes[index.arch]

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
	w.entities[entity.id] = entityIndex{archIdx, newIndex}
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
	oldArch := &w.archetypes[index.arch]
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

	archIdx, ok := w.findArchetype(mask)
	if !ok {
		archIdx = w.createArchetype(newIDs...)
	}
	arch := &w.archetypes[archIdx]
	oldArch = &w.archetypes[index.arch]

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
	w.entities[entity.id] = entityIndex{archIdx, newIndex}
}

func (w *World) findArchetype(mask Mask) (int, bool) {
	length := len(w.archetypes)
	for i := 0; i < length; i++ {
		if w.archetypes[i].mask == mask {
			return i, true
		}
	}
	return 0, false
}

func (w *World) createArchetype(comps ...ID) int {
	sort.Slice(comps, func(i, j int) bool { return comps[i] < comps[j] })
	types := make([]componentType, len(comps))
	for i, id := range comps {
		types[i] = componentType{id, w.registry.types[id]}
	}
	w.archetypes = append(w.archetypes, newArchetype(w.config.CapacityIncrement, types...))
	return len(w.archetypes) - 1
}

// Alive reports whether an entity is still alive.
func (w *World) Alive(entity Entity) bool {
	return w.entityPool.Alive(entity)
}

// componentID returns the ID for a component type, and registers it if not already registered.
func (w *World) componentID(tp reflect.Type) ID {
	return w.registry.ComponentID(tp)
}

// Query creates a [Query] iterator for the given components.
//
// Locks the world to prevent changes to component compositions.
func (w *World) Query(comps ...ID) Query {
	mask := NewMask(comps...)
	arches := []archetypeIter{}
	length := len(w.archetypes)
	count := 0
	for i := 0; i < length; i++ {
		arch := &w.archetypes[i]
		if arch.mask.Contains(mask) {
			arches = append(arches, newArchetypeIter(arch))
			count += int(arch.Len())
		}
	}
	lock := w.bitPool.Get()
	w.locks.Set(ID(lock), true)
	return newQuery(w, arches, count, lock)
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
