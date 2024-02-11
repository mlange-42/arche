// Demonstrates the use of "Resources".
//
// Resources are similar to components, but they are not associated to an Entity.
// Rather, there can only be one resource object per type in the World.
package main

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// TimeStep is a resource holding the model step
type TimeStep struct {
	Step int
}

func main() {
	// Create a World.
	world := ecs.NewWorld()

	// Add a resource to the world
	ecs.AddResource(&world, &TimeStep{0})

	// Run the simulation
	run(&world)
	// Run the simulation using generic access
	runGeneric(&world)
}

// Makes use of the resource by ID access.
func run(w *ecs.World) {
	timeStepID := ecs.ResourceID[TimeStep](w)

	for {
		// Get the the TimeStep resource from the world
		time := w.Resources().Get(timeStepID).(*TimeStep)

		// Use the resource
		time.Step++
		fmt.Println(time.Step)
		if time.Step >= 50 {
			break
		}
	}
}

// Makes use of the resource by generic access.
func runGeneric(w *ecs.World) {
	mapper := generic.NewResource[TimeStep](w)
	for {
		// Get the the TimeStep resource from the world
		time := mapper.Get()

		// Use the resource
		time.Step++
		fmt.Println(time.Step)
		if time.Step >= 100 {
			break
		}
	}
}
