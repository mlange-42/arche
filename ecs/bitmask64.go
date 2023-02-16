package ecs

// bitMask64 is there just for performance comparison with the new 128 bit BitMask.
type bitMask64 uint64

func newBitMask64(ids ...ID) bitMask64 {
	var mask bitMask64
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}
func (e bitMask64) Get(bit ID) bool {
	mask := bitMask64(1 << bit)
	return e&mask == mask
}

func (e *bitMask64) Set(bit ID, value bool) {
	if value {
		*e |= bitMask64(1 << bit)
	} else {
		*e &= bitMask64(^(1 << bit))
	}
}
