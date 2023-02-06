package ecs

import (
	"reflect"
	"sort"
	"unsafe"
)

// World is the interface for the ECS world
type World interface {
	NewEntity() Entity
	RemEntity(entity Entity) bool
	Get(entity Entity, comps ID) unsafe.Pointer
	GetAt(arch, index int, id ID) unsafe.Pointer
	GetEntityAt(arch, index int) Entity
	Has(entity Entity, comps ID) bool
	Add(entity Entity, comps ...ID)
	Remove(entity Entity, comps ...ID)
	Alive(entity Entity) bool
	Registry() *ComponentRegistry
	Query(comps ...ID) Query
	Next(mask Mask, arch, index int) (int, int, bool)
}

// NewWorld creates a new World
func NewWorld() World {
	return newWorld()
}

func newWorld() *world {
	return &world{
		entities:   []entityIndex{},
		entityPool: NewEntityPool(),
		registry:   NewComponentRegistry(),
		archetypes: []Archetype{NewArchetype()},
	}
}

type world struct {
	entities   []entityIndex
	archetypes []Archetype
	entityPool EntityPool
	registry   ComponentRegistry
}

func (w *world) NewEntity() Entity {
	entity := w.entityPool.Get()
	idx := w.archetypes[0].Add(entity)
	if int(entity.id) == len(w.entities) {
		w.entities = append(w.entities, entityIndex{0, idx})
	} else {
		w.entities[entity.id] = entityIndex{0, idx}
	}
	return entity
}

func (w *world) RemEntity(entity Entity) bool {
	if !w.entityPool.Alive(entity) {
		return false
	}

	index := w.entities[entity.id]
	oldArch := &w.archetypes[index.arch]
	swapped := oldArch.Remove(int(index.index))
	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}

	w.entities[entity.id].arch = -1
	return true
}

func (w *world) Get(entity Entity, comp ID) unsafe.Pointer {
	index := w.entities[entity.id]
	arch := w.archetypes[index.arch]

	if !arch.HasComponent(comp) {
		return nil
	}

	return arch.Get(int(index.index), comp)
}

func (w *world) Has(entity Entity, comp ID) bool {
	index := w.entities[entity.id]
	arch := w.archetypes[index.arch]
	return arch.HasComponent(comp)
}

func (w *world) Add(entity Entity, comps ...ID) {
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

	swapped := oldArch.Remove(int(index.index))

	if swapped {
		swapEntity := oldArch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{archIdx, newIndex}
}

func (w *world) Remove(entity Entity, comps ...ID) {
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

	swapped := oldArch.Remove(int(index.index))
	if swapped {
		swapEntity := oldArch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{archIdx, newIndex}
}

func (w *world) findArchetype(mask Mask) (int, bool) {
	for i, a := range w.archetypes {
		if a.mask == mask {
			return i, true
		}
	}
	return 0, false
}

func (w *world) createArchetype(comps ...ID) int {
	sort.Slice(comps, func(i, j int) bool { return comps[i] < comps[j] })
	types := make([]ComponentType, len(comps))
	for i, id := range comps {
		types[i] = ComponentType{id, w.registry.types[id]}
	}
	a := NewArchetype(types...)
	w.archetypes = append(w.archetypes, a)
	return len(w.archetypes) - 1
}

func (w *world) Alive(entity Entity) bool {
	return w.entityPool.Alive(entity)
}

func (w *world) Registry() *ComponentRegistry {
	return &w.registry
}

func (w *world) Query(comps ...ID) Query {
	return NewQuery(w, NewMask(comps...))
}

func (w *world) Next(mask Mask, arch, index int) (int, int, bool) {
	if arch < 0 || index >= int(w.archetypes[arch].Len())-1 {
		arch++
		index = -1
		match := false
		for arch < len(w.archetypes) {
			if w.archetypes[arch].mask.Contains(mask) {
				match = true
				break
			}
			arch++
		}
		if !match {
			return 0, 0, false
		}
	}
	return arch, index + 1, true
}

func (w *world) GetAt(arch, index int, id ID) unsafe.Pointer {
	return w.archetypes[arch].Get(index, id)
}

func (w *world) GetEntityAt(arch, index int) Entity {
	return w.archetypes[arch].GetEntity(index)
}

// RegisterComponent provides a way to register components via generics
func RegisterComponent[T any](w World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.Registry().RegisterComponent(tp)
}

// ComponentID provides a way to get a component's ID via generics
func ComponentID[T any](w World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.Registry().ComponentID(tp)
}
