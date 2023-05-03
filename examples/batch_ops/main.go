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
	// Run the simulation.
	run()
	// Run the simulation using the generic API.
	runGeneric()
}

// Uses the standard API with ID access.
func run() {
	// Create a World.
	world := ecs.NewWorld()

	// Get component IDs.
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	// Create an entity builder with components.
	builder := ecs.NewBuilder(&world, posID, velID)

	// Batch-create entities.
	builder.NewBatch(100)

	// Batch-create entities, and iterate them.
	query := builder.NewQuery(100)
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		pos.X = 1.0
		pos.Y = 1.0
	}

	// Batch-remove all entities with exactly the given components.
	filterExcl := ecs.All(posID, velID).Exclusive()
	world.Batch().RemoveEntities(&filterExcl)

	// Batch-remove all entities with the given components (and potentially further components).
	filter := ecs.All(posID, velID)
	world.Batch().RemoveEntities(&filter)
}

// Uses the type-safe generic API.
func runGeneric() {
	// Create a World.
	world := ecs.NewWorld()

	// Get component mapper.
	mapper := generic.NewMap2[Position, Velocity](&world)

	// Batch-create entities using the mapper.
	mapper.NewBatch(100)

	// Batch-create entities using the mapper, and iterate them.
	query := mapper.NewQuery(100)
	for query.Next() {
		pos, _ := query.Get()
		pos.X = 1.0
		pos.Y = 1.0
	}

	// Batch-remove all entities with exactly the given components.
	mapper.RemoveEntities(true)

	// Batch-remove all entities with the given components (and potentially further components).
	mapper.RemoveEntities(false)
}
