package worldaccess

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// Position component
type Position struct {
	X float64
	Y float64
}

// Heading component
type Heading struct {
	Angle float64
}

func TestGet(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	entity := world.NewEntity(posID, headID)

	pos := (*Position)(world.Get(entity, posID))
	head := (*Heading)(world.Get(entity, headID))

	_, _ = pos, head
}

func TestGetGeneric(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap2[Position, Heading](&world)
	entity := builder.New()

	pos, head := builder.Get(entity)

	_, _ = pos, head
}

func TestHas(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	entity := world.NewEntity(posID, headID)

	hasPos := world.Has(entity, posID)
	_ = hasPos
}

func TestHasGeneric(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap2[Position, Heading](&world)
	entity := builder.New()

	mapper := generic.NewMap[Position](&world)

	hasPos := mapper.Has(entity)
	_ = hasPos
}

func TestGetUnchecked(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	entity := world.NewEntity(posID, headID)

	pos := (*Position)(world.Get(entity, posID))
	head := (*Heading)(world.GetUnchecked(entity, headID))

	_, _ = pos, head
}
