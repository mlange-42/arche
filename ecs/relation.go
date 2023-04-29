package ecs

import "reflect"

var relationType = reflect.TypeOf((*Relation)(nil)).Elem()

// Relation must be embedded as first field into components that represent an entity relation.
//
//	type ChildOf struct {
//		Relation
//	}
type Relation struct{}
