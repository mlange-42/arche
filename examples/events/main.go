// Demonstrates the use of a EntityEvent through listeners.
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

// ChildOf relation component
type ChildOf struct {
	ecs.Relation
}

// EventHandler that prints information about received events.
type EventHandler struct{}

// Notify is called with all subscribed events.
func (h *EventHandler) Notify(world *ecs.World, evt ecs.EntityEvent) {
	fmt.Printf("    Type mask:     %06b\n", evt.EventTypes)
	fmt.Printf("    Entity:        %+v\n", evt.Entity)
	fmt.Printf("    Added/Removed: %+v / %+v\n", evt.AddedIDs, evt.RemovedIDs)
	fmt.Printf("    Relation:      %+v -> %+v\n", evt.OldRelation, evt.NewRelation)

	var target ecs.Entity
	if evt.NewRelation != nil {
		target = world.Relations().Get(evt.Entity, *evt.NewRelation)
	}
	fmt.Printf("    Target:        %+v -> %+v\n", evt.OldTarget, target)
}

func main() {
	// Create a World.
	world := ecs.NewWorld()

	// Get component IDs
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)
	childID := ecs.ComponentID[ChildOf](&world)

	// Create the event handler.
	handler := EventHandler{}

	// Create the listener.
	listener := listener.NewCallback(
		// Function to handle events.
		handler.Notify,
		// Subscribe to event types; could also use `event.All` here.
		event.EntityCreated|event.EntityRemoved|
			event.ComponentAdded|event.ComponentRemoved|
			event.RelationChanged|event.TargetChanged,
		// Optionally, restrict subscription to a list of component types / IDs.
		//posID, childID,
	)

	// Register a listener function.
	world.SetListener(&listener)

	fmt.Println("=== Create entity ===")
	child := world.NewEntity()

	fmt.Println("=== Create entity with normal components ===")
	parent := world.NewEntity(posID, velID)

	fmt.Println("=== Create entity with relation component ===")
	world.NewEntity(childID)

	fmt.Println("=== Add normal component to entity ===")
	world.Add(child, posID)

	fmt.Println("=== Add relation component to entity ===")
	world.Add(child, childID)

	fmt.Println("=== Change relation target ===")
	world.Relations().Set(child, childID, parent)
}
