package ecs

import "math/bits"

// MaskTotalBits is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = 128
const wordSize = 64

// BitMask is a 128 bit bitmask.
type BitMask struct {
	lo uint64
	hi uint64
}

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
func (b BitMask) Get(bit ID) bool {
	if bit < wordSize {
		mask := uint64(1 << bit)
		return b.lo&mask == mask
	}
	mask := uint64(1 << (bit - wordSize))
	return b.hi&mask == mask
}

// Set sets the state of bit index to true or false.
//
// This function has no effect for bit >= [MaskTotalBits].
func (b *BitMask) Set(bit ID, value bool) {
	if bit < wordSize {
		if value {
			b.lo |= uint64(1 << bit)
		} else {
			b.lo &= uint64(^(1 << bit))
		}
	}
	if value {
		b.hi |= uint64(1 << (bit - wordSize))
	} else {
		b.hi &= uint64(^(1 << (bit - wordSize)))
	}
}

// IsZero returns whether no bits are set in the bitmask.
func (b BitMask) IsZero() bool {
	return b.lo == 0 && b.hi == 0
}

// Reset changes the state of all bits to false.
func (b *BitMask) Reset() {
	b.lo = 0
	b.hi = 0
}

// Contains reports if other mask is a subset of this mask.
func (b BitMask) Contains(other BitMask) bool {
	return b.lo&other.lo == other.lo && b.hi&other.hi == other.hi
}

// ContainsAny reports if any bit of other mask is a subset of this mask.
func (b BitMask) ContainsAny(other BitMask) bool {
	return b.lo&other.lo != 0 || b.hi&other.hi != 0
}

// TotalBitsSet returns how many bits are set in this mask.
func (b BitMask) TotalBitsSet() int {
	return bits.OnesCount64(uint64(b.hi)) + bits.OnesCount64(uint64(b.lo))
}
