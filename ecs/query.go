package ecs

import (
	"unsafe"
)

// Query is an iterator to iterate entities
type Query struct {
	archetypes []archetypeIter
	index      int
	done       bool
}

// NewQuery creates a new QueryIter
func NewQuery(arches []archetypeIter) Query {
	return Query{
		archetypes: arches,
		index:      0,
		done:       false,
	}
}

// Next proceeds to the next entity
func (q *Query) Next() bool {
	if q.done {
		return false
	}
	for {
		if q.archetypes[q.index].Next() {
			return true
		}
		q.index++
		if q.index >= len(q.archetypes) {
			q.done = true
			return false
		}
	}
}

// Has returns whether the current entity has the given component
func (q *Query) Has(comp ID) bool {
	return q.archetypes[q.index].Has(comp)
}

// Get returns the pointer to the given component at the iterator's position
func (q *Query) Get(comp ID) unsafe.Pointer {
	return q.archetypes[q.index].Get(comp)
}

// Entity returns the entity at the iterator's position
func (q *Query) Entity() Entity {
	return q.archetypes[q.index].Entity()
}

type archetypeIter struct {
	Archetype *Archetype
	Length    int
	Index     int
}

func newArchetypeIter(arch *Archetype) archetypeIter {
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
