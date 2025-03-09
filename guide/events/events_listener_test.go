package main

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
)

// PositionAddedListener listens to additions of a Position component.
type PositionAddedListener struct {
	subs  event.Subscription
	comps ecs.Mask
}

// NewPositionAddedListener creates a new PositionAddedListener.
func NewPositionAddedListener(world *ecs.World) PositionAddedListener {
	posID := ecs.ComponentID[Position](world)
	return PositionAddedListener{
		subs:  event.ComponentAdded,
		comps: ecs.All(posID),
	}
}

// Notify the listener about a subscribed event.
func (l *PositionAddedListener) Notify(world *ecs.World, evt ecs.EntityEvent) {
	fmt.Println("Position component added to entity ", evt.Entity)
}

// Subscriptions to one or more event types.
func (l *PositionAddedListener) Subscriptions() event.Subscription {
	return l.subs
}

// Components the listener subscribes to. Listening to all components indicated by nil.
func (l *PositionAddedListener) Components() *ecs.Mask {
	return &l.comps
}
