// Demonstrates the use of a EntityEvent listener,
// notifying on changes to the component composition of entities.
package main

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/mlange-42/arche/listener"
)

// Position component
type Position struct {
	X float64
	Y float64
}

// Velocity component
type Velocity struct {
	X float64
	Y float64
}

// Rotation component
type Rotation struct {
	A float64
}

func main() {
	// Create a World.
	world := ecs.NewWorld()
	listener := listener.NewCallback(
		func(world *ecs.World, evt ecs.EntityEvent) {
			fmt.Printf("Events: %08b Added: %v Removed: %v\n", evt.EventTypes, evt.Added, evt.Removed)
		},
		event.Entities|event.Components,
	)

	// Get component IDs
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)
	rotID := ecs.ComponentID[Rotation](&world)

	// Register a listener function.
	world.SetListener(&listener)

	// Create/manipulate/delete entities and observe the listener's output
	e0 := world.NewEntity(posID)
	e1 := world.NewEntity(posID, velID)

	world.Add(e0, velID)
	world.Add(e1, rotID)

	world.Remove(e1, posID, velID)

	world.RemoveEntity(e0)
	world.RemoveEntity(e1)
}
