package filter

import "internal/base"

// Mask is a mask for a combination of components.
type Mask = base.Mask

// MaskFilter is the interface for logic filters
type MaskFilter interface {
	Matches(mask base.BitMask) bool
}

// All matches all the given components.
//
// Like [And] for combining individual components.
func All(comps ...base.ID) Mask {
	return base.NewMask(comps...)
}

// OneOf matches any of the two components.
//
// Like [Or] for combining individual components.
func OneOf(compA base.ID, compB base.ID) Or {
	return Or{base.NewMask(compA), base.NewMask(compB)}
}

// And is a filter for ANDing together components
type And struct {
	a MaskFilter
	b MaskFilter
}

// Or is a filter for ORing together components
type Or struct {
	a MaskFilter
	b MaskFilter
}

// XOr is a filter for XORing together components
type XOr struct {
	a MaskFilter
	b MaskFilter
}

// Not is a filter for excluding components
type Not Mask

// Matches matches a filter against a mask
func (f *And) Matches(mask base.BitMask) bool {
	return f.a.Matches(mask) && f.b.Matches(mask)
}

// Matches matches a filter against a mask
func (f *Or) Matches(mask base.BitMask) bool {
	return f.a.Matches(mask) || f.b.Matches(mask)
}

// Matches matches a filter against a mask
func (f *XOr) Matches(mask base.BitMask) bool {
	return f.a.Matches(mask) != f.b.Matches(mask)
}

// Matches matches a filter against a mask
func (f Not) Matches(mask base.BitMask) bool {
	return !mask.ContainsAny(f.BitMask)
}
