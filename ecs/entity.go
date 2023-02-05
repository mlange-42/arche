package ecs

// Entity identifier
type Entity struct {
	id  uint32
	gen uint16
}

func newEntity(id int) Entity {
	return Entity{uint32(id), 0}
}
