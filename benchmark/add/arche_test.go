package add

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
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

func addArcheGeneric(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()
		b.StartTimer()

		for i := 0; i < count; i++ {
			entity := world.NewEntity()
			generic.Add[position](&world, entity)
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

func BenchmarkArcheAddGeneric_100(b *testing.B) {
	addArcheGeneric(b, 100)
}

func BenchmarkArcheAddGeneric_1000(b *testing.B) {
	addArcheGeneric(b, 1000)
}

func BenchmarkArcheAddGeneric_10000(b *testing.B) {
	addArcheGeneric(b, 10000)
}
