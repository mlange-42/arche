package ecs

import "reflect"

// Eid is the entity identifier/index type
type eid uint32

// ID is the component identifier type
type ID = uint8

// Component is a component ID/pointer pair.
type Component struct {
	ID
	Comp interface{}
}

// componentType is a component ID with a data type
type componentType struct {
	ID
	Type reflect.Type
}
