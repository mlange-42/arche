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
	add         []ecs.ID
	remove      []ecs.ID
	hasRelation bool
	relationID  ecs.ID
	builder     ecs.Builder
	world       *ecs.World
}

// NewExchange creates a new Exchange object.
func NewExchange(w *ecs.World) *Exchange {
	return &Exchange{
		world: w,
	}
}

// WithRelation sets the [Relation] component for the Exchange.
//
// Use in conjunction with the optional target argument of [Exchange.NewEntity], [Exchange.Add], [Exchange.Remove] and [Exchange.Exchange].
//
// See [Relation] for details and examples.
func (m *Exchange) WithRelation(comp Comp) *Exchange {
	m.hasRelation = true
	m.relationID = ecs.TypeID(m.world, comp)

	m.builder = *ecs.NewBuilder(m.world, m.add...).WithRelation(m.relationID)
	return m
}

// Adds sets components to add in calls to [Exchange.Add] and [Exchange.Exchange].
//
// Create the required mask items with [T].
func (m *Exchange) Adds(add ...Comp) *Exchange {
	m.add = toIds(m.world, add)

	b := ecs.NewBuilder(m.world, m.add...)
	if m.hasRelation {
		b = b.WithRelation(m.relationID)
	}
	m.builder = *b
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
// The optional argument can be used to set the target [ecs.Entity] for the Exchange's [ecs.Relation].
// See [Exchange.WithRelation].
//
// See also [ecs.World.NewEntity].
func (m *Exchange) NewEntity(target ...ecs.Entity) ecs.Entity {
	if len(target) > 0 {
		if !m.hasRelation {
			panic("can't set target entity: Exchange has no relation")
		}
		return m.builder.New(target[0])
	}
	return m.world.NewEntity(m.add...)
}

// Add the components set via [Exchange.Adds] to the given entity.
//
// The optional argument can be used to set the target [ecs.Entity] for the Exchange's [ecs.Relation].
// See [Exchange.WithRelation].
//
// See also [ecs.World.Add].
func (m *Exchange) Add(entity ecs.Entity, target ...ecs.Entity) {
	if len(target) > 0 {
		if !m.hasRelation {
			panic("can't set target entity: Exchange has no relation")
		}
		m.world.Relations().Exchange(entity, m.add, nil, m.relationID, target[0])
	} else {
		m.world.Add(entity, m.add...)
	}
}

// Remove the components set via [Exchange.Removes] from the given entity.
//
// The optional argument can be used to set the target [ecs.Entity] for the Exchange's [ecs.Relation].
// See [Exchange.WithRelation].
//
// See also [ecs.World.Remove].
func (m *Exchange) Remove(entity ecs.Entity, target ...ecs.Entity) {
	if len(target) > 0 {
		if !m.hasRelation {
			panic("can't set target entity: Exchange has no relation")
		}
		m.world.Relations().Exchange(entity, nil, m.remove, m.relationID, target[0])
	} else {
		m.world.Remove(entity, m.remove...)
	}
}

// Exchange components on an entity.
//
// Removes the components set via [Exchange.Removes].
// Adds the components set via [Exchange.Adds].
// The optional argument can be used to set the target [ecs.Entity] for the Exchange's [ecs.Relation].
// See [Exchange.WithRelation].
// When a [ecs.Relation] component is removed and another one is added,
// the target entity of the relation is set to zero if no target is given.
//
// See also [ecs.World.Exchange].
func (m *Exchange) Exchange(entity ecs.Entity, target ...ecs.Entity) {
	if len(target) > 0 {
		if !m.hasRelation {
			panic("can't set target entity: Exchange has no relation")
		}
		m.world.Relations().Exchange(entity, m.add, m.remove, m.relationID, target[0])
	} else {
		m.world.Exchange(entity, m.add, m.remove)
	}
}

// ExchangeBatch exchanges components on many entities, matching a filter.
// Returns the number of affected entities.
//
// Removes the components set via [Exchange.Removes].
// Adds the components set via [Exchange.Adds].
// The optional argument can be used to set the target [ecs.Entity] for the Exchange's [ecs.Relation].
// See [Exchange.WithRelation].
// When a [ecs.Relation] component is removed and another one is added,
// the target entity of the relation is set to zero if no target is given.
//
// See also [ecs.Batch.Exchange] and [ecs.Batch.ExchangeQ].
func (m *Exchange) ExchangeBatch(filter ecs.Filter, target ...ecs.Entity) int {
	if len(target) > 0 {
		if !m.hasRelation {
			panic("can't set target entity: Exchange has no relation")
		}
		return m.world.Relations().ExchangeBatch(filter, m.add, m.remove, m.relationID, target[0])
	}
	return m.world.Batch().Exchange(filter, m.add, m.remove)

}
