// Package relations benchmarks different ways to create entity relations:
//   - Children know their parent entity
//   - Parents have an array of child entities
//   - Parents have a slice of child entities
//   - Children form an implicit linked lists, parents know the first child
//
// The task is to sum up a value over all children of each parent.
//
// Each parent has 10 children. Benchmarks are run for 100, 1000 and 10000 parent entities.
package relations

import "github.com/mlange-42/arche/ecs"

const numChildren = 10

// ParentData component
type ParentData struct {
	Value int
}

// ChildData component
type ChildData struct {
	Value int
}

// Child component
type Child struct {
	Parent ecs.Entity
}

// ChildList component
type ChildList struct {
	Next ecs.Entity
}

// ParentArr component
type ParentArr struct {
	Children [numChildren]ecs.Entity
}

// ParentSlice component
type ParentSlice struct {
	Children []ecs.Entity
}

// ParentList component
type ParentList struct {
	FirstChild  ecs.Entity
	NumChildren int
}
