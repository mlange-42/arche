package base

import "reflect"

// MaskTotalBits is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = 64

// ComponentRegistry keeps track of component IDs
type ComponentRegistry struct {
	Components map[reflect.Type]ID
	Types      [MaskTotalBits]reflect.Type
}

// NewComponentRegistry creates a new ComponentRegistry
func NewComponentRegistry() ComponentRegistry {
	return ComponentRegistry{
		Components: map[reflect.Type]ID{},
	}
}

// RegisterComponent registers a components and assigns an ID for it
func (r *ComponentRegistry) RegisterComponent(tp reflect.Type, totalBits int) ID {
	id := ID(len(r.Components))
	if int(id) >= totalBits {
		panic("maximum of 64 component types exceeded")
	}
	r.Components[tp] = id
	r.Types[id] = tp
	return id
}

// ComponentID returns the ID for a component type, and registers it if not already registered
func (r *ComponentRegistry) ComponentID(tp reflect.Type) ID {
	if id, ok := r.Components[tp]; ok {
		return id
	}
	return r.RegisterComponent(tp, MaskTotalBits)
}

// ComponentType returns the type of a component by ID
func (r *ComponentRegistry) ComponentType(id ID) reflect.Type {
	return r.Types[id]
}
