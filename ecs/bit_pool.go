package ecs

// bitPool is an entityPool implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type bitPool struct {
	bits      [MaskTotalBits]uint8
	next      uint8
	length    uint8
	available uint8
}

// newBitPool creates a new, initialized bit pool
func newBitPool() bitPool {
	return bitPool{
		bits:      [MaskTotalBits]uint8{},
		next:      0,
		length:    0,
		available: 0,
	}
}

// Get returns a fresh or recycled bit
func (p *bitPool) Get() uint8 {
	if p.available == 0 {
		if p.length >= MaskTotalBits {
			panic("run out of the maximum of 64 bits")
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

// Recycle hands a bit back for recycling
func (p *bitPool) Recycle(b uint8) {
	p.next, p.bits[b] = b, p.next
	p.available++
}
