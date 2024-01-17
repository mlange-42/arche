package listener

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
)

// Returns whether a listener is interested in an event based on event type and component subscriptions.
//
// Argument trigger should only contain the subscription bits that triggered the event.
// I.e. subscriptions & evenTypes.
func subscribes(trigger event.Subscription, added *ecs.Mask, removed *ecs.Mask, subs *ecs.Mask, oldRel *ecs.ID, newRel *ecs.ID) bool {
	if trigger == 0 {
		return false
	}
	if subs == nil {
		// No component subscriptions
		return true
	}
	if trigger.ContainsAny(event.Relations) {
		// Contains event.RelationChanged and/or event.TargetChanged
		if (oldRel != nil && subs.Get(*oldRel)) || (newRel != nil && subs.Get(*newRel)) {
			return true
		}
	}
	if trigger.ContainsAny(event.EntityCreated | event.ComponentAdded) {
		// Contains additions-like types
		if added != nil && subs.ContainsAny(added) {
			return true
		}
	}
	if trigger.ContainsAny(event.EntityRemoved | event.ComponentRemoved) {
		// Contains additions-like types
		if removed != nil && subs.ContainsAny(removed) {
			return true
		}
	}
	return false
}
