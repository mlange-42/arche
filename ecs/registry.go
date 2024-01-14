package ecs

import (
	"fmt"
	"reflect"
)

// componentRegistry keeps track of component IDs.
type componentRegistry[T uint8] struct {
	Components map[reflect.Type]T
	Types      []reflect.Type
	Used       Mask
	IsRelation Mask
}

// newComponentRegistry creates a new ComponentRegistry.
func newComponentRegistry[T uint8]() componentRegistry[T] {
	return componentRegistry[T]{
		Components: map[reflect.Type]T{},
		Types:      make([]reflect.Type, MaskTotalBits),
		Used:       Mask{},
		IsRelation: Mask{},
	}
}

// ComponentID returns the ID for a component type, and registers it if not already registered.
// The second return value indicates if it is a newly created ID.
func (r *componentRegistry[T]) ComponentID(tp reflect.Type) (T, bool) {
	if id, ok := r.Components[tp]; ok {
		return id, false
	}
	return r.registerComponent(tp, MaskTotalBits), true
}

// ComponentType returns the type of a component by ID.
func (r *componentRegistry[T]) ComponentType(id T) (reflect.Type, bool) {
	return r.Types[id], r.Used.Get(uint8(id))
}

// ComponentType returns the type of a component by ID.
func (r *componentRegistry[T]) Count() int {
	return len(r.Components)
}

// registerComponent registers a components and assigns an ID for it.
func (r *componentRegistry[T]) registerComponent(tp reflect.Type, totalBits int) T {
	val := len(r.Components)
	if val >= totalBits {
		panic(fmt.Sprintf("maximum of %d component types exceeded", totalBits))
	}
	id := T(val)
	r.Components[tp], r.Types[id] = id, tp
	r.Used.Set(uint8(id), true)
	if r.isRelation(tp) {
		r.IsRelation.Set(uint8(id), true)
	}
	return id
}

func (r *componentRegistry[T]) unregisterLastComponent() {
	id := T(len(r.Components) - 1)
	tp, _ := r.ComponentType(id)
	delete(r.Components, tp)
	r.Types[id] = nil
	r.Used.Set(uint8(id), false)
	r.IsRelation.Set(uint8(id), false)
}

func (r *componentRegistry[T]) isRelation(tp reflect.Type) bool {
	if tp.NumField() == 0 {
		return false
	}
	field := tp.Field(0)
	return field.Type == relationType && field.Name == relationType.Name()
}
