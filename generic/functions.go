package generic

import (
	"github.com/mlange-42/arche/ecs"
)

// New1 creates a new entity with two component types.
// Returns the entity and the new components.
//
// See also [ecs.World.NewEntity].
func New1[A any](w *ecs.World, entity ecs.Entity) (ecs.Entity, *A) {
	idA := ecs.ComponentID[A](w)
	entity = w.NewEntity(idA)
	return entity, (*A)(w.Get(entity, idA))
}

// New2 creates a new entity with two component types.
// Returns the entity and the new components.
//
// See also [ecs.World.NewEntity].
func New2[A any, B any](w *ecs.World, entity ecs.Entity) (ecs.Entity, *A, *B) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	entity = w.NewEntity(idA, idB)
	return entity, (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB))
}

// New3 creates a new entity with three component types.
// Returns the entity and the new components.
//
// See also and [ecs.World.NewEntity].
func New3[A any, B any, C any](w *ecs.World, entity ecs.Entity) (ecs.Entity, *A, *B, *C) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	entity = w.NewEntity(idA, idB, idC)
	return entity, (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC))
}

// New4 creates a new entity with four component types.
// Returns the entity and the new components.
//
// See also [ecs.World.NewEntity].
func New4[A any, B any, C any, D any](w *ecs.World, entity ecs.Entity) (ecs.Entity, *A, *B, *C, *D) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	idD := ecs.ComponentID[D](w)
	entity = w.NewEntity(idA, idB, idC, idD)
	return entity, (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD))
}

// New5 creates a new entity with five component types.
// Returns the entity and the new components.
//
// See also [ecs.World.NewEntity].
func New5[A any, B any, C any, D any, E any](w *ecs.World, entity ecs.Entity) (ecs.Entity, *A, *B, *C, *D, *E) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	idD := ecs.ComponentID[D](w)
	idE := ecs.ComponentID[E](w)
	entity = w.NewEntity(idA, idB, idC, idD, idE)
	return entity, (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD)), (*E)(w.Get(entity, idE))
}

// NewWith1 creates a new entity with one component, taken from a pointer.
// Returns the entity and the new components.
//
// See also [ecs.World.NewEntity].
func NewWith1[A any](w *ecs.World, entity ecs.Entity, a *A) (ecs.Entity, *A) {
	idA := ecs.ComponentID[A](w)
	entity = w.NewEntityWith(ecs.Component{ID: idA, Component: a})
	return entity, (*A)(w.Get(entity, idA))
}

// NewWith2 creates a new entity with two components, taken from pointers.
// Returns the entity and the new components.
//
// See also [ecs.World.NewEntity].
func NewWith2[A any, B any](w *ecs.World, entity ecs.Entity, a *A, b *B) (ecs.Entity, *A, *B) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	entity = w.NewEntityWith(ecs.Component{ID: idA, Component: a}, ecs.Component{ID: idB, Component: b})
	return entity, (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB))
}

// NewWith3 creates a new entity with three components, taken from pointers.
// Returns the entity and the new components.
//
// See also [ecs.World.NewEntity].
func NewWith3[A any, B any, C any](w *ecs.World, entity ecs.Entity, a *A, b *B, c *C) (ecs.Entity, *A, *B, *C) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	entity = w.NewEntityWith(ecs.Component{ID: idA, Component: a}, ecs.Component{ID: idB, Component: b}, ecs.Component{ID: idC, Component: c})
	return entity, (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC))
}

// NewWith4 creates a new entity with four components, taken from pointers.
// Returns the entity and the new components.
//
// See also [ecs.World.NewEntity].
func NewWith4[A any, B any, C any, D any](w *ecs.World, entity ecs.Entity, a *A, b *B, c *C, d *D) (ecs.Entity, *A, *B, *C, *D) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	idD := ecs.ComponentID[D](w)
	entity = w.NewEntityWith(ecs.Component{ID: idA, Component: a}, ecs.Component{ID: idB, Component: b}, ecs.Component{ID: idC, Component: c}, ecs.Component{ID: idD, Component: d})
	return entity, (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD))
}

// NewWith5 creates a new entity with five components, taken from pointers.
// Returns the entity and the new components.
//
// See also [ecs.World.NewEntity].
func NewWith5[A any, B any, C any, D any, E any](w *ecs.World, entity ecs.Entity, a *A, b *B, c *C, d *D, e *E) (ecs.Entity, *A, *B, *C, *D, *E) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	idD := ecs.ComponentID[D](w)
	idE := ecs.ComponentID[E](w)
	entity = w.NewEntityWith(ecs.Component{ID: idA, Component: a}, ecs.Component{ID: idB, Component: b}, ecs.Component{ID: idC, Component: c}, ecs.Component{ID: idD, Component: d}, ecs.Component{ID: idE, Component: e})
	return entity, (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD)), (*E)(w.Get(entity, idE))
}

// Add1 adds a component type to an entity, and returns the component.
//
// See also [ecs.World.Add].
func Add1[A any](w *ecs.World, entity ecs.Entity) *A {
	id := ecs.ComponentID[A](w)
	w.Add(entity, id)
	return (*A)(w.Get(entity, id))
}

// Add2 adds two component types to an entity, and returns the components.
//
// See also [ecs.World.Add].
func Add2[A any, B any](w *ecs.World, entity ecs.Entity) (*A, *B) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	w.Add(entity, idA, idB)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB))
}

// Add3 adds three component types to an entity, and returns the components.
//
// See also [ecs.World.Add].
func Add3[A any, B any, C any](w *ecs.World, entity ecs.Entity) (*A, *B, *C) {
	idA := ecs.ComponentID[A](w)
	idB := ecs.ComponentID[B](w)
	idC := ecs.ComponentID[C](w)
	w.Add(entity, idA, idB, idC)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC))
}

// Add4 adds four component types to an entity, and returns the components.
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

// Add5 adds five component types to an entity, and returns the components.
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

// Assign5 adds five components to an entity.
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
