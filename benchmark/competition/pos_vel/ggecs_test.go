package posvel

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
)

const (
	PositionComponentID ecs.ComponentID = iota
	VelocityComponentID
)

func BenchmarkIterGGEcs(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(1024)
	world.Register(ecs.NewComponentRegistry[Position](PositionComponentID))
	world.Register(ecs.NewComponentRegistry[Velocity](VelocityComponentID))

	for i := 0; i < nPos; i++ {
		_ = world.NewEntity(PositionComponentID)
	}
	for i := 0; i < nPosVel; i++ {
		_ = world.NewEntity(PositionComponentID, VelocityComponentID)
	}

	mask := ecs.MakeComponentMask(PositionComponentID, VelocityComponentID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(mask)
		for query.Next() {
			pos := (*Position)(query.Component(PositionComponentID))
			vel := (*Velocity)(query.Component(VelocityComponentID))
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkBuildGGEcs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := ecs.NewWorld(1024)
		world.Register(ecs.NewComponentRegistry[Position](PositionComponentID))
		world.Register(ecs.NewComponentRegistry[Velocity](VelocityComponentID))

		for i := 0; i < nPos; i++ {
			world.NewEntity(PositionComponentID)
		}
		for i := 0; i < nPosVel; i++ {
			world.NewEntity(PositionComponentID, VelocityComponentID)
		}
	}
}
