package ecs

// EntityPool is the interface for entity recycling
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
		next:      newEntity(0),
		available: 0,
	}
}

// Creates a pool with proper initialization of next
// TODO: check for a better solution than using the max generation
// The problem it that we identify dead entities partially by tha fact that their ID does not match their position.
// Depending on the initialization of `next` and the first recycled entity, this assumption may not be valid.
func newImplicitListEntityPool() *implicitListEntityPool {
	return &implicitListEntityPool{
		entities:  []Entity{},
		next:      Entity{0, 1},
		available: 0,
	}
}

// implicitListEntityPool is an EntityPool implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type implicitListEntityPool struct {
	entities  []Entity
	next      Entity
	available uint32
}

// Get returns a fresh or recycled entity
func (p *implicitListEntityPool) Get() Entity {
	if p.available == 0 {
		e := newEntity(len(p.entities))
		p.entities = append(p.entities, e)
		return e
	}
	p.next, p.entities[p.next.id] = p.entities[p.next.id], p.next
	p.available--
	return p.entities[p.next.id]
}

// Recycle hands an entity back for recycling
func (p *implicitListEntityPool) Recycle(e Entity) {
	e.gen++
	p.next, p.entities[e.id] = e, p.next
	p.available++
}

// Alive return whether an entity is currently alive
func (p *implicitListEntityPool) Alive(e Entity) bool {
	return p.entities[e.id].id == e.id && e.gen == p.entities[e.id].gen
}
