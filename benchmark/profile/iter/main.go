package main

// Profiling:
// go build ./benchmark/profile/iter
// iter
// go tool pprof -http=":8000" -nodefraction=0.001 iter cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 iter mem.pprof

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

func run(rounds, iters, entities int) {
	for i := 0; i < rounds; i++ {
		world := ecs.NewWorld(
			ecs.NewConfig().WithCapacityIncrement(1024),
		)

		posID := ecs.ComponentID[position](&world)
		rotID := ecs.ComponentID[rotation](&world)

		for j := 0; j < entities; j++ {
			_ = world.NewEntity(posID, rotID)
		}

		for j := 0; j < iters; j++ {
			query := world.Query(ecs.All(posID, rotID))
			for query.Next() {
				pos := (*position)(query.Get(posID))
				pos.X = 1
			}
		}
	}
}
