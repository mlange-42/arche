package ecs

import (
	"unsafe"
)

// Filter is the interface for logic filters.
// Filters are required to query entities using [World.Query].
//
// See [Mask] and [MaskFilter] for basic filters.
// For type-safe generics queries, see package [github.com/mlange-42/arche/generic].
// For advanced filtering, see package [github.com/mlange-42/arche/filter].
type Filter interface {
	// Matches the filter against a bitmask, i.e. a component composition.
	Matches(bits Mask) bool
}

// MaskFilter is a [Filter] for including and excluding certain components.
// See [All] and [Mask.Without].
type MaskFilter struct {
	Include Mask
	Exclude Mask
}

// Matches matches a filter against a mask.
func (f *MaskFilter) Matches(bits Mask) bool {
	return bits.Contains(f.Include) &&
		(f.Exclude.IsZero() || !bits.ContainsAny(f.Exclude))
}

// Query is an iterator to iterate entities, filtered by a [Filter].
//
// Create queries through the [World] using [World.Query].
//
// See also the generic alternatives [github.com/mlange-42/arche/generic.Query1],
// [github.com/mlange-42/arche/generic.Query2], etc.
// For advanced filtering, see package [github.com/mlange-42/arche/filter]
type Query struct {
	queryIter
	filter Filter
}

// newQuery creates a new Filter
func newQuery(world *World, filter Filter, lockBit uint8) Query {
	return Query{
		queryIter: queryIter{
			world:   world,
			index:   -1,
			lockBit: lockBit,
			count:   -1,
		},
		filter: filter,
	}
}

// Next proceeds to the next [Entity] in the Query.
func (q *Query) Next() bool {
	if q.archetype.Next() {
		return true
	}
	// outline to allow inlining of the fast path
	return q.nextArchetype()
}

// Count counts the entities matching this query.
//
// Involves a small overhead of iterating through archetypes when called the first time.
// However, it is considerable faster then manual counting via iteration.
func (q *Query) Count() int {
	if q.count >= 0 {
		return q.count
	}
	q.count = q.world.count(q.filter)
	return q.count
}

func (q *Query) nextArchetype() bool {
	index, archetype, ok := q.world.nextArchetype(q.filter, q.index)
	q.index = index
	if ok {
		q.archetype = archetype
		return true
	}
	q.world.closeQuery(&q.queryIter)
	return false
}

type queryIter struct {
	world     *World
	archetype archetypeIter
	index     int
	lockBit   uint8
	count     int
}

// Has returns whether the current [Entity] has the given component.
func (q *queryIter) Has(comp ID) bool {
	return q.archetype.Has(comp)
}

// Get returns the pointer to the given component at the iterator's current [Entity].
func (q *queryIter) Get(comp ID) unsafe.Pointer {
	return q.archetype.Get(comp)
}

// Entity returns the [Entity] at the iterator's position
func (q *queryIter) Entity() Entity {
	return q.archetype.Entity()
}

// Mask returns the archetype [BitMask] for the [Entity] at the iterator's current position.
//
// Can be used for fast checks of the entity composition, e.g. using a [Filter].
func (q *queryIter) Mask() Mask {
	return q.archetype.Archetype.Mask
}

// Close closes the Query and unlocks the world.
//
// Automatically called when iteration finishes.
// Needs to be called only if breaking out of the query iteration.
func (q *queryIter) Close() {
	q.world.closeQuery(q)
}

type archetypeIter struct {
	Archetype *archetype
	Length    uint32
	Index     uint32
}

func newArchetypeIter(arch *archetype) archetypeIter {
	return archetypeIter{
		Archetype: arch,
		Length:    arch.Len(),
	}
}

func (it *archetypeIter) Next() bool {
	it.Index++
	return it.Index < it.Length
}

// Has returns whether the current entity has the given component
func (it *archetypeIter) Has(comp ID) bool {
	return it.Archetype.HasComponent(comp)
}

// Get returns the pointer to the given component at the iterator's position
func (it *archetypeIter) Get(comp ID) unsafe.Pointer {
	return it.Archetype.Get(it.Index, comp)
}

// Entity returns the entity at the iterator's position
func (it *archetypeIter) Entity() Entity {
	return it.Archetype.GetEntity(it.Index)
}
