package main

// Profiling:
// go build ./profile/main/profile.go
// profile
// go tool pprof -http=":8000" profile cpu.pprof
// go tool pprof -http=":8000" profile mem.pprof

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
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	count := 1000
	iters := 1000
	rounds := 100

	for i := 0; i < rounds; i++ {
		world := ecs.NewWorld()

		posID := ecs.RegisterComponent[position](&world)
		rotID := ecs.RegisterComponent[rotation](&world)

		for j := 0; j < count; j++ {
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
