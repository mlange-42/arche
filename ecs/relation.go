package ecs

import "reflect"

var relationType = reflect.TypeOf((*Relation)(nil)).Elem()

// Relation must be embedded as first field into components that represent an entity relation.
//
//	type ChildOf struct {
//		Relation
//	}
type Relation struct {
	entity Entity
}

// Target that is the target of the relation.
func (r *Relation) Target() Entity {
	return r.entity
}

// setTarget sets the target entity of the relation.
func (r *Relation) setTarget(e Entity) {
	r.entity = e
}
