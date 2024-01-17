package ecs

import (
	"fmt"
	"math"
)

type number interface {
	int | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

// entityPool is an implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type entityPool struct {
	entities          []Entity
	next              eid
	available         uint32
	capacityIncrement uint32
}

// newEntityPool creates a new, initialized Entity pool.
func newEntityPool(capacityIncrement uint32) entityPool {
	entities := make([]Entity, 1, capacityIncrement)
	entities[0] = Entity{0, math.MaxUint32}
	return entityPool{
		entities:          entities,
		next:              0,
		available:         0,
		capacityIncrement: capacityIncrement,
	}
}

// Get returns a fresh or recycled entity.
func (p *entityPool) Get() Entity {
	if p.available == 0 {
		return p.getNew()
	}
	curr := p.next
	p.next, p.entities[p.next].id = p.entities[p.next].id, p.next
	p.available--
	return p.entities[curr]
}

// Allocates and returns a new entity. For internal use.
func (p *entityPool) getNew() Entity {
	e := newEntity(eid(len(p.entities)))
	if len(p.entities) == cap(p.entities) {
		old := p.entities
		p.entities = make([]Entity, len(p.entities), len(p.entities)+int(p.capacityIncrement))
		copy(p.entities, old)
	}
	p.entities = append(p.entities, e)
	return e
}

// Recycle hands an entity back for recycling.
func (p *entityPool) Recycle(e Entity) {
	if e.id == 0 {
		panic("can't recycle reserved zero entity")
	}
	p.entities[e.id].gen++
	p.next, p.entities[e.id].id = e.id, p.next
	p.available++
}

// Reset recycles all entities. Does NOT free the reserved memory.
func (p *entityPool) Reset() {
	p.entities = p.entities[:1]
	p.next = 0
	p.available = 0
}

// Alive returns whether an entity is still alive, based on the entity's generations.
func (p *entityPool) Alive(e Entity) bool {
	return e.gen == p.entities[e.id].gen
}

// Len returns the current number of used entities.
func (p *entityPool) Len() int {
	return len(p.entities) - 1 - int(p.available)
}

// Cap returns the current capacity (used and recycled entities).
func (p *entityPool) Cap() int {
	return len(p.entities) - 1
}

// TotalCap returns the current capacity in terms of reserved memory.
func (p *entityPool) TotalCap() int {
	return cap(p.entities)
}

// Available returns the current number of available/recycled entities.
func (p *entityPool) Available() int {
	return int(p.available)
}

// bitPool is an entityPool implementation using implicit linked lists.
type bitPool struct {
	bits      [MaskTotalBits]uint8
	next      uint8
	length    uint16
	available uint8
}

// Get returns a fresh or recycled bit.
func (p *bitPool) Get() uint8 {
	if p.available == 0 {
		return p.getNew()
	}
	curr := p.next
	p.next, p.bits[p.next] = p.bits[p.next], p.next
	p.available--
	return p.bits[curr]
}

// Allocates and returns a new bit. For internal use.
func (p *bitPool) getNew() uint8 {
	if p.length >= MaskTotalBits {
		panic(fmt.Sprintf("run out of the maximum of %d bits", MaskTotalBits))
	}
	b := uint8(p.length)
	p.bits[p.length] = b
	p.length++
	return b
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

// entityPool is an implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type intPool[T number] struct {
	pool              []T
	next              T
	available         uint32
	capacityIncrement uint32
}

// newEntityPool creates a new, initialized Entity pool.
func newIntPool[T number](capacityIncrement uint32) intPool[T] {
	return intPool[T]{
		pool:              make([]T, 0, capacityIncrement),
		next:              0,
		available:         0,
		capacityIncrement: capacityIncrement,
	}
}

// Get returns a fresh or recycled entity.
func (p *intPool[T]) Get() T {
	if p.available == 0 {
		return p.getNew()
	}
	curr := p.next
	p.next, p.pool[p.next] = p.pool[p.next], p.next
	p.available--
	return p.pool[curr]
}

// Allocates and returns a new entity. For internal use.
func (p *intPool[T]) getNew() T {
	e := T(len(p.pool))
	if len(p.pool) == cap(p.pool) {
		old := p.pool
		p.pool = make([]T, len(p.pool), len(p.pool)+int(p.capacityIncrement))
		copy(p.pool, old)
	}
	p.pool = append(p.pool, e)
	return e
}

// Recycle hands an entity back for recycling.
func (p *intPool[T]) Recycle(e T) {
	p.next, p.pool[e] = e, p.next
	p.available++
}

// Reset recycles all entities. Does NOT free the reserved memory.
func (p *intPool[T]) Reset() {
	p.pool = p.pool[:0]
	p.next = 0
	p.available = 0
}
