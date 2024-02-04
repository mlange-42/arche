// Demonstrates random sampling of a fixed number of entities from a query using Query.EntityAt().
package main

import (
	"fmt"
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func main() {
	// Create a World.
	world := ecs.NewWorld()

	// Create entities.
	for i := 0; i < 1000; i++ {
		world.NewEntity()
	}

	// Create a generic filter.
	filter := generic.NewFilter0()
	// Register the filter. Optional, but recommended for query random access.
	filter.Register(&world)

	// Get a fresh query iterator.
	query := filter.Query(&world)
	// Get the number of entities in the query
	count := query.Count()
	// Get random indices without replacement
	sample := sample(25, count)
	fmt.Println(sample)

	// Iterate over sampled indices.
	for _, idx := range sample {
		// Get the entity at the random index.
		entity := query.EntityAt(idx)
		// Do something with the entity and/or components at the current iterator position.
		fmt.Println(entity)
	}
}

// Returns a random sample of indices without replacement.Take k of n sample.
func sample(k, n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := rand.Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	m = m[:k]
	return m
}
