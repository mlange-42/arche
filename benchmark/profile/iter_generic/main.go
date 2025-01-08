package main

// Profiling:
// go build ./benchmark/profile/iter_generic
// iter_generic
// go tool pprof -http=":8000" -nodefraction=0.001 iter_generic cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 iter_generic mem.pprof

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
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
		world := ecs.NewWorld(1024)

		mapper := generic.NewMap2[position, velocity](&world)

		for j := 0; j < entities; j++ {
			mapper.New()
		}

		filter := generic.NewFilter2[position, velocity]()

		for j := 0; j < iters; j++ {
			query := filter.Query(&world)
			for query.Next() {
				pos, vel := query.Get()
				pos.X += vel.X
				pos.Y += vel.Y
			}
		}
	}
}
