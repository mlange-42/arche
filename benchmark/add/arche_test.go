package add

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func addArcheWorld(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()

		posID := ecs.ComponentID[position](&world)
		rotID := ecs.ComponentID[rotation](&world)
		b.StartTimer()

		for i := 0; i < count; i++ {
			entity := world.NewEntity()
			world.Add(entity, posID, rotID)
		}
	}
}

func BenchmarkArcheAdd_100(b *testing.B) {
	addArcheWorld(b, 100)
}

func BenchmarkArcheAdd_1000(b *testing.B) {
	addArcheWorld(b, 1000)
}

func BenchmarkArcheAdd_10000(b *testing.B) {
	addArcheWorld(b, 10000)
}
