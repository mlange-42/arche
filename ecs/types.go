package ecs

import "reflect"

// Eid is the entity identifier/index type.
type eid uint32

// ID is the component identifier type.
type ID struct {
	id uint8
}

func id(id uint8) ID {
	return ID{id: id}
}

// ResID is the resource identifier type.
type ResID struct {
	id uint8
}

// Component is a component ID/pointer pair.
//
// It is a helper for [World.Assign], [World.NewEntityWith] and [NewBuilderWith].
// It is not related to how components are implemented in Arche.
type Component struct {
	ID   ID          // Component ID.
	Comp interface{} // The component, as a pointer to a struct.
}

// componentType is a component ID with a data type
type componentType struct {
	ID
	Type reflect.Type
}

// CompInfo provides information about a registered component.
// Returned by [ComponentInfo].
type CompInfo struct {
	ID         ID
	Type       reflect.Type
	IsRelation bool
}

// EntityDump is a dump of the entire entity data of the world.
//
// See [World.DumpEntities] and [World.LoadEntities].
type EntityDump struct {
	Entities  []Entity // Entities in the World's entity pool.
	Alive     []uint32 // IDs of all alive entities in query iteration order.
	Next      uint32   // The next free entity of the World's entity pool.
	Available uint32   // The number of allocated and available entities in the World's entity pool.
}
