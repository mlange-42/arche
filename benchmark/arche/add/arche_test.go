package add

import (
	"testing"

	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func addArcheWorld(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

		posID := ecs.ComponentID[c.Position](&world)
		rotID := ecs.ComponentID[c.Rotation](&world)
		comps := []ecs.ID{posID, rotID}
		b.StartTimer()

		var e ecs.Entity
		for i := 0; i < count; i++ {
			e = world.NewEntity(comps...)
		}
		_ = e
	}
}

func addArcheGeneric(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
		mut := generic.NewMap2[c.Position, c.Rotation](&world)
		b.StartTimer()

		var e ecs.Entity
		for i := 0; i < count; i++ {
			e = mut.NewEntity()
		}
		_ = e
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
