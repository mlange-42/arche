package ecs

import (
	"fmt"
	"reflect"
)

type resources struct {
	registry  componentRegistry[ResID]
	resources []any
}

func newResources() resources {
	return resources{
		registry:  newComponentRegistry(),
		resources: make([]any, MaskTotalBits),
	}
}

// Add adds a resource to the world.
// The resource should always be a pointer.
//
// Panics if there is already a resource of the given type.
func (r *resources) Add(id ResID, res any) {
	if r.resources[id] != nil {
		panic(fmt.Sprintf("Resource of ID %d was already added (type %v)", id, reflect.TypeOf(res)))
	}
	r.resources[id] = res
}

// Remove removes a resource from the world.
//
// Panics if there is no resource of the given type.
func (r *resources) Remove(id ResID) {
	if r.resources[id] == nil {
		panic(fmt.Sprintf("Resource of ID %d is not present", id))
	}
	r.resources[id] = nil
}

// Get returns a pointer to the resource of the given type.
//
// Returns nil if there is no such resource.
func (r *resources) Get(id ResID) interface{} {
	return r.resources[id]
}

// Has returns whether the world has the given resource.
func (r *resources) Has(id ResID) bool {
	return r.resources[id] != nil
}

// Reset removes all resources.
func (r *resources) Reset() {
	for i := 0; i < len(r.resources); i++ {
		r.resources[i] = nil
	}
}
