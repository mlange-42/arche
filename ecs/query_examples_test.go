package ecs_test

import (
	"fmt"
	"math/rand"

	"github.com/mlange-42/arche/ecs"
)

func ExampleQuery() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	filter := ecs.All(posID, velID)
	query := world.Query(filter)
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		vel := (*Velocity)(query.Get(velID))
		pos.X += vel.X
		pos.Y += vel.Y
	}
	// Output:
}

func ExampleQuery_Count() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	world.NewEntity(posID)

	query := world.Query(ecs.All(posID))
	cnt := query.Count()
	fmt.Println(cnt)

	query.Close()
	// Output: 1
}

func ExampleQuery_Close() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	world.NewEntity(posID)

	query := world.Query(ecs.All(posID))
	cnt := query.Count()
	fmt.Println(cnt)

	query.Close()
	// Output: 1
}

func ExampleQuery_EntityAt() {
	// Set up the world.
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)

	// Create entities.
	builder := ecs.NewBuilder(&world, posID)
	builder.NewBatch(100)

	// Create a random generator.
	rng := rand.New(rand.NewSource(42))

	// Register a filter (optional, but recommended).
	filter := world.Cache().Register(ecs.All(posID))

	// Query some entities.
	query := world.Query(&filter)
	// Get the number of entities in the query
	cnt := query.Count()

	// Sample some random entities.
	for i := 0; i < 5; i++ {
		randomEntity := query.EntityAt(rng.Intn(cnt))
		fmt.Println(randomEntity)
	}

	// Close the query.
	query.Close()
	// Output: {6 0}
	// {88 0}
	// {69 0}
	// {51 0}
	// {24 0}
}
