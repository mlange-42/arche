package ecs

import "reflect"

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
