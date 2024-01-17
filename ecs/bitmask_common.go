package ecs

const wordSize = 64

// All creates a new Mask from a list of IDs.
// Matches all entities that have the respective components, and potentially further components.
//
// See also [Mask.Without] and [Mask.Exclusive]
//
// If any [ID] is greater than or equal to [MaskTotalBits], it will not be added to the mask.
func All(ids ...ID) Mask {
	var mask Mask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}

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
