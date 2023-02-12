package generic

import (
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func toIds(w *ecs.World, types []reflect.Type) []ecs.ID {
	ids := make([]ecs.ID, len(types))
	for i, t := range types {
		ids[i] = ecs.TypeID(w, t)
	}
	return ids
}

func toMask(w *ecs.World, types []reflect.Type) ecs.Mask {
	mask := ecs.BitMask(0)
	for _, t := range types {
		mask.Set(ecs.TypeID(w, t), true)
	}
	return ecs.Mask{BitMask: mask}
}

func toMaskOptional(w *ecs.World, include []ecs.ID, optional []reflect.Type) ecs.Mask {
	mask := ecs.NewBitMask(include...)
	for _, t := range optional {
		mask.Set(ecs.TypeID(w, t), false)
	}
	return ecs.Mask{BitMask: mask}
}

// Q0Builder builds a [Q0] query
type Q0Builder struct {
	include []reflect.Type
	exclude []reflect.Type
}

// Query0 creates a generic query for no components.
//
// See also [ecs.World.Query].
func Query0() Q0Builder {
	return Q0Builder{}
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q0Builder) With(mask []reflect.Type) Q0Builder {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q0Builder) Without(mask []reflect.Type) Q0Builder {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q0 query for iteration.
func (q Q0Builder) Build(w *ecs.World) Q0 {
	ids := toIds(w, q.include)
	return Q0{
		w.Query(ecs.MaskPair{
			Mask:    ecs.Mask{BitMask: ecs.NewBitMask(ids...)},
			Exclude: toMask(w, q.exclude),
		}),
	}
}

// Q0 is a generic query for no components.
//
// Create one with [Query0]
type Q0 struct {
	ecs.Query
}

// Q1Builder builds a [Q1] query
type Q1Builder[A any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
}

// Query1 creates a generic query for one component.
//
// See also [ecs.World.Query].
func Query1[A any]() Q1Builder[A] {
	return Q1Builder[A]{
		include: []reflect.Type{typeOf[A]()},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q Q1Builder[A]) Optional(mask []reflect.Type) Q1Builder[A] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q1Builder[A]) With(mask []reflect.Type) Q1Builder[A] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q1Builder[A]) Without(mask []reflect.Type) Q1Builder[A] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q1 query for iteration.
func (q Q1Builder[A]) Build(w *ecs.World) Q1[A] {
	ids := toIds(w, q.include)
	return Q1[A]{
		w.Query(ecs.MaskPair{
			Mask:    toMaskOptional(w, ids, q.optional),
			Exclude: toMask(w, q.exclude),
		}),
		ids[0],
	}
}

// Q1 is a generic query for one component.
//
// Create one with [Query1]
type Q1[A any] struct {
	ecs.Query
	id ecs.ID
}

// Get1 returns the first queried component for the current query position
func (q *Q1[A]) Get1() *A {
	return (*A)(q.Query.Get(q.id))
}

//////////////////////////////////////////////////////////////////////////

// Q2Builder builds a [Q2] query
type Q2Builder[A any, B any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
}

// Query2 creates a generic query for two components.
//
// See also [ecs.World.Query].
func Query2[A any, B any]() Q2Builder[A, B] {
	return Q2Builder[A, B]{
		include: []reflect.Type{typeOf[A](), typeOf[B]()},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q Q2Builder[A, B]) Optional(mask []reflect.Type) Q2Builder[A, B] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q2Builder[A, B]) With(mask []reflect.Type) Q2Builder[A, B] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q2Builder[A, B]) Without(mask []reflect.Type) Q2Builder[A, B] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q2 query for iteration.
func (q Q2Builder[A, B]) Build(w *ecs.World) Q2[A, B] {
	ids := toIds(w, q.include)
	return Q2[A, B]{
		w.Query(ecs.MaskPair{
			Mask:    toMaskOptional(w, ids, q.optional),
			Exclude: toMask(w, q.exclude),
		}),
		ids,
	}
}

// Q2 is a generic query for two components.
//
// Create one with [Query2]
type Q2[A any, B any] struct {
	ecs.Query
	ids []ecs.ID
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

// Q3Builder builds a [Q3] query
type Q3Builder[A any, B any, C any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
}

// Query3 creates a generic query for three components.
//
// See also [ecs.World.Query].
func Query3[A any, B any, C any]() Q3Builder[A, B, C] {
	return Q3Builder[A, B, C]{
		include: []reflect.Type{typeOf[A](), typeOf[B](), typeOf[C]()},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q Q3Builder[A, B, C]) Optional(mask []reflect.Type) Q3Builder[A, B, C] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q3Builder[A, B, C]) With(mask []reflect.Type) Q3Builder[A, B, C] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q3Builder[A, B, C]) Without(mask []reflect.Type) Q3Builder[A, B, C] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q3 query for iteration.
func (q Q3Builder[A, B, C]) Build(w *ecs.World) Q3[A, B, C] {
	ids := toIds(w, q.include)
	return Q3[A, B, C]{
		w.Query(ecs.MaskPair{
			Mask:    toMaskOptional(w, ids, q.optional),
			Exclude: toMask(w, q.exclude),
		}),
		ids,
	}
}

// Q3 is a generic query for three components.
//
// Create one with [Query3]
type Q3[A any, B any, C any] struct {
	ecs.Query
	ids []ecs.ID
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

// Q4Builder builds a [Q4] query
type Q4Builder[A any, B any, C any, D any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
}

// Query4 creates a generic query for four components.
//
// See also [ecs.World.Query].
func Query4[A any, B any, C any, D any]() Q4Builder[A, B, C, D] {
	return Q4Builder[A, B, C, D]{
		include: []reflect.Type{typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D]()},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q Q4Builder[A, B, C, D]) Optional(mask []reflect.Type) Q4Builder[A, B, C, D] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q4Builder[A, B, C, D]) With(mask []reflect.Type) Q4Builder[A, B, C, D] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q4Builder[A, B, C, D]) Without(mask []reflect.Type) Q4Builder[A, B, C, D] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q4 query for iteration.
func (q Q4Builder[A, B, C, D]) Build(w *ecs.World) Q4[A, B, C, D] {
	ids := toIds(w, q.include)
	return Q4[A, B, C, D]{
		w.Query(ecs.MaskPair{
			Mask:    toMaskOptional(w, ids, q.optional),
			Exclude: toMask(w, q.exclude),
		}),
		ids,
	}
}

// Q4 is a generic query for four components.
//
// Create one with [Query4]
type Q4[A any, B any, C any, D any] struct {
	ecs.Query
	ids []ecs.ID
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

// Q5Builder builds a [Q5] query
type Q5Builder[A any, B any, C any, D any, E any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
}

// Query5 creates a generic query for five components.
//
// See also [ecs.World.Query].
func Query5[A any, B any, C any, D any, E any]() Q5Builder[A, B, C, D, E] {
	return Q5Builder[A, B, C, D, E]{
		include: []reflect.Type{
			typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
			typeOf[E](),
		},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q Q5Builder[A, B, C, D, E]) Optional(mask []reflect.Type) Q5Builder[A, B, C, D, E] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q5Builder[A, B, C, D, E]) With(mask []reflect.Type) Q5Builder[A, B, C, D, E] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q5Builder[A, B, C, D, E]) Without(mask []reflect.Type) Q5Builder[A, B, C, D, E] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q5 query for iteration.
func (q Q5Builder[A, B, C, D, E]) Build(w *ecs.World) Q5[A, B, C, D, E] {
	ids := toIds(w, q.include)
	return Q5[A, B, C, D, E]{
		w.Query(ecs.MaskPair{
			Mask:    toMaskOptional(w, ids, q.optional),
			Exclude: toMask(w, q.exclude),
		}),
		ids,
	}
}

// Q5 is a generic query for five components.
//
// Create one with [Query5]
type Q5[A any, B any, C any, D any, E any] struct {
	ecs.Query
	ids []ecs.ID
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

// Q6Builder builds a [Q6] query
type Q6Builder[A any, B any, C any, D any, E any, F any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
}

// Query6 creates a generic query for six components.
//
// See also [ecs.World.Query].
func Query6[A any, B any, C any, D any, E any, F any]() Q6Builder[A, B, C, D, E, F] {
	return Q6Builder[A, B, C, D, E, F]{
		include: []reflect.Type{
			typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
			typeOf[E](), typeOf[F](),
		},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q Q6Builder[A, B, C, D, E, F]) Optional(mask []reflect.Type) Q6Builder[A, B, C, D, E, F] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q6Builder[A, B, C, D, E, F]) With(mask []reflect.Type) Q6Builder[A, B, C, D, E, F] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q6Builder[A, B, C, D, E, F]) Without(mask []reflect.Type) Q6Builder[A, B, C, D, E, F] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q6 query for iteration.
func (q Q6Builder[A, B, C, D, E, F]) Build(w *ecs.World) Q6[A, B, C, D, E, F] {
	ids := toIds(w, q.include)
	return Q6[A, B, C, D, E, F]{
		w.Query(ecs.MaskPair{
			Mask:    toMaskOptional(w, ids, q.optional),
			Exclude: toMask(w, q.exclude),
		}),
		ids,
	}
}

// Q6 is a generic query for six components.
//
// Create one with [Query6]
type Q6[A any, B any, C any, D any, E any, F any] struct {
	ecs.Query
	ids []ecs.ID
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

// Q7Builder builds a [Q7] query
type Q7Builder[A any, B any, C any, D any, E any, F any, G any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
}

// Query7 creates a generic query for seven components.
//
// See also [ecs.World.Query].
func Query7[A any, B any, C any, D any, E any, F any, G any]() Q7Builder[A, B, C, D, E, F, G] {
	return Q7Builder[A, B, C, D, E, F, G]{
		include: []reflect.Type{
			typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
			typeOf[E](), typeOf[F](), typeOf[G](),
		},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q Q7Builder[A, B, C, D, E, F, G]) Optional(mask []reflect.Type) Q7Builder[A, B, C, D, E, F, G] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q7Builder[A, B, C, D, E, F, G]) With(mask []reflect.Type) Q7Builder[A, B, C, D, E, F, G] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q7Builder[A, B, C, D, E, F, G]) Without(mask []reflect.Type) Q7Builder[A, B, C, D, E, F, G] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q7 query for iteration.
func (q Q7Builder[A, B, C, D, E, F, G]) Build(w *ecs.World) Q7[A, B, C, D, E, F, G] {
	ids := toIds(w, q.include)
	return Q7[A, B, C, D, E, F, G]{
		w.Query(ecs.MaskPair{
			Mask:    toMaskOptional(w, ids, q.optional),
			Exclude: toMask(w, q.exclude),
		}),
		ids,
	}
}

// Q7 is a generic query for seven components.
//
// Create one with [Query7]
type Q7[A any, B any, C any, D any, E any, F any, G any] struct {
	ecs.Query
	ids []ecs.ID
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

// Q8Builder builds a [Q8] query
type Q8Builder[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
}

// Query8 creates a generic query for eight components.
//
// See also [ecs.World.Query].
func Query8[A any, B any, C any, D any, E any, F any, G any, H any]() Q8Builder[A, B, C, D, E, F, G, H] {
	return Q8Builder[A, B, C, D, E, F, G, H]{
		include: []reflect.Type{
			typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
			typeOf[E](), typeOf[F](), typeOf[G](), typeOf[H](),
		},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q Q8Builder[A, B, C, D, E, F, G, H]) Optional(mask []reflect.Type) Q8Builder[A, B, C, D, E, F, G, H] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q8Builder[A, B, C, D, E, F, G, H]) With(mask []reflect.Type) Q8Builder[A, B, C, D, E, F, G, H] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q Q8Builder[A, B, C, D, E, F, G, H]) Without(mask []reflect.Type) Q8Builder[A, B, C, D, E, F, G, H] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q8 query for iteration.
func (q Q8Builder[A, B, C, D, E, F, G, H]) Build(w *ecs.World) Q8[A, B, C, D, E, F, G, H] {
	ids := toIds(w, q.include)
	return Q8[A, B, C, D, E, F, G, H]{
		w.Query(ecs.MaskPair{
			Mask:    toMaskOptional(w, ids, q.optional),
			Exclude: toMask(w, q.exclude),
		}),
		ids,
	}
}

// Q8 is a generic query for seven components.
//
// Create one with [Query8]
type Q8[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	ecs.Query
	ids []ecs.ID
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

// Mask1 creates a component type list for one component type.
func Mask1[A any]() []reflect.Type {
	return []reflect.Type{typeOf[A]()}
}

// Mask2 creates a component type list for two component types.
func Mask2[A any, B any]() []reflect.Type {
	return []reflect.Type{typeOf[A](), typeOf[B]()}
}

// Mask3 creates a component type list for three component types.
func Mask3[A any, B any, C any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](),
	}
}

// Mask4 creates a component type list for four component types.
func Mask4[A any, B any, C any, D any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
	}
}

// Mask5 creates a component type list for five component types.
func Mask5[A any, B any, C any, D any, E any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](),
	}
}

// Mask6 creates a component type list for six component types.
func Mask6[A any, B any, C any, D any, E any, F any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](),
	}
}

// Mask7 creates a component type list for seven component types.
func Mask7[A any, B any, C any, D any, E any, F any, G any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](), typeOf[G](),
	}
}

// Mask8 creates a component type list for eight component types.
func Mask8[A any, B any, C any, D any, E any, F any, G any, H any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](), typeOf[G](), typeOf[H](),
	}
}
