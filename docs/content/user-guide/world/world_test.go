package world

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func TestWorldSimple(t *testing.T) {
	world := ecs.NewWorld()
	_ = world
}

func TestWorldConfig(t *testing.T) {
	config := ecs.NewConfig().WithCapacityIncrement(1024)
	world := ecs.NewWorld(config)
	_ = world
}

func TestWorldConfigRelations(t *testing.T) {
	config := ecs.NewConfig().
		WithCapacityIncrement(1024).
		WithRelationCapacityIncrement(128)

	world := ecs.NewWorld(config)
	_ = world
}

func TestWorldReset(t *testing.T) {
	world := ecs.NewWorld()
	// ... do something with the world

	world.Reset()
	// ... start over again
}
