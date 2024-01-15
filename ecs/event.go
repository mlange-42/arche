package ecs

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
	Entity                   Entity       // The entity that was changed.
	OldMask                  Mask         // The old component masks. Get the new mask with [World.Mask].
	Added, Removed           []ID         // Components added and removed. DO NOT MODIFY! Get the current components with [World.Ids].
	OldRelation, NewRelation *ID          // Old and new relation component ID. No relation is indicated by nil.
	OldTarget                Entity       // Old relation target entity. Get the new target with [World.Relations] and [Relations.Get].
	EventTypes               Subscription // Bit mask of event types. See [Subscription].
}

// EntityAdded reports whether the entity was newly added.
func (e *EntityEvent) EntityAdded() bool {
	return e.EventTypes.Contains(EntityCreated)
}

// EntityRemoved reports whether the entity was removed.
func (e *EntityEvent) EntityRemoved() bool {
	return e.EventTypes.Contains(EntityRemoved)
}

// Subscription bits for a [Listener]
type Subscription uint8

const (
	// EntityCreated subscription bit
	EntityCreated Subscription = 0b00000001
	// EntityRemoved subscription bit
	EntityRemoved Subscription = 0b00000010
	// ComponentAdded subscription bit
	ComponentAdded Subscription = 0b00000100
	// ComponentRemoved subscription bit
	ComponentRemoved Subscription = 0b000001000
	// RelationChanged subscription bit
	RelationChanged Subscription = 0b000010000
	// TargetChanged subscription bit
	TargetChanged Subscription = 0b000100000
)

func subscription(entityCreated, entityRemoved, componentAdded, componentRemoved, relationChanged, targetChanged bool) Subscription {
	var bits Subscription = 0
	if entityCreated {
		bits |= EntityCreated
	}
	if entityRemoved {
		bits |= EntityRemoved
	}
	if componentAdded {
		bits |= ComponentAdded
	}
	if componentRemoved {
		bits |= ComponentRemoved
	}
	if relationChanged {
		bits |= RelationChanged
	}
	if targetChanged {
		bits |= TargetChanged
	}
	return bits
}

// Contains checks whether the argument is in this Subscription.
func (s Subscription) Contains(bit Subscription) bool {
	return (bit & s) == bit
}

// ContainsAny checks whether any of the argument's bits is in this Subscription.
func (s Subscription) ContainsAny(bit Subscription) bool {
	return (bit & s) != 0
}

// Listener interface
type Listener interface {
	Notify(e EntityEvent)
	Subscriptions() Subscription
}

// CallbackListener for [EntityEvent]s.
type CallbackListener struct {
	Callback  func(e EntityEvent)
	Subscribe Subscription
}

// NewCallbackListener creates a new [CallbackListener] that subscribes to all event types.
func NewCallbackListener(callback func(e EntityEvent)) CallbackListener {
	return CallbackListener{
		Callback:  callback,
		Subscribe: EntityCreated | EntityRemoved | ComponentAdded | ComponentRemoved | RelationChanged | TargetChanged,
	}
}

// Notify the listener
func (l *CallbackListener) Notify(e EntityEvent) {
	l.Callback(e)
}

// Subscriptions of the listener
func (l *CallbackListener) Subscriptions() Subscription {
	return l.Subscribe
}
