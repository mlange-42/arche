package ecs

import "math/bits"

// MaskTotalBits is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = 128
const wordSize = 64

// Mask is a 128 bit bitmask.
// It is also a [Filter] for including certain components (see [All] and [Mask.Without]).
//
// # Example
//
//	world := NewWorld()
//	posID := ComponentID[Position](&world)
//	velID := ComponentID[Velocity](&world)
//
//	filter := All(posID, velID)
//	query := world.Query(filter)
type Mask struct {
	Lo uint64 // First 64 bits of the mask
	Hi uint64 // Second 64 bits of the mask
}

// All creates a new Mask from a list of IDs.
// Matches al entities that have the respective components, and potentially further components.
//
// See also [Mask.Without] and [Mask.Exact]
//
// If any [ID] is bigger or equal [MaskTotalBits], it'll not be added to the mask.
func All(ids ...ID) Mask {
	var mask Mask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}

// Matches matches a filter against a bitmask.
func (b Mask) Matches(bits Mask) bool {
	return bits.Contains(b)
}

// Without creates a [MaskFilter] which filters for including the mask's components,
// and excludes the components given as arguments.
//
// # Example
//
//	world := NewWorld()
//	posID := ComponentID[Position](&world)
//	velID := ComponentID[Velocity](&world)
//
//	filter := All(posID).Without(velID)
//	query := world.Query(filter)
func (b Mask) Without(comps ...ID) MaskFilter {
	return MaskFilter{
		Include: b,
		Exclude: All(comps...),
	}
}

// Exact creates a [MaskFilter] which filters for exactly the mask's components.
// Matches only entities that have exactly the given components, and no other.
//
// # Example
//
//	world := NewWorld()
//	posID := ComponentID[Position](&world)
//	velID := ComponentID[Velocity](&world)
//
//	filter := All(posID, velID).Exact()
//	query := world.Query(filter)
func (b Mask) Exact() MaskFilter {
	return MaskFilter{
		Include: b,
		Exclude: b.Not(),
	}
}

// Get reports if bit index defined by [ID] is true or false.
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

// Not returns the inversion of this mask.
func (b *Mask) Not() Mask {
	return Mask{
		Lo: ^b.Lo,
		Hi: ^b.Hi,
	}
}

// IsZero returns whether no bits are set in the bitmask.
func (b *Mask) IsZero() bool {
	return b.Lo == 0 && b.Hi == 0
}

// Reset changes the state of all bits to false.
func (b *Mask) Reset() {
	b.Lo, b.Lo = 0, 0
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
