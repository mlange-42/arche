package apis

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

type Position struct{}

func TestIDs(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	filter := ecs.All(posID)

	query := world.Query(filter)
	for query.Next() {
		// do something
	}
}

func TestGeneric(t *testing.T) {
	world := ecs.NewWorld()

	filter := generic.NewFilter1[Position]()

	query := filter.Query(&world)
	for query.Next() {
		// do something
	}
}
