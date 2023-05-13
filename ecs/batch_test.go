package ecs_test

import "github.com/mlange-42/arche/ecs"

func ExampleBatch() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	builder := ecs.NewBuilder(&world, posID, velID)
	builder.NewBatch(100)

	world.Batch().Remove(ecs.All(posID, velID), velID)
	world.Batch().RemoveEntities(ecs.All(posID))
}

func ExampleBatch_Add() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	builder := ecs.NewBuilder(&world, posID)
	builder.NewBatch(100)

	filter := ecs.All(posID)
	world.Batch().Add(filter, velID)
}

func ExampleBatch_Remove() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	builder := ecs.NewBuilder(&world, posID, velID)
	builder.NewBatch(100)

	filter := ecs.All(posID, velID)
	world.Batch().Remove(filter, velID)
}

func ExampleBatch_Exchange() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	builder := ecs.NewBuilder(&world, posID)
	builder.NewBatch(100)

	filter := ecs.All(posID)
	world.Batch().Exchange(
		filter,          // Filter
		[]ecs.ID{velID}, // Add components
		[]ecs.ID{posID}, // Remove components
	)
}

func ExampleBatch_RemoveEntities() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)

	builder := ecs.NewBuilder(&world, posID)
	builder.NewBatch(100)

	filter := ecs.All(posID)
	world.Batch().RemoveEntities(filter)
}

func ExampleBatch_SetRelation() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	childID := ecs.ComponentID[ChildOf](&world)

	target := world.NewEntity()

	builder := ecs.NewBuilder(&world, posID, childID)
	builder.NewBatch(100)

	filter := ecs.All(childID)
	world.Batch().SetRelation(filter, childID, target, nil)
}
