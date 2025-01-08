package main

// Profiling:
// go build ./benchmark/profile/create_many
// ./create_many
// go tool pprof -http=":8000" -nodefraction=0.001 ./create_many cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 ./create_many mem.pprof

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/pkg/profile"
)

type position struct {
	X int
	Y int
}

type rotation struct {
	Angle int
}

func main() {

	iters := 100
	entities := 1_000_000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(iters, entities)
	stop.Stop()

	stop = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("."))
	run(iters, entities)
	stop.Stop()
}

func run(iters, entityCount int) {
	for i := 0; i < iters; i++ {
		world := ecs.NewWorld(1024)

		posID := ecs.ComponentID[position](&world)
		rotID := ecs.ComponentID[rotation](&world)

		for range entityCount {
			world.NewEntity(posID, rotID)
		}
	}
}
