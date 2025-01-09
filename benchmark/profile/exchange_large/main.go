package main

// Profiling:
// go build ./benchmark/profile/exchange_large
// ./exchange_large
// go tool pprof -http=":8000" -nodefraction=0.001 ./exchange_large cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 ./exchange_large mem.pprof

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

type c1 struct {
	X int
	Y int
}

type c2 struct {
	X int
	Y int
}

type c3 struct {
	X int
	Y int
}

type c4 struct {
	X int
	Y int
}

type c5 struct {
	X int
	Y int
}

type c6 struct {
	X int
	Y int
}

type c7 struct {
	X int
	Y int
}

type c8 struct {
	X int
	Y int
}

type c9 struct {
	X int
	Y int
}

type c10 struct {
	X int
	Y int
}

func main() {

	count := 250
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

		allIDs := []ecs.ID{
			posID, velID,
			ecs.ComponentID[c1](&world),
			ecs.ComponentID[c2](&world),
			ecs.ComponentID[c3](&world),
			ecs.ComponentID[c4](&world),
			ecs.ComponentID[c5](&world),
			ecs.ComponentID[c6](&world),
			ecs.ComponentID[c7](&world),
			ecs.ComponentID[c8](&world),
			ecs.ComponentID[c9](&world),
			ecs.ComponentID[c10](&world),
		}

		world.Batch().New(numEntities, allIDs...)

		filter := ecs.All(posID)

		for j := 0; j < iters; j++ {
			add := []ecs.ID{rotID}
			remove := []ecs.ID{velID}
			if j%2 != 0 {
				add, remove = remove, add
			}
			world.Batch().Exchange(&filter, add, remove)
		}
	}
}
