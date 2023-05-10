// Demonstrates entity relations
package main

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// ChildOf component for child relation to parent.
type ChildOf struct {
	ecs.Relation
}

func main() {
	run()
	runGeneric()
}

// Entity relations using the code API
func run() {
	// Create a new world.
	world := ecs.NewWorld()

	// Get component IDs.
	childID := ecs.ComponentID[ChildOf](&world)

	// For creating entities with relations, a Builder with a relation is required.
	childBuilder := ecs.NewBuilder(&world, childID).WithRelation(childID)

	// Create parent entities.
	parent1 := world.NewEntity()
	parent2 := world.NewEntity()

	// Create an entity with a ChildOf relation to a parent entity.
	child := childBuilder.New(parent1)
	// Change the child's relation target.
	world.Relations().Set(child, childID, parent2)
	// Get the child's relation target.
	fmt.Println(world.Relations().Get(child, childID))

	// Create entities with a relation in batches.
	childBuilder.NewBatch(10, parent1)
	childBuilder.NewBatch(10, parent2)

	// Create a filter for a relation target.
	filter := ecs.NewRelationFilter(ecs.All(childID), parent1)

	// Create a query from it.
	query := world.Query(&filter)

	// Iterate the query.
	for query.Next() {
		// Get the relation target of the current entity.
		target := query.Relation(childID)
		fmt.Println(target)
	}
}

// Entity relations using the generic API
func runGeneric() {

	// Create a new world.
	world := ecs.NewWorld()

	// For creating entities with relations, a generic Map with a relation is required.
	childBuilder := generic.NewMap1[ChildOf](&world, generic.T[ChildOf]())

	// Create parent entities.
	parent1 := world.NewEntity()
	parent2 := world.NewEntity()

	// Create an entity with a ChildOf relation to a parent entity.
	child := childBuilder.New(parent1)
	// To set or get the relation target via the world, a Map is required.
	relationMap := generic.NewMap[ChildOf](&world)
	// Change the child's relation target.
	relationMap.SetRelation(child, parent2)
	// Get the child's relation target.
	fmt.Println(relationMap.GetRelation(child))

	// Create entities with a relation in batches.
	childBuilder.NewBatch(10, parent1)
	childBuilder.NewBatch(10, parent2)

	// Create a filter for a relation target.
	filter := generic.NewFilter1[ChildOf]().WithRelation(generic.T[ChildOf](), parent1)

	// Create a query from it.
	query := filter.Query(&world)

	// Iterate the query.
	for query.Next() {
		// Get the relation target of the current entity.
		target := query.Relation()
		fmt.Println(target)
	}
}
