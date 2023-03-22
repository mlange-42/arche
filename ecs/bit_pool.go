package ecs

// bitPool is an entityPool implementation using implicit linked lists.
type bitPool struct {
	bits      [MaskTotalBits]uint8
	next      uint8
	length    uint8
	available uint8
}

// newBitPool creates a new, initialized bit pool.
func newBitPool() bitPool {
	return bitPool{
		bits:      [MaskTotalBits]uint8{},
		next:      0,
		length:    0,
		available: 0,
	}
}

// Get returns a fresh or recycled bit.
func (p *bitPool) Get() uint8 {
	if p.available == 0 {
		if p.length >= MaskTotalBits {
			panic("run out of the maximum of 128 bits")
		}
		b := p.length
		p.bits[p.length] = b
		p.length++
		return b
	}
	curr := p.next
	p.next, p.bits[p.next] = p.bits[p.next], p.next
	p.available--
	return p.bits[curr]
}

// Recycle hands a bit back for recycling.
func (p *bitPool) Recycle(b uint8) {
	p.next, p.bits[b] = b, p.next
	p.available++
}

// Reset recycles all bits.
func (p *bitPool) Reset() {
	p.next = 0
	p.length = 0
	p.available = 0
}
