package posvel

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func BenchmarkIterArche(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	posMap := generic.NewMap1[Position](&world)
	posVelMap := generic.NewMap2[Position, Velocity](&world)

	for i := 0; i < nPos; i++ {
		posMap.NewEntity()
	}
	for i := 0; i < nPosVel; i++ {
		posVelMap.NewEntity()
	}

	filter := generic.NewFilter2[Position, Velocity]()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := filter.Query(&world)
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkBuildArche(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

		posMap := generic.NewMap1[Position](&world)
		posVelMap := generic.NewMap2[Position, Velocity](&world)

		for i := 0; i < nPos; i++ {
			posMap.NewEntity()
		}
		for i := 0; i < nPosVel; i++ {
			posVelMap.NewEntity()
		}
	}
}
