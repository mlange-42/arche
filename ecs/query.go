package ecs

import "unsafe"

// Query is an iterator to iterate entities
type Query struct {
	Mask      Mask
	World     *World
	archetype int
	index     int
}

// NewQuery creates a new QueryIter
func NewQuery(w *World, mask Mask) Query {
	return Query{
		Mask:      mask,
		World:     w,
		archetype: -1,
		index:     -1,
	}
}

// Next proceeds to the next entity
func (q *Query) Next() bool {
	a, i, hasNext := q.World.Next(q.Mask, q.archetype, q.index)
	q.archetype, q.index = a, i
	return hasNext
}

// Get returns the pointer to the given component at the iterator's position
func (q *Query) Get(comp ID) unsafe.Pointer {
	return q.World.GetAt(q.archetype, q.index, comp)
}

// Entity returns the entity at the iterator's position
func (q *Query) Entity() Entity {
	return q.World.GetEntityAt(q.archetype, q.index)
}
