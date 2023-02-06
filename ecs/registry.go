package ecs

import "reflect"

// componentRegistry keeps track of component IDs
type componentRegistry struct {
	components map[reflect.Type]ID
	types      [MaskTotalBits]reflect.Type
}

// newComponentRegistry creates a new ComponentRegistry
func newComponentRegistry() componentRegistry {
	return componentRegistry{
		components: map[reflect.Type]ID{},
	}
}

// RegisterComponent registers a components and assigns an ID for it
func (r *componentRegistry) RegisterComponent(tp reflect.Type) ID {
	id := ID(len(r.components))
	r.components[tp] = id
	r.types[id] = tp
	return id
}

// ComponentID returns the ID for a component type, and registers it if not already registered
func (r *componentRegistry) ComponentID(tp reflect.Type) ID {
	if id, ok := r.components[tp]; ok {
		return id
	}
	return r.RegisterComponent(tp)
}

// ComponentType returns the type of a component by ID
func (r *componentRegistry) ComponentType(id ID) reflect.Type {
	return r.types[id]
}
