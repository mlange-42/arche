package main

// Profiling:
// go build ./profile/main/profile.go
// profile
// go tool pprof -http=":8000" -nodefraction=0.001 profile cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 profile mem.pprof

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

	count := 100
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
		world := ecs.NewConfig().
			WithCapacityIncrement(1024).
			Build()

		posID := ecs.ComponentID[position](&world)
		rotID := ecs.ComponentID[rotation](&world)

		for j := 0; j < entities; j++ {
			entity := world.NewEntity()
			world.Add(entity, posID, rotID)
		}

		for j := 0; j < iters; j++ {
			query := world.Query(posID, rotID)
			for query.Next() {
				pos := (*position)(query.Get(posID))
				_ = pos
			}
		}
	}
}
