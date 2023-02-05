package ecs

// World is the interface for the ECS world
type World interface {
	NewEntity() Entity
	RemEntity(entity Entity)
}

// NewWorld creates a new World
func NewWorld() World {
	return &world{
		entityPool: NewEntityPool(),
	}
}

type world struct {
	entityPool EntityPool
}

func (w *world) NewEntity() Entity {
	return w.entityPool.Get()
}

func (w *world) RemEntity(entity Entity) {
	w.entityPool.Recycle(entity)
}
