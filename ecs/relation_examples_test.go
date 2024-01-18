package ecs_test

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
)

type ChildOf struct {
	ecs.Relation
}

func ExampleRelation() {
	world := ecs.NewWorld()
	childID := ecs.ComponentID[ChildOf](&world)

	// Create a target/parent entity for the relation.
	parent := world.NewEntity()

	// Create a child entity with a relation to the parent.
	childBuilder := ecs.NewBuilder(&world, childID).WithRelation(childID)
	child := childBuilder.New(parent)

	// Get the relation target of the child.
	_ = world.Relations().Get(child, childID)

	// Filter for the relation.
	filter := ecs.NewRelationFilter(ecs.All(childID), parent)

	query := world.Query(&filter)
	for query.Next() {
		fmt.Println(
			query.Entity(),
			query.Relation(childID),
		)
	}
	// Output: {2 0} {1 0}
}
