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
	mask      bitMask
	archetype archetypeIter
	index     int
	done      bool
	lockBit   uint8
}

// newQuery creates a new Query
func newQuery(world *World, mask bitMask, lockBit uint8) Query {
	return Query{
		world:   world,
		mask:    mask,
		index:   -1,
		lockBit: lockBit,
	}
}

// Next proceeds to the next [Entity] in the Query.
func (q *Query) Next() bool {
	if q.archetype.Next() {
		return true
	}
	i, a, ok := q.world.nextArchetype(q.mask, q.index)
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
	return it.Archetype.Get(it.Index, comp)
}

// Entity returns the entity at the iterator's position
func (it *archetypeIter) Entity() Entity {
	return it.Archetype.GetEntity(it.Index)
}

// Q1 is a generic query for one component
type Q1[A any] struct {
	Query
	id ID
}

// Query1 creates a generic query for one component
func Query1[A any](w *World) Q1[A] {
	id := ComponentID[A](w)
	return Q1[A]{
		Query: w.Query(id),
		id:    id,
	}
}

// Get returns the first queried component for the current query position
func (q *Q1[A]) Get1() *A {
	return (*A)(q.Query.Get(q.id))
}

// Q2 is a generic query for two components
type Q2[A any, B any] struct {
	Query
	ids [2]ID
}

// Query2 creates a generic query for two components
func Query2[A any, B any](w *World) Q2[A, B] {
	ids := [2]ID{ComponentID[A](w), ComponentID[B](w)}
	return Q2[A, B]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Get returns all queried component for the current query position
func (q *Q2[A, B]) GetAll() (*A, *B) {
	return (*A)(q.Query.Get(q.ids[0])), (*B)(q.Query.Get(q.ids[1]))
}

// Get returns the first queried component for the current query position
func (q *Q2[A, B]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get returns the second queried component for the current query position
func (q *Q2[A, B]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Q3 is a generic query for three components
type Q3[A any, B any, C any] struct {
	Query
	ids [3]ID
}

// Query3 creates a generic query for three components
func Query3[A any, B any, C any](w *World) Q3[A, B, C] {
	ids := [3]ID{ComponentID[A](w), ComponentID[B](w), ComponentID[C](w)}
	return Q3[A, B, C]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Get returns all queried component for the current query position
func (q *Q3[A, B, C]) GetAll() (*A, *B, *C) {
	return (*A)(q.Query.Get(q.ids[0])), (*B)(q.Query.Get(q.ids[1])), (*C)(q.Query.Get(q.ids[2]))
}

// Get returns the first queried component for the current query position
func (q *Q3[A, B, C]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get returns the second queried component for the current query position
func (q *Q3[A, B, C]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Get returns the third queried component for the current query position
func (q *Q3[A, B, C]) Get3() *C {
	return (*C)(q.Query.Get(q.ids[2]))
}

// Q4 is a generic query for four components
type Q4[A any, B any, C any, D any] struct {
	Query
	ids [4]ID
}

// Query4 creates a generic query for four components
func Query4[A any, B any, C any, D any](w *World) Q4[A, B, C, D] {
	ids := [4]ID{ComponentID[A](w), ComponentID[B](w), ComponentID[C](w), ComponentID[D](w)}
	return Q4[A, B, C, D]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Get returns all queried component for the current query position
func (q *Q4[A, B, C, D]) GetAll() (*A, *B, *C, *D) {
	return (*A)(q.Query.Get(q.ids[0])),
		(*B)(q.Query.Get(q.ids[1])),
		(*C)(q.Query.Get(q.ids[2])),
		(*D)(q.Query.Get(q.ids[3]))
}

// Get returns the first queried component for the current query position
func (q *Q4[A, B, C, D]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get returns the second queried component for the current query position
func (q *Q4[A, B, C, D]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Get returns the third queried component for the current query position
func (q *Q4[A, B, C, D]) Get3() *C {
	return (*C)(q.Query.Get(q.ids[2]))
}

// Get returns the fourth queried component for the current query position
func (q *Q4[A, B, C, D]) Get4() *D {
	return (*D)(q.Query.Get(q.ids[3]))
}

// Q5 is a generic query for five components
type Q5[A any, B any, C any, D any, E any] struct {
	Query
	ids [5]ID
}

// Query5 creates a generic query for four components
func Query5[A any, B any, C any, D any, E any](w *World) Q5[A, B, C, D, E] {
	ids := [5]ID{ComponentID[A](w), ComponentID[B](w), ComponentID[C](w), ComponentID[D](w), ComponentID[E](w)}
	return Q5[A, B, C, D, E]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Get returns all queried component for the current query position
func (q *Q5[A, B, C, D, E]) GetAll() (*A, *B, *C, *D, *E) {
	return (*A)(q.Query.Get(q.ids[0])),
		(*B)(q.Query.Get(q.ids[1])),
		(*C)(q.Query.Get(q.ids[2])),
		(*D)(q.Query.Get(q.ids[3])),
		(*E)(q.Query.Get(q.ids[4]))
}

// Get returns the first queried component for the current query position
func (q *Q5[A, B, C, D, E]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get returns the second queried component for the current query position
func (q *Q5[A, B, C, D, E]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Get returns the third queried component for the current query position
func (q *Q5[A, B, C, D, E]) Get3() *C {
	return (*C)(q.Query.Get(q.ids[2]))
}

// Get returns the fourth queried component for the current query position
func (q *Q5[A, B, C, D, E]) Get4() *D {
	return (*D)(q.Query.Get(q.ids[3]))
}

// Get returns the fifth queried component for the current query position
func (q *Q5[A, B, C, D, E]) Get5() *E {
	return (*E)(q.Query.Get(q.ids[4]))
}
