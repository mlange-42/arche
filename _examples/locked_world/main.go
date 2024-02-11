// Demonstrates how to create, remove or alter entities,
// despite the World is locked during query iteration.
package main

import (
	"fmt"
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// Position component
type Position struct {
	X float64
	Y float64
}

// Energy component
type Energy struct {
	Value float64
}

// Reproduction is a helper for delaying reproduction events
type Reproduction struct {
	X         float64
	Y         float64
	Offspring int
}

func main() {

	// Create a World.
	world := ecs.NewWorld()

	// Create a component mapper.
	mapper := generic.NewMap2[Position, Energy](&world)

	// Create entities.
	for i := 0; i < 1000; i++ {
		// Create a new Entity with components.
		e := mapper.New()
		pos, en := mapper.Get(e)

		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100
		en.Value = rand.Float64()
	}

	// Filter for generating queries.
	filter := generic.NewFilter2[Position, Energy]()

	// Slices to collect events to perform after the world is unlocked.
	toDie := []ecs.Entity{}
	toReproduce := []Reproduction{}

	// Time loop.
	for t := 0; t < 100; t++ {
		// Get a fresh query iterator.
		query := filter.Query(&world)
		// We want to count our entities.
		alive := 0
		born := 0
		died := 0
		// Iterate it.
		// During iteration, the world is locked.
		// Thus, it is not possible to add or remove entities or components.
		for query.Next() {
			// Component access through the Query.
			entity := query.Entity()
			pos, ene := query.Get()
			// Do model logic.
			ene.Value += rand.NormFloat64() * 0.05

			// Entities with zero energy die.
			if ene.Value <= 0.0 {
				// Put the entity into a slice for later processing.
				toDie = append(toDie, entity)
				died++
				continue
			}

			alive++

			// Entities with energy > 0 reproduce.
			if ene.Value >= 1.0 {
				ene.Value = 0.5
				offspring := 1 + rand.Intn(3)
				born += offspring
				// Put the reproduction information into a slice for later processing.
				toReproduce = append(toReproduce, Reproduction{pos.X, pos.Y, offspring})
			}
		}
		// Here, the world lock is released automatically.

		// Remove entities that are to die.
		for _, entity := range toDie {
			world.RemoveEntity(entity)
		}

		// Create entities through reproduction.
		for _, repro := range toReproduce {
			// Create offspring entities.
			for i := 0; i < repro.Offspring; i++ {
				// Create a new Entity and initialize it.
				e := mapper.New()
				pos, en := mapper.Get(e)

				pos.X = repro.X + rand.NormFloat64()
				pos.Y = repro.Y + rand.NormFloat64()
				en.Value = rand.Float64() * 0.5
			}
		}

		// Clear the slices.
		toDie = toDie[:0]
		toReproduce = toReproduce[:0]

		fmt.Printf("Alive: %4d, Born: %2d, Died: %2d\n", alive, born, died)
	}
}
