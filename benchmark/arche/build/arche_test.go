package iterate

import (
	"testing"

	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
)

func BenchmarkBuild1kArch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()
		ids := c.RegisterAll(&world)
		b.StartTimer()

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
