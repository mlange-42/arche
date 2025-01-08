// Demonstrates outer parallelism by running multiple simulations concurrently.
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// Run a total of 1000 simulations
const totalRuns = 1000

func main() {
	runSequential()
	runParallel()
}

// Run simulations sequentially, with a reused world.
func runSequential() {
	start := time.Now()

	// Create an empty world. It will be reused for all simulations.
	w := ecs.NewWorld()

	// Run many simulations.
	for i := 0; i < totalRuns; i++ {
		runModel(&w)
	}

	fmt.Printf("Sequential:   %s\n", time.Since(start))
}

// Run simulations in parallel, using workers.
func runParallel() {
	start := time.Now()

	// As many workers as processors available.
	workers := runtime.NumCPU()
	// Channel for sending jobs to workers (buffered!).
	jobs := make(chan int, totalRuns)
	// Channel for retrieving results / done messages (buffered!).
	results := make(chan int, totalRuns)

	// Start the workers.
	for w := 0; w < workers; w++ {
		go worker(jobs, results)
	}

	// Send the jobs. Does not block due to buffered channel.
	for j := 0; j < totalRuns; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect done messages.
	for j := 0; j < totalRuns; j++ {
		<-results
	}

	fmt.Printf("Parallel (%d): %s\n", workers, time.Since(start))
}

// Worker for running simulations on a reused world.
// Each worker needs its own world.
func worker(jobs <-chan int, results chan<- int) {
	// Create the worker's world. Will be reused for all jobs of the worker.
	w := ecs.NewWorld()

	// Process incoming jobs.
	for j := range jobs {
		// Run the model.
		runModel(&w)
		// Send done message. Does not block due to buffered channel.
		results <- j
	}
}

// A simulation that creates 1000 entities, and then removes them after 1000 steps, one each step.
// The argument is an ECS world that is reused for multiple simulations.
func runModel(w *ecs.World) {
	// Reset the world.
	w.Reset()

	// Create 10k entities, each with a countdown.
	mapper := generic.NewMap1[Countdown](w)
	query := mapper.NewBatchQ(1000)

	// Initialize the just created entities.
	cnt := 0
	for query.Next() {
		cd := query.Get()
		cd.Remaining = cnt + 1000
		cnt++
	}

	// List of entities to remove in each step.
	toRemove := []ecs.Entity{}

	// Filter for querying entities with Countdown.
	filter := generic.NewFilter1[Countdown]()

	// Run until there are no more entities.
	for {
		query = filter.Query(w)

		// If no more entities, break out of the loop.
		if query.Count() == 0 {
			query.Close()
			break
		}

		// Iterate entities.
		for query.Next() {
			cd := query.Get()
			// Countdown.
			cd.Remaining--
			// List entity for removal if countdown hits zero.
			if cd.Remaining <= 0 {
				toRemove = append(toRemove, query.Entity())
			}
		}

		// Remove all entities with countdown zero.
		for _, e := range toRemove {
			w.RemoveEntity(e)
		}

		// Clear list of entities to remove.
		toRemove = toRemove[:0]
	}
}

// Countdown component
type Countdown struct {
	Remaining int // Remaining ticks.
}
