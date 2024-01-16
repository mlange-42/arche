// Package event contains a mask type and bit switches for listener subscriptions.
//
// See also ecs.Listener and ecs.EntityEvent.
package event

// Subscription bits for individual events
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

// Subscription bits for groups of events
const (
	// Entities subscription for entity creation or removal
	Entities Subscription = EntityCreated | EntityRemoved
	// Components subscription for component addition or removal
	Components Subscription = ComponentAdded | ComponentRemoved
	// Relations subscription for relation and target changes
	Relations Subscription = RelationChanged | TargetChanged
	// All subscriptions
	All Subscription = Entities | Components | Relations
)

// Subscription bits for an ecs.Listener
type Subscription uint8

// Contains checks whether all the argument's bits are contained in this Subscription.
func (s Subscription) Contains(bits Subscription) bool {
	return (bits & s) == bits
}

// ContainsAny checks whether any of the argument's bits are contained in this Subscription.
func (s Subscription) ContainsAny(bits Subscription) bool {
	return (bits & s) != 0
}
