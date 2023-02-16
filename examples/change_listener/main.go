// Demonstrates the use of a ChangeEvent listener.
package main

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
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

// Listen is called on entity changes.
func Listen(evt ecs.ChangeEvent) {
	// Just prints out what the event is about.
	// This could also be a method of a type that manages events.
	// Could use e.g. filters to distribute events to interested/registered systems.
	if evt.EntityAdded() {
		fmt.Printf("Entity added, has components %v\n", evt.Current)
	} else if evt.EntityRemoved() {
		fmt.Printf("Entity removed, had components %v\n", evt.Current)
	} else {
		fmt.Printf("Entity changed, has components %v\n", evt.Current)
	}
}

func main() {
	// Create a World.
	world := ecs.NewWorld()

	// Get component IDs
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)
	rotID := ecs.ComponentID[Rotation](&world)

	// Register a listener function.
	world.RegisterListener(Listen)

	// Create/manipulate/delete entities and observe the listener's output
	e0 := world.NewEntity(posID)
	e1 := world.NewEntity(posID, velID)

	world.Add(e0, velID)
	world.Add(e1, rotID)

	world.Remove(e1, posID, velID)

	world.RemoveEntity(e0)
	world.RemoveEntity(e1)
}
