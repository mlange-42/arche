// Demonstrates the use of a ChangeEvent listener.
package main

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
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

	// Register a listener function.
	world.RegisterListener(Listen)

	// Create/manipulate/delete entities and observe the listener's output
	e0, _ := generic.NewEntity1[Position](&world)
	e1, _, _ := generic.NewEntity2[Position, Velocity](&world)

	generic.Add1[Velocity](&world, e0)
	generic.Add1[Rotation](&world, e1)

	generic.Remove2[Position, Velocity](&world, e1)

	world.RemEntity(e0)
	world.RemEntity(e1)
}
