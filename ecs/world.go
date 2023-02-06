package ecs

import "reflect"

// World is the interface for the ECS world
type World interface {
	NewEntity() Entity
	RemEntity(entity Entity) bool
	Alive(entity Entity) bool
	Registry() *ComponentRegistry
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
		w.entities = append(w.entities, entityIndex{&w.archetypes[0], idx})
	} else {
		w.entities[entity.id] = entityIndex{&w.archetypes[0], idx}
	}
	return entity
}

func (w *world) RemEntity(entity Entity) bool {
	if !w.entityPool.Alive(entity) {
		return false
	}

	index := w.entities[entity.id]
	swapped := index.arch.Remove(int(index.index))
	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := index.arch.GetEntity(int(index.index))
		w.entities[swapEntity.id].index = index.index
	}
	return true
}

func (w *world) Alive(entity Entity) bool {
	return w.entityPool.Alive(entity)
}

func (w *world) Registry() *ComponentRegistry {
	return &w.registry
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
