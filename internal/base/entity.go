package base

// Entity identifier.
// Holds an entity ID and it's generation for recycling.
//
// Entities should only be created via the [World], using [World.NewEntity].
type Entity struct {
	ID  Eid
	gen uint16
}

// NewEntity creates a new Entity.
func NewEntity(id Eid) Entity {
	return Entity{id, 0}
}

// NewEntityGen creates a new Entity with a given generation.
func NewEntityGen(id Eid, gen uint16) Entity {
	return Entity{id, gen}
}

// IsZero returns whether this entity is the reserved zero entity.
func (e Entity) IsZero() bool {
	return e.ID == 0
}

// EntityIndex indicates where an entity is currently stored
type EntityIndex struct {
	Arch  *Archetype
	Index uint32
}
