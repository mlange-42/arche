package listener_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/mlange-42/arche/listener"
	"github.com/stretchr/testify/assert"
)

func TestCallback(t *testing.T) {
	w := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&w)
	velID := ecs.ComponentID[Velocity](&w)

	evt := []ecs.EntityEvent{}
	ls := listener.NewCallback(
		func(w *ecs.World, e ecs.EntityEvent) {
			evt = append(evt, e)
		},
		event.All,
		posID,
	)
	w.SetListener(&ls)

	assert.Equal(t, event.All, ls.Subscriptions())
	assert.Equal(t, ecs.All(posID), *ls.Components())

	w.NewEntity(posID)
	assert.Equal(t, 1, len(evt))

	w.NewEntity(velID)
	assert.Equal(t, 1, len(evt))
}

func ExampleCallback() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)

	ls := listener.NewCallback(
		func(w *ecs.World, e ecs.EntityEvent) {
			// Print the EventType bits of the event.
			fmt.Printf("   EventType: %08b\n", e.EventTypes)
		},
		// Subscribe to all events.
		event.All,
	)
	world.SetListener(&ls)

	fmt.Println("Create entity")
	e := world.NewEntity()

	fmt.Println("Add component")
	world.Add(e, posID)

	fmt.Println("Remove component")
	world.Remove(e, posID)

	fmt.Println("Remove entity")
	world.RemoveEntity(e)

	fmt.Println("Create entity with component(s)")
	e = world.NewEntity(posID)

	fmt.Println("Remove entity with component(s)")
	world.RemoveEntity(e)
	// Output: Create entity
	//    EventType: 00000001
	// Add component
	//    EventType: 00000100
	// Remove component
	//    EventType: 00001000
	// Remove entity
	//    EventType: 00000010
	// Create entity with component(s)
	//    EventType: 00000101
	// Remove entity with component(s)
	//    EventType: 00001010
}
