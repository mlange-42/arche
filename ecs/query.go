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
	count  int
}

// newQuery creates a new Filter
func newQuery(world *World, filter Filter, arches []*archetype, count int, lockBit uint8) Query {
	var arch *archetype
	len := len(arches)
	archLen := 0
	if len > 0 {
		arch = arches[0]
		archLen = int(arch.Len())
	}
	return Query{
		queryIter: queryIter{
			world:      world,
			archetypes: arches,
			archetype:  arch,
			index:      -1,
			length:     archLen,
			lockBit:    lockBit,
		},
		count:  count,
		filter: filter,
	}
}

// Count returns the number of entities in this query.
func (q *Query) Count() int {
	return q.count
}

// Next proceeds to the next [Entity] in the Query.
func (q *Query) Next() bool {
	q.index++
	if q.index < q.length {
		return true
	}
	// outline to allow inlining of the fast path
	return q.nextArchetype()
}

func (q *Query) nextArchetype() bool {
	q.archIndex++
	if q.archIndex < len(q.archetypes) {
		q.archetype = q.archetypes[q.archIndex]
		q.index = 0
		q.length = int(q.archetype.Len())
		return true
	}
	q.world.closeQuery(&q.queryIter)
	return false
}

type queryIter struct {
	world      *World
	archetypes []*archetype
	archetype  *archetype
	archIndex  int
	length     int
	index      int
	lockBit    uint8
}

// Has returns whether the current [Entity] has the given component.
func (q *queryIter) Has(comp ID) bool {
	return q.archetype.HasComponent(comp)
}

// Get returns the pointer to the given component at the iterator's current [Entity].
func (q *queryIter) Get(comp ID) unsafe.Pointer {
	return q.archetype.Get(uint32(q.index), comp)
}

// Entity returns the [Entity] at the iterator's position
func (q *queryIter) Entity() Entity {
	return q.archetype.GetEntity(uint32(q.index))
}

// Mask returns the archetype [BitMask] for the [Entity] at the iterator's current position.
//
// Can be used for fast checks of the entity composition, e.g. using a [Filter].
func (q *queryIter) Mask() Mask {
	return q.archetype.Mask
}

// Close closes the Query and unlocks the world.
//
// Automatically called when iteration finishes.
// Needs to be called only if breaking out of the query iteration.
func (q *queryIter) Close() {
	q.world.closeQuery(q)
}
