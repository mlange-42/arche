package ecs_test

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
)

// ChildOf demonstrates how to define a relation component.
// There may be more fields _after_ the embed.
type ChildOf struct {
	ecs.Relation
}

func ExampleRelation() {
	world := ecs.NewWorld()
	childID := ecs.ComponentID[ChildOf](&world)

	// Create a target/parent entity for the relation.
	parent := world.NewEntity()

	// Create a builder with a relation.
	childBuilder := ecs.NewBuilder(&world, childID).WithRelation(childID)
	// Create a child entity with a relation to the parent.
	child := childBuilder.New(parent)

	// Create a child entity with a relation to the zero entity.
	child2 := childBuilder.New(parent)

	// Get the relation target of the child.
	_ = world.Relations().Get(child, childID)

	// Set the relation target.
	world.Relations().Set(child2, childID, parent)

	// Filter for the relation with a given target.
	filter := ecs.NewRelationFilter(ecs.All(childID), parent)

	query := world.Query(&filter)
	for query.Next() {
		fmt.Println(
			query.Entity(),
			query.Relation(childID), // Get the relation target in a query.
		)
	}
	// Output: {2 0} {1 0}
	// {3 0} {1 0}
}
