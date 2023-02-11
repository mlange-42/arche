package ecs

// Mask is a mask for a combination of components.
type Mask struct {
	mask bitMask
}

// Q0 is a generic query for no components.
//
// Create one with [Query0]
type Q0 struct {
	Query
}

// Query0 creates a generic query for no components.
//
// See also [World.Query].
func Query0(w *World) Q0 {
	return Q0{
		Query: w.query(0, 0),
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q0) Not(mask Mask) Q0 {
	q.exclude = mask.mask
	return q
}

// Q1 is a generic query for one component.
//
// Create one with [Query1]
type Q1[A any] struct {
	Query
	id ID
}

// Query1 creates a generic query for one component.
//
// See also [World.Query].
func Query1[A any](w *World) Q1[A] {
	id := ComponentID[A](w)
	return Q1[A]{
		Query: w.Query(id),
		id:    id,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q1[A]) Not(mask Mask) Q1[A] {
	q.exclude = mask.mask
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
	Query
	ids [2]ID
}

// Query2 creates a generic query for two components.
//
// See also [World.Query].
func Query2[A any, B any](w *World) Q2[A, B] {
	ids := [2]ID{ComponentID[A](w), ComponentID[B](w)}
	return Q2[A, B]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q2[A, B]) Not(mask Mask) Q2[A, B] {
	q.exclude = mask.mask
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
	Query
	ids [3]ID
}

// Query3 creates a generic query for three components.
//
// See also [World.Query].
func Query3[A any, B any, C any](w *World) Q3[A, B, C] {
	ids := [3]ID{ComponentID[A](w), ComponentID[B](w), ComponentID[C](w)}
	return Q3[A, B, C]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q3[A, B, C]) Not(mask Mask) Q3[A, B, C] {
	q.exclude = mask.mask
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
	Query
	ids [4]ID
}

// Query4 creates a generic query for four components.
//
// See also [World.Query].
func Query4[A any, B any, C any, D any](w *World) Q4[A, B, C, D] {
	ids := [4]ID{ComponentID[A](w), ComponentID[B](w), ComponentID[C](w), ComponentID[D](w)}
	return Q4[A, B, C, D]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q4[A, B, C, D]) Not(mask Mask) Q4[A, B, C, D] {
	q.exclude = mask.mask
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
	Query
	ids [5]ID
}

// Query5 creates a generic query for five components.
//
// See also [World.Query].
func Query5[A any, B any, C any, D any, E any](w *World) Q5[A, B, C, D, E] {
	ids := [5]ID{ComponentID[A](w), ComponentID[B](w), ComponentID[C](w), ComponentID[D](w), ComponentID[E](w)}
	return Q5[A, B, C, D, E]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q5[A, B, C, D, E]) Not(mask Mask) Q5[A, B, C, D, E] {
	q.exclude = mask.mask
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
	Query
	ids [6]ID
}

// Query6 creates a generic query for six components.
//
// See also [World.Query].
func Query6[A any, B any, C any, D any, E any, F any](w *World) Q6[A, B, C, D, E, F] {
	ids := [6]ID{
		ComponentID[A](w), ComponentID[B](w), ComponentID[C](w),
		ComponentID[D](w), ComponentID[E](w), ComponentID[F](w),
	}
	return Q6[A, B, C, D, E, F]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q6[A, B, C, D, E, F]) Not(mask Mask) Q6[A, B, C, D, E, F] {
	q.exclude = mask.mask
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
	Query
	ids [7]ID
}

// Query7 creates a generic query for seven components.
//
// See also [World.Query].
func Query7[A any, B any, C any, D any, E any, F any, G any](w *World) Q7[A, B, C, D, E, F, G] {
	ids := [7]ID{
		ComponentID[A](w), ComponentID[B](w), ComponentID[C](w),
		ComponentID[D](w), ComponentID[E](w), ComponentID[F](w),
		ComponentID[G](w),
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
func (q Q7[A, B, C, D, E, F, G]) Not(mask Mask) Q7[A, B, C, D, E, F, G] {
	q.exclude = mask.mask
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
	Query
	ids [8]ID
}

// Query8 creates a generic query for seven components.
//
// See also [World.Query].
func Query8[A any, B any, C any, D any, E any, F any, G any, H any](w *World) Q8[A, B, C, D, E, F, G, H] {
	ids := [8]ID{
		ComponentID[A](w), ComponentID[B](w), ComponentID[C](w),
		ComponentID[D](w), ComponentID[E](w), ComponentID[F](w),
		ComponentID[G](w), ComponentID[H](w),
	}
	return Q8[A, B, C, D, E, F, G, H]{
		Query: w.Query(ids[:]...),
		ids:   ids,
	}
}

// Not excludes entities with the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q8[A, B, C, D, E, F, G, H]) Not(mask Mask) Q8[A, B, C, D, E, F, G, H] {
	q.exclude = mask.mask
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
func Mask1[A any](w *World) Mask {
	return Mask{newMask(ComponentID[A](w))}
}

// Mask2 creates a component [Mask] for two component types.
func Mask2[A any, B any](w *World) Mask {
	return Mask{newMask(
		ComponentID[A](w), ComponentID[B](w),
	)}
}

// Mask3 creates a component [Mask] for three component types.
func Mask3[A any, B any, C any](w *World) Mask {
	return Mask{newMask(
		ComponentID[A](w), ComponentID[B](w), ComponentID[C](w),
	)}
}

// Mask4 creates a component [Mask] for four component types.
func Mask4[A any, B any, C any, D any](w *World) Mask {
	return Mask{newMask(
		ComponentID[A](w), ComponentID[B](w), ComponentID[C](w),
		ComponentID[D](w),
	)}
}

// Mask5 creates a component [Mask] for five component types.
func Mask5[A any, B any, C any, D any, E any](w *World) Mask {
	return Mask{newMask(
		ComponentID[A](w), ComponentID[B](w), ComponentID[C](w),
		ComponentID[D](w), ComponentID[E](w),
	)}
}

// Mask6 creates a component [Mask] for six component types.
func Mask6[A any, B any, C any, D any, E any, F any](w *World) Mask {
	return Mask{newMask(
		ComponentID[A](w), ComponentID[B](w), ComponentID[C](w),
		ComponentID[D](w), ComponentID[E](w), ComponentID[F](w),
	)}
}

// Mask7 creates a component [Mask] for seven component types.
func Mask7[A any, B any, C any, D any, E any, F any, G any](w *World) Mask {
	return Mask{newMask(
		ComponentID[A](w), ComponentID[B](w), ComponentID[C](w),
		ComponentID[D](w), ComponentID[E](w), ComponentID[F](w),
		ComponentID[G](w),
	)}
}

// Mask8 creates a component [Mask] for eight component types.
func Mask8[A any, B any, C any, D any, E any, F any, G any, H any](w *World) Mask {
	return Mask{newMask(
		ComponentID[A](w), ComponentID[B](w), ComponentID[C](w),
		ComponentID[D](w), ComponentID[E](w), ComponentID[F](w),
		ComponentID[G](w), ComponentID[H](w),
	)}
}
