package benchmark

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
)

func BenchmarkGameEngineEcs(b *testing.B) {
	comps := []ecs.ComponentConfig{
		{ID: 0, Component: position{}},
		{ID: 1, Component: rotation{}},
	}
	world := ecs.NewWorld(comps...)

	for i := 0; i < 1000; i++ {
		entity := world.NewEntity()
		world.Assign(entity, 0, 1)
	}

	filter := world.NewFilter(0, 1)

	for i := 0; i < b.N; i++ {
		for _, e := range filter.Entities() {
			pos := (*position)(world.Component(e, 0))
			_ = pos
		}
	}
}
