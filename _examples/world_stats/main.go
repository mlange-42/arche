// Demonstrates using world and archetype statistics.
package main

import (
	"fmt"
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
	world := ecs.NewWorld(4096)

	ids := []ecs.ID{
		ecs.ComponentID[Position](&world),
		ecs.ComponentID[Velocity](&world),
		ecs.ComponentID[Rotation](&world),
	}

	currIDs := make([]ecs.ID, 0, len(ids))
	// Create 1 million entities.
	for i := 0; i < 1_000_000; i++ {
		// Create a new Entity.
		entity := world.NewEntity()
		// Select some random components
		for _, id := range ids {
			if rand.Float64() < 0.7 {
				currIDs = append(currIDs, id)
			}
		}
		// Add the components
		world.Add(entity, currIDs...)
		// Clear the slice of IDs
		currIDs = currIDs[:0]
	}

	// Get and print world statistics
	stats := world.Stats()
	fmt.Println(stats.String())
}
