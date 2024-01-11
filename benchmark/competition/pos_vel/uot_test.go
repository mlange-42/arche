package posvel

import (
	"testing"

	"github.com/unitoftime/ecs"
)

func BenchmarkIterUot(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld()

	for i := 0; i < nPos; i++ {
		id := world.NewId()
		ecs.Write(world, id,
			ecs.C(Position{0, 0}),
		)
	}
	for i := 0; i < nPosVel; i++ {
		id := world.NewId()
		ecs.Write(world, id,
			ecs.C(Position{0, 0}),
			ecs.C(Velocity{0, 0}),
		)
	}
	query := ecs.Query2[Position, Velocity](world)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query.MapId(func(id ecs.Id, pos *Position, vel *Velocity) {
			pos.X += vel.X
			pos.Y += vel.Y
		})
	}
}

func BenchmarkBuildUot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := ecs.NewWorld()

		for i := 0; i < nPos; i++ {
			id := world.NewId()
			ecs.Write(world, id,
				ecs.C(Position{0, 0}),
			)
		}
		for i := 0; i < nPosVel; i++ {
			id := world.NewId()
			ecs.Write(world, id,
				ecs.C(Position{0, 0}),
				ecs.C(Velocity{0, 0}),
			)
		}
	}
}
