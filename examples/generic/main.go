// Demonstrates the generic API, which provides type-safety and convenience over the ID-based core API.
package main

import (
	"math/rand"

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

// Rotation component
type Rotation struct {
	A float64
}

// Elevation component
type Elevation struct {
	E float64
}

func main() {
	// Create a World.
	world := ecs.NewWorld()

	// Create entities.
	for i := 0; i < 1000; i++ {
		// Create a new Entity.
		entity := world.NewEntity()
		// Add components to it.
		pos, vel, _, _ := generic.Add4[Position, Velocity, Rotation, Elevation](&world, entity)

		// Initialize component fields.
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100

		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Generic queries support up to 8 components.
	// For more components, use World.Query() directly.
	filter := generic.NewFilter2[Position, Velocity]()

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query iterator.
		query := filter.Query(&world)
		// Iterate it.
		for query.Next() {
			// Component access through the Query.
			_, pos, vel := query.GetAll()
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}

	// A more complex generic query using optional and excluded components:
	filter =
		generic.
			NewFilter2[Position, Velocity]().    // Components provided through Get... methods
			Optional(generic.Mask1[Velocity]()). // but those may be nil
			With(generic.Mask1[Elevation]()).    // additional required components
			Without(generic.Mask1[Rotation]())   // entities with any of these are excluded.

	q := filter.Query(&world)

	for q.Next() {
		_, pos, vel := q.GetAll()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}
