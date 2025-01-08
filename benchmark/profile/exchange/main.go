package main

// Profiling:
// go build ./benchmark/profile/exchange
// exchange
// go tool pprof -http=":8000" -nodefraction=0.001 exchange cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 exchange mem.pprof

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

	count := 100
	iters := 1000
	entities := 1000

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

		entities := make([]ecs.Entity, numEntities)
		for j := 0; j < numEntities; j++ {
			entity := world.NewEntity(posID, velID)
			entities[j] = entity
		}

		for j := 0; j < iters; j++ {
			add := []ecs.ID{rotID}
			remove := []ecs.ID{velID}
			if j%2 != 0 {
				add, remove = remove, add
			}
			for _, entity := range entities {
				world.Exchange(entity, add, remove)
			}
		}
	}
}
