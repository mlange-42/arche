package main

// Profiling:
// go build ./benchmark/profile/exchange_batch
// ./exchange_batch
// go tool pprof -http=":8000" -nodefraction=0.001 ./exchange_batch cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 ./exchange_batch mem.pprof

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/pkg/profile"
)

type position struct {
	X int
	Y int
}

type velocity struct {
	X int
	Y int
}

type rotation struct {
	Angle int
}

func main() {

	count := 20
	iters := 100000
	entities := 32

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()

	stop = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()
}

func run(rounds, iters, numEntities int) {
	for i := 0; i < rounds; i++ {
		world := ecs.NewWorld(1024)

		posID := ecs.ComponentID[position](&world)
		velID := ecs.ComponentID[velocity](&world)
		rotID := ecs.ComponentID[rotation](&world)

		for j := 0; j < numEntities; j++ {
			_ = world.NewEntity(posID, velID)
		}
		add := []ecs.ID{rotID}
		remove := []ecs.ID{velID}
		filter := ecs.All(posID)

		for j := 0; j < iters; j++ {
			world.Batch().Exchange(&filter, add, remove)
			add, remove = remove, add
		}
	}
}
