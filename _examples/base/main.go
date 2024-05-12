// Demonstrates the core API that uses component IDs for access.
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

// Rotation component
type Rotation struct {
	A float64
}

func main() {
	// Create a World.
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)
	rotID := ecs.ComponentID[Rotation](&world)

	// Create entities.
	for i := 0; i < 1000; i++ {
		// Create a new Entity with components.
		entity := world.NewEntity(posID, velID)
		// Get the components.
		pos := (*Position)(world.Get(entity, posID))
		vel := (*Velocity)(world.Get(entity, velID))

		// Initialize component fields.
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100

		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Create a filter, demanding and excluding components.
	filter := ecs.All(posID, velID).Without(rotID)

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query iterator.
		query := world.Query(&filter)
		// Iterate it.
		for query.Next() {
			// Component access through the Query.
			pos := (*Position)(query.Get(posID))
			vel := (*Velocity)(query.Get(velID))
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
