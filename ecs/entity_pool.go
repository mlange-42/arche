package ecs

import (
	"math"
)

// newEntityPool creates a new, initialized Entity pool
func newEntityPool(capacityIncrement int) entityPool {
	return entityPool{
		entities:          []Entity{{0, math.MaxUint16}},
		next:              0,
		available:         0,
		capacityIncrement: capacityIncrement,
	}
}

// entityPool is an implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type entityPool struct {
	entities          []Entity
	next              eid
	available         uint32
	capacityIncrement int
}

// Get returns a fresh or recycled entity
func (p *entityPool) Get() Entity {
	if p.available == 0 {
		e := newEntity(eid(len(p.entities)))
		if len(p.entities) == cap(p.entities) {
			old := p.entities
			p.entities = make([]Entity, len(p.entities), len(p.entities)+p.capacityIncrement)
			copy(p.entities, old)
		}
		p.entities = append(p.entities, e)
		return e
	}
	curr := p.next
	p.next, p.entities[p.next].id = p.entities[p.next].id, p.next
	p.available--
	return p.entities[curr]
}

// Recycle hands an entity back for recycling
func (p *entityPool) Recycle(e Entity) {
	if e.id == 0 {
		panic("can't recycle reserved zero entity")
	}
	p.entities[e.id].gen++
	p.next, p.entities[e.id].id = e.id, p.next
	p.available++
}

// Alive return whether an entity is still alive, based on the entity's generations
func (p *entityPool) Alive(e Entity) bool {
	return e.id != 0 && e.gen == p.entities[e.id].gen
}

func (p *entityPool) Len() int {
	return len(p.entities) - 1 - int(p.available)
}
func (p *entityPool) Cap() int {
	return len(p.entities) - 1
}
func (p *entityPool) Available() int {
	return int(p.available)
}
