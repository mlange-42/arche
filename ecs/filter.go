package ecs

// Mask is a mask for a combination of components.
type Mask struct {
	mask bitMask
}

// All matches all the given components.
//
// Like [And] for combining individual components.
func All(comps ...ID) Mask {
	return Mask{newMask(comps...)}
}

// OneOf matches any of the two components.
//
// Like [Or] for combining individual components.
func OneOf(compA ID, compB ID) Or {
	return Or{Mask{newMask(compA)}, Mask{newMask(compB)}}
}

type filter interface {
	Matches(mask Mask) bool
}

// And is a filter for ANDing together components
type And struct {
	a filter
	b filter
}

// Or is a filter for ORing together components
type Or struct {
	a filter
	b filter
}

// XOr is a filter for XORing together components
type XOr struct {
	a filter
	b filter
}

// Not is a filter for excluding components
type Not Mask

// Matches matches a filter against a mask
func (f Mask) Matches(mask Mask) bool {
	return mask.mask.Contains(f.mask)
}

// Matches matches a filter against a mask
func (f *And) Matches(mask Mask) bool {
	return f.a.Matches(mask) && f.b.Matches(mask)
}

// Matches matches a filter against a mask
func (f *Or) Matches(mask Mask) bool {
	return f.a.Matches(mask) || f.b.Matches(mask)
}

// Matches matches a filter against a mask
func (f *XOr) Matches(mask Mask) bool {
	return f.a.Matches(mask) != f.b.Matches(mask)
}

// Matches matches a filter against a mask
func (f Not) Matches(mask Mask) bool {
	return !mask.mask.ContainsAny(f.mask)
}
