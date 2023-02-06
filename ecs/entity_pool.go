package ecs

// newEntityPool creates a new, initialized EntityPool
func newEntityPool() entityPool {
	return entityPool{
		entities:  []Entity{},
		next:      0,
		available: 0,
	}
}

// entityPool is an entityPool implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type entityPool struct {
	entities  []Entity
	next      ID
	available uint32
}

// Get returns a fresh or recycled entity
func (p *entityPool) Get() Entity {
	if p.available == 0 {
		e := newEntity(len(p.entities))
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
	p.entities[e.id].gen++
	p.next, p.entities[e.id].id = e.id, p.next
	p.available++
}

// Alive return whether an entity is still alive, based on the entity's generations
func (p *entityPool) Alive(e Entity) bool {
	return e.gen == p.entities[e.id].gen
}
