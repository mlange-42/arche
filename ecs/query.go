package ecs

import (
	"unsafe"
)

// Query is an iterator to iterate entities.
//
// Create queries through the [World] using [World.Query].
// See also the generic variants [Query1], [Query2], [Query3], ...
//
// # Example:
//
//	query := world.Query(posID, rotID)
//	for query.Next() {
//	    pos := (*position)(query.Get(posID))
//	    pos.X += 1.0
//	}
type Query struct {
	world     *World
	mask      bitMask
	exclude   bitMask
	archetype archetypeIter
	index     int
	done      bool
	lockBit   uint8
}

// newQuery creates a new Query
func newQuery(world *World, mask, exclude bitMask, lockBit uint8) Query {
	return Query{
		world:   world,
		mask:    mask,
		exclude: exclude,
		index:   -1,
		lockBit: lockBit,
	}
}

// Not excludes components from the query.
// Entities with these components will be skipped.
//
// # Example:
//
//	query := world.Query(idA, isB).Not(idC, isD)
func (q Query) Not(comps ...ID) Query {
	q.exclude = newMask(comps...)
	return q
}

// Next proceeds to the next [Entity] in the Query.
func (q *Query) Next() bool {
	if q.archetype.Next() {
		return true
	}
	i, a, ok := q.world.nextArchetype(q.mask, q.exclude, q.index)
	q.index = i
	if ok {
		q.archetype = a
		return true
	}
	q.done = true
	q.world.closeQuery(q)
	return false
}

// Has returns whether the current [Entity] has the given component
func (q *Query) Has(comp ID) bool {
	return q.archetype.Has(comp)
}

// Get returns the pointer to the given component at the iterator's current [Entity]
func (q *Query) Get(comp ID) unsafe.Pointer {
	return q.archetype.Get(comp)
}

// Entity returns the [Entity] at the iterator's position
func (q *Query) Entity() Entity {
	return q.archetype.Entity()
}

// Close closes the Query and unlocks the world.
//
// Automatically called when iteration finishes.
// Needs to be called only if breaking out of the query iteration.
func (q *Query) Close() {
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
