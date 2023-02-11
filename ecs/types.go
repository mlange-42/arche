package ecs

import (
	"internal/base"
	"reflect"
	"unsafe"
)

// eid is the entity identifier/index type
type eid = base.Eid

// ID is the component identifier type
type ID = base.ID

// Component is a Component ID/Component pointer pair
type Component struct {
	ID
	Component interface{}
}

// componentType is a component ID with a data type
type componentType struct {
	ID
	Type reflect.Type
}

// componentPointer is a component ID with a pointer in a storage
type componentPointer struct {
	ID
	Pointer unsafe.Pointer
}
