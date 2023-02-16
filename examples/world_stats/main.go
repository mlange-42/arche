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
	world := ecs.NewWorld()

	ids := []ecs.ID{
		ecs.ComponentID[Position](&world),
		ecs.ComponentID[Velocity](&world),
		ecs.ComponentID[Rotation](&world),
	}

	// Create entities.
	for i := 0; i < 100000; i++ {
		// Create a new Entity.
		entity := world.NewEntity()
		// Add some random components
		for _, id := range ids {
			if rand.Float64() < 0.7 {
				world.Add(entity, id)
			}
		}
	}

	fmt.Println(world.Stats().String())
}
