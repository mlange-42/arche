package base

import (
	"reflect"
	"unsafe"
)

// ID is the component identifier type
type ID uint8

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
