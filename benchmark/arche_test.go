package benchmark

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func BenchmarkArche(b *testing.B) {
	world := ecs.NewWorld()

	comps := []ecs.Component{
		{ID: 0, Component: position{}},
		{ID: 1, Component: rotation{}},
	}

	arch := ecs.NewArchetype(comps...)

	for i := 0; i < 1000; i++ {
		arch.Add(
			world.NewEntity(),
			ecs.Component{ID: 0, Component: &position{1, 2}},
			ecs.Component{ID: 1, Component: &rotation{3}},
		)
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			pos := (*position)(arch.Get(i, ecs.ID(0)))
			_ = pos
		}
	}
}
