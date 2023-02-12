package main

import (
	"math/rand"

	"github.com/mlange-42/arche/ecs"
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

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	// Create entities
	for i := 0; i < 1000; i++ {
		// Create a new Entity.
		entity := world.NewEntity()
		// Add components to it.
		world.Add(entity, posID, velID)
		pos := (*Position)(world.Get(entity, posID))
		vel := (*Position)(world.Get(entity, velID))

		// Initialize component fields.
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100

		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query iterator.
		query := world.Query(ecs.All(posID, velID))
		// Iterate it
		for query.Next() {
			// Component access through a Query.
			pos := (*Position)(query.Get(posID))
			vel := (*Position)(query.Get(velID))
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
