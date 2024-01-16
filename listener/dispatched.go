package listener

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
)

// Dispatched event listener.
type Dispatched struct {
	listeners []ecs.Listener
	events    event.Subscription
}

// NewDispatched returns a new [Dispatched] listener with sub-listeners.
func NewDispatched(listeners ...ecs.Listener) Dispatched {
	var events event.Subscription = 0
	for _, l := range listeners {
		events |= l.Subscriptions()
	}
	return Dispatched{
		listeners: listeners,
		events:    events,
	}
}

// AddListener adds a listener to this listener.
func (l *Dispatched) AddListener(ls ecs.Listener) {
	l.listeners = append(l.listeners, ls)
	l.events |= ls.Subscriptions()
}

// Notify the listener.
func (l *Dispatched) Notify(e ecs.EntityEvent) {
	for _, ls := range l.listeners {
		if ls.Subscriptions().ContainsAny(e.EventTypes) {
			ls.Notify(e)
		}
	}
}

// Subscriptions of the listener.
func (l *Dispatched) Subscriptions() event.Subscription {
	return l.events
}
