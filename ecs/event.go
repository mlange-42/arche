package ecs

// EntityEvent contains information about component changes to an [Entity].
//
// To receive change events, register a function func(e *EntityEvent) with [World.SetListener].
//
// Events notified are entity creation, removal and changes to the component composition.
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
//
// Note that the event pointer received by the listener function should not be stored,
// as the instance behind the pointer might be reused for further notifications.
type EntityEvent struct {
	Entity                   Entity // The entity that was changed.
	OldMask                  Mask   // The old and new component masks.
	Added, Removed           []ID   // Components added and removed. DO NOT MODIFY!
	OldRelation, NewRelation *ID    // Old and new relation component ID. No relation id indicated by nil.
	OldTarget                Entity // Old and new target entity.
	AddedRemoved             int8   // Whether the entity itself was added (> 0), removed (< 0), or only changed (= 0).
	RelationChanged          bool   // Whether the relation component has changed.
	TargetChanged            bool   // Whether the relation target has changed. Will be false if the relation component changes, but the target does not.
}

// EntityAdded reports whether the entity was newly added.
func (e *EntityEvent) EntityAdded() bool {
	return e.AddedRemoved > 0
}

// EntityRemoved reports whether the entity was removed.
func (e *EntityEvent) EntityRemoved() bool {
	return e.AddedRemoved < 0
}
