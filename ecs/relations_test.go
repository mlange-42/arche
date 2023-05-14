package ecs_test

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
)

func ExampleRelations() {
	world := ecs.NewWorld()

	relID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()
	child := world.NewEntity(relID)

	world.Relations().Set(child, relID, parent)
	fmt.Println(world.Relations().Get(child, relID))
	// Output: {1 0}
}

func ExampleRelations_SetBatch() {
	world := ecs.NewWorld()

	relID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()
	ecs.NewBuilder(&world, relID).NewBatch(100)

	filter := ecs.All(relID)
	world.Relations().SetBatch(filter, relID, parent)
	// Output:
}

func ExampleRelations_SetBatchQ() {
	world := ecs.NewWorld()

	relID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()
	ecs.NewBuilder(&world, relID).NewBatch(100)

	filter := ecs.All(relID)
	query := world.Relations().SetBatchQ(filter, relID, parent)
	fmt.Println(query.Count())
	query.Close()
	// Output: 100
}
