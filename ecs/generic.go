package ecs

import "reflect"

// ComponentID returns the ID for a component type. Registers the type if it is not already registered.
func ComponentID[T any](w *World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.componentID(tp)
}

// Get gets a component for an entity. THIS IS SLOW!!!
//
// Using [World.Get] is about 5 times faster,
// and using queries in about 10 times faster.
func Get[A any](w *World, entity Entity) *A {
	id := ComponentID[A](w)
	return (*A)(w.Get(entity, id))
}

// Has returns whether a component attached to an entity. THIS IS SLOW!!!
//
// Using [World.Has] is more than 5 times faster.
func Has[A any](w *World, entity Entity) bool {
	id := ComponentID[A](w)
	return w.Has(entity, id)
}

// Add adds a component type to an entity.
//
// See also [World.Add].
func Add[A any](w *World, entity Entity) *A {
	id := ComponentID[A](w)
	w.Add(entity, id)
	return (*A)(w.Get(entity, id))
}

// Add2 adds two component types to an entity.
//
// See also [World.Add].
func Add2[A any, B any](w *World, entity Entity) (*A, *B) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	w.Add(entity, idA, idB)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB))
}

// Add3 adds three component types to an entity.
//
// See also [World.Add].
func Add3[A any, B any, C any](w *World, entity Entity) (*A, *B, *C) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	idC := ComponentID[C](w)
	w.Add(entity, idA, idB, idC)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC))
}

// Add4 adds four component types to an entity.
//
// See also [World.Add].
func Add4[A any, B any, C any, D any](w *World, entity Entity) (*A, *B, *C, *D) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	idC := ComponentID[C](w)
	idD := ComponentID[D](w)
	w.Add(entity, idA, idB, idC, idD)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD))
}

// Add5 adds five component types to an entity.
//
// See also [World.Add].
func Add5[A any, B any, C any, D any, E any](w *World, entity Entity) (*A, *B, *C, *D, *E) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	idC := ComponentID[C](w)
	idD := ComponentID[D](w)
	idE := ComponentID[E](w)
	w.Add(entity, idA, idB, idC, idD, idE)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD)), (*E)(w.Get(entity, idE))
}

// Assign adds a components to an entity.
//
// See also [World.Assign] and [World.AssignN].
func Assign[A any](w *World, entity Entity, a *A) *A {
	idA := ComponentID[A](w)
	w.Assign(entity, idA, a)
	return (*A)(w.Get(entity, idA))
}

// Assign2 adds two components to an entity.
//
// See also [World.Assign] and [World.AssignN].
func Assign2[A any, B any](w *World, entity Entity, a *A, b *B) (*A, *B) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	w.AssignN(entity, Component{idA, a}, Component{idB, b})
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB))
}

// Assign3 adds three components to an entity.
//
// See also [World.Assign] and [World.AssignN].
func Assign3[A any, B any, C any](w *World, entity Entity, a *A, b *B, c *C) (*A, *B, *C) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	idC := ComponentID[C](w)
	w.AssignN(entity, Component{idA, a}, Component{idB, b}, Component{idC, c})
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC))
}

// Assign4 adds four components to an entity.
//
// See also [World.Assign] and [World.AssignN].
func Assign4[A any, B any, C any, D any](w *World, entity Entity, a *A, b *B, c *C, d *D) (*A, *B, *C, *D) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	idC := ComponentID[C](w)
	idD := ComponentID[D](w)
	w.AssignN(entity, Component{idA, a}, Component{idB, b}, Component{idC, c}, Component{idD, d})
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD))
}

// Assign5 adds four components to an entity.
//
// See also [World.Assign] and [World.AssignN].
func Assign5[A any, B any, C any, D any, E any](w *World, entity Entity, a *A, b *B, c *C, d *D, e *E) (*A, *B, *C, *D, *E) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	idC := ComponentID[C](w)
	idD := ComponentID[D](w)
	idE := ComponentID[E](w)
	w.AssignN(entity, Component{idA, a}, Component{idB, b}, Component{idC, c}, Component{idD, d}, Component{idE, e})
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD)), (*E)(w.Get(entity, idE))
}

// Remove removes a component from an entity.
//
// See also [World.Remove].
func Remove[A any](w *World, entity Entity) {
	w.Remove(entity, ComponentID[A](w))
}

// Remove2 removes two components from an entity.
//
// See also [World.Remove].
func Remove2[A any, B any](w *World, entity Entity) {
	w.Remove(entity, ComponentID[A](w), ComponentID[B](w))
}

// Remove3 removes three components from an entity.
//
// See also [World.Remove].
func Remove3[A any, B any, C any](w *World, entity Entity) {
	w.Remove(entity, ComponentID[A](w), ComponentID[B](w), ComponentID[C](w))
}

// Remove4 removes four components from an entity.
//
// See also [World.Remove].
func Remove4[A any, B any, C any, D any](w *World, entity Entity) {
	w.Remove(entity, ComponentID[A](w), ComponentID[B](w), ComponentID[C](w), ComponentID[D](w))
}

// Remove5 removes five components from an entity.
//
// See also [World.Remove].
func Remove5[A any, B any, C any, D any, E any](w *World, entity Entity) {
	w.Remove(entity, ComponentID[A](w), ComponentID[B](w), ComponentID[C](w), ComponentID[D](w), ComponentID[E](w))
}
