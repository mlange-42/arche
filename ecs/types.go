package ecs

import (
	"reflect"
	"unsafe"
)

// ID defines the format for the components identifier
type ID uint32

// component is a component ID with a reference struct
type component struct {
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
