package ecs_test

import (
	"fmt"

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
