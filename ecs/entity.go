package ecs

// Entity identifier
// TODO: Store ID and generation in a single uint64?
type Entity struct {
	id  ID
	gen uint16
}

func newEntity(id int) Entity {
	return Entity{ID(id), 0}
}

// IsZero returns whether this entity is the reserved zero entity
func (e Entity) IsZero() bool {
	return e.id == 0
}

type entityIndex struct {
	arch  int
	index uint32
}
