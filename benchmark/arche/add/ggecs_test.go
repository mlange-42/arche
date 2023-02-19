package add

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
	c "github.com/mlange-42/arche/benchmark/arche/common"
)

func addGameEngineEcs(b *testing.B, count int) {
	b.StopTimer()
	comps := []ecs.ComponentConfig{
		{ID: 0, Component: c.Position{}},
		{ID: 1, Component: c.Rotation{}},
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld(comps...)
		filter := world.NewFilter(1, 1)
		_ = filter
		b.StartTimer()

		for i := 0; i < count; i++ {
			entity := world.NewEntity()
			world.Assign(entity, 0, 1)
		}
	}
}

func BenchmarkGGEcsAdd_100(b *testing.B) {
	addGameEngineEcs(b, 100)
}

func BenchmarkGGEcsAdd_1000(b *testing.B) {
	addGameEngineEcs(b, 1000)
}

func BenchmarkGGEcsAdd_10000(b *testing.B) {
	addGameEngineEcs(b, 10000)
}
