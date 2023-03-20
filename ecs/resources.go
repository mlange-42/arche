package ecs

import (
	"fmt"
	"reflect"
)

type resources struct {
	registry  componentRegistry
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
func (r *resources) Add(res any) {
	tp := reflect.TypeOf(res).Elem()
	id := r.registry.ComponentID(tp)
	if r.resources[id] != nil {
		panic(fmt.Sprintf("Resource of type %v was already added", tp))
	}
	r.resources[id] = res
}

// Get returns a pointer to the resource of the given type.
func (r *resources) Get(id ID) interface{} {
	res := r.resources[id]
	if res != nil {
		return res
	}
	panic(fmt.Sprintf("No resource for id %d", id))
}

// Has returns whether the world has the given resource.
func (r *resources) Has(id ID) bool {
	return r.resources[id] != nil
}
