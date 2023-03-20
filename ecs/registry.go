package ecs

import "reflect"

// componentRegistry keeps track of component IDs
type componentRegistry[T uint8] struct {
	Components map[reflect.Type]T
	Types      []reflect.Type
}

// newComponentRegistry creates a new ComponentRegistry
func newComponentRegistry[T uint8]() componentRegistry[T] {
	return componentRegistry[T]{
		Components: map[reflect.Type]T{},
		Types:      make([]reflect.Type, MaskTotalBits),
	}
}

// RegisterComponent registers a components and assigns an ID for it
func (r *componentRegistry[T]) RegisterComponent(tp reflect.Type, totalBits int) T {
	id := T(len(r.Components))
	if int(id) >= totalBits {
		panic("maximum of 128 component types exceeded")
	}
	r.Components[tp] = id
	r.Types[id] = tp
	return id
}

// ComponentID returns the ID for a component type, and registers it if not already registered
func (r *componentRegistry[T]) ComponentID(tp reflect.Type) T {
	if id, ok := r.Components[tp]; ok {
		return id
	}
	return r.RegisterComponent(tp, MaskTotalBits)
}

// ComponentType returns the type of a component by ID
func (r *componentRegistry[T]) ComponentType(id T) reflect.Type {
	return r.Types[id]
}
