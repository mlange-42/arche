// Demonstrates the use of a EntityEvent listener,
// notifying on changes to the component composition of entities.
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

// Listener type
type Listener struct {
	World *ecs.World
}

// Listen is called on entity changes.
func (l *Listener) Listen(evt ecs.EntityEvent) {
	// Just prints out what the event is about.
	// This could also be a method of a type that manages events.
	// Could use e.g. filters to distribute events to interested/registered systems.
	if evt.EntityAdded() {
		fmt.Printf("Entity added, has components %v\n", l.World.Ids(evt.Entity))
	} else if evt.EntityRemoved() {
		fmt.Printf("Entity removed, had components %v\n", l.World.Ids(evt.Entity))
	} else {
		fmt.Printf("Entity changed, has components %v\n", l.World.Ids(evt.Entity))
	}
}

func main() {
	// Create a World.
	world := ecs.NewWorld()
	listener := Listener{World: &world}
	wrapper := ecs.NewListener(listener.Listen)

	// Get component IDs
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)
	rotID := ecs.ComponentID[Rotation](&world)

	// Register a listener function.
	world.SetListener(wrapper)

	// Create/manipulate/delete entities and observe the listener's output
	e0 := world.NewEntity(posID)
	e1 := world.NewEntity(posID, velID)

	world.Add(e0, velID)
	world.Add(e1, rotID)

	world.Remove(e1, posID, velID)

	world.RemoveEntity(e0)
	world.RemoveEntity(e1)
}
