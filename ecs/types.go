package ecs

import (
	"reflect"
	"unsafe"
)

// ID defines the format for the components identifier
type ID uint32

// Component is a component ID with a reference struct
type Component struct {
	ID
	Component interface{}
}

// ComponentType is a component ID with a data type
type ComponentType struct {
	ID
	Type reflect.Type
}

// ComponentPointer is a component ID with a pointer in a storage
type ComponentPointer struct {
	ID
	Pointer unsafe.Pointer
}
