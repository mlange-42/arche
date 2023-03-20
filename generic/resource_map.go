package generic

import "github.com/mlange-42/arche/ecs"

// ResMap provides a type-safe way to access a world resources.
//
// Create one with [NewResMap].
type ResMap[T any] struct {
	id    ecs.ID
	world *ecs.World
}

// NewResMap creates a new [ResMap] for a resource type.
//
// ResMap provides a type-safe way to access a world resources.
//
// See also [ecs.World.GetResource], [ecs.World.HasResource] and [ecs.World.AddResource].
func NewResMap[T any](w *ecs.World) ResMap[T] {
	return ResMap[T]{
		id:    ecs.ResourceID[T](w),
		world: w,
	}
}

// ID returns the resource ID for this Map.
func (g *ResMap[T]) ID() ecs.ID {
	return g.id
}

// Get gets the resource of the given type.
//
// See also [ecs.World.Get].
func (g *ResMap[T]) Get() *T {
	return g.world.GetResource(g.id).(*T)
}

// Has returns whether the world has the resource.
//
// See also [ecs.World.HasResource].
func (g *ResMap[T]) Has() bool {
	return g.world.HasResource(g.id)
}
