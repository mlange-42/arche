package generic

import (
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

type compiledQuery struct {
	mask     ecs.Mask
	exclude  ecs.Mask
	Ids      []ecs.ID
	compiled bool
}

func (q *compiledQuery) Compile(w *ecs.World, include, optional, exclude []reflect.Type) {
	if q.compiled {
		return
	}
	q.Ids = toIds(w, include)
	q.mask = toMaskOptional(w, q.Ids, optional)
	q.exclude = toMask(w, exclude)
	q.compiled = true
}

func (q *compiledQuery) MaskPair() ecs.MaskPair {
	return ecs.MaskPair{
		Mask:    q.mask,
		Exclude: q.exclude,
	}
}

// Filter0 builds a [Query0] query
type Filter0 struct {
	include  []reflect.Type
	exclude  []reflect.Type
	compiled compiledQuery
}

// NewFilter0 creates a generic filter for no components.
//
// See also [ecs.World.Query].
func NewFilter0() *Filter0 {
	return &Filter0{}
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Filter0) With(mask []reflect.Type) *Filter0 {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Filter0) Without(mask []reflect.Type) *Filter0 {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a [Query0] query for iteration.
func (q *Filter0) Build(w *ecs.World) Query0 {
	q.compiled.Compile(w, q.include, []reflect.Type{}, q.exclude)
	return Query0{
		w.Query(q.compiled.MaskPair()),
	}
}

// Query0 is a generic query iterator for no components.
//
// Create one with [NewFilter0] and [Filter0.Build].
type Query0 struct {
	ecs.Query
}

// Filter1 builds a [Query1] query
type Filter1[A any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
	compiled compiledQuery
}

// NewFilter1 creates a generic filter for one component.
//
// See also [ecs.World.Query].
func NewFilter1[A any]() *Filter1[A] {
	return &Filter1[A]{
		include: []reflect.Type{typeOf[A]()},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q *Filter1[A]) Optional(mask []reflect.Type) *Filter1[A] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Filter1[A]) With(mask []reflect.Type) *Filter1[A] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Filter1[A]) Without(mask []reflect.Type) *Filter1[A] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a [Query1] query for iteration.
func (q *Filter1[A]) Build(w *ecs.World) Query1[A] {
	q.compiled.Compile(w, q.include, q.optional, q.exclude)
	return Query1[A]{
		w.Query(q.compiled.MaskPair()),
		q.compiled.Ids[0],
	}
}

// Query1 is a generic query iterator for one component.
//
// Create one with [NewFilter1] and [Filter1.Build]
type Query1[A any] struct {
	ecs.Query
	id ecs.ID
}

// Get1 returns the first queried component for the current query position
func (q *Query1[A]) Get1() *A {
	return (*A)(q.Query.Get(q.id))
}

//////////////////////////////////////////////////////////////////////////

// Filter2 builds a [Query2] query
type Filter2[A any, B any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
	compiled compiledQuery
}

// NewQuery2 creates a generic filter for two components.
//
// See also [ecs.World.Query].
func NewQuery2[A any, B any]() *Filter2[A, B] {
	return &Filter2[A, B]{
		include: []reflect.Type{typeOf[A](), typeOf[B]()},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q *Filter2[A, B]) Optional(mask []reflect.Type) *Filter2[A, B] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Filter2[A, B]) With(mask []reflect.Type) *Filter2[A, B] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Filter2[A, B]) Without(mask []reflect.Type) *Filter2[A, B] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a [Query2] query for iteration.
func (q *Filter2[A, B]) Build(w *ecs.World) Query2[A, B] {
	q.compiled.Compile(w, q.include, q.optional, q.exclude)
	return Query2[A, B]{
		w.Query(q.compiled.MaskPair()),
		q.compiled.Ids,
	}
}

// Query2 is a generic query iterator for two components.
//
// Create one with [NewFilter2] and [Filter2.Build]
type Query2[A any, B any] struct {
	ecs.Query
	ids []ecs.ID
}

// GetAll returns all queried components for the current query position
func (q *Query2[A, B]) GetAll() (*A, *B) {
	return (*A)(q.Query.Get(q.ids[0])), (*B)(q.Query.Get(q.ids[1]))
}

// Get1 returns the first queried component for the current query position
func (q *Query2[A, B]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get2 returns the second queried component for the current query position
func (q *Query2[A, B]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

//////////////////////////////////////////////////////////////////////////

// Filter3 builds a [Query3] query
type Filter3[A any, B any, C any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
	compiled compiledQuery
}

// NewFilter3 creates a generic filter for three components.
//
// See also [ecs.World.Query].
func NewFilter3[A any, B any, C any]() *Filter3[A, B, C] {
	return &Filter3[A, B, C]{
		include: []reflect.Type{typeOf[A](), typeOf[B](), typeOf[C]()},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q *Filter3[A, B, C]) Optional(mask []reflect.Type) *Filter3[A, B, C] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Filter3[A, B, C]) With(mask []reflect.Type) *Filter3[A, B, C] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Filter3[A, B, C]) Without(mask []reflect.Type) *Filter3[A, B, C] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a [Query3] query for iteration.
func (q *Filter3[A, B, C]) Build(w *ecs.World) Query3[A, B, C] {
	q.compiled.Compile(w, q.include, q.optional, q.exclude)
	return Query3[A, B, C]{
		w.Query(q.compiled.MaskPair()),
		q.compiled.Ids,
	}
}

// Query3 is a generic query iterator for three components.
//
// Create one with [NewFilter3] and [Filter3.Build]
type Query3[A any, B any, C any] struct {
	ecs.Query
	ids []ecs.ID
}

// GetAll returns all queried components for the current query position
func (q *Query3[A, B, C]) GetAll() (*A, *B, *C) {
	return (*A)(q.Query.Get(q.ids[0])), (*B)(q.Query.Get(q.ids[1])), (*C)(q.Query.Get(q.ids[2]))
}

// Get1 returns the first queried component for the current query position
func (q *Query3[A, B, C]) Get1() *A {
	return (*A)(q.Query.Get(q.ids[0]))
}

// Get2 returns the second queried component for the current query position
func (q *Query3[A, B, C]) Get2() *B {
	return (*B)(q.Query.Get(q.ids[1]))
}

// Get3 returns the third queried component for the current query position
func (q *Query3[A, B, C]) Get3() *C {
	return (*C)(q.Query.Get(q.ids[2]))
}

//////////////////////////////////////////////////////////////////////////

// Query4 builds a [Q4] query
type Query4[A any, B any, C any, D any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
	compiled compiledQuery
}

// NewQuery4 creates a generic query for four components.
//
// See also [ecs.World.Query].
func NewQuery4[A any, B any, C any, D any]() *Query4[A, B, C, D] {
	return &Query4[A, B, C, D]{
		include: []reflect.Type{typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D]()},
	}
}

// Optional makes some of the query's components optional.
//
// Create the required mask with [Mask1], [Mask2], etc.
//
// Only affects component types that were specified in the query.
func (q *Query4[A, B, C, D]) Optional(mask []reflect.Type) *Query4[A, B, C, D] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Query4[A, B, C, D]) With(mask []reflect.Type) *Query4[A, B, C, D] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Query4[A, B, C, D]) Without(mask []reflect.Type) *Query4[A, B, C, D] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q4 query for iteration.
func (q *Query4[A, B, C, D]) Build(w *ecs.World) Q4[A, B, C, D] {
	q.compiled.Compile(w, q.include, q.optional, q.exclude)
	return Q4[A, B, C, D]{
		w.Query(q.compiled.MaskPair()),
		q.compiled.Ids,
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

// Query5 builds a [Q5] query
type Query5[A any, B any, C any, D any, E any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
	compiled compiledQuery
}

// NewQuery5 creates a generic query for five components.
//
// See also [ecs.World.Query].
func NewQuery5[A any, B any, C any, D any, E any]() *Query5[A, B, C, D, E] {
	return &Query5[A, B, C, D, E]{
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
func (q *Query5[A, B, C, D, E]) Optional(mask []reflect.Type) *Query5[A, B, C, D, E] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Query5[A, B, C, D, E]) With(mask []reflect.Type) *Query5[A, B, C, D, E] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Query5[A, B, C, D, E]) Without(mask []reflect.Type) *Query5[A, B, C, D, E] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q5 query for iteration.
func (q *Query5[A, B, C, D, E]) Build(w *ecs.World) Q5[A, B, C, D, E] {
	q.compiled.Compile(w, q.include, q.optional, q.exclude)
	return Q5[A, B, C, D, E]{
		w.Query(q.compiled.MaskPair()),
		q.compiled.Ids,
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

// Query6 builds a [Q6] query
type Query6[A any, B any, C any, D any, E any, F any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
	compiled compiledQuery
}

// NewQuery6 creates a generic query for six components.
//
// See also [ecs.World.Query].
func NewQuery6[A any, B any, C any, D any, E any, F any]() *Query6[A, B, C, D, E, F] {
	return &Query6[A, B, C, D, E, F]{
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
func (q *Query6[A, B, C, D, E, F]) Optional(mask []reflect.Type) *Query6[A, B, C, D, E, F] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Query6[A, B, C, D, E, F]) With(mask []reflect.Type) *Query6[A, B, C, D, E, F] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Query6[A, B, C, D, E, F]) Without(mask []reflect.Type) *Query6[A, B, C, D, E, F] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q6 query for iteration.
func (q *Query6[A, B, C, D, E, F]) Build(w *ecs.World) Q6[A, B, C, D, E, F] {
	q.compiled.Compile(w, q.include, q.optional, q.exclude)
	return Q6[A, B, C, D, E, F]{
		w.Query(q.compiled.MaskPair()),
		q.compiled.Ids,
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

// Query7 builds a [Q7] query
type Query7[A any, B any, C any, D any, E any, F any, G any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
	compiled compiledQuery
}

// NewQuery7 creates a generic query for seven components.
//
// See also [ecs.World.Query].
func NewQuery7[A any, B any, C any, D any, E any, F any, G any]() *Query7[A, B, C, D, E, F, G] {
	return &Query7[A, B, C, D, E, F, G]{
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
func (q *Query7[A, B, C, D, E, F, G]) Optional(mask []reflect.Type) *Query7[A, B, C, D, E, F, G] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Query7[A, B, C, D, E, F, G]) With(mask []reflect.Type) *Query7[A, B, C, D, E, F, G] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Query7[A, B, C, D, E, F, G]) Without(mask []reflect.Type) *Query7[A, B, C, D, E, F, G] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q7 query for iteration.
func (q *Query7[A, B, C, D, E, F, G]) Build(w *ecs.World) Q7[A, B, C, D, E, F, G] {
	q.compiled.Compile(w, q.include, q.optional, q.exclude)
	return Q7[A, B, C, D, E, F, G]{
		w.Query(q.compiled.MaskPair()),
		q.compiled.Ids,
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

// Query8 builds a [Q8] query
type Query8[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	include  []reflect.Type
	optional []reflect.Type
	exclude  []reflect.Type
	compiled compiledQuery
}

// NewQuery8 creates a generic query for eight components.
//
// See also [ecs.World.Query].
func NewQuery8[A any, B any, C any, D any, E any, F any, G any, H any]() *Query8[A, B, C, D, E, F, G, H] {
	return &Query8[A, B, C, D, E, F, G, H]{
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
func (q *Query8[A, B, C, D, E, F, G, H]) Optional(mask []reflect.Type) *Query8[A, B, C, D, E, F, G, H] {
	q.optional = append(q.optional, mask...)
	return q
}

// With adds more required components that are not accessible using Get... methods.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Query8[A, B, C, D, E, F, G, H]) With(mask []reflect.Type) *Query8[A, B, C, D, E, F, G, H] {
	q.include = append(q.include, mask...)
	return q
}

// Without excludes entities with any of the given components from the query.
//
// Create the required mask with [Mask1], [Mask2], etc.
func (q *Query8[A, B, C, D, E, F, G, H]) Without(mask []reflect.Type) *Query8[A, B, C, D, E, F, G, H] {
	q.exclude = append(q.exclude, mask...)
	return q
}

// Build builds a Q8 query for iteration.
func (q *Query8[A, B, C, D, E, F, G, H]) Build(w *ecs.World) Q8[A, B, C, D, E, F, G, H] {
	q.compiled.Compile(w, q.include, q.optional, q.exclude)
	return Q8[A, B, C, D, E, F, G, H]{
		w.Query(q.compiled.MaskPair()),
		q.compiled.Ids,
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
