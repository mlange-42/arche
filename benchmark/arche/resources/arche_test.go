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
	w.AddResource(&c.Position{X: 1, Y: 2})
	posID := ecs.ResourceID[c.Position](&w)

	b.StartTimer()

	var res *c.Position
	for i := 0; i < b.N; i++ {
		res = w.GetResource(posID).(*c.Position)
	}

	_ = res
}

func BenchmarkArcheGetResourceGeneric(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	w.AddResource(&c.Position{X: 1, Y: 2})
	mapper := generic.NewResource[c.Position](&w)

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
	w.AddResource(&c.Position{X: 1, Y: 2})
	posID := ecs.ResourceID[c.Position](&w)

	b.StartTimer()

	var res bool
	for i := 0; i < b.N; i++ {
		res = w.HasResource(posID)
	}

	_ = res
}

func BenchmarkArcheHasResourceGeneric(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	w.AddResource(&c.Position{X: 1, Y: 2})
	mapper := generic.NewResource[c.Position](&w)

	b.StartTimer()

	var res bool
	for i := 0; i < b.N; i++ {
		res = mapper.Has()
	}

	_ = res
}
