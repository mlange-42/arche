package main

// Profiling:
// go build ./benchmark/profile/many_arch
// many_arch
// go tool pprof -http=":8000" -nodefraction=0.001 many_arch cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 many_arch mem.pprof

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
		ids := c.RegisterAll(&world)

		perArch := 2 * entities / 1000

		for i := 0; i < 1024; i++ {
			mask := i
			add := make([]ecs.ID, 0, 10)
			for j := 0; j < 10; j++ {
				id := ids[j]
				m := 1 << j
				if mask&m == m {
					add = append(add, id)
				}
			}
			for j := 0; j < perArch; j++ {
				entity := world.NewEntity()
				world.Add(entity, add...)
			}
		}

		filter := ecs.All(ids[6])
		for j := 0; j < iters; j++ {
			query := world.Query(&filter)
			for query.Next() {
				pos := (*c.TestStruct6)(query.Get(ids[6]))
				pos.Val++
			}
		}
	}
}
