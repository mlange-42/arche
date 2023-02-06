package ecs

import "reflect"

// ComponentRegistry keeps track of component IDs
type ComponentRegistry struct {
	components map[reflect.Type]ID
	types      [MaskTotalBits]reflect.Type
}

// NewComponentRegistry creates a new ComponentRegistry
func NewComponentRegistry() ComponentRegistry {
	return ComponentRegistry{
		components: map[reflect.Type]ID{},
	}
}

// RegisterComponent registers a components and assigns an ID for it
func (r *ComponentRegistry) RegisterComponent(tp reflect.Type) ID {
	id := ID(len(r.components))
	r.components[tp] = id
	r.types[id] = tp
	return id
}

// ComponentID returns the ID for a component type, and registers it if not already registered
func (r *ComponentRegistry) ComponentID(tp reflect.Type) ID {
	if id, ok := r.components[tp]; ok {
		return id
	}
	return r.RegisterComponent(tp)
}

// ComponentType returns the type of a component by ID
func (r *ComponentRegistry) ComponentType(id ID) reflect.Type {
	return r.types[id]
}
