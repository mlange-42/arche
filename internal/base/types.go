package base

// Eid is the entity identifier/index type
type Eid uint32

// ID is the component identifier type
type ID uint8

// Mask is a mask for a combination of components.
type Mask struct {
	BitMask
}

// NewMask creates a new Mask from a list of IDs.
//
// If any ID is bigger or equal [MaskTotalBits], it'll not be added to the mask.
func NewMask(ids ...ID) Mask {
	var mask BitMask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return Mask{mask}
}

// Matches matches a filter against a mask
func (f Mask) Matches(mask BitMask) bool {
	return mask.Contains(f.BitMask)
}
