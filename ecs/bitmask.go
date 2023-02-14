package ecs

import "math/bits"

// MaskTotalBits is the size of BitMask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = 64

// BitMask is a 64 bit bitmask.
type BitMask uint64

// NewBitMask creates a new bitmask from a list of IDs.
//
// If any ID is bigger or equal [MaskTotalBits], it'll not be added to the mask.
func NewBitMask(ids ...ID) BitMask {
	var mask BitMask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}

// Get reports if bit index defined by ID is true or false.
//
// The return will be always false for bit >= [MaskTotalBits].
func (e BitMask) Get(bit ID) bool {
	mask := BitMask(1 << bit)
	return e&mask == mask
}

// Set sets the state of bit index to true or false.
//
// This function has no effect for bit >= [MaskTotalBits].
func (e *BitMask) Set(bit ID, value bool) {
	if value {
		*e |= BitMask(1 << bit)
	} else {
		*e &= BitMask(^(1 << bit))
	}
}

// Reset changes the state of all bits to false.
func (e *BitMask) Reset() {
	*e = 0
}

// Contains reports if other mask is a subset of this mask.
func (e BitMask) Contains(other BitMask) bool {
	return e&other == other
}

// ContainsAny reports if any bit of other mask is a subset of this mask.
func (e BitMask) ContainsAny(other BitMask) bool {
	return e&other != 0
}

// TotalBitsSet returns how many bits are set in this mask.
func (e BitMask) TotalBitsSet() int {
	return bits.OnesCount64(uint64(e))
}
