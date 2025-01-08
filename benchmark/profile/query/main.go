package main

// Profiling:
// go build ./benchmark/profile/query
// query
// go tool pprof -http=":8000" -nodefraction=0.001 query cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 query mem.pprof

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/pkg/profile"
)

type position struct {
	X float64
	Y float64
}

type velocity struct {
	X float64
	Y float64
}

func main() {

	count := 250
	iters := 1000000
	entities := 1000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()

	stop = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()
}

func run(rounds, iters, entities int) {
	for i := 0; i < rounds; i++ {
		world := ecs.NewWorld(1024)

		posID := ecs.ComponentID[position](&world)
		velID := ecs.ComponentID[velocity](&world)

		for j := 0; j < entities; j++ {
			_ = world.NewEntity(posID, velID)
		}

		filter := ecs.All(posID, velID)
		var query ecs.Query

		for j := 0; j < iters; j++ {
			query = world.Query(&filter)
			query.Close()
		}
	}
}
