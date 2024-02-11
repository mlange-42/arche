// The minimal example from the README using generic access.
//
// For automatic testing by the GitHub CI.
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

func main() {
	// Create a World.
	world := ecs.NewWorld()

	// Create a component mapper.
	mapper := generic.NewMap2[Position, Velocity](&world)

	// Create entities.
	for i := 0; i < 1000; i++ {
		// Create a new Entity with components.
		entity := mapper.New()
		// Get the components
		pos, vel := mapper.Get(entity)
		// Initialize component fields.
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100
		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Create a generic filter.
	filter := generic.NewFilter2[Position, Velocity]()

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query.
		query := filter.Query(&world)
		// Iterate it
		for query.Next() {
			// Component access through the Query.
			pos, vel := query.Get()
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
