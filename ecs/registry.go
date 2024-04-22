package ecs

import (
	"fmt"
	"reflect"
)

// componentRegistry keeps track of component IDs.
type componentRegistry struct {
	Components map[reflect.Type]uint8
	Types      []reflect.Type
	IDs        []uint8
	Used       Mask
	IsRelation Mask
}

// newComponentRegistry creates a new ComponentRegistry.
func newComponentRegistry() componentRegistry {
	return componentRegistry{
		Components: map[reflect.Type]uint8{},
		Types:      make([]reflect.Type, MaskTotalBits),
		Used:       Mask{},
		IsRelation: Mask{},
		IDs:        []uint8{},
	}
}

// ComponentID returns the ID for a component type, and registers it if not already registered.
// The second return value indicates if it is a newly created ID.
func (r *componentRegistry) ComponentID(tp reflect.Type) (uint8, bool) {
	if id, ok := r.Components[tp]; ok {
		return id, false
	}
	return r.registerComponent(tp, MaskTotalBits), true
}

// ComponentType returns the type of a component by ID.
func (r *componentRegistry) ComponentType(id uint8) (reflect.Type, bool) {
	return r.Types[id], r.Used.Get(ID{id: id})
}

// ComponentType returns the type of a component by ID.
func (r *componentRegistry) Count() int {
	return len(r.Components)
}

// registerComponent registers a components and assigns an ID for it.
func (r *componentRegistry) registerComponent(tp reflect.Type, totalBits int) uint8 {
	val := len(r.Components)
	if val >= totalBits {
		panic(fmt.Sprintf("exceeded the maximum of %d component types or resource types", totalBits))
	}
	newID := uint8(val)
	id := id(newID)
	r.Components[tp], r.Types[newID] = newID, tp
	r.Used.Set(id, true)
	if r.isRelation(tp) {
		r.IsRelation.Set(id, true)
	}
	r.IDs = append(r.IDs, newID)
	return newID
}

func (r *componentRegistry) unregisterLastComponent() {
	newID := uint8(len(r.Components) - 1)
	id := id(newID)
	tp, _ := r.ComponentType(newID)
	delete(r.Components, tp)
	r.Types[newID] = nil
	r.Used.Set(id, false)
	r.IsRelation.Set(id, false)
	r.IDs = r.IDs[:len(r.IDs)-1]
}

func (r *componentRegistry) isRelation(tp reflect.Type) bool {
	if tp.Kind() != reflect.Struct || tp.NumField() == 0 {
		return false
	}
	field := tp.Field(0)
	return field.Type == relationType && field.Name == relationType.Name()
}
