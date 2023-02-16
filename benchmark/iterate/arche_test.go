package iterate

import (
	"testing"

	c "github.com/mlange-42/arche/benchmark/common"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/filter"
	"github.com/mlange-42/arche/generic"
)

func runArcheQuery(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(ecs.All(posID, rotID))
		b.StartTimer()
		for query.Next() {
			pos := (*c.Position)(query.Get(posID))
			_ = pos
		}
	}
}

func runArcheFilter(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter.All(posID, rotID))
		b.StartTimer()
		for query.Next() {
			pos := (*c.Position)(query.Get(posID))
			_ = pos
		}
	}
}

func runArcheQueryGeneric(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(posID, rotID)
	}
	query := generic.NewFilter1[c.Position]()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		q := query.Query(&world)
		b.StartTimer()
		for q.Next() {
			_, pos := q.Get()
			_ = pos
		}
	}
}

func runArcheQuery5C(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	id0 := ecs.ComponentID[c.TestStruct0](&world)
	id1 := ecs.ComponentID[c.TestStruct1](&world)
	id2 := ecs.ComponentID[c.TestStruct2](&world)
	id3 := ecs.ComponentID[c.TestStruct3](&world)
	id4 := ecs.ComponentID[c.TestStruct4](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(id0, id1, id2, id3, id4)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(ecs.All(id0, id1, id2, id3, id4))
		b.StartTimer()
		for query.Next() {
			t1 := (*c.TestStruct0)(query.Get(id0))
			t2 := (*c.TestStruct1)(query.Get(id1))
			t3 := (*c.TestStruct2)(query.Get(id2))
			t4 := (*c.TestStruct3)(query.Get(id3))
			t5 := (*c.TestStruct4)(query.Get(id4))
			_, _, _, _, _ = t1, t2, t3, t4, t5
		}
	}
}

func runArcheQueryGeneric5C(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	id0 := ecs.ComponentID[c.TestStruct0](&world)
	id1 := ecs.ComponentID[c.TestStruct1](&world)
	id2 := ecs.ComponentID[c.TestStruct2](&world)
	id3 := ecs.ComponentID[c.TestStruct3](&world)
	id4 := ecs.ComponentID[c.TestStruct4](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(id0, id1, id2, id3, id4)
	}

	query := generic.NewFilter5[c.TestStruct0, c.TestStruct1, c.TestStruct2, c.TestStruct3, c.TestStruct4]()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		q := query.Query(&world)
		b.StartTimer()
		for q.Next() {
			_, _, _, _, _, _ = q.Get()
		}
	}
}

func runArcheQuery1kArch(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	c.RegisterAll(&world)

	perArch := 2 * count / 1000

	for i := 0; i < 1024; i++ {
		mask := i
		add := make([]ecs.ID, 0, 10)
		for j := 0; j < 10; j++ {
			id := ecs.ID(j)
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			entity := world.NewEntity()
			world.Add(entity, add...)
		}
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(ecs.All(6))
		b.StartTimer()
		for query.Next() {
			pos := (*c.Position)(query.Get(6))
			_ = pos
		}
	}
}

func runArcheFilter1kArch(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	c.RegisterAll(&world)

	perArch := 2 * count / 1000

	for i := 0; i < 1024; i++ {
		mask := i
		add := make([]ecs.ID, 0, 10)
		for j := 0; j < 10; j++ {
			id := ecs.ID(j)
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			entity := world.NewEntity()
			world.Add(entity, add...)
		}
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter.All(6))
		b.StartTimer()
		for query.Next() {
			pos := (*c.Position)(query.Get(6))
			_ = pos
		}
	}
}

func runArcheWorld(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(ecs.All(posID, rotID))
		b.StartTimer()
		for query.Next() {
			entity := query.Entity()
			pos := (*c.Position)(world.Get(entity, posID))
			_ = pos
		}
	}
}

func runArcheWorldGeneric(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	get := generic.NewMap[c.Position](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(ecs.All(posID, rotID))
		b.StartTimer()
		for query.Next() {
			entity := query.Entity()
			pos := get.Get(entity)
			_ = pos
		}
	}
}

func BenchmarkArcheIterQueryID_1_000(b *testing.B) {
	runArcheQuery(b, 1000)
}

func BenchmarkArcheIterQueryID_10_000(b *testing.B) {
	runArcheQuery(b, 10000)
}

func BenchmarkArcheIterQueryID_100_000(b *testing.B) {
	runArcheQuery(b, 100000)
}

func BenchmarkArcheIterFilter_1_000(b *testing.B) {
	runArcheFilter(b, 1000)
}

func BenchmarkArcheIterFilter_10_000(b *testing.B) {
	runArcheFilter(b, 10000)
}

func BenchmarkArcheIterFilter_100_000(b *testing.B) {
	runArcheFilter(b, 100000)
}

func BenchmarkArcheIterQueryGeneric_1_000(b *testing.B) {
	runArcheQueryGeneric(b, 1000)
}

func BenchmarkArcheIterQueryGeneric_10_000(b *testing.B) {
	runArcheQueryGeneric(b, 10000)
}

func BenchmarkArcheIterQueryGeneric_100_000(b *testing.B) {
	runArcheQueryGeneric(b, 100000)
}

func BenchmarkArcheIterQueryID_5C_1_000(b *testing.B) {
	runArcheQuery5C(b, 1000)
}

func BenchmarkArcheIterQueryID_5C_10_000(b *testing.B) {
	runArcheQuery5C(b, 10000)
}

func BenchmarkArcheIterQueryID_5C_100_000(b *testing.B) {
	runArcheQuery5C(b, 100000)
}

func BenchmarkArcheIterQueryGeneric_5C_1_000(b *testing.B) {
	runArcheQueryGeneric5C(b, 1000)
}

func BenchmarkArcheIterQueryGeneric_5C_10_000(b *testing.B) {
	runArcheQueryGeneric5C(b, 10000)
}

func BenchmarkArcheIterQueryGeneric_5C_100_000(b *testing.B) {
	runArcheQueryGeneric5C(b, 100000)
}

func BenchmarkArcheIterWorldID_1_000(b *testing.B) {
	runArcheWorld(b, 1000)
}

func BenchmarkArcheIterWorldID_10_000(b *testing.B) {
	runArcheWorld(b, 10000)
}

func BenchmarkArcheIterWorldID_100_000(b *testing.B) {
	runArcheWorld(b, 100000)
}

func BenchmarkArcheIterWorldGeneric_1_000(b *testing.B) {
	runArcheWorldGeneric(b, 1000)
}

func BenchmarkArcheIterWorldGeneric_10_000(b *testing.B) {
	runArcheWorldGeneric(b, 10000)
}

func BenchmarkArcheIterWorldGeneric_100_000(b *testing.B) {
	runArcheWorldGeneric(b, 100000)
}

func BenchmarkArcheIter1kArchID_1_000(b *testing.B) {
	runArcheQuery1kArch(b, 1000)
}

func BenchmarkArcheIter1kArchID_10_000(b *testing.B) {
	runArcheQuery1kArch(b, 10000)
}

func BenchmarkArcheIter1kArchID_100_000(b *testing.B) {
	runArcheQuery1kArch(b, 100000)
}

func BenchmarkArcheFilter1kArchID_1_000(b *testing.B) {
	runArcheFilter1kArch(b, 1000)
}

func BenchmarkArcheFilter1kArchID_10_000(b *testing.B) {
	runArcheFilter1kArch(b, 10000)
}

func BenchmarkArcheFilter1kArchID_100_000(b *testing.B) {
	runArcheFilter1kArch(b, 100000)
}
