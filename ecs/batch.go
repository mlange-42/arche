package ecs

// Batch is a helper to perform batched operations on the world.
//
// Create using [World.Batch].
type Batch struct {
	world *World
}

// Add adds components to many entities, matching a filter.
//
// Panics:
//   - when called with components that can't be added because they are already present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.AddQ] and [World.Add].
func (b *Batch) Add(filter Filter, comps ...ID) {
	b.world.exchangeBatch(filter, comps, nil)
}

// AddQ adds components to many entities, matching a filter.
// It returns a query over the affected entities.
//
// Panics:
//   - when called with components that can't be added because they are already present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.Add] and [World.Add].
func (b *Batch) AddQ(filter Filter, comps ...ID) Query {
	return b.world.exchangeBatchQuery(filter, comps, nil)
}

// Remove removes components from many entities, matching a filter.
//
// Panics:
//   - when called with components that can't be removed because they are not present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.RemoveQ] and [World.Remove].
func (b *Batch) Remove(filter Filter, comps ...ID) {
	b.world.exchangeBatch(filter, nil, comps)
}

// RemoveQ removes components from many entities, matching a filter.
// It returns a query over the affected entities.
//
// Panics:
//   - when called with components that can't be removed because they are not present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.Remove] and [World.Remove].
func (b *Batch) RemoveQ(filter Filter, comps ...ID) Query {
	return b.world.exchangeBatchQuery(filter, nil, comps)
}

// SetRelation sets the [Relation] target for many entities, matching a filter.
//
// Panics:
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Relations.Set] and [Relations.SetBatch].
func (b *Batch) SetRelation(filter Filter, comp ID, target Entity) {
	b.world.setRelationBatch(filter, comp, target)
}

// SetRelationQ sets the [Relation] target for many entities, matching a filter.
// It returns a query over the affected entities.
//
// Panics:
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Relations.Set] and [Relations.SetBatch].
func (b *Batch) SetRelationQ(filter Filter, comp ID, target Entity) Query {
	return b.world.setRelationBatchQuery(filter, comp, target)
}

// Exchange exchanges components for many entities, matching a filter.
//
// Panics:
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.ExchangeQ] and [World.Exchange].
func (b *Batch) Exchange(filter Filter, add []ID, rem []ID) {
	b.world.exchangeBatch(filter, add, rem)
}

// ExchangeQ exchanges components for many entities, matching a filter.
// It returns a query over the affected entities.
//
// Panics:
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.Exchange] and [World.Exchange].
func (b *Batch) ExchangeQ(filter Filter, add []ID, rem []ID) Query {
	return b.world.exchangeBatchQuery(filter, add, rem)
}

// RemoveEntities removes and recycles all entities matching a filter.
//
// Returns the number of removed entities.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// See also [World.RemoveEntity]
func (b *Batch) RemoveEntities(filter Filter) int {
	return b.world.removeEntities(filter)
}
