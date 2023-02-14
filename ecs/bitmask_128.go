package ecs

import "math/bits"

// MaskTotalBits128 is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits128 = 128

const wordSize = 64

// BitMask128 is a 128 bit bitmask.
type BitMask128 struct {
	lo uint64
	hi uint64
}

// Get reports if bit index defined by ID is true or false.
//
// The return will be always false for bit >= [MaskTotalBits].
func (b BitMask128) Get(bit ID) bool {
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
func (b *BitMask128) Set(bit ID, value bool) {
	if bit < wordSize {
		if value {
			b.lo |= uint64(1 << bit)
		} else {
			b.hi &= uint64(^(1 << bit))
		}
	}
	if value {
		b.lo |= uint64(1 << (bit - wordSize))
	} else {
		b.hi &= uint64(^(1 << (bit - wordSize)))
	}
}

// Reset changes the state of all bits to false.
func (b *BitMask128) Reset() {
	b.lo = 0
	b.hi = 0
}

// Contains reports if other mask is a subset of this mask.
func (b BitMask128) Contains(other BitMask128) bool {
	return b.lo&other.lo == other.lo && b.hi&other.hi == other.hi
}

// ContainsAny reports if any bit of other mask is a subset of this mask.
func (b BitMask128) ContainsAny(other BitMask128) bool {
	return b.lo&other.lo != 0 || b.hi&other.hi != 0
}

// And is bitwise AND
func (b BitMask128) And(o BitMask128) BitMask128 {
	return BitMask128{b.lo & o.lo, b.hi & o.hi}
}

// Or is bitwise OR
func (b BitMask128) Or(o BitMask128) BitMask128 {
	return BitMask128{b.lo | o.lo, b.hi | o.hi}
}

// XOr is bitwise VOR
func (b BitMask128) XOr(o BitMask128) BitMask128 {
	return BitMask128{b.lo ^ o.lo, b.hi ^ o.hi}
}

// TotalBitsSet returns how many bits are set in this mask.
func (b BitMask128) TotalBitsSet() int {
	return bits.OnesCount64(uint64(b.hi)) + bits.OnesCount64(uint64(b.lo))
}
