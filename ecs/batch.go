package ecs

// Batch is a helper to perform batched operations on the world.
//
// Create using [World.Batch].
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
func (b *Batch) NewEntitiesWithQuery(count int, comps ...Component) Query {
	return b.world.newEntitiesWithQuery(count, comps...)
}

// RemoveEntities removes and recycles all entities matching a filter.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
func (b *Batch) RemoveEntities(filter Filter) {
	b.world.removeEntities(filter)
}
