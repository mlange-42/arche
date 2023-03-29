package ecs

import (
	"fmt"
	"reflect"
)

// Resources manage a world's resources.
//
// Access it using [World.Resources].
type Resources struct {
	registry  componentRegistry[ResID]
	resources []any
}

// newResources creates a new Resources manager.
func newResources() Resources {
	return Resources{
		registry:  newComponentRegistry(),
		resources: make([]any, MaskTotalBits),
	}
}

// Add adds a resource to the world.
// The resource should always be a pointer.
//
// Panics if there is already a resource of the given type.
//
// See also [github.com/mlange-42/arche/generic.Resource.Add] for a generic variant.
func (r *Resources) Add(id ResID, res any) {
	if r.resources[id] != nil {
		panic(fmt.Sprintf("Resource of ID %d was already added (type %v)", id, reflect.TypeOf(res)))
	}
	r.resources[id] = res
}

// Remove removes a resource from the world.
//
// Panics if there is no resource of the given type.
//
// See also [github.com/mlange-42/arche/generic.Resource.Remove] for a generic variant.
func (r *Resources) Remove(id ResID) {
	if r.resources[id] == nil {
		panic(fmt.Sprintf("Resource of ID %d is not present", id))
	}
	r.resources[id] = nil
}

// Get returns a pointer to the resource of the given type.
//
// Returns nil if there is no such resource.
//
// See also [github.com/mlange-42/arche/generic.Resource.Get] for a generic variant.
func (r *Resources) Get(id ResID) interface{} {
	return r.resources[id]
}

// Has returns whether the world has the given resource.
//
// See also [github.com/mlange-42/arche/generic.Resource.Has] for a generic variant.
func (r *Resources) Has(id ResID) bool {
	return r.resources[id] != nil
}

// reset removes all resources.
func (r *Resources) reset() {
	for i := 0; i < len(r.resources); i++ {
		r.resources[i] = nil
	}
}
