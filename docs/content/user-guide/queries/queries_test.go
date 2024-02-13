package userguide

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/mlange-42/arche/ecs"
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

func TestQueryIterate(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	filter := ecs.All(posID, velID)
	query := world.Query(filter)
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		vel := (*Velocity)(query.Get(velID))
		pos.X += vel.X
		pos.Y += vel.Y
	}
}

func TestQueryIterateGeneric(t *testing.T) {
	world := ecs.NewWorld()

	filter := generic.NewFilter2[Position, Velocity]()
	query := filter.Query(&world)
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}

func TestQueryRemoveEntities(t *testing.T) {
	world := ecs.NewWorld()

	// Create some entities.
	for i := 0; i < 100; i++ {
		world.NewEntity()
	}

	// A slice that we (re)-use to defer entity removal.
	var toRemove []ecs.Entity

	// A time loop.
	for time := 0; time < 100; time++ {
		// Query... the world gets locked.
		query := world.Query(ecs.All())
		// Iterate, and collect entities to remove.
		for query.Next() {
			if rand.Float64() < 0.1 {
				toRemove = append(toRemove, query.Entity())
			}
		}
		// The world is unlocked again.
		// Actually remove the collected entities.
		for _, e := range toRemove {
			world.RemoveEntity(e)
		}
		// Empty the slice, so we can reuse it in the next time step.
		toRemove = toRemove[:0]
	}
}

func TestQueryCount(t *testing.T) {
	world := ecs.NewWorld()

	query := world.Query(ecs.All())

	cnt := query.Count()
	fmt.Println(cnt)

	query.Close()
}

func TestQueryEntityAt(t *testing.T) {
	world := ecs.NewWorld()

	// Create some entities.
	for i := 0; i < 100; i++ {
		world.NewEntity()
	}

	// Query and count entities.
	query := world.Query(ecs.All())
	cnt := query.Count()

	// Draw random entities.
	for i := 0; i < 10; i++ {
		entity := query.EntityAt(rand.Intn(cnt))
		fmt.Println(entity)
	}

	query.Close()
}
