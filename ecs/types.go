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

// bitMask is a bitmask.
type bitMask = base.BitMask

// MaskTotalBits is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = base.MaskTotalBits

// Mask is a mask for a combination of components.
type Mask = base.Mask

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
