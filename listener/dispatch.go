package listener

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
)

// Dispatch event listener.
//
// Dispatches events to sub-listeners and manages subscription automatically, based on their settings.
// Sub-listeners should not alter their subscriptions or components after being added.
//
// To make it possible for systems to add listeners, Dispatch can be added to the [ecs.World] as a resource.
type Dispatch struct {
	listeners     []ecs.Listener     // Sub-listeners to dispatch events to.
	events        event.Subscription // Subscribed event types.
	components    ecs.Mask           // Subscribed components.
	hasComponents bool               // Whether there is a restriction to components.
}

// NewDispatch returns a new [Dispatch] listener with the given sub-listeners.
func NewDispatch(listeners ...ecs.Listener) Dispatch {
	var events event.Subscription = 0
	var components ecs.Mask
	hasComponents := true
	for _, l := range listeners {
		events |= l.Subscriptions()
		cmp := l.Components()
		if cmp == nil {
			hasComponents = false
		} else {
			components = components.Or(cmp)
		}
	}
	return Dispatch{
		listeners:     listeners,
		events:        events,
		components:    components,
		hasComponents: hasComponents,
	}
}

// AddListener adds a sub-listener to this listener.
func (l *Dispatch) AddListener(ls ecs.Listener) {
	l.listeners = append(l.listeners, ls)
	l.events |= ls.Subscriptions()

	cmp := ls.Components()
	if cmp == nil {
		l.hasComponents = false
	} else {
		l.components = l.components.Or(cmp)
	}
}

// Notify the listener.
func (l *Dispatch) Notify(world *ecs.World, evt ecs.EntityEvent) {
	for _, ls := range l.listeners {
		trigger := ls.Subscriptions() & evt.EventTypes
		if trigger != 0 && subscribes(trigger, &evt.Added, &evt.Removed, ls.Components(), evt.OldRelation, evt.NewRelation) {
			ls.Notify(world, evt)
		}
	}
}

// Subscriptions of the listener.
func (l *Dispatch) Subscriptions() event.Subscription {
	return l.events
}

// Components the listener subscribes to.
func (l *Dispatch) Components() *ecs.Mask {
	if l.hasComponents {
		return &l.components
	}
	return nil
}
