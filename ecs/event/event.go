// Package event contains a mask type and bit switches for listener subscriptions.
//
// See also ecs.Listener and ecs.EntityEvent.
package event

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
const (
	// All subscriptions
	All Subscription = EntityCreated | EntityRemoved | ComponentAdded | ComponentRemoved | RelationChanged | TargetChanged
	// Entities subscription for entity creation or removal
	Entities Subscription = EntityCreated | EntityRemoved
	// Components subscription for component addition or removal
	Components Subscription = ComponentAdded | ComponentRemoved
	// Relations subscription for relation and target changes
	Relations Subscription = RelationChanged | TargetChanged
)

// Subscription bits for an ecs.Listener
type Subscription uint8

// Contains checks whether the argument is in this Subscription.
func (s Subscription) Contains(bit Subscription) bool {
	return (bit & s) == bit
}

// ContainsAny checks whether any of the argument's bits is in this Subscription.
func (s Subscription) ContainsAny(bit Subscription) bool {
	return (bit & s) != 0
}
