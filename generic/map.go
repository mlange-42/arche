package generic

import "github.com/mlange-42/arche/ecs"

// Map provides a type-safe way to access a component type by entity ID.
//
// Create one with [NewMap].
type Map[T any] struct {
	id    ecs.ID
	world *ecs.World
}

// NewMap creates a new [Map] for a component type.
//
// Map provides a type-safe way to access a component type by entity ID.
//
// See also [ecs.World.Get], [ecs.World.Has] and [ecs.World.Set].
func NewMap[T any](w *ecs.World) Map[T] {
	return Map[T]{
		id:    ecs.ComponentID[T](w),
		world: w,
	}
}

// ID returns the component ID for this Map.
func (m *Map[T]) ID() ecs.ID {
	return m.id
}

// Get returns a pointer to the component of the given entity.
//
// ⚠️ Important: The obtained pointer should not be stored persistently!
//
// See [Map.HasUnchecked] for an optimized version for static entities.
// See also [ecs.World.Get].
func (m *Map[T]) Get(entity ecs.Entity) *T {
	return (*T)(m.world.Get(entity, m.id))
}

// GetUnchecked returns a pointer to the component of the given entity.
//
// ⚠️ Important: The obtained pointer should not be stored persistently!
//
// GetUnchecked is an optimized version of [Map.Get],
// for cases where entities are static or checked with [ecs.World.Alive] in user code.
//
// See also [ecs.World.GetUnchecked].
func (m *Map[T]) GetUnchecked(entity ecs.Entity) *T {
	return (*T)(m.world.GetUnchecked(entity, m.id))
}

// Has returns whether the entity has the component.
//
// See [Map.HasUnchecked] for an optimized version for static entities.
// See also [ecs.World.Has].
func (m *Map[T]) Has(entity ecs.Entity) bool {
	return m.world.Has(entity, m.id)
}

// HasUnchecked returns whether the entity has the component.
//
// HasUnchecked is an optimized version of [Map.Has],
// for cases where entities are static or checked with [ecs.World.Alive] in user code.
//
// See also [ecs.World.HasUnchecked].
func (m *Map[T]) HasUnchecked(entity ecs.Entity) bool {
	return m.world.HasUnchecked(entity, m.id)
}

// Set overwrites the component for the given entity.
//
// Panics if the entity does not have a component of that type.
//
// See also [ecs.World.Set].
func (m *Map[T]) Set(entity ecs.Entity, comp *T) *T {
	p := (*T)(m.world.Get(entity, m.id))
	if p == nil {
		panic("can't copy component into entity that has no such component type")
	}
	*p = *comp
	return p
}

// GetRelation returns the target entity for the given entity and the Map's relation component.
//
// Panics:
//   - if the entity does not have a component of that type.
//   - if the component is not an [ecs.Relation].
//   - if the entity has been removed/recycled.
//
// See also [ecs.World.GetRelation].
func (m *Map[T]) GetRelation(entity ecs.Entity) ecs.Entity {
	return m.world.Relations().Get(entity, m.id)
}

// GetRelation returns the target entity for the given entity and the Map's relation component.
//
// Returns the zero entity if the entity does not have the given component,
// or if the component is not an [ecs.Relation].
//
// GetRelationUnchecked is an optimized version of [Map.GetRelation].
// Does not check if the entity is alive or that the component ID is applicable.
//
// See also [ecs.World.GetRelationUnchecked].
func (m *Map[T]) GetRelationUnchecked(entity ecs.Entity) ecs.Entity {
	return m.world.Relations().GetUnchecked(entity, m.id)
}

// SetRelation sets the target entity for the given entity and the Map's component.
//
// Panics if the entity does not have a component of that type.
// Panics if the component is not a relation.
//
// See also [ecs.World.SetRelation].
func (m *Map[T]) SetRelation(entity, target ecs.Entity) {
	m.world.Relations().Set(entity, m.id, target)
}

// SetRelationBatch sets the target entity for many entities and the Map's component.
// Returns the number of affected entities.
//
// Panics if the entity does not have a component of that type.
// Panics if the component is not a relation.
//
// See also [ecs.Batch.SetRelation].
func (m *Map[T]) SetRelationBatch(filter ecs.Filter, target ecs.Entity) int {
	return m.world.Batch().SetRelation(filter, m.id, target)
}

// SetRelationBatch sets the target entity for many entities and the Map's component,
// and returns a query over them.
//
// Panics if the entity does not have a component of that type.
// Panics if the component is not a relation.
//
// See also [ecs.Batch.SetRelation].
func (m *Map[T]) SetRelationBatchQ(filter ecs.Filter, target ecs.Entity) Query1[T] {
	query := m.world.Batch().SetRelationQ(filter, m.id, target)
	return Query1[T]{
		Query:       query,
		id0:         m.id,
		hasRelation: true,
		relation:    m.id,
	}
}
