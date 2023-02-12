package filter

import (
	"github.com/mlange-42/arche/ecs"
)

// Mask is a mask for a combination of components.
type Mask = ecs.Mask

// ALL matches entities that have all the given components.
type ALL Mask

// All matches entities that have all the given components.
//
// Like [And] for combining individual components.
func All(comps ...ecs.ID) ALL {
	return ALL(ecs.All(comps...))
}

// Not inverts this filter to exclude entities with all the given components
func (f ALL) Not() NoneOF {
	return NoneOF(f)
}

// Matches matches a filter against a bitmask
func (f ALL) Matches(bits ecs.BitMask) bool {
	return bits.Contains(f.BitMask)
}

// ANY matches entities that have any of the given components.
type ANY Mask

// Any matches entities that have any of the given components.
func Any(comps ...ecs.ID) ANY {
	return ANY(ecs.All(comps...))
}

// Not inverts this filter to exclude entities with any of the given components
func (f ANY) Not() AnyNOT {
	return AnyNOT(f)
}

// Matches matches a filter against a bitmask
func (f ANY) Matches(bits ecs.BitMask) bool {
	return bits.ContainsAny(f.BitMask)
}

// NoneOF matches entities that are missing all the given components.
type NoneOF Mask

// NoneOf matches entities that are missing all the given components.
func NoneOf(comps ...ecs.ID) NoneOF {
	return NoneOF(ecs.All(comps...))
}

// Matches matches a filter against a bitmask
func (f NoneOF) Matches(bits ecs.BitMask) bool {
	return !bits.ContainsAny(f.BitMask)
}

// AnyNOT matches entities that are missing any of the given components.
type AnyNOT Mask

// AnyNot matches entities that are missing any of the given components.
func AnyNot(comps ...ecs.ID) AnyNOT {
	return AnyNOT(ecs.All(comps...))
}

// Matches matches a filter against a bitmask
func (f AnyNOT) Matches(bits ecs.BitMask) bool {
	return !bits.Contains(f.BitMask)
}

// AND combines two filters using AND.
type AND struct {
	L ecs.Filter
	R ecs.Filter
}

// And combines two filters using AND.
func And(l, r ecs.Filter) *AND {
	return &AND{L: l, R: r}
}

// Matches matches a filter against a bitmask
func (f *AND) Matches(bits ecs.BitMask) bool {
	return f.L.Matches(bits) && f.R.Matches(bits)
}

// OR combines two filters using OR.
type OR struct {
	L ecs.Filter
	R ecs.Filter
}

// Or combines two filters using OR.
func Or(l, r ecs.Filter) *OR {
	return &OR{L: l, R: r}
}

// Matches matches a filter against a bitmask
func (f *OR) Matches(bits ecs.BitMask) bool {
	return f.L.Matches(bits) || f.R.Matches(bits)
}

// XOR combines two filters using XOR.
type XOR struct {
	L ecs.Filter
	R ecs.Filter
}

// XOr combines two filters using XOR.
func XOr(l, r ecs.Filter) *XOR {
	return &XOR{L: l, R: r}
}

// Matches matches a filter against a bitmask
func (f *XOR) Matches(bits ecs.BitMask) bool {
	return f.L.Matches(bits) != f.R.Matches(bits)
}

// NOT inverts a filter. It matches if the inner filter does not.
type NOT struct {
	f ecs.Filter
}

// Not inverts a filter. It matches if the inner filter does not.
func Not(f ecs.Filter) *NOT {
	return &NOT{f: f}
}

// Matches matches a filter against a bitmask
func (f *NOT) Matches(bits ecs.BitMask) bool {
	return !f.f.Matches(bits)
}
