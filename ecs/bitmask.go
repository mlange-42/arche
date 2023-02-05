package ecs

// Mask is a bitmask
type Mask uint64

// MaskTotalBits is the size of Mask in bits
const MaskTotalBits = 64

var nibbleToBitsSet = [16]uint{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4}

// NewMask creates a new bitmask from a list of IDs
// If any ID is bigger or equal MaskTotalBits, it'll not be added to the mask
// Implementation taken from https://github.com/marioolofo/go-gameengine-ecs
func NewMask(ids ...ID) Mask {
	var mask Mask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}

// Get reports if bit index defined by ID is true or false
// The return will be always false for bit >= MaskTotalBits
func (e Mask) Get(bit ID) bool {
	mask := Mask(1 << bit)
	return e&mask == mask
}

// Set sets the state of bit index to true or false
// This function has no effect for bit >= MaskTotalBits
func (e *Mask) Set(bit ID, value bool) {
	if value {
		*e |= Mask(1 << bit)
	} else {
		*e &= Mask(^(1 << bit))
	}
}

// Reset change the state of all bits to false
func (e *Mask) Reset() {
	*e = 0
}

// Contains reports if other mask is a subset of this mask
func (e Mask) Contains(other Mask) bool {
	return e&other == other
}

// TotalBitsSet returns how many bits are set in this mask
func (e Mask) TotalBitsSet() uint {
	var count uint

	for e != 0 {
		count += nibbleToBitsSet[e&0xf]
		e >>= 4
	}
	return count
}
