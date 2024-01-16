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

// TestListener type
type TestListener struct {
	World *ecs.World
}

// Notify is called on entity changes.
func (l *TestListener) Notify(evt ecs.EntityEvent) {
	// Just prints out what the event is about.
	// This could also be a method of a type that manages events.
	// Could use e.g. filters to distribute events to interested/registered systems.
	if evt.Contains(event.EntityCreated) {
		fmt.Printf("Entity added, has components %v\n", l.World.Ids(evt.Entity))
	} else if evt.Contains(event.EntityRemoved) {
		fmt.Printf("Entity removed, had components %v\n", l.World.Ids(evt.Entity))
	} else {
		fmt.Printf("Entity changed, has components %v\n", l.World.Ids(evt.Entity))
	}
}

func main() {
	// Create a World.
	world := ecs.NewWorld()
	ls := TestListener{World: &world}
	wrapper := listener.NewCallback(event.Entities|event.Components, ls.Notify)

	// Get component IDs
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)
	rotID := ecs.ComponentID[Rotation](&world)

	// Register a listener function.
	world.SetListener(&wrapper)

	// Create/manipulate/delete entities and observe the listener's output
	e0 := world.NewEntity(posID)
	e1 := world.NewEntity(posID, velID)

	world.Add(e0, velID)
	world.Add(e1, rotID)

	world.Remove(e1, posID, velID)

	world.RemoveEntity(e0)
	world.RemoveEntity(e1)
}
