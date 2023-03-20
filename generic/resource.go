package generic

import "github.com/mlange-42/arche/ecs"

// Resource provides a type-safe way to access a world resources.
//
// Create one with [NewResource].
type Resource[T any] struct {
	id    ecs.ID
	world *ecs.World
}

// NewResource creates a new [Resource] for a resource type.
//
// [Resource] provides a type-safe way to access a world resources.
//
// See also [ecs.World.GetResource], [ecs.World.HasResource] and [ecs.World.AddResource].
func NewResource[T any](w *ecs.World) Resource[T] {
	return Resource[T]{
		id:    ecs.ResourceID[T](w),
		world: w,
	}
}

// ID returns the resource ID for this Map.
func (g *Resource[T]) ID() ecs.ID {
	return g.id
}

// Get gets the resource of the given type.
//
// Returns nil if there is no such resource.
//
// See also [ecs.World.GetResource].
func (g *Resource[T]) Get() *T {
	return g.world.GetResource(g.id).(*T)
}

// Has returns whether the world has the resource.
//
// See also [ecs.World.HasResource].
func (g *Resource[T]) Has() bool {
	return g.world.HasResource(g.id)
}
