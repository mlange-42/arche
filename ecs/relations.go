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
//
// Panics:
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [Relations.Set], [Relations.SetBatchQ] and [Batch.SetRelation].
func (r *Relations) SetBatch(filter Filter, comp ID, target Entity) {
	r.world.setRelationBatch(filter, comp, target)
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
