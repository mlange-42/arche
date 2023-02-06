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

type entityIndex struct {
	arch  *Archetype
	index uint32
}
