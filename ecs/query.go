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
	return bits.Contains(f.Include) && (f.Exclude.IsZero() || !bits.ContainsAny(f.Exclude))
}

// dummyFilter is a filter that returns either true or false, irrespective of the matched mask.
//
// Used for the [Query] returned by entity batch creation methods.
type dummyFilter struct{ Value bool }

// Matches matches a filter against a mask.
func (f dummyFilter) Matches(bits Mask) bool {
	return f.Value
}

// Query is an iterator to iterate entities, filtered by a [Filter].
//
// Create queries through the [World] using [World.Query].
//
// See also the generic alternatives [github.com/mlange-42/arche/generic.Query1],
// [github.com/mlange-42/arche/generic.Query2], etc.
// For advanced filtering, see package [github.com/mlange-42/arche/filter]
type Query struct {
	archetypeIter
	filter     Filter
	world      *World
	archetypes archetypes
	index      int
	lockBit    uint8
	count      int
}

// newQuery creates a new Filter
func newQuery(world *World, filter Filter, lockBit uint8, archetypes archetypes) Query {
	return Query{
		filter:     filter,
		world:      world,
		archetypes: archetypes,
		index:      -1,
		lockBit:    lockBit,
		count:      -1,
	}
}

// newQuery creates a query on a single archetype
func newArchQuery(world *World, lockBit uint8, archetype *archetype, start uint32) Query {
	if start > 0 {
		iter := newArchetypeIter(archetype)
		iter.Index = uintptr(start - 1)
		return Query{
			filter:        dummyFilter{true},
			world:         world,
			archetypes:    singleArchetype{archetype},
			index:         0,
			lockBit:       lockBit,
			count:         int(archetype.Len() - start),
			archetypeIter: iter,
		}
	}
	return Query{
		filter:     dummyFilter{true},
		world:      world,
		archetypes: singleArchetype{archetype},
		index:      -1,
		lockBit:    lockBit,
		count:      int(archetype.Len()),
	}
}

// Next proceeds to the next [Entity] in the Query.
func (q *Query) Next() bool {
	if q.archetypeIter.Next() {
		return true
	}
	// outline to allow inlining of the fast path
	return q.nextArchetype()
}

// Step advances the query iterator by the given number of entities.
//
// Query.Step(1) is equivalent to [Query.Next]().
//
// This method, used together with [Query.Count], can be useful for the selection of random entities.
func (q *Query) Step(step int) bool {
	if step <= 0 {
		panic("step size must be positive")
	}
	var ok bool
	for {
		step, ok = q.archetypeIter.Step(uint32(step))
		if ok {
			return true
		}
		if !q.nextArchetype() {
			return false
		}
		if step == 0 {
			return true
		}
	}
}

// Count counts the entities matching this query.
//
// Involves a small overhead of iterating through archetypes when called the first time.
// However, this is still much faster than manual counting via iteration.
func (q *Query) Count() int {
	if q.count >= 0 {
		return q.count
	}
	q.count = q.countEntities()
	return q.count
}

func (q *Query) countEntities() int {
	len := int(q.archetypes.Len())
	count := uint32(0)
	for i := 0; i < len; i++ {
		a := q.archetypes.Get(i)
		if q.filter.Matches(a.Mask) {
			count += a.Len()
		}
	}
	return int(count)
}

// Mask returns the archetype [BitMask] for the [Entity] at the iterator's current position.
//
// Can be used for fast checks of the entity composition, e.g. using a [Filter].
func (q *Query) Mask() Mask {
	return q.Access.Mask
}

// Close closes the Query and unlocks the world.
//
// Automatically called when iteration finishes.
// Needs to be called only if breaking out of the query iteration.
func (q *Query) Close() {
	q.world.closeQuery(q)
}

// nextArchetype proceeds to the next archetype, and returns whether this was successful/possible.
func (q *Query) nextArchetype() bool {
	if mask, ok := q.filter.(Mask); ok {
		return q.nextArchetypeMask(mask)
	}
	if mask, ok := q.filter.(*MaskFilter); ok {
		return q.nextArchetypeMaskFilter(mask)
	}
	return q.nextArchetypeFilter()
}

func (q *Query) nextArchetypeMask(f Mask) bool {
	len := int(q.archetypes.Len())
	for i := q.index + 1; i < len; i++ {
		a := q.archetypes.Get(i)
		if f.Matches(a.Mask) && a.Len() > 0 {
			q.index = i
			q.archetypeIter = newArchetypeIter(a)
			return true
		}
	}
	q.index = len
	q.world.closeQuery(q)
	return false
}

func (q *Query) nextArchetypeMaskFilter(f *MaskFilter) bool {
	len := int(q.archetypes.Len())
	for i := q.index + 1; i < len; i++ {
		a := q.archetypes.Get(i)
		if f.Matches(a.Mask) && a.Len() > 0 {
			q.index = i
			q.archetypeIter = newArchetypeIter(a)
			return true
		}
	}
	q.index = len
	q.world.closeQuery(q)
	return false
}

func (q *Query) nextArchetypeFilter() bool {
	len := int(q.archetypes.Len())
	for i := q.index + 1; i < len; i++ {
		a := q.archetypes.Get(i)
		if q.filter.Matches(a.Mask) && a.Len() > 0 {
			q.index = i
			q.archetypeIter = newArchetypeIter(a)
			return true
		}
	}
	q.index = len
	q.world.closeQuery(q)
	return false
}

// archetypeIter is an iterator ovr a single archetype.
type archetypeIter struct {
	Access *archetypeAccess
	Length uintptr
	Index  uintptr
}

// newArchetypeIter creates a new archetypeIter.
func newArchetypeIter(arch *archetype) archetypeIter {
	return archetypeIter{
		Access: &arch.archetypeAccess,
		Length: uintptr(arch.Len()),
	}
}

// Next proceeds to the next entity in the archetype, and returns whether this was successful/possible.
func (it *archetypeIter) Next() bool {
	it.Index++
	return it.Index < it.Length
}

// Step proceeds/steps by the given number of entities.
func (it *archetypeIter) Step(count uint32) (int, bool) {
	if it.Length == 0 {
		return int(count - 1), false
	}
	it.Index += uintptr(count)
	if it.Index < it.Length {
		return 0, true
	}
	return int(it.Index) - int(it.Length), false
}

// Has returns whether the current entity has the given component.
func (it *archetypeIter) Has(comp ID) bool {
	return it.Access.HasComponent(comp)
}

// Get returns the pointer to the given component at the iterator's position.
func (it *archetypeIter) Get(comp ID) unsafe.Pointer {
	return it.Access.Get(it.Index, comp)
}

// Entity returns the entity at the iterator's position.
func (it *archetypeIter) Entity() Entity {
	return it.Access.GetEntity(it.Index)
}
