package ecs

import (
	"encoding/json"
	"reflect"
)

// Reflection type of an [Entity].
var entityType = reflect.TypeOf(Entity{})

// Size of an [Entity] in memory, in bytes.
var entitySize uint32 = uint32(entityType.Size())

// Size of an [entityIndex] in memory.
var entityIndexSize uint32 = uint32(reflect.TypeOf(entityIndex{}).Size())

// Entity identifier.
// Holds an entity ID and its generation for recycling.
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

// ID returns the entity's ID.
func (e Entity) ID() uint32 {
	return uint32(e.id)
}

// Generation returns the entity's generation.
func (e Entity) Generation() uint32 {
	return e.gen
}

// MarshalJSON returns a JSON representation of the entity, for serialization purposes.
//
// The JSON representation of an entity is a two-element array of entity ID and generation.
func (e Entity) MarshalJSON() ([]byte, error) {
	arr := [2]uint32{uint32(e.id), e.gen}
	jsonValue, _ := json.Marshal(arr) // Ignore the error, as we can be sure this works.
	return jsonValue, nil
}

// UnmarshalJSON into an entity.
//
// For serialization purposes only. Do not use this to create entities!
func (e *Entity) UnmarshalJSON(data []byte) error {
	arr := [2]uint32{}
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	e.id = eid(arr[0])
	e.gen = arr[1]

	return nil
}

// entityIndex indicates where an entity is currently stored.
type entityIndex struct {
	arch  *archetype // Entity's current archetype
	index uint32     // Entity's current index in the archetype
}
