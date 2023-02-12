package filter

import (
	"github.com/mlange-42/arche/ecs"
)

// Mask is a mask for a combination of components.
type Mask = ecs.Mask

// OneOf matches any of the two components.
//
// Like [Or] for combining individual components.
func OneOf(compA ecs.ID, compB ecs.ID) *OR {
	return &OR{ecs.All(compA), ecs.All(compB)}
}

// ALL is a filter including entities with all the given components
type ALL Mask

// All matches all the given components.
//
// Like [And] for combining individual components.
func All(comps ...ecs.ID) ALL {
	return ALL(ecs.All(comps...))
}

// Not inverts this filter to exclude entities with all the given components
func (f ALL) Not() NOT {
	return NOT(f)
}

// Matches matches a filter against a bitmask
func (f ALL) Matches(bits ecs.BitMask) bool {
	return bits.Contains(f.BitMask)
}

// ANY is a filter including entities with any of the given components
type ANY Mask

// Any constructs a NotANY filter
func Any(comps ...ecs.ID) ANY {
	return ANY(ecs.All(comps...))
}

// Not inverts this filter to exclude entities with any of the given components
func (f ANY) Not() NotANY {
	return NotANY(f)
}

// Matches matches a filter against a bitmask
func (f ANY) Matches(bits ecs.BitMask) bool {
	return bits.ContainsAny(f.BitMask)
}

// NOT is a filter for excluding entities with all given components
type NOT Mask

// Not constructs a NOT filter
func Not(comps ...ecs.ID) NOT {
	return NOT(ecs.All(comps...))
}

// Matches matches a filter against a bitmask
func (f NOT) Matches(bits ecs.BitMask) bool {
	return !bits.Contains(f.BitMask)
}

// NotANY is a filter for excluding entities with any of the given components
type NotANY Mask

// NotAny constructs a NotANY filter
func NotAny(comps ...ecs.ID) NotANY {
	return NotANY(ecs.All(comps...))
}

// Matches matches a filter against a bitmask
func (f NotANY) Matches(bits ecs.BitMask) bool {
	return !bits.ContainsAny(f.BitMask)
}

// AND is a filter for ANDing together components
type AND struct {
	L ecs.Filter
	R ecs.Filter
}

// And constructs a pointer to a AND filter
func And(l, r ecs.Filter) *AND {
	return &AND{L: l, R: r}
}

// Matches matches a filter against a bitmask
func (f *AND) Matches(bits ecs.BitMask) bool {
	return f.L.Matches(bits) && f.R.Matches(bits)
}

// OR is a filter for ORing together components
type OR struct {
	L ecs.Filter
	R ecs.Filter
}

// Or constructs a pointer to a OR filter
func Or(l, r ecs.Filter) *OR {
	return &OR{L: l, R: r}
}

// Matches matches a filter against a bitmask
func (f *OR) Matches(bits ecs.BitMask) bool {
	return f.L.Matches(bits) || f.R.Matches(bits)
}

// XOR is a filter for XORing together components
type XOR struct {
	L ecs.Filter
	R ecs.Filter
}

// XOr constructs a pointer to a XOR filter
func XOr(l, r ecs.Filter) *XOR {
	return &XOR{L: l, R: r}
}

// Matches matches a filter against a bitmask
func (f *XOR) Matches(bits ecs.BitMask) bool {
	return f.L.Matches(bits) != f.R.Matches(bits)
}
