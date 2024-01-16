package listener

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
)

// Callback listener for ecs.EntityEvent.
//
// Calls a function on events that are contained in the subscription mask.
type Callback struct {
	callback      func(w *ecs.World, e ecs.EntityEvent)
	events        event.Subscription
	components    ecs.Mask
	hasComponents bool
}

// NewCallback creates a new Callback listener for the given events.
//
// Subscribes to the specified events with changes on the specified components.
// If no component IDs are given, is subscribes to all components.
func NewCallback(callback func(*ecs.World, ecs.EntityEvent), events event.Subscription, components ...ecs.ID) Callback {
	return Callback{
		callback:      callback,
		events:        events,
		components:    ecs.All(components...),
		hasComponents: len(components) > 0,
	}
}

// Notify the listener.
func (l *Callback) Notify(w *ecs.World, e ecs.EntityEvent) {
	l.callback(w, e)
}

// Subscriptions of the listener.
func (l *Callback) Subscriptions() event.Subscription {
	return l.events
}

// Components the listener subscribes to.
func (l *Callback) Components() *ecs.Mask {
	if l.hasComponents {
		return &l.components
	}
	return nil
}
