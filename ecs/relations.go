package ecs

// Relations provides access to entity [Relation] targets.
//
// Access it using [World.Relations].
type Relations struct {
	world *World
}

// Get returns the target entity for an entity relation.
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//
// See [Relation] for details and examples.
func (r *Relations) Get(entity Entity, comp ID) Entity {
	return r.world.getRelation(entity, comp)
}

// GetUnchecked returns the target entity for an entity relation.
//
// GetUnchecked is an optimized version of [Relations.Get].
// Does not check if the entity is alive or that the component ID is applicable.
//
// Panics when called for a removed entity, but not for a recycled entity.
func (r *Relations) GetUnchecked(entity Entity, comp ID) Entity {
	return r.world.getRelationUnchecked(entity, comp)
}

// Set sets the target entity for an entity relation.
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - when called for a removed (and potentially recycled) target.
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See [Relation] for details and examples.
func (r *Relations) Set(entity Entity, comp ID, target Entity) {
	r.world.setRelation(entity, comp, target)
}

// SetBatch sets the [Relation] target for many entities, matching a filter.
// Returns the number of affected entities.
//
// Panics:
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Relations.Set], [Relations.SetBatchQ] and [Batch.SetRelation].
func (r *Relations) SetBatch(filter Filter, comp ID, target Entity) int {
	return r.world.setRelationBatch(filter, comp, target)
}

// SetBatchQ sets the [Relation] target for many entities, matching a filter.
// Returns a query over all affected entities.
//
// Panics:
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Relations.Set], [Relations.SetBatch] and [Batch.SetRelation].
func (r *Relations) SetBatchQ(filter Filter, comp ID, target Entity) Query {
	return r.world.setRelationBatchQuery(filter, comp, target)
}

// Exchange adds and removes components in one pass.
// This is more efficient than subsequent use of [World.Add] and [World.Remove].
// In contrast to [World.Exchange], it allows to also set a relation target.
//
// When a [Relation] component is removed and another one is added,
// the target entity of the relation is set to zero if no target is given.
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called for a missing relation component.
//   - when called for a component that is not a relation.
//   - when called without any components to add or remove. Use [World.Relations] instead.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also the generic variants under [github.com/mlange-42/arche/generic.Exchange].
func (r *Relations) Exchange(entity Entity, add []ID, rem []ID, relation ID, target Entity) {
	r.world.exchange(entity, add, rem, relation, true, target)
}

// ExchangeBatch exchanges components for many entities, matching a filter.
// Returns the number of affected entities.
// In contrast to [Batch.Exchange], it allows to also set a relation target.
//
// When a [Relation] component is removed and another one is added,
// the target entity of the relation is set to zero if no target is given.
//
// Panics:
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called for a missing relation component.
//   - when called for a component that is not a relation.
//   - when called without any components to add or remove. Use [Batch.SetRelation] instead.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.Exchange], [Batch.ExchangeQ], [Relations.ExchangeBatch] and [World.Exchange].
func (r *Relations) ExchangeBatch(filter Filter, add []ID, rem []ID, relation ID, target Entity) int {
	return r.world.exchangeBatch(filter, add, rem, relation, true, target)
}

// ExchangeBatchQ exchanges components for many entities, matching a filter.
// It returns a query over the affected entities.
// In contrast to [Batch.ExchangeQ], it allows to also set a relation target.
//
// When a [Relation] component is removed and another one is added,
// the target entity of the relation is set to zero if no target is given.
//
// Panics:
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called for a missing relation component.
//   - when called for a component that is not a relation.
//   - when called without any components to add or remove. Use [Batch.SetRelationQ] instead.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Batch.Exchange], [Batch.ExchangeQ], [Relations.ExchangeBatch] and [World.Exchange].
func (r *Relations) ExchangeBatchQ(filter Filter, add []ID, rem []ID, relation ID, target Entity) Query {
	return r.world.exchangeBatchQuery(filter, add, rem, relation, true, target)
}
