package ecs

import "reflect"

// World is the interface for the ECS world
type World interface {
	NewEntity() Entity
	RemEntity(entity Entity)
	Registry() *ComponentRegistry
}

// NewWorld creates a new World
func NewWorld() World {
	return &world{
		entityPool: NewEntityPool(),
		registry:   NewComponentRegistry(),
	}
}

type world struct {
	entityPool EntityPool
	registry   ComponentRegistry
}

func (w *world) NewEntity() Entity {
	return w.entityPool.Get()
}

func (w *world) RemEntity(entity Entity) {
	w.entityPool.Recycle(entity)
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
