package ecs

import (
	"unsafe"
)

// Filter is the interface for logic filters
type Filter interface {
	// Matches the filter against a bitmask, i.e. a component composition.
	Matches(bits BitMask) bool
}

// Mask for a combination of components.
type Mask struct {
	BitMask BitMask
}

// All creates a new Mask from a list of IDs.
//
// If any ID is bigger or equal [MaskTotalBits], it'll not be added to the mask.
func All(ids ...ID) Mask {
	var mask BitMask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return Mask{mask}
}

// Matches matches a filter against a bitmask
func (f Mask) Matches(bits BitMask) bool {
	return bits.Contains(f.BitMask)
}

// Without excludes the given components.
func (f Mask) Without(comps ...ID) MaskPair {
	return MaskPair{
		Mask:    f,
		Exclude: All(comps...),
	}
}

// MaskPair is a filter for including an excluding components
type MaskPair struct {
	Mask    Mask
	Exclude Mask
}

// Matches matches a filter against a mask
func (f MaskPair) Matches(bits BitMask) bool {
	ex := f.Exclude.BitMask
	return bits.Contains(f.Mask.BitMask) && (ex.IsZero() || !bits.ContainsAny(ex))
}

// QueryIter is the interface for iterable queries
type QueryIter interface {
	// Next proceeds to the next [Entity] in the Query.
	Next() bool
	// Has returns whether the current [Entity] has the given component
	Has(comp ID) bool
	// Get returns the pointer to the given component at the iterator's current [Entity]
	Get(comp ID) unsafe.Pointer
	// Entity returns the [Entity] at the iterator's position
	Entity() Entity
	// Close closes the Query and unlocks the world.
	//
	// Automatically called when iteration finishes.
	// Needs to be called only if breaking out of the query iteration.
	Close()
}

// Query is an advanced iterator to iterate entities.
//
// Create queries through the [World] using [World.Query].
//
// See also the generic alternatives [github.com/mlange-42/arche/generic.Query1],
// [github.com/mlange-42/arche/generic.Query2], etc.
//
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
		},
		filter: filter,
	}
}

// Next proceeds to the next [Entity] in the Query.
func (q *Query) Next() bool {
	if q.archetype.Next() {
		return true
	}
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
func (q *queryIter) Mask() BitMask {
	return q.archetype.Archetype.Mask
}

// IDs returns the archetype's component IDs for the [Entity] at the iterator's current position.
//
// Makes a copy of the slice for immutability, so there is a certain overhead involved.
func (q *queryIter) IDs(entity Entity) []ID {
	var ids []ID
	return append(ids, q.archetype.Archetype.Ids...)
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
	return it.Archetype.GetUnsafe(it.Index, comp)
}

// Entity returns the entity at the iterator's position
func (it *archetypeIter) Entity() Entity {
	return it.Archetype.GetEntity(it.Index)
}
