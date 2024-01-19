package ecs

const wordSize = 64

// Matches the mask as filter against another mask.
func (b Mask) Matches(bits *Mask) bool {
	return bits.Contains(&b)
}

// Without creates a [MaskFilter] which filters for including the mask's components,
// and excludes the components given as arguments.
func (b Mask) Without(comps ...ID) MaskFilter {
	return MaskFilter{
		Include: b,
		Exclude: All(comps...),
	}
}

// Exclusive creates a [MaskFilter] which filters for exactly the mask's components.
// Matches only entities that have exactly the given components, and no other.
func (b Mask) Exclusive() MaskFilter {
	return MaskFilter{
		Include: b,
		Exclude: b.Not(),
	}
}
