package generic

import "github.com/mlange-42/arche/ecs"

// Resource provides a type-safe way to access world resources.
//
// Create one with [NewResource].
type Resource[T any] struct {
	id    ecs.ResID
	world *ecs.World
}

// NewResource creates a new [Resource] mapper for a resource type.
// This does not add a resource to the world, but only creates a mapper for resource access!
//
// [Resource] provides a type-safe way to access a world resources.
//
// See also [ecs.World.Resources].
func NewResource[T any](w *ecs.World) Resource[T] {
	return Resource[T]{
		id:    ecs.ResourceID[T](w),
		world: w,
	}
}

// ID returns the resource [ecs.ResID] for this Resource mapper.
func (g *Resource[T]) ID() ecs.ResID {
	return g.id
}

// Add adds a resource to the world.
//
// Panics if there is already a resource of the given type.
//
// See also [ecs.Resources.Add].
func (g *Resource[T]) Add(res *T) {
	g.world.Resources().Add(g.id, res)
}

// Remove removes a resource from the world.
//
// Panics if there is no resource of the given type.
//
// See also [ecs.Resources.Remove].
func (g *Resource[T]) Remove() {
	g.world.Resources().Remove(g.id)
}

// Get gets the resource of the given type.
//
// Returns nil if there is no such resource.
//
// See also [ecs.Resources.Get].
func (g *Resource[T]) Get() *T {
	return g.world.Resources().Get(g.id).(*T)
}

// Has returns whether the world has the resource type.
//
// See also [ecs.Resources.Has].
func (g *Resource[T]) Has() bool {
	return g.world.Resources().Has(g.id)
}
