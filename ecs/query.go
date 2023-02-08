package ecs

import (
	"unsafe"
)

// Query is an iterator to iterate entities.
//
// Create queries through the [World] using [World.Query].
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
	mask      Mask
	archetype archetypeIter
	index     int
	done      bool
	lockBit   uint8
}

// newQuery creates a new Query
func newQuery(world *World, mask Mask, lockBit uint8) Query {
	return Query{
		world:   world,
		mask:    mask,
		index:   -1,
		lockBit: lockBit,
	}
}

// Next proceeds to the next [Entity] in the Query.
func (q *Query) Next() bool {
	if q.done {
		panic("Query is used up. Create a new Query!")
	}
	for {
		if q.archetype.Next() {
			return true
		}
		if i, a, ok := q.world.nextArchetype(q.mask, q.index); ok {
			q.index = i
			q.archetype = a
			q.archetype.Next()
			return true
		}
		q.done = true
		q.world.closeQuery(q)
		return false
	}
}

// Has returns whether the current [Entity] has the given component
func (q *Query) Has(comp ID) bool {
	if q.done {
		panic("Query is used up. Create a new Query!")
	}
	return q.archetype.Has(comp)
}

// Get returns the pointer to the given component at the iterator's current [Entity]
func (q *Query) Get(comp ID) unsafe.Pointer {
	if q.done {
		panic("Query is used up. Create a new Query!")
	}
	return q.archetype.Get(comp)
}

// Entity returns the [Entity] at the iterator's position
func (q *Query) Entity() Entity {
	if q.done {
		panic("Query is used up. Create a new Query!")
	}
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
		Index:     -1,
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
