package ecs

import "reflect"

// componentRegistry keeps track of component IDs.
type componentRegistry[T uint8] struct {
	Components map[reflect.Type]T
	Types      []reflect.Type
	Used       Mask
}

// newComponentRegistry creates a new ComponentRegistry.
func newComponentRegistry[T uint8]() componentRegistry[T] {
	return componentRegistry[T]{
		Components: map[reflect.Type]T{},
		Types:      make([]reflect.Type, MaskTotalBits),
		Used:       Mask{},
	}
}

// ComponentID returns the ID for a component type, and registers it if not already registered.
func (r *componentRegistry[T]) ComponentID(tp reflect.Type) T {
	if id, ok := r.Components[tp]; ok {
		return id
	}
	return r.registerComponent(tp, MaskTotalBits)
}

// ComponentType returns the type of a component by ID.
func (r *componentRegistry[T]) ComponentType(id T) (reflect.Type, bool) {
	return r.Types[id], r.Used.Get(uint8(id))
}

// registerComponent registers a components and assigns an ID for it.
func (r *componentRegistry[T]) registerComponent(tp reflect.Type, totalBits int) T {
	id := T(len(r.Components))
	if int(id) >= totalBits {
		panic("maximum of 128 component types exceeded")
	}
	r.Components[tp], r.Types[id] = id, tp
	r.Used.Set(uint8(id), true)
	return id
}
