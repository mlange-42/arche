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

	// Create a filter
	filter := ecs.All(posID, velID).Without(rotID)

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query iterator.
		q := world.Query(filter)
		// Iterate it
		for q.Next() {
			// Component access through a Query.
			pos := (*Position)(q.Get(posID))
			vel := (*Position)(q.Get(velID))
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
