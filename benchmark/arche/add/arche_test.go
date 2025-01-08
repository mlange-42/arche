package add

import (
	"testing"

	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func addArche(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()

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

func addArcheBatch(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()

		posID := ecs.ComponentID[c.Position](&world)
		rotID := ecs.ComponentID[c.Rotation](&world)
		comps := []ecs.ID{posID, rotID}
		b.StartTimer()

		ecs.NewBuilder(&world, comps...).NewBatch(count)
	}
}

func addSetArche(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()

		posID := ecs.ComponentID[c.Position](&world)
		rotID := ecs.ComponentID[c.Rotation](&world)
		comps := []ecs.ID{posID, rotID}
		b.StartTimer()

		var e ecs.Entity
		for i := 0; i < count; i++ {
			e = world.NewEntity(comps...)
			pos := (*c.Position)(world.Get(e, posID))
			rot := (*c.Rotation)(world.Get(e, rotID))
			pos.X = 1
			rot.Angle = 1
		}
		_ = e
	}
}

func addSetArcheBatch(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()

		posID := ecs.ComponentID[c.Position](&world)
		rotID := ecs.ComponentID[c.Rotation](&world)
		comps := []ecs.ID{posID, rotID}
		b.StartTimer()

		query := ecs.NewBuilder(&world, comps...).NewBatchQ(count)
		for query.Next() {
			pos := (*c.Position)(query.Get(posID))
			rot := (*c.Rotation)(query.Get(rotID))
			pos.X = 1
			rot.Angle = 1
		}
	}
}

func addArcheGeneric(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()
		mut := generic.NewMap2[c.Position, c.Rotation](&world)
		b.StartTimer()

		var e ecs.Entity
		for i := 0; i < count; i++ {
			e = mut.New()
		}
		_ = e
	}
}

func addArcheGenericBatch(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()
		mut := generic.NewMap2[c.Position, c.Rotation](&world)
		b.StartTimer()

		mut.NewBatch(count)
	}
}

func addSetArcheGeneric(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()
		mut := generic.NewMap2[c.Position, c.Rotation](&world)
		b.StartTimer()

		var e ecs.Entity
		for i := 0; i < count; i++ {
			e = mut.New()
			pos, rot := mut.Get(e)
			pos.X = 1
			rot.Angle = 1
		}
		_ = e
	}
}

func addSetArcheGenericBatch(b *testing.B, count int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld()
		mut := generic.NewMap2[c.Position, c.Rotation](&world)
		b.StartTimer()

		query := mut.NewBatchQ(count)

		for query.Next() {
			pos, rot := query.Get()
			pos.X = 1
			rot.Angle = 1
		}
	}
}

func BenchmarkArcheAdd_100(b *testing.B) {
	addArche(b, 100)
}

func BenchmarkArcheAdd_1000(b *testing.B) {
	addArche(b, 1000)
}

func BenchmarkArcheAdd_10000(b *testing.B) {
	addArche(b, 10000)
}

func BenchmarkArcheAdd_Batch_100(b *testing.B) {
	addArcheBatch(b, 100)
}

func BenchmarkArcheAdd_Batch_1000(b *testing.B) {
	addArcheBatch(b, 1000)
}

func BenchmarkArcheAdd_Batch_10000(b *testing.B) {
	addArcheBatch(b, 10000)
}

func BenchmarkArcheAddSet_100(b *testing.B) {
	addSetArche(b, 100)
}

func BenchmarkArcheAddSet_1000(b *testing.B) {
	addSetArche(b, 1000)
}

func BenchmarkArcheAddSet_10000(b *testing.B) {
	addSetArche(b, 10000)
}

func BenchmarkArcheAddSet_Batch_100(b *testing.B) {
	addSetArcheBatch(b, 100)
}

func BenchmarkArcheAddSet_Batch_1000(b *testing.B) {
	addSetArcheBatch(b, 1000)
}

func BenchmarkArcheAddSet_Batch_10000(b *testing.B) {
	addSetArcheBatch(b, 10000)
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

func BenchmarkArcheAddGeneric_Batch_100(b *testing.B) {
	addArcheGenericBatch(b, 100)
}

func BenchmarkArcheAddGeneric_Batch_1000(b *testing.B) {
	addArcheGenericBatch(b, 1000)
}

func BenchmarkArcheAddGeneric_Batch_10000(b *testing.B) {
	addArcheGenericBatch(b, 10000)
}

func BenchmarkArcheAddSetGeneric_100(b *testing.B) {
	addSetArcheGeneric(b, 100)
}

func BenchmarkArcheAddSetGeneric_1000(b *testing.B) {
	addSetArcheGeneric(b, 1000)
}

func BenchmarkArcheAddSetGeneric_10000(b *testing.B) {
	addSetArcheGeneric(b, 10000)
}

func BenchmarkArcheAddSetGeneric_Batch_100(b *testing.B) {
	addSetArcheGenericBatch(b, 100)
}

func BenchmarkArcheAddSetGeneric_Batch_1000(b *testing.B) {
	addSetArcheGenericBatch(b, 1000)
}

func BenchmarkArcheAddSetGeneric_Batch_10000(b *testing.B) {
	addSetArcheGenericBatch(b, 10000)
}
