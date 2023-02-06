package benchmark

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func runArche(b *testing.B, count int) {
	world := ecs.NewWorld()

	comps := []ecs.Component{
		{ID: 0, Component: position{}},
		{ID: 1, Component: rotation{}},
	}

	arch := ecs.NewArchetype(comps...)

	for i := 0; i < count; i++ {
		arch.Add(
			world.NewEntity(),
			ecs.Component{ID: 0, Component: &position{1, 2}},
			ecs.Component{ID: 1, Component: &rotation{3}},
		)
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < count; j++ {
			pos := (*position)(arch.Get(i, ecs.ID(0)))
			_ = pos
		}
	}
}

func BenchmarkArche100(b *testing.B) {
	runArche(b, 100)
}

func BenchmarkArche1000(b *testing.B) {
	runArche(b, 1000)
}

func BenchmarkArche10000(b *testing.B) {
	runArche(b, 10000)
}
