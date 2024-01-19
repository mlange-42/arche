//go:build tiny

package ecs

import (
	"math/bits"
)

// MaskTotalBits is the size of a [Mask] in bits.
// It is the maximum number of component types that may exist in any [World].
//
// ⚠️ This build uses the build tag `tiny`. Remove the tag for 256 bit masks.
const MaskTotalBits = 64

// Mask is a 64 bit bitmask.
// It is also a [Filter] for including certain components.
//
// Use [All] to create a mask for a list of component IDs.
// A mask can be further specified using [Mask.Without] or [Mask.Exclusive].
//
// ⚠️ This build uses the build tag `tiny`. Remove the tag for 256 bit masks.
type Mask struct {
	bits uint64 // 64 bits of the mask
}

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

// Get reports whether the bit at the given index [ID] is set.
//
// Returns false for bit >= [MaskTotalBits].
func (b *Mask) Get(bit ID) bool {
	mask := uint64(1 << bit.id)
	return b.bits&mask == mask
}

// Set sets the state of the bit at the given index.
//
// Has no effect for bit >= [MaskTotalBits].
func (b *Mask) Set(bit ID, value bool) {
	if value {
		b.bits |= (1 << bit.id)
	} else {
		b.bits &= ^(1 << bit.id)
	}
}

// Not returns the inversion of this mask.
func (b *Mask) Not() Mask {
	return Mask{
		bits: ^b.bits,
	}
}

// IsZero returns whether no bits are set in the mask.
func (b *Mask) IsZero() bool {
	return b.bits == 0
}

// Reset the mask setting all bits to false.
func (b *Mask) Reset() {
	b.bits = 0
}

// Contains reports if the other mask is a subset of this mask.
func (b *Mask) Contains(other *Mask) bool {
	return b.bits&other.bits == other.bits
}

// ContainsAny reports if any bit of the other mask is in this mask.
func (b *Mask) ContainsAny(other *Mask) bool {
	return b.bits&other.bits != 0
}

// And returns the bitwise AND of two masks.
func (b *Mask) And(other *Mask) Mask {
	return Mask{
		bits: b.bits & other.bits,
	}
}

// Or returns the bitwise OR of two masks.
func (b *Mask) Or(other *Mask) Mask {
	return Mask{
		bits: b.bits | other.bits,
	}
}

// Xor returns the bitwise XOR of two masks.
func (b *Mask) Xor(other *Mask) Mask {
	return Mask{
		bits: b.bits ^ other.bits,
	}
}

// TotalBitsSet returns how many bits are set in this mask.
func (b *Mask) TotalBitsSet() int {
	return bits.OnesCount64(b.bits)
}

func (b *Mask) toTypes(reg *componentRegistry) []componentType {
	if b.bits == 0 {
		return []componentType{}
	}

	count := int(b.TotalBitsSet())
	types := make([]componentType, count)

	idx := 0
	for j := 0; j < wordSize; j++ {
		id := ID{id: uint8(j)}
		if b.Get(id) {
			types[idx] = componentType{ID: id, Type: reg.Types[id.id]}
			idx++
		}
	}
	return types
}
