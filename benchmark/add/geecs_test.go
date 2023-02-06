package add

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
)

func addGameEngineEcs(b *testing.B, count int) {
	comps := []ecs.ComponentConfig{
		{ID: 0, Component: position{}},
		{ID: 1, Component: rotation{}},
	}

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

func BenchmarkAddGEEcs100(b *testing.B) {
	addGameEngineEcs(b, 100)
}

func BenchmarkAddGEEcs1000(b *testing.B) {
	addGameEngineEcs(b, 1000)
}

func BenchmarkAddGEEcs10000(b *testing.B) {
	addGameEngineEcs(b, 10000)
}
