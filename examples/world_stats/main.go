// Demonstrates the core API that uses component IDs for access.
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
	world :=
		ecs.NewConfig().
			WithCapacityIncrement(4096).
			Build()

	ids := []ecs.ID{
		ecs.ComponentID[Position](&world),
		ecs.ComponentID[Velocity](&world),
		ecs.ComponentID[Rotation](&world),
	}

	currIDs := make([]ecs.ID, 0, len(ids))
	// Create 1 million entities.
	for i := 0; i < 1000000; i++ {
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
		currIDs = currIDs[:0]
	}

	fmt.Println(world.Stats().String())
}
