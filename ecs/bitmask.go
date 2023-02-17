package ecs

import "math/bits"

// MaskTotalBits is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = 128
const wordSize = 64

// BitMask is a 128 bit bitmask.
type BitMask struct {
	Lo uint64
	Hi uint64
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
		return b.Lo&mask == mask
	}
	mask := uint64(1 << (bit - wordSize))
	return b.Hi&mask == mask
}

// Set sets the state of bit index to true or false.
//
// This function has no effect for bit >= [MaskTotalBits].
func (b *BitMask) Set(bit ID, value bool) {
	if bit < wordSize {
		if value {
			b.Lo |= uint64(1 << bit)
		} else {
			b.Lo &= uint64(^(1 << bit))
		}
	}
	if value {
		b.Hi |= uint64(1 << (bit - wordSize))
	} else {
		b.Hi &= uint64(^(1 << (bit - wordSize)))
	}
}

// IsZero returns whether no bits are set in the bitmask.
func (b BitMask) IsZero() bool {
	return b.Lo == 0 && b.Hi == 0
}

// Reset changes the state of all bits to false.
func (b *BitMask) Reset() {
	b.Lo = 0
	b.Hi = 0
}

// Contains reports if other mask is a subset of this mask.
func (b BitMask) Contains(other BitMask) bool {
	return b.Lo&other.Lo == other.Lo && b.Hi&other.Hi == other.Hi
}

// ContainsAny reports if any bit of other mask is in this mask.
func (b BitMask) ContainsAny(other BitMask) bool {
	return b.Lo&other.Lo != 0 || b.Hi&other.Hi != 0
}

// TotalBitsSet returns how many bits are set in this mask.
func (b BitMask) TotalBitsSet() int {
	return bits.OnesCount64(b.Hi) + bits.OnesCount64(b.Lo)
}
