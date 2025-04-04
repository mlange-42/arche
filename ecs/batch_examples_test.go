package ecs_test

import "github.com/mlange-42/arche/ecs"

func ExampleBatch() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	world.Batch().New(100, posID, velID)

	world.Batch().Remove(ecs.All(posID, velID), velID)
	world.Batch().RemoveEntities(ecs.All(posID))
	// Output:
}

func ExampleBatch_New() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	world.Batch().New(100, posID, velID)
	// Output:
}

func ExampleBatch_NewQ() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	query := world.Batch().NewQ(100, posID, velID)

	for query.Next() {
		pos := (*Position)(query.Get(posID))
		pos.X = 100
	}
	// Output:
}

func ExampleBatch_Add() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	world.Batch().New(100, posID)

	filter := ecs.All(posID)
	world.Batch().Add(filter, velID)
	// Output:
}

func ExampleBatch_AddQ() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	world.Batch().New(100, posID)

	filter := ecs.All(posID)
	query := world.Batch().AddQ(filter, velID)

	for query.Next() {
		pos := (*Position)(query.Get(posID))
		pos.X = 100
	}
	// Output:
}

func ExampleBatch_Remove() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	world.Batch().New(100, posID, velID)

	filter := ecs.All(posID, velID)
	world.Batch().Remove(filter, velID)
	// Output:
}

func ExampleBatch_RemoveQ() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	world.Batch().New(100, posID, velID)

	filter := ecs.All(posID, velID)
	query := world.Batch().RemoveQ(filter, velID)

	for query.Next() {
		pos := (*Position)(query.Get(posID))
		pos.X = 100
	}
	// Output:
}

func ExampleBatch_Exchange() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	world.Batch().New(100, posID)

	filter := ecs.All(posID)
	world.Batch().Exchange(
		filter,          // Filter
		[]ecs.ID{velID}, // Add components
		[]ecs.ID{posID}, // Remove components
	)
	// Output:
}

func ExampleBatch_ExchangeQ() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	world.Batch().New(100, posID)

	filter := ecs.All(posID)
	query := world.Batch().ExchangeQ(
		filter,          // Filter
		[]ecs.ID{velID}, // Add components
		[]ecs.ID{posID}, // Remove components
	)

	for query.Next() {
		vel := (*Velocity)(query.Get(velID))
		vel.X = 100
	}
	// Output:
}

func ExampleBatch_RemoveEntities() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)

	world.Batch().New(100, posID)

	filter := ecs.All(posID)
	world.Batch().RemoveEntities(filter)
	// Output:
}

func ExampleBatch_SetRelation() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	childID := ecs.ComponentID[ChildOf](&world)

	target := world.NewEntity()

	world.Batch().New(100, posID, childID)

	filter := ecs.All(childID)
	world.Batch().SetRelation(filter, childID, target)
	// Output:
}

func ExampleBatch_SetRelationQ() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	childID := ecs.ComponentID[ChildOf](&world)

	target := world.NewEntity()

	world.Batch().New(100, posID, childID)

	filter := ecs.All(childID)
	query := world.Batch().SetRelationQ(filter, childID, target)

	for query.Next() {
		pos := (*Position)(query.Get(posID))
		pos.X = 100
	}
	// Output:
}
