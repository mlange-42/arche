package listener

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
)

// Returns whether a listener is interested in an event based on event type and component subscriptions.
//
// Argument trigger should only contain the subscription bits that triggered the event.
// I.e. subscriptions & evenTypes.
func subscribes(trigger event.Subscription, changed *ecs.Mask, subs *ecs.Mask, oldRel *ecs.ID, newRel *ecs.ID) bool {
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
	if trigger.ContainsAny(event.Entities | event.Components) {
		// Contains any other than event.RelationChanged and/or event.TargetChanged
		return subs.ContainsAny(changed)
	}
	return false
}
