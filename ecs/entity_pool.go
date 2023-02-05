package ecs

// EntityPool is the interface for entity recycling
// TODO: Add a way to shrink/compact the pool?
type EntityPool interface {
	// Get returns a fresh or recycled entity
	Get() Entity
	// Recycle hands an entity back for recycling
	Recycle(Entity)
	// Alive return whether an entity is currently alive
	Alive(Entity) bool
}

// NewEntityPool creates a new, initialized EntityPool
func NewEntityPool() EntityPool {
	return &implicitListEntityPool{
		entities:  []Entity{},
		next:      0,
		available: 0,
	}
}

// Creates a pool with proper initialization of next
func newImplicitListEntityPool() *implicitListEntityPool {
	return &implicitListEntityPool{
		entities:  []Entity{},
		next:      0,
		available: 0,
	}
}

// implicitListEntityPool is an EntityPool implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type implicitListEntityPool struct {
	entities  []Entity
	next      ID
	available uint32
}

// Get returns a fresh or recycled entity
func (p *implicitListEntityPool) Get() Entity {
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
func (p *implicitListEntityPool) Recycle(e Entity) {
	p.entities[e.id].gen++
	p.next, p.entities[e.id].id = e.id, p.next
	p.available++
}

// Alive return whether an entity is still alive, based on the entity's generations
func (p *implicitListEntityPool) Alive(e Entity) bool {
	return e.gen == p.entities[e.id].gen
}
