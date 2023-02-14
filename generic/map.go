package generic

import "github.com/mlange-42/arche/ecs"

// Map provides a type-safe way to access a component type by entity ID.
//
// Create one with [NewMap].
type Map[T any] struct {
	id    ecs.ID
	world *ecs.World
}

// NewMap creates a new [Map] for a component type.
//
// Map provides a type-safe way to access a component type by entity ID.
//
// See also [ecs.World.Get], [ecs.World.Has] and [ecs.World.Set].
func NewMap[T any](w *ecs.World) Map[T] {
	return Map[T]{
		id:    ecs.ComponentID[T](w),
		world: w,
	}
}

// ID returns the component ID for this Map.
func (g *Map[T]) ID() ecs.ID {
	return g.id
}

// Get gets the component for the given entity.
//
// See also [ecs.World.Get].
func (g *Map[T]) Get(entity ecs.Entity) *T {
	return (*T)(g.world.Get(entity, g.id))
}

// Has returns whether the entity has the component.
//
// See also [ecs.World.Has].
func (g *Map[T]) Has(entity ecs.Entity) bool {
	return g.world.Has(entity, g.id)
}

// Set overwrites the component for the given entity.
//
// Panics if the entity does not have a component of that type.
//
// See also [ecs.World.Set].
func (g *Map[T]) Set(entity ecs.Entity, comp *T) *T {
	return (*T)(g.world.Set(entity, g.id, comp))
}
