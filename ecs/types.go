package ecs

import "reflect"

// Eid is the entity identifier/index type.
type eid uint32

// ID is the component identifier type.
type ID = uint8

// Component is a component ID/pointer pair.
//
// It is a helper for [World.Assign] and [World.NewEntityWith].
// It is not related to how components are implemented in Arche.
type Component struct {
	ID   ID
	Comp interface{}
}

// componentType is a component ID with a data type
type componentType struct {
	ID
	Type reflect.Type
}
