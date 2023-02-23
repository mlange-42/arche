package add

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
	c "github.com/mlange-42/arche/benchmark/arche/common"
)

const (
	PositionComponentID ecs.ComponentID = iota
	RotationComponentID
)

func addGameEngineEcs(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld(1024)
		world.Register(ecs.NewComponentRegistry[c.Position](PositionComponentID))
		world.Register(ecs.NewComponentRegistry[c.Rotation](RotationComponentID))
		b.StartTimer()

		for i := 0; i < count; i++ {
			_ = world.NewEntity(PositionComponentID, RotationComponentID)
		}
	}
}

func BenchmarkGGEcsAdd_100(b *testing.B) {
	addGameEngineEcs(b, 100)
}

func BenchmarkGGEcsAdd_1000(b *testing.B) {
	addGameEngineEcs(b, 1000)
}

func BenchmarkGGEcsAdd_10000(b *testing.B) {
	addGameEngineEcs(b, 10000)
}
