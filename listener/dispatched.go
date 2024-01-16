package listener

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
)

// Dispatched event listener.
//
// Dispatches events to sub-listeners and manages subscription automatically.
//
// Sub-listeners should not alter their subscriptions or components after adding them.
type Dispatched struct {
	listeners     []ecs.Listener
	events        event.Subscription
	components    ecs.Mask
	hasComponents bool
}

// NewDispatched returns a new [Dispatched] listener with the given sub-listeners.
func NewDispatched(listeners ...ecs.Listener) Dispatched {
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
	return Dispatched{
		listeners:     listeners,
		events:        events,
		components:    components,
		hasComponents: hasComponents,
	}
}

// AddListener adds a sub-listener to this listener.
func (l *Dispatched) AddListener(ls ecs.Listener) {
	l.listeners = append(l.listeners, ls)
	l.events |= ls.Subscriptions()

	cmp := l.Components()
	if cmp == nil {
		l.hasComponents = false
	} else {
		l.components = l.components.Or(cmp)
	}
}

// Notify the listener.
func (l *Dispatched) Notify(evt ecs.EntityEvent) {
	for _, ls := range l.listeners {
		if ls.Subscriptions().ContainsAny(evt.EventTypes) &&
			subscribes(evt.EventTypes, &evt.Changed, ls.Components(), evt.OldRelation, evt.NewRelation) {
			ls.Notify(evt)
		}
	}
}

// Subscriptions of the listener.
func (l *Dispatched) Subscriptions() event.Subscription {
	return l.events
}

// Components the listener subscribes to.
func (l *Dispatched) Components() *ecs.Mask {
	if l.hasComponents {
		return &l.components
	}
	return nil
}

func subscribes(evtType event.Subscription, changed *ecs.Mask, subs *ecs.Mask, oldRel *ecs.ID, newRel *ecs.ID) bool {
	if event.Relations.Contains(evtType) {
		return subs == nil || (oldRel != nil && subs.Get(*oldRel)) || (newRel != nil && subs.Get(*newRel))
	}
	return subs == nil || subs.ContainsAny(changed)
}
