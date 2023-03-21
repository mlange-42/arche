package posvel

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func BenchmarkIterArche(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	for i := 0; i < nPos; i++ {
		world.NewEntity(posID)
	}
	for i := 0; i < nPosVel; i++ {
		world.NewEntity(posID, velID)
	}

	filter := ecs.All(posID, velID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		for query.Next() {
			pos := (*Position)(query.Get(posID))
			vel := (*Velocity)(query.Get(velID))
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkBuildArche(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

		posID := ecs.ComponentID[Position](&world)
		velID := ecs.ComponentID[Velocity](&world)

		for i := 0; i < nPos; i++ {
			world.NewEntity(posID)
		}
		for i := 0; i < nPosVel; i++ {
			world.NewEntity(posID, velID)
		}
	}
}
