package apis

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

type Position struct {
	X float64
	Y float64
}
type Velocity struct {
	X float64
	Y float64
}

func TestIDs(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)
	filter := ecs.All(posID, velID)

	query := world.Query(filter)
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		vel := (*Velocity)(query.Get(posID))
		pos.X += vel.X
		pos.Y += vel.Y
	}
}

func TestGeneric(t *testing.T) {
	world := ecs.NewWorld()

	filter := generic.NewFilter2[Position, Velocity]()

	query := filter.Query(&world)
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}
