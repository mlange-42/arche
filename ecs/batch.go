package ecs

// Batch is a helper to perform batched operations on the world.
//
// Create using [World.Batch].
//
// # Example
//
//	world := NewWorld()
//	posID := ComponentID[Position](&world)
//	velID := ComponentID[Velocity](&world)
//
//	world.Batch().NewEntities(10_000, posID, velID)
type Batch struct {
	world *World
}

// NewEntities creates the given number of entities [Entity].
// The given component types are added to all entities.
//
// See also [Batch.NewEntitiesQuery].
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
//
// # Example
//
//	world := NewWorld()
//	posID := ComponentID[Position](&world)
//	velID := ComponentID[Velocity](&world)
//
//	world.Batch().NewEntities(10_000, posID, velID)
func (b *Batch) NewEntities(count int, comps ...ID) {
	b.world.newEntities(count, comps...)
}

// NewEntitiesQuery creates the given number of entities [Entity].
// The given component types are added to all entities.
//
// Returns a [Query] for iterating the created entities.
// The [Query] must be closed if it is not used!
// Listener notification is delayed until the query is closed of fully iterated.
// See also [Batch.NewEntities].
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
//
// # Example
//
//	world := NewWorld()
//	posID := ComponentID[Position](&world)
//	velID := ComponentID[Velocity](&world)
//
//	query := world.Batch().NewEntitiesQuery(10_000, posID, velID)
//
//	for query.Next() {
//		// initialize components of the newly created entities
//	}
func (b *Batch) NewEntitiesQuery(count int, comps ...ID) Query {
	return b.world.newEntitiesQuery(count, comps...)
}

// NewEntitiesWith creates the given number of entities [Entity].
// The given component values are assigned to all entity.
//
// See also [Batch.NewEntitiesWithQuery].
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
//
// # Example
//
//	world := NewWorld()
//	posID := ComponentID[Position](&world)
//	velID := ComponentID[Velocity](&world)
//
//	components := []Component{
//		{ID: posID, Comp: &Position{X: 0, Y: 0}},
//		{ID: velID, Comp: &Velocity{X: 10, Y: 2}},
//	}
//
//	world.Batch().NewEntitiesWith(10_000, components...)
func (b *Batch) NewEntitiesWith(count int, comps ...Component) {
	b.world.newEntitiesWith(count, comps...)
}

// NewEntitiesWithQuery creates the given number of entities [Entity].
// The given component values are assigned to all entity.
//
// Returns a [Query] for iterating the created entities.
// The [Query] must be closed if it is not used!
// Listener notification is delayed until the query is closed of fully iterated.
// See also [Batch.NewEntitiesWith].
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Map1], etc.
//
// # Example
//
//	world := NewWorld()
//	posID := ComponentID[Position](&world)
//	velID := ComponentID[Velocity](&world)
//
//	components := []Component{
//		{ID: posID, Comp: &Position{X: 0, Y: 0}},
//		{ID: velID, Comp: &Velocity{X: 10, Y: 2}},
//	}
//
//	query := world.Batch().NewEntitiesWithQuery(10_000, components...)
//
//	for query.Next() {
//		// initialize components of the newly created entities
//	}
func (b *Batch) NewEntitiesWithQuery(count int, comps ...Component) Query {
	return b.world.newEntitiesWithQuery(count, comps...)
}

// RemoveEntities removes and recycles all entities matching a filter.
//
// Returns the number of removed entities.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// # Example
//
//	world := NewWorld()
//	posID := ComponentID[Position](&world)
//	velID := ComponentID[Velocity](&world)
//
//	world.Batch().NewEntities(10_000, posID, velID)
//
//	filter := All(posID, velID).Exclusive()
//	world.Batch().RemoveEntities(filter)
func (b *Batch) RemoveEntities(filter Filter) int {
	return b.world.removeEntities(filter)
}
