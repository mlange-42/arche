package ecs

import (
	"unsafe"
)

// Query is an iterator to iterate entities.
//
// Create queries through the [World] using [World.Query].
//
// ## Example:
//
//	query := world.Query(posID, rotID)
//	for query.Next() {
//	    pos := (*position)(query.Get(posID))
//	    pos.X += 1.0
//	}
type Query struct {
	world      *World
	archetypes []archetypeIter
	index      int
	done       bool
	count      int
	lockBit    uint8
}

// newQuery creates a new Query
func newQuery(world *World, arches []archetypeIter, count int, lockBit uint8) Query {
	return Query{
		world:      world,
		archetypes: arches,
		index:      0,
		done:       false,
		count:      count,
		lockBit:    lockBit,
	}
}

// Next proceeds to the next [Entity] in the Query.
func (q *Query) Next() bool {
	if q.done {
		panic("Query is used up. Create a new Query!")
	}
	for {
		if q.archetypes[q.index].Next() {
			return true
		}
		q.index++
		if q.index >= len(q.archetypes) {
			q.done = true
			q.world.closeQuery(q)
			return false
		}
	}
}

// Has returns whether the current [Entity] has the given component
func (q *Query) Has(comp ID) bool {
	if q.done {
		panic("Query is used up. Create a new Query!")
	}
	return q.archetypes[q.index].Has(comp)
}

// Get returns the pointer to the given component at the iterator's current [Entity]
func (q *Query) Get(comp ID) unsafe.Pointer {
	if q.done {
		panic("Query is used up. Create a new Query!")
	}
	return q.archetypes[q.index].Get(comp)
}

// Entity returns the [Entity] at the iterator's position
func (q *Query) Entity() Entity {
	if q.done {
		panic("Query is used up. Create a new Query!")
	}
	return q.archetypes[q.index].Entity()
}

// Close closes the Query and unlocks the world.
//
// Automatically called when iteration finishes.
// Needs to be called only if breaking out of the query iteration.
func (q *Query) Close() {
	q.done = true
	q.world.closeQuery(q)
}

// Count returns the number of matching entities
func (q *Query) Count() int {
	return q.count
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
