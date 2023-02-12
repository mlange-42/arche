package filter

import (
	"github.com/mlange-42/arche/ecs"
)

// Mask is a mask for a combination of components.
type Mask = ecs.Mask

// All matches all the given components.
//
// Like [And] for combining individual components.
func All(comps ...ecs.ID) Mask {
	return ecs.NewMask(comps...)
}

// OneOf matches any of the two components.
//
// Like [Or] for combining individual components.
func OneOf(compA ecs.ID, compB ecs.ID) *OR {
	return &OR{ecs.NewMask(compA), ecs.NewMask(compB)}
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

// Matches matches a filter against a mask
func (f *AND) Matches(mask ecs.BitMask) bool {
	return f.L.Matches(mask) && f.R.Matches(mask)
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

// Matches matches a filter against a mask
func (f *OR) Matches(mask ecs.BitMask) bool {
	return f.L.Matches(mask) || f.R.Matches(mask)
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

// Matches matches a filter against a mask
func (f *XOR) Matches(mask ecs.BitMask) bool {
	return f.L.Matches(mask) != f.R.Matches(mask)
}

// NOT is a filter for excluding entities with all given components
type NOT Mask

// Not constructs a NOT filter
func Not(comps ...ecs.ID) NOT {
	return NOT(ecs.NewMask(comps...))
}

// Matches matches a filter against a mask
func (f NOT) Matches(mask ecs.BitMask) bool {
	return !mask.Contains(f.BitMask)
}

// NotANY is a filter for excluding entities with any of the the given components
type NotANY Mask

// NotAny constructs a NotANY filter
func NotAny(comps ...ecs.ID) NotANY {
	return NotANY(ecs.NewMask(comps...))
}

// Matches matches a filter against a mask
func (f NotANY) Matches(mask ecs.BitMask) bool {
	return !mask.ContainsAny(f.BitMask)
}
