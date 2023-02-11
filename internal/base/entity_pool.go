package base

import (
	"math"
)

// NewEntityPool creates a new, initialized Entity pool
func NewEntityPool(capacityIncrement int) EntityPool {
	return EntityPool{
		entities:          []Entity{{0, math.MaxUint16}},
		next:              0,
		available:         0,
		capacityIncrement: capacityIncrement,
	}
}

// EntityPool is an implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type EntityPool struct {
	entities          []Entity
	next              Eid
	available         uint32
	capacityIncrement int
}

// Get returns a fresh or recycled entity
func (p *EntityPool) Get() Entity {
	if p.available == 0 {
		e := NewEntity(Eid(len(p.entities)))
		if len(p.entities) == cap(p.entities) {
			old := p.entities
			p.entities = make([]Entity, len(p.entities), len(p.entities)+p.capacityIncrement)
			copy(p.entities, old)
		}
		p.entities = append(p.entities, e)
		return e
	}
	curr := p.next
	p.next, p.entities[p.next].ID = p.entities[p.next].ID, p.next
	p.available--
	return p.entities[curr]
}

// Recycle hands an entity back for recycling
func (p *EntityPool) Recycle(e Entity) {
	if e.ID == 0 {
		panic("can't recycle reserved zero entity")
	}
	p.entities[e.ID].gen++
	p.next, p.entities[e.ID].ID = e.ID, p.next
	p.available++
}

// Alive return whether an entity is still alive, based on the entity's generations
func (p *EntityPool) Alive(e Entity) bool {
	return e.ID != 0 && e.gen == p.entities[e.ID].gen
}
