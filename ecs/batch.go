package ecs

// Batch is a helper to perform batched operations on the world.
//
// Create using [World.Batch].
type Batch struct {
	world *World
}

// Add adds components to many entities, matching a filter.
// Returns the number of affected entities.
//
// Panics:
//   - when called with components that can't be added because they are already present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.AddQ] and [World.Add].
func (b *Batch) Add(filter Filter, comps ...ID) int {
	return b.world.exchangeBatch(filter, comps, nil, ID{}, false, Entity{})
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
	return b.world.exchangeBatchQuery(filter, comps, nil, ID{}, false, Entity{})
}

// Remove removes components from many entities, matching a filter.
// Returns the number of affected entities.
//
// Panics:
//   - when called with components that can't be removed because they are not present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.RemoveQ] and [World.Remove].
func (b *Batch) Remove(filter Filter, comps ...ID) int {
	return b.world.exchangeBatch(filter, nil, comps, ID{}, false, Entity{})
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
	return b.world.exchangeBatchQuery(filter, nil, comps, ID{}, false, Entity{})
}

// SetRelation sets the [Relation] target for many entities, matching a filter.
// Returns the number of affected entities.
//
// Entities that match the filter but already have the desired target entity are not processed,
// and no events are emitted for them.
//
// Panics:
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Relations.Set] and [Relations.SetBatch].
func (b *Batch) SetRelation(filter Filter, comp ID, target Entity) int {
	return b.world.setRelationBatch(filter, comp, target)
}

// SetRelationQ sets the [Relation] target for many entities, matching a filter.
// It returns a query over the affected entities.
//
// Entities that match the filter but already have the desired target entity are not processed,
// not included in the query, and no events are emitted for them.
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
// Returns the number of affected entities.
//
// When a [Relation] component is removed and another one is added,
// the target entity of the relation is reset to zero.
//
// Panics:
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.ExchangeQ] and [World.Exchange].
// For batch-exchange with a relation target, see [Relations.ExchangeBatch].
func (b *Batch) Exchange(filter Filter, add []ID, rem []ID) int {
	return b.world.exchangeBatch(filter, add, rem, ID{}, false, Entity{})
}

// ExchangeQ exchanges components for many entities, matching a filter.
// It returns a query over the affected entities.
//
// When a [Relation] component is removed and another one is added,
// the target entity of the relation is reset to zero.
//
// Panics:
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.Exchange] and [World.Exchange].
// For batch-exchange with a relation target, see [Relations.ExchangeBatchQ].
func (b *Batch) ExchangeQ(filter Filter, add []ID, rem []ID) Query {
	return b.world.exchangeBatchQuery(filter, add, rem, ID{}, false, Entity{})
}

// RemoveEntities removes and recycles all entities matching a filter.
// Returns the number of removed entities.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
//
// Unlike with the other batch operations, it is not easily possible to provide a query version RemoveEntitiesQ.
// However, one can simply query with the same filter before calling RemoveEntities.
//
// See also [World.RemoveEntity]
func (b *Batch) RemoveEntities(filter Filter) int {
	return b.world.removeEntities(filter)
}
