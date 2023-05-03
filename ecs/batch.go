package ecs

// Batch is a helper to perform batched operations on the world.
//
// Create using [World.Batch].
type Batch struct {
	world *World
}

// Add adds components to many entities, matching a filter.
//
// If the callback argument is given, it is called with a [Query] over the affected entities,
// one Query for each affected archetype.
//
// Panics:
//   - when called with components that can't be added because they are already present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [World.Exchange].
func (b *Batch) Add(filter Filter, callback func(Query), comps ...ID) {
	b.world.exchangeBatch(filter, comps, nil, callback)
}

// Remove removes components from many entities, matching a filter.
//
// If the callback argument is given, it is called with a [Query] over the affected entities,
// one Query for each affected archetype.
//
// Panics:
//   - when called with components that can't be removed because they are not present.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [World.Exchange].
func (b *Batch) Remove(filter Filter, callback func(Query), comps ...ID) {
	b.world.exchangeBatch(filter, nil, comps, callback)
}

// SetRelation sets the [Relation] target for many entities, matching a filter.
//
// If the callback argument is given, it is called with a [Query] over the affected entities,
// one Query for each affected archetype.
//
// Panics:
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [World.SetRelation].
func (b *Batch) SetRelation(filter Filter, comp ID, target Entity, callback func(Query)) {
	b.world.setRelationBatch(filter, comp, target, callback)
}

// Exchange exchanges components for many entities, matching a filter.
//
// If the callback argument is given, it is called with a [Query] over the affected entities,
// one Query for each affected archetype.
//
// Panics:
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [World.Exchange].
func (b *Batch) Exchange(filter Filter, add []ID, rem []ID, callback func(Query)) {
	b.world.exchangeBatch(filter, add, rem, callback)
}

// RemoveEntities removes and recycles all entities matching a filter.
//
// Returns the number of removed entities.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
func (b *Batch) RemoveEntities(filter Filter) int {
	return b.world.removeEntities(filter)
}
