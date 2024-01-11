package ecs

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Reflection type of an [Entity].
var entityType = reflect.TypeOf(Entity{})

// Size of an [Entity] in memory, in bytes.
var entitySize uint32 = uint32(entityType.Size())

// Size of an [entityIndex] in memory.
var entityIndexSize uint32 = uint32(reflect.TypeOf(entityIndex{}).Size())

// Entity identifier.
// Holds an entity ID and it's generation for recycling.
//
// Entities are only created via the [World], using [World.NewEntity] or [World.NewEntityWith].
// Batch creation of entities is possible via [Builder].
//
// ⚠️ Important:
// Entities are intended to be stored and passed around via copy, not via pointers!
// The zero value should be used to indicate "nil", and can be checked with [Entity.IsZero].
type Entity struct {
	id  eid    // Entity ID
	gen uint32 // Entity generation
}

// newEntity creates a new Entity.
func newEntity(id eid) Entity {
	return Entity{id, 0}
}

// newEntityGen creates a new Entity with a given generation.
func newEntityGen(id eid, gen uint32) Entity {
	return Entity{id, gen}
}

// IsZero returns whether this entity is the reserved zero entity.
func (e Entity) IsZero() bool {
	return e.id == 0
}

// ID of the entity. For serialization purposes.
func (e Entity) ID() uint32 {
	return uint32(e.id)
}

// Gen returns the generation of the entity. For serialization purposes.
func (e Entity) Gen() uint32 {
	return e.gen
}

// MarshalJSON returns a JSON representation of the entity, for serialization purposes.
func (e Entity) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"ID\": %d, \"Gen\": %d}", e.id, e.gen)), nil
}

type entityHelper struct {
	ID  eid
	Gen uint32
}

// UnmarshalJSON into an entity.
func (e *Entity) UnmarshalJSON(data []byte) error {
	helper := entityHelper{}
	if err := json.Unmarshal(data, &helper); err != nil {
		return err
	}
	e.id = helper.ID
	e.gen = helper.Gen
	return nil
}

// entityIndex indicates where an entity is currently stored.
type entityIndex struct {
	arch  *archetype // Entity's current archetype
	index uint32     // Entity's current index in the archetype
}
