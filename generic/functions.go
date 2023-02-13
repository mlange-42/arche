package generic

import (
	"github.com/mlange-42/arche/ecs"
)

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

// Add1 adds a component type to an entity.
//
// See also [ecs.World.Add].
func Add1[A any](w *ecs.World, entity ecs.Entity) *A {
	id := ecs.ComponentID[A](w)
	w.Add(entity, id)
	return (*A)(w.Get(entity, id))
}

// Add2 adds two component types to an entity.
//
// See also [ecs.World.Add].
func Add2[A any, B any](w *ecs.World, entity ecs.Entity) (*A, *B) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	w.Add(entity, idA, idB)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB))
}

// Add3 adds three component types to an entity.
//
// See also [ecs.World.Add].
func Add3[A any, B any, C any](w *ecs.World, entity ecs.Entity) (*A, *B, *C) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	w.Add(entity, idA, idB, idC)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC))
}

// Add4 adds four component types to an entity.
//
// See also [ecs.World.Add].
func Add4[A any, B any, C any, D any](w *ecs.World, entity ecs.Entity) (*A, *B, *C, *D) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	idD := ecs.ComponentID[D](w)
	w.Add(entity, idA, idB, idC, idD)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD))
}

// Add5 adds five component types to an entity.
//
// See also [ecs.World.Add].
func Add5[A any, B any, C any, D any, E any](w *ecs.World, entity ecs.Entity) (*A, *B, *C, *D, *E) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	idD := ecs.ComponentID[D](w)
	idE := ecs.ComponentID[E](w)
	w.Add(entity, idA, idB, idC, idD, idE)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD)), (*E)(w.Get(entity, idE))
}

// Assign1 adds a components to an entity.
//
// See also [ecs.World.Assign] and [ecs.World.AssignN].
func Assign1[A any](w *ecs.World, entity ecs.Entity, a *A) *A {
	idA := ecs.ComponentID[A](w)
	w.Assign(entity, idA, a)
	return (*A)(w.Get(entity, idA))
}

// Assign2 adds two components to an entity.
//
// See also [ecs.World.Assign] and [ecs.World.AssignN].
func Assign2[A any, B any](w *ecs.World, entity ecs.Entity, a *A, b *B) (*A, *B) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	w.AssignN(entity, ecs.Component{ID: idA, Component: a}, ecs.Component{ID: idB, Component: b})
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB))
}

// Assign3 adds three components to an entity.
//
// See also [ecs.World.Assign] and [ecs.World.AssignN].
func Assign3[A any, B any, C any](w *ecs.World, entity ecs.Entity, a *A, b *B, c *C) (*A, *B, *C) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	w.AssignN(entity, ecs.Component{ID: idA, Component: a}, ecs.Component{ID: idB, Component: b}, ecs.Component{ID: idC, Component: c})
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC))
}

// Assign4 adds four components to an entity.
//
// See also [ecs.World.Assign] and [ecs.World.AssignN].
func Assign4[A any, B any, C any, D any](w *ecs.World, entity ecs.Entity, a *A, b *B, c *C, d *D) (*A, *B, *C, *D) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	idD := ecs.ComponentID[D](w)
	w.AssignN(entity, ecs.Component{ID: idA, Component: a}, ecs.Component{ID: idB, Component: b}, ecs.Component{ID: idC, Component: c}, ecs.Component{ID: idD, Component: d})
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD))
}

// Assign5 adds four components to an entity.
//
// See also [ecs.World.Assign] and [ecs.World.AssignN].
func Assign5[A any, B any, C any, D any, E any](w *ecs.World, entity ecs.Entity, a *A, b *B, c *C, d *D, e *E) (*A, *B, *C, *D, *E) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	idD := ecs.ComponentID[D](w)
	idE := ecs.ComponentID[E](w)
	w.AssignN(entity, ecs.Component{ID: idA, Component: a}, ecs.Component{ID: idB, Component: b}, ecs.Component{ID: idC, Component: c}, ecs.Component{ID: idD, Component: d}, ecs.Component{ID: idE, Component: e})
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD)), (*E)(w.Get(entity, idE))
}

// Remove1 removes a component from an entity.
//
// See also [ecs.World.Remove].
func Remove1[A any](w *ecs.World, entity ecs.Entity) {
	w.Remove(entity, ecs.ComponentID[A](w))
}

// Remove2 removes two components from an entity.
//
// See also [ecs.World.Remove].
func Remove2[A any, B any](w *ecs.World, entity ecs.Entity) {
	w.Remove(entity, ecs.ComponentID[A](w), ecs.ComponentID[B](w))
}

// Remove3 removes three components from an entity.
//
// See also [ecs.World.Remove].
func Remove3[A any, B any, C any](w *ecs.World, entity ecs.Entity) {
	w.Remove(entity, ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w))
}

// Remove4 removes four components from an entity.
//
// See also [ecs.World.Remove].
func Remove4[A any, B any, C any, D any](w *ecs.World, entity ecs.Entity) {
	w.Remove(entity, ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w), ecs.ComponentID[D](w))
}

// Remove5 removes five components from an entity.
//
// See also [ecs.World.Remove].
func Remove5[A any, B any, C any, D any, E any](w *ecs.World, entity ecs.Entity) {
	w.Remove(entity, ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w), ecs.ComponentID[D](w), ecs.ComponentID[E](w))
}
