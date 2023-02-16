package ecs

import "reflect"

// Size of an Entity in memory.
var entitySize = reflect.TypeOf(Entity{}).Size()

// Size of an entitySize in memory.
var entityIndexSize = reflect.TypeOf(entityIndex{}).Size()

// Entity identifier.
// Holds an entity ID and it's generation for recycling.
//
// Entities should only be created via the [World], using [World.NewEntity].
type Entity struct {
	id  eid
	gen uint16
}

// newEntity creates a new Entity.
func newEntity(id eid) Entity {
	return Entity{id, 0}
}

// newEntityGen creates a new Entity with a given generation.
func newEntityGen(id eid, gen uint16) Entity {
	return Entity{id, gen}
}

// IsZero returns whether this entity is the reserved zero entity.
func (e Entity) IsZero() bool {
	return e.id == 0
}

// entityIndex indicates where an entity is currently stored
type entityIndex struct {
	arch  *archetype
	index uint32
}
