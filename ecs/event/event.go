// Package event contains a mask type and bit switches for listener subscriptions.
//
// See also [github.com/mlange-42/arche/ecs.Listener] and [github.com/mlange-42/arche/ecs.EntityEvent].
package event

// Subscription bits for an [github.com/mlange-42/arche/ecs.Listener]
type Subscription uint8

// Contains checks whether all the argument's bits are contained in this Subscription.
func (s Subscription) Contains(bits Subscription) bool {
	return (bits & s) == bits
}

// ContainsAny checks whether any of the argument's bits are contained in this Subscription.
func (s Subscription) ContainsAny(bits Subscription) bool {
	return (bits & s) != 0
}

// Subscription bits for individual events.
const (
	// EntityCreated subscription bit.
	//
	// Without component subscription:
	//   - Creation of an entity with or without any components
	// With component subscription:
	//   - Creation of an entity with any of the given components
	EntityCreated Subscription = 1

	// EntityRemoved subscription bit.
	//
	// Without component subscription:
	//   - Removal of an entity, with or without any components
	// With component subscription:
	//   - Removal of an entity with any of the given components
	EntityRemoved Subscription = 1 << 1

	// ComponentAdded subscription bit.
	//
	// Without component subscription:
	//   - Addition of any component(s) to an entity
	//   - Creation of an entity with any components
	// With component subscription:
	//   - Addition of any of the given components to an entity
	//   - Creation of an entity with any of the given components
	ComponentAdded Subscription = 1 << 2

	// ComponentRemoved subscription bit.
	//
	// Without component subscription:
	//   - Removal of any component(s) from an entity
	//   - Removal of an entity with any components
	// With component subscription:
	//   - Removal of any of the given components from an entity
	//   - Removal of an entity with any of the given components
	ComponentRemoved Subscription = 1 << 3

	// RelationChanged subscription bit.
	//
	// Without component subscription:
	//   - Addition of a relation component
	//   - Removal of a relation component
	//   - Exchange of a relation component with another
	// With component subscription:
	//   - Addition of any of the given relation components
	//   - Removal of any of the given relation components
	//   - Exchange if any of the two is among the given components
	RelationChanged Subscription = 1 << 4

	// TargetChanged subscription bit.
	//
	// Without component subscription:
	//   - Whenever RelationChanged is triggered
	//   - Change of the target entity of any relation component
	// With component subscription:
	//   - Whenever RelationChanged is triggered
	//   - Change of the target entity of any of the given (relation) components
	TargetChanged Subscription = 1 << 5
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
