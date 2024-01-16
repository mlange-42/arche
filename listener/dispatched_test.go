package listener_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/mlange-42/arche/listener"
	"github.com/stretchr/testify/assert"
)

func TestDispatched(t *testing.T) {
	entityEvents := []ecs.EntityEvent{}
	componentEvents := []ecs.EntityEvent{}

	l1 := listener.NewCallback(event.Entities, func(evt ecs.EntityEvent) { entityEvents = append(entityEvents, evt) })
	l2 := listener.NewCallback(event.Components, func(evt ecs.EntityEvent) { componentEvents = append(componentEvents, evt) })

	ls := listener.NewDispatched(&l1)
	ls.AddListener(&l2)

	assert.Equal(t, event.Entities|event.Components, ls.Subscriptions())

	ls.Notify(ecs.EntityEvent{EventTypes: event.EntityCreated | event.ComponentAdded})
	ls.Notify(ecs.EntityEvent{EventTypes: event.EntityCreated | event.RelationChanged})
	ls.Notify(ecs.EntityEvent{EventTypes: event.Relations})

	assert.Equal(t, 2, len(entityEvents))
	assert.Equal(t, 1, len(componentEvents))
}

type Position struct {
	X float64
	Y float64
}

func ExampleDispatched() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)

	entityListener := listener.NewCallback(
		event.EntityCreated|event.EntityRemoved,
		func(evt ecs.EntityEvent) { fmt.Println("Entity event") },
	)
	componentListener := listener.NewCallback(
		event.ComponentAdded|event.ComponentRemoved,
		func(evt ecs.EntityEvent) { fmt.Println("Component event") },
	)

	mainListener := listener.NewDispatched(
		&entityListener,
		&componentListener,
	)

	world.SetListener(&mainListener)

	// Triggers event.EntityCreated
	e := world.NewEntity()
	// Triggers event.ComponentAdded
	world.Add(e, posID)

	// Triggers event.EntityCreated and event.ComponentAdded
	_ = world.NewEntity(posID)
	// Output: Entity event
	// Component event
	// Entity event
	// Component event
}
