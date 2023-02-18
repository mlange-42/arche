package ecs

import "math/bits"

// MaskTotalBits is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = 128
const wordSize = 64

// Mask is a 128 bit bitmask.
// It is also a [Filter] for including certain components (see [All]).
type Mask struct {
	Lo uint64
	Hi uint64
}

// All creates a new Mask from a list of IDs.
//
// If any ID is bigger or equal [MaskTotalBits], it'll not be added to the mask.
func All(ids ...ID) Mask {
	var mask Mask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}

// Without creates a [MaskPair] which filters for including the mask's components,
// and excludes the components given as arguments.
func (b Mask) Without(comps ...ID) MaskPair {
	return MaskPair{
		Include: b,
		Exclude: All(comps...),
	}
}

// Get reports if bit index defined by ID is true or false.
//
// The return will be always false for bit >= [MaskTotalBits].
func (b *Mask) Get(bit ID) bool {
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
func (b *Mask) Set(bit ID, value bool) {
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
func (b *Mask) IsZero() bool {
	return b.Lo == 0 && b.Hi == 0
}

// Reset changes the state of all bits to false.
func (b *Mask) Reset() {
	b.Lo = 0
	b.Hi = 0
}

// Contains reports if other mask is a subset of this mask.
func (b *Mask) Contains(other Mask) bool {
	return b.Lo&other.Lo == other.Lo && b.Hi&other.Hi == other.Hi
}

// ContainsAny reports if any bit of other mask is in this mask.
func (b *Mask) ContainsAny(other Mask) bool {
	return b.Lo&other.Lo != 0 || b.Hi&other.Hi != 0
}

// TotalBitsSet returns how many bits are set in this mask.
func (b *Mask) TotalBitsSet() int {
	return bits.OnesCount64(b.Hi) + bits.OnesCount64(b.Lo)
}

// MaskPair is a filter for including and excluding components.
// It is a [Filter] for including and excluding certain components (see [Mask.Without]).
type MaskPair struct {
	Include Mask
	Exclude Mask
}

// Matches matches a filter against a mask.
func (f *MaskPair) Matches(bits Mask) bool {
	return bits.Contains(f.Include) &&
		(f.Exclude.IsZero() || !bits.ContainsAny(f.Exclude))
}

// Matches matches a filter against a bitmask.
func (b Mask) Matches(bits Mask) bool {
	return bits.Contains(b)
}
