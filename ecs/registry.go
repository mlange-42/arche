package ecs

import "reflect"

// componentRegistry keeps track of component IDs
type componentRegistry struct {
	Components map[reflect.Type]ID
	Types      [MaskTotalBits]reflect.Type
}

// newComponentRegistry creates a new ComponentRegistry
func newComponentRegistry() componentRegistry {
	return componentRegistry{
		Components: map[reflect.Type]ID{},
	}
}

// RegisterComponent registers a components and assigns an ID for it
func (r *componentRegistry) RegisterComponent(tp reflect.Type, totalBits int) ID {
	id := ID(len(r.Components))
	if int(id) >= totalBits {
		panic("maximum of 128 component types exceeded")
	}
	r.Components[tp] = id
	r.Types[id] = tp
	return id
}

// ComponentID returns the ID for a component type, and registers it if not already registered
func (r *componentRegistry) ComponentID(tp reflect.Type) ID {
	if id, ok := r.Components[tp]; ok {
		return id
	}
	return r.RegisterComponent(tp, MaskTotalBits)
}

// ComponentType returns the type of a component by ID
func (r *componentRegistry) ComponentType(id ID) reflect.Type {
	return r.Types[id]
}
