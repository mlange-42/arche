package main

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/mlange-42/arche/listener"
)

// Position component.
type Position struct {
	X float64
	Y float64
}

// Heading component
type Heading struct {
	Angle float64
}

func TestSubscriptions(t *testing.T) {
	_ = event.EntityCreated
	_ = event.EntityRemoved
	_ = event.ComponentAdded
	_ = event.ComponentRemoved
	_ = event.RelationChanged
	_ = event.TargetChanged
}

func TestCombineSubscriptions(t *testing.T) {
	subs := event.EntityCreated | event.EntityRemoved
	_ = subs
}

func TestComponentSubscriptions(t *testing.T) {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	compSubs := ecs.All(posID, headID)
	_ = compSubs
}

func TestListeners(t *testing.T) {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	// Listener for all entity creation and entity removal events.
	entityListener := listener.NewCallback(
		// A function we want to call on notification.
		func(w *ecs.World, ee ecs.EntityEvent) { /* Do something here */ },
		// Subscription to event types.
		event.EntityCreated|event.EntityRemoved,
	)

	posOrHeadAddedListener := listener.NewCallback(
		// A function we want to call on notification.
		func(w *ecs.World, ee ecs.EntityEvent) { /* Do something here */ },
		// Subscription to event types.
		event.ComponentAdded,
		// Subscription is restricted to these component types.
		posID, headID,
	)

	// Create the dispatch listener from both sub-listeners.
	dispatch := listener.NewDispatch(
		&entityListener,
		&posOrHeadAddedListener,
	)
	// Set it as the world's listener.
	world.SetListener(&dispatch)
}
