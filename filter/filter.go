package filter

import (
	"github.com/mlange-42/arche/ecs"
)

// All matches entities that have all the given components.
// Synonym for [ecs.All].
//
// Like [AND] for combining individual components.
//
// See also [ecs.All], [ecs.Mask], [ecs.Mask.Without] and [ecs.Mask.Exclusive].
func All(comps ...ecs.ID) ecs.Mask {
	return ecs.All(comps...)
}

// ANY matches entities that have any of the given components.
type ANY ecs.Mask

// Any matches entities that have any of the given components.
func Any(comps ...ecs.ID) ANY {
	return ANY(ecs.All(comps...))
}

// Matches the filter against a mask.
func (f ANY) Matches(bits *ecs.Mask) bool {
	m := ecs.Mask(f)
	return bits.ContainsAny(&m)
}

// NoneOF matches entities that are missing all the given components.
type NoneOF ecs.Mask

// NoneOf matches entities that are missing all the given components.
func NoneOf(comps ...ecs.ID) NoneOF {
	return NoneOF(ecs.All(comps...))
}

// Matches the filter against a mask.
func (f NoneOF) Matches(bits *ecs.Mask) bool {
	m := ecs.Mask(f)
	return !bits.ContainsAny(&m)
}

// AnyNOT matches entities that are missing any of the given components.
type AnyNOT ecs.Mask

// AnyNot matches entities that are missing any of the given components.
func AnyNot(comps ...ecs.ID) AnyNOT {
	return AnyNOT(ecs.All(comps...))
}

// Matches the filter against a mask.
func (f AnyNOT) Matches(bits *ecs.Mask) bool {
	m := ecs.Mask(f)
	return !bits.Contains(&m)
}

// AND combines two filters using AND.
// Matches if both filters match.
//
// Ignores relation target in wrapped ecs.RelationFilter.
type AND struct {
	L ecs.Filter
	R ecs.Filter
}

// And creates an [AND] logic filter and returns a pointer to it.
func And(l, r ecs.Filter) *AND {
	return &AND{L: l, R: r}
}

// Matches the filter against a mask.
func (f *AND) Matches(bits *ecs.Mask) bool {
	return f.L.Matches(bits) && f.R.Matches(bits)
}

// OR combines two filters using OR.
// Matches if any of the filters matches.
//
// Ignores relation target in wrapped ecs.RelationFilter.
type OR struct {
	L ecs.Filter
	R ecs.Filter
}

// Or creates an [OR] logic filter and returns a pointer to it.
func Or(l, r ecs.Filter) *OR {
	return &OR{L: l, R: r}
}

// Matches the filter against a mask.
func (f *OR) Matches(bits *ecs.Mask) bool {
	return f.L.Matches(bits) || f.R.Matches(bits)
}

// XOR combines two filters using XOR.
// Matches if exactly one of the filters matches.
//
// Ignores relation target in wrapped ecs.RelationFilter.
type XOR struct {
	L ecs.Filter
	R ecs.Filter
}

// XOr creates an [XOR] logic filter and returns a pointer to it.
func XOr(l, r ecs.Filter) *XOR {
	return &XOR{L: l, R: r}
}

// Matches the filter against a mask.
func (f *XOR) Matches(bits *ecs.Mask) bool {
	return f.L.Matches(bits) != f.R.Matches(bits)
}

// NOT inverts a filter. It matches if the inner filter does not.
//
// Ignores relation target in wrapped ecs.RelationFilter.
//
// To invert a simple [ecs.Mask] filter, it is more efficient to use [ecs.Mask.Not].
type NOT struct {
	F ecs.Filter
}

// Not creates an [NOT] logic filter and returns a pointer to it.
func Not(f ecs.Filter) *NOT {
	return &NOT{F: f}
}

// Matches the filter against a mask.
func (f *NOT) Matches(bits *ecs.Mask) bool {
	return !f.F.Matches(bits)
}
