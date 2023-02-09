package ecs

import "reflect"

// ComponentID returns the ID for a component type. Registers the type if it is not already registered.
func ComponentID[T any](w *World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.componentID(tp)
}

// Add adds a component type to an entity
func Add[A any](w *World, entity Entity) *A {
	id := ComponentID[A](w)
	w.Add(entity, id)
	return (*A)(w.Get(entity, id))
}

// Add2 adds two component type to an entity
func Add2[A any, B any](w *World, entity Entity) (*A, *B) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	w.Add(entity, idA, idB)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB))
}

// Add3 adds three component type to an entity
func Add3[A any, B any, C any](w *World, entity Entity) (*A, *B, *C) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	idC := ComponentID[C](w)
	w.Add(entity, idA, idB, idC)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC))
}

// Add4 adds four component type to an entity
func Add4[A any, B any, C any, D any](w *World, entity Entity) (*A, *B, *C, *D) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	idC := ComponentID[C](w)
	idD := ComponentID[D](w)
	w.Add(entity, idA, idB, idC, idD)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD))
}

// Add5 adds five component type to an entity
func Add5[A any, B any, C any, D any, E any](w *World, entity Entity) (*A, *B, *C, *D, *E) {
	idA := ComponentID[A](w)
	idB := ComponentID[B](w)
	idC := ComponentID[C](w)
	idD := ComponentID[D](w)
	idE := ComponentID[E](w)
	w.Add(entity, idA, idB, idC, idD, idE)
	return (*A)(w.Get(entity, idA)), (*B)(w.Get(entity, idB)), (*C)(w.Get(entity, idC)), (*D)(w.Get(entity, idD)), (*E)(w.Get(entity, idE))
}

// Remove removes a component type to an entity
func Remove[A any](w *World, entity Entity) {
	w.Remove(entity, ComponentID[A](w))
}

// Remove2 removes two component type to an entity
func Remove2[A any, B any](w *World, entity Entity) {
	w.Remove(entity, ComponentID[A](w), ComponentID[B](w))
}

// Remove3 removes three component type to an entity
func Remove3[A any, B any, C any](w *World, entity Entity) {
	w.Remove(entity, ComponentID[A](w), ComponentID[B](w), ComponentID[C](w))
}

// Remove4 removes four component type to an entity
func Remove4[A any, B any, C any, D any](w *World, entity Entity) {
	w.Remove(entity, ComponentID[A](w), ComponentID[B](w), ComponentID[C](w), ComponentID[D](w))
}

// Remove5 removes five component type to an entity
func Remove5[A any, B any, C any, D any, E any](w *World, entity Entity) {
	w.Remove(entity, ComponentID[A](w), ComponentID[B](w), ComponentID[C](w), ComponentID[D](w), ComponentID[E](w))
}
