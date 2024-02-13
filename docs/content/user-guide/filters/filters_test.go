package filters

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

func TestMask(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	mask := ecs.All(posID, headID)

	query := world.Query(mask)
	query.Close()
}

func TestMaskGeneric(t *testing.T) {
	world := ecs.NewWorld()

	filter := generic.NewFilter2[Position, Heading]()
	query := filter.Query(&world)
	query.Close()
}

func TestMaskWithout(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	mask := ecs.All(posID).Without(headID)
	_ = mask
}

func TestMaskWithoutGeneric(t *testing.T) {
	filter := generic.NewFilter1[Position]().
		Without(generic.T[Heading]())

	_ = filter
}

func TestMaskExclusive(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	mask := ecs.All(posID, headID).Exclusive()
	_ = mask
}

func TestMaskExclusiveGeneric(t *testing.T) {
	// TODO
	//filter := generic.NewFilter2[Position, Heading]().
	//	Exclusive()
	//_ = filter
}
