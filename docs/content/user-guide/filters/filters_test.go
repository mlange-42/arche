package filters

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/filter"
	"github.com/mlange-42/arche/generic"
)

// Position component
type Position struct {
	X float64
	Y float64
}

// Velocity component
type Velocity struct {
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
	filter := generic.NewFilter2[Position, Heading]().
		Exclusive()
	_ = filter
}

func TestGenericOptional(t *testing.T) {
	world := ecs.NewWorld()

	filter := generic.NewFilter2[Position, Heading]().
		Optional(generic.T[Heading]())

	query := filter.Query(&world)
	for query.Next() {
		_, head := query.Get()
		if head == nil {
			// Optional component Heading not present
			fmt.Println("Heading not present in entity ", query.Entity())
		}
	}
}

func TestGenericWith(t *testing.T) {
	world := ecs.NewWorld()

	filter := generic.NewFilter1[Position]().
		With(generic.T[Heading]())

	query := filter.Query(&world)
	for query.Next() {
		pos := query.Get()
		_ = pos
	}
}

func TestLogicFilters(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)
	headID := ecs.ComponentID[Heading](&world)

	// Either Position and Velocity, or Position and Heading.
	_ = filter.OR{
		L: ecs.All(posID, velID),
		R: ecs.All(posID, headID),
	}

	// Same as above, expressed with a different logic.
	_ = filter.AND{
		L: ecs.All(posID),
		R: filter.Any(velID, headID),
	}

	// Missing any of Position or Velocity.
	_ = filter.AnyNot(posID, velID)
}

func TestRegister(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)

	mask := ecs.All(posID)
	filter := world.Cache().Register(mask)

	// Use the registered filter in queries!
	query := world.Query(&filter)
	query.Close()
}

func TestRegisterGeneric(t *testing.T) {
	world := ecs.NewWorld()

	filter := generic.NewFilter1[Position]()
	filter.Register(&world)

	query := filter.Query(&world)
	query.Close()
}
