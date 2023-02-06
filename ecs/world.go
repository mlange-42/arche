package ecs

import (
	"reflect"
	"sort"
	"unsafe"
)

// NewWorld creates a new World
func NewWorld() World {
	return World{
		entities:         []entityIndex{},
		entityPool:       NewEntityPool(),
		registry:         NewComponentRegistry(),
		archetypes:       []Archetype{NewArchetype()},
		archetypeLengths: []int{0},
	}
}

// World holds all ECS data
type World struct {
	entities         []entityIndex
	archetypes       []Archetype
	archetypeLengths []int
	entityPool       EntityPool
	registry         ComponentRegistry
}

// NewEntity creates a new or recycled entity
func (w *World) NewEntity() Entity {
	entity := w.entityPool.Get()
	idx := w.archetypes[0].Add(entity)
	w.archetypeLengths[0]++
	if int(entity.id) == len(w.entities) {
		w.entities = append(w.entities, entityIndex{0, idx})
	} else {
		w.entities[entity.id] = entityIndex{0, idx}
	}
	return entity
}

// RemEntity recycles an entity
func (w *World) RemEntity(entity Entity) bool {
	if !w.entityPool.Alive(entity) {
		return false
	}

	index := w.entities[entity.id]
	oldArch := &w.archetypes[index.arch]
	swapped := oldArch.Remove(int(index.index))
	w.archetypeLengths[index.arch]--

	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}

	w.entities[entity.id].arch = -1
	return true
}

// Get returns a component for an entity
func (w *World) Get(entity Entity, comp ID) unsafe.Pointer {
	index := w.entities[entity.id]
	arch := w.archetypes[index.arch]

	if !arch.HasComponent(comp) {
		return nil
	}

	return arch.Get(int(index.index), comp)
}

// Has returns whether an entity has a component
func (w *World) Has(entity Entity, comp ID) bool {
	index := w.entities[entity.id]
	arch := w.archetypes[index.arch]
	return arch.HasComponent(comp)
}

// Add adds components to an entity
func (w *World) Add(entity Entity, comps ...ID) {
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

	allComps := make([]ComponentPointer, 0, len(oldIDs)+len(comps))
	for _, id := range oldIDs {
		comp := oldArch.Get(int(index.index), id)
		allComps = append(allComps, ComponentPointer{id, comp})
	}
	for _, id := range comps {
		allComps = append(allComps, ComponentPointer{id, nil})
	}

	newIndex := arch.AddPointer(entity, allComps...)
	w.archetypeLengths[archIdx]++

	swapped := oldArch.Remove(int(index.index))
	w.archetypeLengths[index.arch]--

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{archIdx, newIndex}
}

// Remove removes components from an entity
func (w *World) Remove(entity Entity, comps ...ID) {
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

	allComps := make([]ComponentPointer, 0, len(newIDs))
	for _, id := range newIDs {
		comp := oldArch.Get(int(index.index), id)
		allComps = append(allComps, ComponentPointer{id, comp})
	}

	newIndex := arch.AddPointer(entity, allComps...)
	w.archetypeLengths[archIdx]++

	swapped := oldArch.Remove(int(index.index))
	w.archetypeLengths[index.arch]--

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{archIdx, newIndex}
}

func (w *World) findArchetype(mask Mask) (int, bool) {
	for i, a := range w.archetypes {
		if a.mask == mask {
			return i, true
		}
	}
	return 0, false
}

func (w *World) createArchetype(comps ...ID) int {
	sort.Slice(comps, func(i, j int) bool { return comps[i] < comps[j] })
	types := make([]ComponentType, len(comps))
	for i, id := range comps {
		types[i] = ComponentType{id, w.registry.types[id]}
	}
	a := NewArchetype(types...)
	w.archetypes = append(w.archetypes, a)
	w.archetypeLengths = append(w.archetypeLengths, 0)
	return len(w.archetypes) - 1
}

// Alive reports whether an entity is still alive
func (w *World) Alive(entity Entity) bool {
	return w.entityPool.Alive(entity)
}

// Registry returns the world's ComponentRegistry
func (w *World) Registry() *ComponentRegistry {
	return &w.registry
}

// Query creates a query iterator for the given components
func (w *World) Query(comps ...ID) Query {
	mask := NewMask(comps...)
	arches := []archetypeIter{}
	for _, arch := range w.archetypes {
		if arch.mask.Contains(mask) {
			arches = append(arches, newArchetypeIter(&arch))
		}
	}
	return NewQuery(arches)
}

// RegisterComponent provides a way to register components via generics
func RegisterComponent[T any](w *World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.Registry().RegisterComponent(tp)
}

// ComponentID provides a way to get a component's ID via generics
func ComponentID[T any](w *World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.Registry().ComponentID(tp)
}
