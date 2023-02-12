package generic

import (
	"github.com/mlange-42/arche/ecs"
)

// Q0 is a generic query for no components.
//
// Create one with [Query0]
type Q0 struct {
	ecs.Query
}

// Query0 creates a generic query for no components.
//
// See also [World.Query].
func Query0(w *ecs.World) Q0 {
	return Q0{
		Query: w.Query(),
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q0) Not(mask ecs.Mask) Q0 {
	q.Exclude = q.Exclude | mask.BitMask
	return q
}

// Q1 is a generic query for one component.
//
// Create one with [Query1]
type Q1[A any] struct {
	ecs.Query
	id ecs.ID
}

// Query1 creates a generic query for one component.
//
// See also [World.Query].
func Query1[A any](w *ecs.World) Q1[A] {
	id := ecs.ComponentID[A](w)
	return Q1[A]{
		Query: w.Query(id),
		id:    id,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q1[A]) Not(mask ecs.Mask) Q1[A] {
	q.Exclude = q.Exclude | mask.BitMask
	return q
}

// Get1 returns the first queried component for the current query position
func (q *Q1[A]) Get1() *A {
	return (*A)(q.Query.Get(q.id))
}

//////////////////////////////////////////////////////////////////////////

// Q2 is a generic query for two components.
//
// Create one with [Query2]
type Q2[A any, B any] struct {
	ecs.Query
	ids [2]ecs.ID
}

// Query2 creates a generic query for two components.
//
// See also [World.Query].
func Query2[A any, B any](w *ecs.World) Q2[A, B] {
	ids := [2]ecs.ID{ecs.ComponentID[A](w), ecs.ComponentID[B](w)}
	return Q2[A, B]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q2[A, B]) Not(mask ecs.Mask) Q2[A, B] {
	q.Exclude = q.Exclude | mask.BitMask
	return q
}

// GetAll returns all queried components for the current query position
func (q *Q2[A, B]) GetAll() (*A, *B) {
	return (*A)(q.Query.Get(q.ids[0])), (*B)(q.Query.Get(q.ids[1]))
}

// Get1 returns the first queried component for the current query position
func (q *Q2[A, B]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get2 returns the second queried component for the current query position
func (q *Q2[A, B]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

//////////////////////////////////////////////////////////////////////////

// Q3 is a generic query for three components.
//
// Create one with [Query3]
type Q3[A any, B any, C any] struct {
	ecs.Query
	ids [3]ecs.ID
}

// Query3 creates a generic query for three components.
//
// See also [World.Query].
func Query3[A any, B any, C any](w *ecs.World) Q3[A, B, C] {
	ids := [3]ecs.ID{ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w)}
	return Q3[A, B, C]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q3[A, B, C]) Not(mask ecs.Mask) Q3[A, B, C] {
	q.Exclude = q.Exclude | mask.BitMask
	return q
}

// GetAll returns all queried components for the current query position
func (q *Q3[A, B, C]) GetAll() (*A, *B, *C) {
	return (*A)(q.Query.Get(q.ids[0])), (*B)(q.Query.Get(q.ids[1])), (*C)(q.Query.Get(q.ids[2]))
}

// Get1 returns the first queried component for the current query position
func (q *Q3[A, B, C]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get2 returns the second queried component for the current query position
func (q *Q3[A, B, C]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Get3 returns the third queried component for the current query position
func (q *Q3[A, B, C]) Get3() *C {
	return (*C)(q.Query.Get(q.ids[2]))
}

//////////////////////////////////////////////////////////////////////////

// Q4 is a generic query for four components.
//
// Create one with [Query4]
type Q4[A any, B any, C any, D any] struct {
	ecs.Query
	ids [4]ecs.ID
}

// Query4 creates a generic query for four components.
//
// See also [World.Query].
func Query4[A any, B any, C any, D any](w *ecs.World) Q4[A, B, C, D] {
	ids := [4]ecs.ID{ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w), ecs.ComponentID[D](w)}
	return Q4[A, B, C, D]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q4[A, B, C, D]) Not(mask ecs.Mask) Q4[A, B, C, D] {
	q.Exclude = q.Exclude | mask.BitMask
	return q
}

// GetAll returns all queried components for the current query position
func (q *Q4[A, B, C, D]) GetAll() (*A, *B, *C, *D) {
	return (*A)(q.Query.Get(q.ids[0])),
		(*B)(q.Query.Get(q.ids[1])),
		(*C)(q.Query.Get(q.ids[2])),
		(*D)(q.Query.Get(q.ids[3]))
}

// Get1 returns the first queried component for the current query position
func (q *Q4[A, B, C, D]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get2 returns the second queried component for the current query position
func (q *Q4[A, B, C, D]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Get3 returns the third queried component for the current query position
func (q *Q4[A, B, C, D]) Get3() *C {
	return (*C)(q.Query.Get(q.ids[2]))
}

// Get4 returns the fourth queried component for the current query position
func (q *Q4[A, B, C, D]) Get4() *D {
	return (*D)(q.Query.Get(q.ids[3]))
}

//////////////////////////////////////////////////////////////////////////

// Q5 is a generic query for five components.
//
// Create one with [Query5]
type Q5[A any, B any, C any, D any, E any] struct {
	ecs.Query
	ids [5]ecs.ID
}

// Query5 creates a generic query for five components.
//
// See also [World.Query].
func Query5[A any, B any, C any, D any, E any](w *ecs.World) Q5[A, B, C, D, E] {
	ids := [5]ecs.ID{ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w), ecs.ComponentID[D](w), ecs.ComponentID[E](w)}
	return Q5[A, B, C, D, E]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q5[A, B, C, D, E]) Not(mask ecs.Mask) Q5[A, B, C, D, E] {
	q.Exclude = q.Exclude | mask.BitMask
	return q
}

// GetAll returns all queried components for the current query position
func (q *Q5[A, B, C, D, E]) GetAll() (*A, *B, *C, *D, *E) {
	return (*A)(q.Query.Get(q.ids[0])),
		(*B)(q.Query.Get(q.ids[1])),
		(*C)(q.Query.Get(q.ids[2])),
		(*D)(q.Query.Get(q.ids[3])),
		(*E)(q.Query.Get(q.ids[4]))
}

// Get1 returns the first queried component for the current query position
func (q *Q5[A, B, C, D, E]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get2 returns the second queried component for the current query position
func (q *Q5[A, B, C, D, E]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Get3 returns the third queried component for the current query position
func (q *Q5[A, B, C, D, E]) Get3() *C {
	return (*C)(q.Query.Get(q.ids[2]))
}

// Get4 returns the fourth queried component for the current query position
func (q *Q5[A, B, C, D, E]) Get4() *D {
	return (*D)(q.Query.Get(q.ids[3]))
}

// Get5 returns the fifth queried component for the current query position
func (q *Q5[A, B, C, D, E]) Get5() *E {
	return (*E)(q.Query.Get(q.ids[4]))
}

//////////////////////////////////////////////////////////////////////////

// Q6 is a generic query for six components.
//
// Create one with [Query6]
type Q6[A any, B any, C any, D any, E any, F any] struct {
	ecs.Query
	ids [6]ecs.ID
}

// Query6 creates a generic query for six components.
//
// See also [World.Query].
func Query6[A any, B any, C any, D any, E any, F any](w *ecs.World) Q6[A, B, C, D, E, F] {
	ids := [6]ecs.ID{
		ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w),
		ecs.ComponentID[D](w), ecs.ComponentID[E](w), ecs.ComponentID[F](w),
	}
	return Q6[A, B, C, D, E, F]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q6[A, B, C, D, E, F]) Not(mask ecs.Mask) Q6[A, B, C, D, E, F] {
	q.Exclude = q.Exclude | mask.BitMask
	return q
}

// GetAll returns all queried components for the current query position
func (q *Q6[A, B, C, D, E, F]) GetAll() (*A, *B, *C, *D, *E, *F) {
	return (*A)(q.Query.Get(q.ids[0])),
		(*B)(q.Query.Get(q.ids[1])),
		(*C)(q.Query.Get(q.ids[2])),
		(*D)(q.Query.Get(q.ids[3])),
		(*E)(q.Query.Get(q.ids[4])),
		(*F)(q.Query.Get(q.ids[5]))
}

// Get1 returns the first queried component for the current query position
func (q *Q6[A, B, C, D, E, F]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get2 returns the second queried component for the current query position
func (q *Q6[A, B, C, D, E, F]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Get3 returns the third queried component for the current query position
func (q *Q6[A, B, C, D, E, F]) Get3() *C {
	return (*C)(q.Query.Get(q.ids[2]))
}

// Get4 returns the fourth queried component for the current query position
func (q *Q6[A, B, C, D, E, F]) Get4() *D {
	return (*D)(q.Query.Get(q.ids[3]))
}

// Get5 returns the fifth queried component for the current query position
func (q *Q6[A, B, C, D, E, F]) Get5() *E {
	return (*E)(q.Query.Get(q.ids[4]))
}

// Get6 returns the sixth queried component for the current query position
func (q *Q6[A, B, C, D, E, F]) Get6() *F {
	return (*F)(q.Query.Get(q.ids[5]))
}

//////////////////////////////////////////////////////////////////////////

// Q7 is a generic query for seven components.
//
// Create one with [Query7]
type Q7[A any, B any, C any, D any, E any, F any, G any] struct {
	ecs.Query
	ids [7]ecs.ID
}

// Query7 creates a generic query for seven components.
//
// See also [World.Query].
func Query7[A any, B any, C any, D any, E any, F any, G any](w *ecs.World) Q7[A, B, C, D, E, F, G] {
	ids := [7]ecs.ID{
		ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w),
		ecs.ComponentID[D](w), ecs.ComponentID[E](w), ecs.ComponentID[F](w),
		ecs.ComponentID[G](w),
	}
	return Q7[A, B, C, D, E, F, G]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// GetAll returns all queried components for the current query position
func (q *Q7[A, B, C, D, E, F, G]) GetAll() (*A, *B, *C, *D, *E, *F, *G) {
	return (*A)(q.Query.Get(q.ids[0])),
		(*B)(q.Query.Get(q.ids[1])),
		(*C)(q.Query.Get(q.ids[2])),
		(*D)(q.Query.Get(q.ids[3])),
		(*E)(q.Query.Get(q.ids[4])),
		(*F)(q.Query.Get(q.ids[5])),
		(*G)(q.Query.Get(q.ids[6]))
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q7[A, B, C, D, E, F, G]) Not(mask ecs.Mask) Q7[A, B, C, D, E, F, G] {
	q.Exclude = q.Exclude | mask.BitMask
	return q
}

// Get1 returns the first queried component for the current query position
func (q *Q7[A, B, C, D, E, F, G]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get2 returns the second queried component for the current query position
func (q *Q7[A, B, C, D, E, F, G]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Get3 returns the third queried component for the current query position
func (q *Q7[A, B, C, D, E, F, G]) Get3() *C {
	return (*C)(q.Query.Get(q.ids[2]))
}

// Get4 returns the fourth queried component for the current query position
func (q *Q7[A, B, C, D, E, F, G]) Get4() *D {
	return (*D)(q.Query.Get(q.ids[3]))
}

// Get5 returns the fifth queried component for the current query position
func (q *Q7[A, B, C, D, E, F, G]) Get5() *E {
	return (*E)(q.Query.Get(q.ids[4]))
}

// Get6 returns the sixth queried component for the current query position
func (q *Q7[A, B, C, D, E, F, G]) Get6() *F {
	return (*F)(q.Query.Get(q.ids[5]))
}

// Get7 returns the seventh queried component for the current query position
func (q *Q7[A, B, C, D, E, F, G]) Get7() *G {
	return (*G)(q.Query.Get(q.ids[6]))
}

//////////////////////////////////////////////////////////////////////////

// Q8 is a generic query for seven components.
//
// Create one with [Query8]
type Q8[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	ecs.Query
	ids [8]ecs.ID
}

// Query8 creates a generic query for seven components.
//
// See also [World.Query].
func Query8[A any, B any, C any, D any, E any, F any, G any, H any](w *ecs.World) Q8[A, B, C, D, E, F, G, H] {
	ids := [8]ecs.ID{
		ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w),
		ecs.ComponentID[D](w), ecs.ComponentID[E](w), ecs.ComponentID[F](w),
		ecs.ComponentID[G](w), ecs.ComponentID[H](w),
	}
	return Q8[A, B, C, D, E, F, G, H]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q8[A, B, C, D, E, F, G, H]) Not(mask ecs.Mask) Q8[A, B, C, D, E, F, G, H] {
	q.Exclude = q.Exclude | mask.BitMask
	return q
}

// GetAll returns all queried components for the current query position
func (q *Q8[A, B, C, D, E, F, G, H]) GetAll() (*A, *B, *C, *D, *E, *F, *G, *H) {
	return (*A)(q.Query.Get(q.ids[0])),
		(*B)(q.Query.Get(q.ids[1])),
		(*C)(q.Query.Get(q.ids[2])),
		(*D)(q.Query.Get(q.ids[3])),
		(*E)(q.Query.Get(q.ids[4])),
		(*F)(q.Query.Get(q.ids[5])),
		(*G)(q.Query.Get(q.ids[6])),
		(*H)(q.Query.Get(q.ids[7]))
}

// Get1 returns the first queried component for the current query position
func (q *Q8[A, B, C, D, E, F, G, H]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get2 returns the second queried component for the current query position
func (q *Q8[A, B, C, D, E, F, G, H]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Get3 returns the third queried component for the current query position
func (q *Q8[A, B, C, D, E, F, G, H]) Get3() *C {
	return (*C)(q.Query.Get(q.ids[2]))
}

// Get4 returns the fourth queried component for the current query position
func (q *Q8[A, B, C, D, E, F, G, H]) Get4() *D {
	return (*D)(q.Query.Get(q.ids[3]))
}

// Get5 returns the fifth queried component for the current query position
func (q *Q8[A, B, C, D, E, F, G, H]) Get5() *E {
	return (*E)(q.Query.Get(q.ids[4]))
}

// Get6 returns the sixth queried component for the current query position
func (q *Q8[A, B, C, D, E, F, G, H]) Get6() *F {
	return (*F)(q.Query.Get(q.ids[5]))
}

// Get7 returns the seventh queried component for the current query position
func (q *Q8[A, B, C, D, E, F, G, H]) Get7() *G {
	return (*G)(q.Query.Get(q.ids[6]))
}

// Get8 returns the eighth queried component for the current query position
func (q *Q8[A, B, C, D, E, F, G, H]) Get8() *H {
	return (*H)(q.Query.Get(q.ids[7]))
}

// Mask1 creates a component [Mask] for one component type.
func Mask1[A any](w *ecs.World) ecs.Mask {
	return ecs.NewMask(ecs.ComponentID[A](w))
}

// Mask2 creates a component [Mask] for two component types.
func Mask2[A any, B any](w *ecs.World) ecs.Mask {
	return ecs.NewMask(
		ecs.ComponentID[A](w), ecs.ComponentID[B](w),
	)
}

// Mask3 creates a component [Mask] for three component types.
func Mask3[A any, B any, C any](w *ecs.World) ecs.Mask {
	return ecs.NewMask(
		ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w),
	)
}

// Mask4 creates a component [Mask] for four component types.
func Mask4[A any, B any, C any, D any](w *ecs.World) ecs.Mask {
	return ecs.NewMask(
		ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w),
		ecs.ComponentID[D](w),
	)
}

// Mask5 creates a component [Mask] for five component types.
func Mask5[A any, B any, C any, D any, E any](w *ecs.World) ecs.Mask {
	return ecs.NewMask(
		ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w),
		ecs.ComponentID[D](w), ecs.ComponentID[E](w),
	)
}

// Mask6 creates a component [Mask] for six component types.
func Mask6[A any, B any, C any, D any, E any, F any](w *ecs.World) ecs.Mask {
	return ecs.NewMask(
		ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w),
		ecs.ComponentID[D](w), ecs.ComponentID[E](w), ecs.ComponentID[F](w),
	)
}

// Mask7 creates a component [Mask] for seven component types.
func Mask7[A any, B any, C any, D any, E any, F any, G any](w *ecs.World) ecs.Mask {
	return ecs.NewMask(
		ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w),
		ecs.ComponentID[D](w), ecs.ComponentID[E](w), ecs.ComponentID[F](w),
		ecs.ComponentID[G](w),
	)
}

// Mask8 creates a component [Mask] for eight component types.
func Mask8[A any, B any, C any, D any, E any, F any, G any, H any](w *ecs.World) ecs.Mask {
	return ecs.NewMask(
		ecs.ComponentID[A](w), ecs.ComponentID[B](w), ecs.ComponentID[C](w),
		ecs.ComponentID[D](w), ecs.ComponentID[E](w), ecs.ComponentID[F](w),
		ecs.ComponentID[G](w), ecs.ComponentID[H](w),
	)
}
