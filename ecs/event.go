package ecs

import "github.com/mlange-42/arche/ecs/event"

// EntityEvent contains information about ECS operations like component and relation changes to an [Entity].
//
// To receive change events, register a [Listener] with [World.SetListener].
//
// # Event types & subscriptions
//
// Events notified are entity creation and removal, component addition and removal,
// and change of relations and their targets.
//
// Event types that are subscribed are determined by [Listener].Subscriptions.
// Events that cover multiple types (e.g. entity creation and component addition) are only notified once.
// Field EventTypes contains the [event.Subscription] bits of covered event types.
//
// See sub-package [event] and the [event.Subscription] constants for event types and subscription logic.
//
// # Event scheduling
//
// Events are emitted immediately after the change is applied.
//
// Except for removed entities, events are always fired when the [World] is in an unlocked state.
// Events for removed entities are fired right before removal of the entity,
// to allow for inspection of its components.
// Therefore, the [World] is in a locked state during entity removal events.
//
// Events for batch-creation of entities using a [Builder] are fired after all entities are created.
// For batch methods that return a [Query], events are fired after the [Query] is closed (or fully iterated).
// This allows the [World] to be in an unlocked state, and notifies after potential entity initialization.
type EntityEvent struct {
	OldRelation, NewRelation *ID                // Old and new relation component ID. No relation is indicated by nil.
	AddedIDs, RemovedIDs     []ID               // Components added and removed. DO NOT MODIFY! Get the current components with [World.Ids].
	Added, Removed           Mask               // Masks indicating changed components (additions and removals).
	Entity                   Entity             // The entity that was changed.
	OldTarget                Entity             // Old relation target entity. Get the new target with [World.Relations] and [Relations.Get].
	EventTypes               event.Subscription // Bit mask of event types. See [event.Subscription].
}

// Contains returns whether the event's types contain the given type/subscription bit.
func (e *EntityEvent) Contains(bit event.Subscription) bool {
	return e.EventTypes.Contains(bit)
}

// Listener interface for listening to [EntityEvent] notifications
// on ECS operations like entity creation and removal, component addition and removal, and relation changes.
//
// A listener can be added to a [World] with [World.SetListener].
//
// # Subscriptions
//
// Listeners can subscribe to one or more event types via method Subscriptions.
// Further, subscriptions can be restricted to one or more components via method Components.
//
// See sub-package [event] and the [event.Subscription] constants for event types and subscription logic.
//
// # See also
//
// See [EntityEvent] for more details.
// See package [github.com/mlange-42/arche/listener] for Listener implementations.
type Listener interface {
	// Notify the listener about a subscribed event.
	Notify(world *World, evt EntityEvent)
	// Subscriptions to one or more event types.
	Subscriptions() event.Subscription
	// Components the listener subscribes to. Listening to all components indicated by nil.
	Components() *Mask
}

// testListener for [EntityEvent]s.
type testListener struct {
	Callback  func(world *World, e EntityEvent)
	Subscribe event.Subscription
}

// newTestListener creates a new [CallbackListener] that subscribes to all event types.
func newTestListener(callback func(world *World, e EntityEvent)) testListener {
	return testListener{
		Callback:  callback,
		Subscribe: event.EntityCreated | event.EntityRemoved | event.ComponentAdded | event.ComponentRemoved | event.RelationChanged | event.TargetChanged,
	}
}

// Notify the listener.
func (l *testListener) Notify(world *World, e EntityEvent) {
	l.Callback(world, e)
}

// Subscriptions of the listener in terms of event types.
func (l *testListener) Subscriptions() event.Subscription {
	return l.Subscribe
}

// Components the listener subscribes to.
// Will be notified about changes on any (not all!) of the components in the mask.
func (l *testListener) Components() *Mask {
	return nil
}
