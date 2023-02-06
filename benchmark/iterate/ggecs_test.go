package iterate

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
)

func runGameEngineEcs(b *testing.B, count int) {
	b.StopTimer()
	comps := []ecs.ComponentConfig{
		{ID: 0, Component: position{}},
		{ID: 1, Component: rotation{}},
	}
	world := ecs.NewWorld(comps...)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Assign(entity, 0, 1)
	}
	b.StartTimer()

	filter := world.NewFilter(0, 1)

	for i := 0; i < b.N; i++ {
		for _, e := range filter.Entities() {
			pos := (*position)(world.Component(e, 0))
			_ = pos
		}
	}
}

func BenchmarkGGEcsIter_1_000(b *testing.B) {
	runGameEngineEcs(b, 1000)
}

func BenchmarkGGEcsIter_10_000(b *testing.B) {
	runGameEngineEcs(b, 10000)
}

func BenchmarkGGEcsIter_100_000(b *testing.B) {
	runGameEngineEcs(b, 100000)
}
