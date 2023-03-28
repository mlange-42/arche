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
	Include Mask // Components to include.
	Exclude Mask // Components to exclude.
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

// CachedFilter is a filter that is cached by the world.
type CachedFilter struct {
	filter Filter
	id     ID
}

// Matches matches a filter against a mask.
func (f *CachedFilter) Matches(bits Mask) bool {
	return f.filter.Matches(bits)
}

// Query is an iterator to iterate entities, filtered by a [Filter].
//
// Create queries through the [World] using [World.Query].
//
// See also the generic alternatives [github.com/mlange-42/arche/generic.Query1],
// [github.com/mlange-42/arche/generic.Query2], etc.
// For advanced filtering, see package [github.com/mlange-42/arche/filter]
type Query struct {
	filter         Filter
	world          *World
	archetypes     archetypes
	access         *archetypeAccess
	index          int
	lockBit        uint8
	count          int
	entityIndex    uintptr
	entityIndexMax uintptr
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
		return Query{
			filter:         dummyFilter{true},
			world:          world,
			archetypes:     batchArchetype{archetype, start},
			access:         &archetype.archetypeAccess,
			index:          0,
			lockBit:        lockBit,
			count:          int(archetype.Len() - start),
			entityIndex:    uintptr(start - 1),
			entityIndexMax: uintptr(archetype.Len()) - 1,
		}
	}
	return Query{
		filter:     dummyFilter{true},
		world:      world,
		archetypes: batchArchetype{archetype, start},
		index:      -1,
		lockBit:    lockBit,
		count:      int(archetype.Len()),
	}
}

// Next proceeds to the next [Entity] in the Query.
func (q *Query) Next() bool {
	if q.entityIndex < q.entityIndexMax {
		q.entityIndex++
		return true
	}
	// outline to allow inlining of the fast path
	return q.nextArchetype()
}

// Has returns whether the current entity has the given component.
func (q *Query) Has(comp ID) bool {
	return q.access.HasComponent(comp)
}

// Get returns the pointer to the given component at the iterator's position.
func (q *Query) Get(comp ID) unsafe.Pointer {
	return q.access.Get(q.entityIndex, comp)
}

// Entity returns the entity at the iterator's position.
func (q *Query) Entity() Entity {
	return q.access.GetEntity(q.entityIndex)
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
		step, ok = q.stepArchetype(uint32(step))
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

// Mask returns the archetype [BitMask] for the [Entity] at the iterator's current position.
//
// Can be used for fast checks of the entity composition, e.g. using a [Filter].
func (q *Query) Mask() Mask {
	return q.access.Mask
}

// Close closes the Query and unlocks the world.
//
// Automatically called when iteration finishes.
// Needs to be called only if breaking out of the query iteration.
func (q *Query) Close() {
	q.world.closeQuery(q)
}

func (q *Query) stepArchetype(step uint32) (int, bool) {
	q.entityIndex += uintptr(step)
	if q.entityIndex <= q.entityIndexMax {
		return 0, true
	}
	return int(q.entityIndex) - int(q.entityIndexMax) - 1, false
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

// nextArchetype proceeds to the next archetype, and returns whether this was successful/possible.
func (q *Query) nextArchetype() bool {
	switch q.filter.(type) {
	case *CachedFilter:
		return q.nextArchetypeCached()
	default:
		return q.nextArchetypeFilter()
	}
}

// nextArchetypeCached is called if the query is cached.
func (q *Query) nextArchetypeCached() bool {
	len := int(q.archetypes.Len()) - 1
	for q.index < len {
		q.index++
		a := q.archetypes.Get(q.index)
		aLen := a.Len()
		if aLen > 0 {
			q.access = &a.archetypeAccess
			q.entityIndex = 0
			q.entityIndexMax = uintptr(aLen) - 1
			return true
		}
	}
	q.world.closeQuery(q)
	return false
}

// nextArchetypeFilter is called if the query is not cached.
func (q *Query) nextArchetypeFilter() bool {
	len := int(q.archetypes.Len()) - 1
	for q.index < len {
		q.index++
		a := q.archetypes.Get(q.index)
		aLen := a.Len()
		if q.filter.Matches(a.Mask) && aLen > 0 {
			q.access = &a.archetypeAccess
			q.entityIndex = 0
			q.entityIndexMax = uintptr(aLen) - 1
			return true
		}
	}
	q.world.closeQuery(q)
	return false
}
