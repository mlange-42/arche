package main

// Profiling:
// go build ./benchmark/profile/world
// world
// go tool pprof -http=":8000" -nodefraction=0.001 world cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 world mem.pprof

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

	count := 250
	iters := 10000
	entities := 1000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()

	stop = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()
}

func run(rounds, iters, entityCount int) {
	for i := 0; i < rounds; i++ {
		world := ecs.NewWorld(
			ecs.NewConfig().WithCapacityIncrement(1024),
		)

		posID := ecs.ComponentID[position](&world)
		rotID := ecs.ComponentID[rotation](&world)

		entities := make([]ecs.Entity, entityCount)
		for i := 0; i < entityCount; i++ {
			entities[i] = world.NewEntity(posID, rotID)
		}

		for j := 0; j < iters; j++ {
			for _, e := range entities {
				pos := (*position)(world.Get(e, posID))
				pos.X = 1
			}
		}
	}
}
