// Demonstrates random sampling of a fixed number of entities from a query.
package main

import (
	"fmt"
	"math/rand"
	"sort"

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

	// Get a fresh query iterator.
	query := filter.Query(&world)
	// Get the number of entities in the query
	count := query.Count()
	// Get sorted random indices
	sample := sortedSample(25, count)
	fmt.Println(sample)

	// Iterate, only visiting the entities at the given indices.
	last := -1
	for _, idx := range sample {
		// Calculate the step size (can be 0)
		step := idx - last
		// Advance the query iterator
		query.Step(step)
		// Do something with the entity and/or components at the current iterator position.
		fmt.Println(query.Entity())
		// Required for calculating the step size.
		last = idx
	}
}

// Returns a sorted random sample of indices of size k, for a "thing" of length n.
func sortedSample(k, n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := rand.Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	m = m[:k]
	sort.Ints(m)
	return m
}
