package ecs

// ID defines the format for the components identifier
type ID uint32

// Component is a component ID with a reference struct
type Component struct {
	ID
	reference interface{}
}
