package posvel

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func BenchmarkIterArche(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	world.Batch().NewEntities(nPos, posID)
	world.Batch().NewEntities(nPosVel, posID, velID)

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

func BenchmarkIterArcheGeneric(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	posMapper := generic.NewMap1[Position](&world)
	posVelMapper := generic.NewMap2[Position, Velocity](&world)

	for i := 0; i < nPos; i++ {
		posMapper.NewEntity()
	}
	for i := 0; i < nPosVel; i++ {
		posVelMapper.NewEntity()
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

func BenchmarkBuildArcheGeneric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

		posMapper := generic.NewMap1[Position](&world)
		posVelMapper := generic.NewMap2[Position, Velocity](&world)

		for i := 0; i < nPos; i++ {
			posMapper.NewEntity()
		}
		for i := 0; i < nPosVel; i++ {
			posVelMapper.NewEntity()
		}
	}
}
