package main

// Profiling:
// go build ./benchmark/profile/many_targets
// many_targets
// go tool pprof -http=":8000" -nodefraction=0.001 many_targets cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 many_targets mem.pprof

import (
	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
	"github.com/pkg/profile"
)

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
		world := ecs.NewWorld()
		posID := ecs.ComponentID[c.TestStruct0](&world)
		relID := ecs.ComponentID[c.ChildOf](&world)

		perArch := entities / 1000

		builder := ecs.NewBuilder(&world)
		targetQuery := builder.NewBatchQ(1000)
		targets := make([]ecs.Entity, 0, 1000)
		for targetQuery.Next() {
			targets = append(targets, targetQuery.Entity())
		}

		childBuilder := ecs.NewBuilder(&world, posID, relID).WithRelation(relID)
		for _, target := range targets {
			childBuilder.NewBatch(perArch, target)
		}

		filter := ecs.All(posID, relID)
		for j := 0; j < iters; j++ {
			query := world.Query(&filter)
			for query.Next() {
				pos := (*c.TestStruct0)(query.Get(posID))
				pos.Val++
			}
		}
	}
}
