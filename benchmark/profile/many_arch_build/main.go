package main

// Profiling:
// go build ./benchmark/profile/many_arch_build
// many_arch_build
// go tool pprof -http=":8000" -nodefraction=0.001 many_arch_build cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 many_arch_build mem.pprof

import (
	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
	"github.com/pkg/profile"
)

func main() {

	count := 1000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count)
	stop.Stop()

	stop = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("."))
	run(count)
	stop.Stop()
}

func run(rounds int) {
	for i := 0; i < rounds; i++ {
		world := ecs.NewWorld()
		ids := c.RegisterAll(&world)

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
			entity := world.NewEntity()
			world.Add(entity, add...)
		}
	}
}
