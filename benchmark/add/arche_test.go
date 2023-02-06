package add

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func addArcheWorld(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()

		posID := ecs.RegisterComponent[position](&world)
		rotID := ecs.RegisterComponent[rotation](&world)
		b.StartTimer()

		for i := 0; i < count; i++ {
			entity := world.NewEntity()
			world.Add(entity, posID, rotID)
		}
	}
}

func BenchmarkAddArche_100(b *testing.B) {
	addArcheWorld(b, 100)
}

func BenchmarkAddArche_1000(b *testing.B) {
	addArcheWorld(b, 1000)
}

func BenchmarkAddArche_10000(b *testing.B) {
	addArcheWorld(b, 10000)
}
