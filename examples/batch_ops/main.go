// Demonstrates batch-creation and batch-removal of entities.
//
// Batch operations are an optimization for creating and removing many entities in one go.
package main

import (
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

	// Run the simulation
	run(&world)
	// Run the simulation using generic access
	runGeneric(&world)
}

// Makes use of the resource by ID access.
func run(w *ecs.World) {
	// Get component IDs
	posID := ecs.ComponentID[Position](w)
	velID := ecs.ComponentID[Velocity](w)

	// Batch-create entities with components
	w.Batch().NewEntities(100, posID, velID)

	// Batch-create entities with components, and iterate them
	query := w.Batch().NewEntitiesQuery(100, posID, velID)
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		pos.X = 1.0
		pos.Y = 1.0
	}

	// Batch-remove all entities with exactly the given components
	filterExact := ecs.All(posID, velID).Exact()
	w.Batch().RemoveEntities(&filterExact)

	// Batch-remove all entities with the given components (and potentially further components)
	filter := ecs.All(posID, velID)
	w.Batch().RemoveEntities(&filter)
}

// Makes use of the resource by generic access.
func runGeneric(w *ecs.World) {
	// Get component mapper
	mapper := generic.NewMap2[Position, Velocity](w)

	// Batch-create entities using the mapper
	mapper.NewEntities(100)

	// Batch-create entities using the mapper, and iterate them
	query := mapper.NewEntitiesQuery(100)
	for query.Next() {
		pos, _ := query.Get()
		pos.X = 1.0
		pos.Y = 1.0
	}

	// Batch-remove all entities with exactly the given components
	mapper.RemoveEntities(true)

	// Batch-remove all entities with the given components (and potentially further components)
	mapper.RemoveEntities(true)
}
