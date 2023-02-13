package main

import (
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	fi "github.com/mlange-42/arche/filter"
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

// Elevation component
type Elevation struct {
	E float64
}

func main() {
	// Create a World.
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)
	rotID := ecs.ComponentID[Rotation](&world)
	eleID := ecs.ComponentID[Elevation](&world)

	// Create entities
	for i := 0; i < 1000; i++ {
		// Create a new Entity.
		entity := world.NewEntity()
		// Add components to it.
		comps := []ecs.ID{posID, velID}
		if i%2 == 0 {
			comps = append(comps, rotID)
		}
		if i%3 == 0 {
			comps = append(comps, eleID)
		}
		world.Add(entity, comps...)

		pos := (*Position)(world.Get(entity, posID))
		vel := (*Position)(world.Get(entity, velID))

		// Initialize component fields.
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100

		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Create a logic filter
	filter := fi.And(
		fi.All(posID, velID),
		fi.Or(
			fi.All(rotID),
			fi.NoneOf(eleID),
		),
	)

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query iterator.
		query := world.Query(filter)
		// Iterate it.
		for query.Next() {
			// Component access through a Query.
			pos := (*Position)(query.Get(posID))
			vel := (*Velocity)(query.Get(velID))
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
