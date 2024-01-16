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
		func(e ecs.EntityEvent) {
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

	ls := listener.NewCallback(
		func(e ecs.EntityEvent) {
			fmt.Println(e)
		},
		event.Entities|event.Components,
	)
	world.SetListener(&ls)

	world.NewEntity()
}
