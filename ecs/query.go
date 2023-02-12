package ecs

import (
	"unsafe"

	"github.com/mlange-42/arche/internal/base"
)

// Filter is the interface for logic filters
type Filter interface {
	Matches(mask base.BitMask) bool
}

// Mask is a mask for a combination of components.
type Mask struct {
	BitMask bitMask
}

// NewMask creates a new Mask from a list of IDs.
//
// If any ID is bigger or equal [MaskTotalBits], it'll not be added to the mask.
func NewMask(ids ...ID) Mask {
	var mask bitMask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return Mask{mask}
}

// Matches matches a filter against a mask
func (f Mask) Matches(mask bitMask) bool {
	return mask.Contains(f.BitMask)
}

// All matches all the given components.
func All(comps ...base.ID) Mask {
	return NewMask(comps...)
}

// Not excludes the given components.
func (f Mask) Not(comps ...base.ID) MaskPair {
	return MaskPair{
		Mask:    f,
		Exclude: NewMask(comps...),
	}
}

// MaskPair is a filter for including an excluding components
type MaskPair struct {
	Mask    Mask
	Exclude Mask
}

// Matches matches a filter against a mask
func (f MaskPair) Matches(mask base.BitMask) bool {
	return mask.Contains(f.Mask.BitMask) && !mask.Contains(f.Exclude.BitMask)
}

// EntityIter is the interface for iterable queries
type EntityIter interface {
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
// See also the generic alternatives [generic.Query1],  [generic.Query2], etc.
//
// For advanced filtering, see package [filter]
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
	i, a, ok := q.world.nextArchetypeFilter(q.filter, q.index)
	q.index = i
	if ok {
		q.archetype = a
		return true
	}
	q.done = true
	q.world.closeQuery(&q.queryIter)
	return false
}

type queryIter struct {
	world     *World
	archetype archetypeIter
	index     int
	done      bool
	lockBit   uint8
}

// Has returns whether the current [Entity] has the given component
func (q *queryIter) Has(comp ID) bool {
	return q.archetype.Has(comp)
}

// Get returns the pointer to the given component at the iterator's current [Entity]
func (q *queryIter) Get(comp ID) unsafe.Pointer {
	return q.archetype.Get(comp)
}

// Entity returns the [Entity] at the iterator's position
func (q *queryIter) Entity() Entity {
	return q.archetype.Entity()
}

// Close closes the Query and unlocks the world.
//
// Automatically called when iteration finishes.
// Needs to be called only if breaking out of the query iteration.
func (q *queryIter) Close() {
	q.done = true
	q.world.closeQuery(q)
}

type archetypeIter struct {
	Archetype *archetype
	Length    int
	Index     int
}

func newArchetypeIter(arch *archetype) archetypeIter {
	return archetypeIter{
		Archetype: arch,
		Length:    int(arch.Len()),
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
