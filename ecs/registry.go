package ecs

import (
	"fmt"
	"reflect"
)

// componentRegistry keeps track of component or resource IDs.
type registry struct {
	Components map[reflect.Type]uint8
	Types      []reflect.Type
	IDs        []uint8
	Used       Mask
}

// newComponentRegistry creates a new ComponentRegistry.
func newRegistry() registry {
	return registry{
		Components: map[reflect.Type]uint8{},
		Types:      make([]reflect.Type, MaskTotalBits),
		Used:       Mask{},
		IDs:        []uint8{},
	}
}

// ComponentID returns the ID for a component type, and registers it if not already registered.
// The second return value indicates if it is a newly created ID.
func (r *registry) ComponentID(tp reflect.Type) (uint8, bool) {
	if id, ok := r.Components[tp]; ok {
		return id, false
	}
	return r.registerComponent(tp, MaskTotalBits), true
}

// ComponentType returns the type of a component by ID.
func (r *registry) ComponentType(id uint8) (reflect.Type, bool) {
	return r.Types[id], r.Used.Get(ID{id: id})
}

// Count returns the total number of reserved IDs. It is the maximum ID plus 1.
func (r *registry) Count() int {
	return len(r.Components)
}

// Reset clears the registry.
func (r *registry) Reset() {
	for t := range r.Components {
		delete(r.Components, t)
	}
	for i := range r.Types {
		r.Types[i] = nil
	}
	r.Used.Reset()
	r.IDs = r.IDs[:0]
}

// registerComponent registers a components and assigns an ID for it.
func (r *registry) registerComponent(tp reflect.Type, totalBits int) uint8 {
	val := len(r.Components)
	if val >= totalBits {
		panic(fmt.Sprintf("exceeded the maximum of %d component types or resource types", totalBits))
	}
	newID := uint8(val)
	id := id(newID)
	r.Components[tp], r.Types[newID] = newID, tp
	r.Used.Set(id, true)
	r.IDs = append(r.IDs, newID)
	return newID
}

func (r *registry) unregisterLastComponent() {
	newID := uint8(len(r.Components) - 1)
	id := id(newID)
	tp, _ := r.ComponentType(newID)
	delete(r.Components, tp)
	r.Types[newID] = nil
	r.Used.Set(id, false)
	r.IDs = r.IDs[:len(r.IDs)-1]
}

// componentRegistry keeps track of component IDs.
type componentRegistry struct {
	registry
	IsRelation Mask
	IsPointer  Mask
}

// newComponentRegistry creates a new ComponentRegistry.
func newComponentRegistry() componentRegistry {
	return componentRegistry{
		registry:   newRegistry(),
		IsRelation: Mask{},
		IsPointer:  Mask{},
	}
}

// Reset clears the registry.
func (r *componentRegistry) Reset() {
	r.registry.Reset()
	r.IsRelation.Reset()
	r.IsPointer.Reset()
}

// registerComponent registers a components and assigns an ID for it.
func (r *componentRegistry) registerComponent(tp reflect.Type, totalBits int) uint8 {
	newID := r.registry.registerComponent(tp, totalBits)
	if r.isRelation(tp) {
		r.IsRelation.Set(id(newID), true)
	}
	if r.isPointer(tp) {
		r.IsPointer.Set(id(newID), true)
	}
	return newID
}

func (r *componentRegistry) unregisterLastComponent() {
	newID := uint8(len(r.Components) - 1)
	r.registry.unregisterLastComponent()
	r.IsRelation.Set(id(newID), false)
	r.IsPointer.Set(id(newID), false)
}

func (r *componentRegistry) isRelation(tp reflect.Type) bool {
	if tp.Kind() != reflect.Struct || tp.NumField() == 0 {
		return false
	}
	field := tp.Field(0)
	return field.Type == relationType && field.Name == relationType.Name()
}

// isPointer determines whether an object contains pointers that need proper garbage collection.
func (r *componentRegistry) isPointer(tp reflect.Type) bool {
	switch tp.Kind() {
	case reflect.Pointer, reflect.Interface:
		elem := tp.Elem()
		return r.isPointerRecursive(elem)
	default:
		return r.isPointerRecursive(tp)
	}
}

// isPointerRecursive determines whether an object contains pointers that need proper garbage collection.
func (r *componentRegistry) isPointerRecursive(tp reflect.Type) bool {
	switch tp.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return true
	case reflect.Struct:
		for i := 0; i < tp.NumField(); i++ {
			if r.isPointerRecursive(tp.Field(i).Type) {
				return true
			}
		}
		return false
	default:
		return false
	}
}
