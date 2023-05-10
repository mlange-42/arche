package ecs

// Data structure for compact storage of booleans.
type bitSet struct {
	data []uint64
}

// Get a value.
func (b *bitSet) Get(bit eid) bool {
	chunk, bit := bit/wordSize, bit%wordSize
	mask := uint64(1 << bit)
	return b.data[chunk]&mask == mask
}

// Set a value.
func (b *bitSet) Set(bit eid, value bool) {
	chunk, bit := bit/wordSize, bit%wordSize
	if value {
		b.data[chunk] |= uint64(1 << bit)
	} else {
		b.data[chunk] &= uint64(^(1 << bit))
	}
}

// Reset all values.
func (b *bitSet) Reset() {
	for i := range b.data {
		b.data[i] = 0
	}
}

// Extend to hold at least the given bits.
func (b *bitSet) ExtendTo(length int) {
	chunks, bit := length/wordSize, length%wordSize
	if bit > 0 {
		chunks++
	}
	if len(b.data) >= chunks {
		return
	}

	old := b.data
	b.data = make([]uint64, chunks)

	copy(b.data, old)
}
