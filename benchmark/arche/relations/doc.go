// Package relations benchmarks different ways to create entity relations:
//   - Children know their parent entity
//   - Parents have an array of child entities
//   - Parents have a slice of child entities
//   - Children form an implicit linked lists, parents know the first child
//
// The task is to sum up a value over all children of each parent.
//
// The benchmarks are implemented in a way that only one components needs to be fetched for each parent and each child.
//
// Each parent has 10 children. Benchmarks are run for 100, 1000 and 10k and 100k parent entities.
package relations

import "github.com/mlange-42/arche/ecs"

const numArrChildren = 10

// Child component
type Child struct {
	Parent ecs.Entity
	Value  int
}

// ChildRelation component
type ChildRelation struct {
	ecs.Relation
	Value int
}

// ChildList component
type ChildList struct {
	Next  ecs.Entity
	Value int
}

// ParentArr component
type ParentArr struct {
	Children [numArrChildren]ecs.Entity
	Value    int
}

// ParentSlice component
type ParentSlice struct {
	Children []ecs.Entity
	Value    int
}

// ParentList component
type ParentList struct {
	FirstChild ecs.Entity
	Value      int
}
