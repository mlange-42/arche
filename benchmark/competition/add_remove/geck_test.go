package addremove

import (
	"testing"

	"github.com/mlange-42/arche/benchmark/competition/add_remove/geckecs"
)

func BenchmarkIterGECK(b *testing.B) {
	b.StopTimer()
	world := geckecs.NewWorld()

	p := geckecs.Position{}
	v := geckecs.Velocity{}

	posEntities := world.Entities(nEntities)
	world.SetPositions(p, posEntities...)

	// Iterate once for more fairness
	world.SetVelocities(v, posEntities...)

	iter := world.PositionReadIter()
	for iter.HasNext() {
		entity := iter.NextEntity()
		entity.RemoveVelocity()
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		world.SetVelocities(v, posEntities...)
		world.RemoveVelocities(posEntities...)
	}
}

func BenchmarkBuildGECK(b *testing.B) {
	pos := geckecs.Position{}
	for i := 0; i < b.N; i++ {
		world := geckecs.NewWorld()
		entities := world.Entities(nEntities)
		world.SetPositions(pos, entities...)
	}
}
