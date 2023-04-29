// Demonstrates generic access to arbitrary/random entities.
package main

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// Linked component
type Linked struct {
	Prev ecs.Entity
	Next ecs.Entity
}

func main() {
	// Create a World.
	world := ecs.NewWorld()

	// Create a component mapper.
	mapper := generic.NewMap1[Linked](&world)

	var prevEntity ecs.Entity
	// Create entities.
	for i := 0; i < 1000; i++ {
		// Create a new Entity with components.
		e := mapper.New()
		link := mapper.Get(e)
		// Make it an implicit linked list.
		link.Prev = prevEntity
	}

	// Create a generic filter.
	filter := generic.NewFilter1[Linked]()

	// Get a fresh query iterator.
	query := filter.Query(&world)
	for query.Next() {
		// Get entity and link component.
		entity := query.Entity()
		link := query.Get()
		if link.Prev.IsZero() {
			continue
		}

		// Get a component from another entity than the one of the current iteration.
		prevLink := mapper.Get(link.Prev)
		// Make it a double-linked list.
		prevLink.Next = entity
	}
}
