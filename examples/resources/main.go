// Demonstrates the use of "Resources".
//
// Resources are similar to components, but they are not associated to an Entity.
// Rather, there can only be one resource object per type in the World.
package main

import "github.com/mlange-42/arche/ecs"

// TimeStep is a resource holding the model step
type TimeStep struct {
	Step int
}

func main() {
	// Create a World.
	world := ecs.NewWorld()

	// Add a resource to the world
	world.AddResource(&TimeStep{0})

	// Run the simulation
	run(&world)
}

func run(w *ecs.World) {
	for {
		// Get the the TimeStep resource from the world
		time := ecs.GetResource[TimeStep](w)

		// Use the resource
		time.Step++
		if time.Step >= 1000 {
			break
		}
	}
}
