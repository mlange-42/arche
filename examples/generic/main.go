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

	// Create entities
	for i := 0; i < 1000; i++ {
		// Create a new Entity.
		entity := world.NewEntity()
		// Add components to it.
		pos, vel := generic.Add2[Position, Velocity](&world, entity)

		// Initialize component fields.
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100

		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query.
		// Generic queries support up to 8 components.
		// For more components, use World.Query()
		query := generic.Query2[Position, Velocity](&world)
		// Iterate it
		for query.Next() {
			// Component access through a Query.
			pos, vel := query.GetAll()
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
