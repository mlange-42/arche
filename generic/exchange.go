package generic

import "github.com/mlange-42/arche/ecs"

// Exchange is a helper for adding, removing and exchanging components.
//
// Set it up like this:
//
//	ex := NewExchange(&world).
//	          Adds(T3[A, B, C]()...).
//	          Removes(T2[D, E]()...)
//
// For only adding or removing components, see also [Map1], [Map2] etc.
type Exchange struct {
	add    []ecs.ID
	remove []ecs.ID
	world  *ecs.World
}

// NewExchange creates a new Exchange object.
func NewExchange(w *ecs.World) *Exchange {
	return &Exchange{
		world: w,
	}
}

// Adds sets components to add in calls to [Exchange.Add] and [Exchange.Exchange].
//
// Create the required mask items with [T].
func (m *Exchange) Adds(add ...Comp) *Exchange {
	m.add = toIds(m.world, add)
	return m
}

// Removes sets components to remove in calls to [Exchange.Remove] and [Exchange.Exchange].
//
// Create the required mask items with [T].
func (m *Exchange) Removes(remove ...Comp) *Exchange {
	m.remove = toIds(m.world, remove)
	return m
}

// NewEntity creates a new [ecs.Entity] with the components set via [Exchange.Adds].
//
// See also [ecs.World.NewEntity].
func (m *Exchange) NewEntity() ecs.Entity {
	entity := m.world.NewEntity(m.add...)
	return entity
}

// Add the components set via [Exchange.Adds] to the given entity.
//
// See also [ecs.World.Add].
func (m *Exchange) Add(entity ecs.Entity) {
	m.world.Add(entity, m.add...)
}

// Remove the components set via [Exchange.Removes] from the given entity.
//
// See also [ecs.World.Remove].
func (m *Exchange) Remove(entity ecs.Entity) {
	m.world.Remove(entity, m.remove...)
}

// Exchange components on an entity.
//
// Removes the components set via [Exchange.Removes].
// Adds the components set via [Exchange.Adds].
//
// When a [Relation] component is removed and another one is added,
// the target entity of the relation remains unchanged.
//
// See also [ecs.World.Exchange].
func (m *Exchange) Exchange(entity ecs.Entity) {
	m.world.Exchange(entity, m.add, m.remove)
}
