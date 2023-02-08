package ecs

// bitMask is a bitmask.
type bitMask uint64

// MaskTotalBits is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = 64

var nibbleToBitsSet = [16]uint{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4}

// newMask creates a new bitmask from a list of IDs.
//
// If any ID is bigger or equal [MaskTotalBits], it'll not be added to the mask.
//
// Implementation taken from https://github.com/marioolofo/go-gameengine-ecs.
func newMask(ids ...ID) bitMask {
	var mask bitMask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}

// Get reports if bit index defined by ID is true or false.
//
// The return will be always false for bit >= [MaskTotalBits].
func (e bitMask) Get(bit ID) bool {
	mask := bitMask(1 << bit)
	return e&mask == mask
}

// Set sets the state of bit index to true or false.
//
// This function has no effect for bit >= [MaskTotalBits].
func (e *bitMask) Set(bit ID, value bool) {
	if value {
		*e |= bitMask(1 << bit)
	} else {
		*e &= bitMask(^(1 << bit))
	}
}

// Reset changes the state of all bits to false.
func (e *bitMask) Reset() {
	*e = 0
}

// Contains reports if other mask is a subset of this mask.
func (e bitMask) Contains(other bitMask) bool {
	return e&other == other
}

// TotalBitsSet returns how many bits are set in this mask.
func (e bitMask) TotalBitsSet() uint {
	var count uint

	for e != 0 {
		count += nibbleToBitsSet[e&0xf]
		e >>= 4
	}
	return count
}
