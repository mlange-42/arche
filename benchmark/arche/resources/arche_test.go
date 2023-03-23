package iterate

import (
	"testing"

	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func BenchmarkArcheGetResource(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	resources := w.Resources()

	posID := ecs.ResourceID[c.Position](&w)
	resources.Add(posID, &c.Position{X: 1, Y: 2})

	b.StartTimer()

	var res *c.Position
	for i := 0; i < b.N; i++ {
		res = resources.Get(posID).(*c.Position)
	}

	_ = res
}

func BenchmarkArcheGetResourceGeneric(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	mapper := generic.NewResource[c.Position](&w)
	mapper.Add(&c.Position{X: 1, Y: 2})

	b.StartTimer()

	var res *c.Position
	for i := 0; i < b.N; i++ {
		res = mapper.Get()
	}

	_ = res
}

func BenchmarkArcheHasResource(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	resources := w.Resources()

	posID := ecs.ResourceID[c.Position](&w)
	resources.Add(posID, &c.Position{X: 1, Y: 2})

	b.StartTimer()

	var res bool
	for i := 0; i < b.N; i++ {
		res = resources.Has(posID)
	}

	_ = res
}

func BenchmarkArcheHasResourceGeneric(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	mapper := generic.NewResource[c.Position](&w)
	mapper.Add(&c.Position{X: 1, Y: 2})

	b.StartTimer()

	var res bool
	for i := 0; i < b.N; i++ {
		res = mapper.Has()
	}

	_ = res
}
