package posvel

import (
	"testing"

	"github.com/mlange-42/arche/benchmark/competition/pos_vel/geckecs"
)

func BenchmarkIterGECK(b *testing.B) {
	b.StopTimer()
	world := geckecs.NewWorld()

	posEntities := world.Entities(nPos)
	world.SetPositions(geckecs.Position{}, posEntities...)

	posVelEntities := world.Entities(nPosVel)
	world.SetPositions(geckecs.Position{}, posVelEntities...)
	world.SetVelocities(geckecs.Velocity{}, posVelEntities...)

	iter := world.PositionVelocitySet.NewIterator()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		iter.Reset()
		for iter.HasNext() {
			_, pos, vel := iter.Next()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkBuildGECK(b *testing.B) {
	var (
		p = geckecs.Position{}
		v = geckecs.Velocity{}
	)
	for i := 0; i < b.N; i++ {
		world := geckecs.NewWorld()

		posEntities := world.Entities(nPos)
		world.SetPositions(p, posEntities...)

		posVelEntities := world.Entities(nPosVel)
		world.SetPositions(p, posVelEntities...)
		world.SetVelocities(v, posVelEntities...)
	}
}
