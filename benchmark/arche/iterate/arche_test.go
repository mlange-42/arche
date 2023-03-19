package iterate

import (
	"testing"

	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/filter"
	"github.com/mlange-42/arche/generic"
)

func runArcheIter(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(posID, rotID)
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(ecs.All(posID, rotID))
		cnt := 0
		b.StartTimer()
		for query.Next() {
			cnt++
		}
		_ = cnt
	}
}

func runArcheGet(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(posID, rotID)
	}

	query := world.Query(ecs.All(posID, rotID))
	for query.Next() {
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < count; i++ {
				pos := (*c.Position)(query.Get(posID))
				pos.X = 1.0
			}
		}
		b.StopTimer()
		query.Close()
		if true {
			break
		}
	}
}

func runArcheGetEntity(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(posID, rotID)
	}

	query := world.Query(ecs.All(posID, rotID))
	for query.Next() {
		b.StartTimer()
		var e ecs.Entity
		for i := 0; i < b.N; i++ {
			for i := 0; i < count; i++ {
				e = query.Entity()
			}
		}
		b.StopTimer()
		_ = e
		query.Close()
		if true {
			break
		}
	}
}

func runArcheQuery(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	for i := 0; i < count; i++ {
		_ = world.NewEntity(posID, rotID)
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(ecs.All(posID, rotID))
		b.StartTimer()
		for query.Next() {
			pos := (*c.Position)(query.Get(posID))
			pos.X = 1.0
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

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter.All(posID, rotID))
		b.StartTimer()
		for query.Next() {
			pos := (*c.Position)(query.Get(posID))
			pos.X = 1.0
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

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		q := query.Query(&world)
		b.StartTimer()
		for q.Next() {
			pos := q.Get()
			pos.X = 1.0
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
			t1.Val, t2.Val, t3.Val, t4.Val, t5.Val = 1, 1, 1, 1, 1
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

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		q := query.Query(&world)
		b.StartTimer()
		for q.Next() {
			t1, t2, t3, t4, t5 := q.Get()
			t1.Val, t2.Val, t3.Val, t4.Val, t5.Val = 1, 1, 1, 1, 1
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

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(ecs.All(6))
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct6)(query.Get(6))
			pos.Val = 1
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

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter.All(6))
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct6)(query.Get(6))
			pos.Val = 1
		}
	}
}

func runArcheWorldGet(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	entities := make([]ecs.Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = world.NewEntity(posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			pos := (*c.Position)(world.Get(e, posID))
			pos.X = 1
		}
	}
}

func runArcheWorldGetGeneric(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	get := generic.NewMap[c.Position](&world)

	entities := make([]ecs.Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = world.NewEntity(posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			pos := get.Get(e)
			pos.X = 1
		}
	}
}

func BenchmarkArcheIter_1_000(b *testing.B) {
	runArcheIter(b, 1000)
}

func BenchmarkArcheIter_10_000(b *testing.B) {
	runArcheIter(b, 10000)
}

func BenchmarkArcheIter_100_000(b *testing.B) {
	runArcheIter(b, 100000)
}

func BenchmarkArcheGet_1_000(b *testing.B) {
	runArcheGet(b, 1000)
}

func BenchmarkArcheGet_10_000(b *testing.B) {
	runArcheGet(b, 10000)
}

func BenchmarkArcheGet_100_000(b *testing.B) {
	runArcheGet(b, 100000)
}

func BenchmarkArcheGetEntity_1_000(b *testing.B) {
	runArcheGetEntity(b, 1000)
}

func BenchmarkArcheGetEntity_10_000(b *testing.B) {
	runArcheGetEntity(b, 10000)
}

func BenchmarkArcheGetEntity_100_000(b *testing.B) {
	runArcheGetEntity(b, 100000)
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
	runArcheWorldGet(b, 1000)
}

func BenchmarkArcheIterWorldID_10_000(b *testing.B) {
	runArcheWorldGet(b, 10000)
}

func BenchmarkArcheIterWorldID_100_000(b *testing.B) {
	runArcheWorldGet(b, 100000)
}

func BenchmarkArcheIterWorldGeneric_1_000(b *testing.B) {
	runArcheWorldGetGeneric(b, 1000)
}

func BenchmarkArcheIterWorldGeneric_10_000(b *testing.B) {
	runArcheWorldGetGeneric(b, 10000)
}

func BenchmarkArcheIterWorldGeneric_100_000(b *testing.B) {
	runArcheWorldGetGeneric(b, 100000)
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
