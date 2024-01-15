package listener

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
)

// Callback listener for [EntityEvent]s.
//
// Calls function Callback on events that are contained in the Subscribe mask.
type Callback struct {
	callback func(e ecs.EntityEvent)
	events   event.Subscription
}

// NewCallback creates a new Callback listener for the given events.
func NewCallback(events event.Subscription, callback func(ecs.EntityEvent)) Callback {
	return Callback{
		callback: callback,
		events:   events,
	}
}

// Notify the listener
func (l *Callback) Notify(e ecs.EntityEvent) {
	l.callback(e)
}

// Subscriptions of the listener
func (l *Callback) Subscriptions() event.Subscription {
	return l.events
}
