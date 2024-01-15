package ecs

import "github.com/mlange-42/arche/ecs/event"

// EntityEvent contains information about component and relation changes to an [Entity].
//
// To receive change events, register a function func(e *EntityEvent) with [World.SetListener].
//
// Events notified are entity creation, removal, changes to the component composition and change of relation targets.
// Events are emitted immediately after the change is applied.
//
// Except for removed entities, events are always fired when the [World] is in an unlocked state.
// Events for removed entities are fired right before removal of the entity,
// to allow for inspection of it's components.
// Therefore, the [World] is in a locked state during entity removal events.
//
// Events for batch-creation of entities using a [Builder] are fired after all entities are created.
// For batch methods that return a [Query], events are fired after the [Query] is closed (or fully iterated).
// This allows the [World] to be in an unlocked state, and notifies after potential entity initialization.
type EntityEvent struct {
	Entity                   Entity             // The entity that was changed.
	OldMask                  Mask               // The old component masks. Get the new mask with [World.Mask].
	Added, Removed           []ID               // Components added and removed. DO NOT MODIFY! Get the current components with [World.Ids].
	OldRelation, NewRelation *ID                // Old and new relation component ID. No relation is indicated by nil.
	OldTarget                Entity             // Old relation target entity. Get the new target with [World.Relations] and [Relations.Get].
	EventTypes               event.Subscription // Bit mask of event types. See [Subscription].
}

// EntityAdded reports whether the entity was newly added.
func (e *EntityEvent) EntityAdded() bool {
	return e.EventTypes.Contains(event.EntityCreated)
}

// EntityRemoved reports whether the entity was removed.
func (e *EntityEvent) EntityRemoved() bool {
	return e.EventTypes.Contains(event.EntityRemoved)
}

func subscription(entityCreated, entityRemoved, componentAdded, componentRemoved, relationChanged, targetChanged bool) event.Subscription {
	var bits event.Subscription = 0
	if entityCreated {
		bits |= event.EntityCreated
	}
	if entityRemoved {
		bits |= event.EntityRemoved
	}
	if componentAdded {
		bits |= event.ComponentAdded
	}
	if componentRemoved {
		bits |= event.ComponentRemoved
	}
	if relationChanged {
		bits |= event.RelationChanged
	}
	if targetChanged {
		bits |= event.TargetChanged
	}
	return bits
}

// Listener interface
type Listener interface {
	Notify(e EntityEvent)
	Subscriptions() event.Subscription
}

// testListener for [EntityEvent]s.
type testListener struct {
	Callback  func(e EntityEvent)
	Subscribe event.Subscription
}

// newTestListener creates a new [CallbackListener] that subscribes to all event types.
func newTestListener(callback func(e EntityEvent)) testListener {
	return testListener{
		Callback:  callback,
		Subscribe: event.EntityCreated | event.EntityRemoved | event.ComponentAdded | event.ComponentRemoved | event.RelationChanged | event.TargetChanged,
	}
}

// Notify the listener
func (l *testListener) Notify(e EntityEvent) {
	l.Callback(e)
}

// Subscriptions of the listener
func (l *testListener) Subscriptions() event.Subscription {
	return l.Subscribe
}
