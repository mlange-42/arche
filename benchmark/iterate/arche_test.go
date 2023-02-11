package iterate

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func runArcheQuery(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[position](&world)
	rotID := ecs.ComponentID[rotation](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(posID, rotID)
		b.StartTimer()
		for query.Next() {
			pos := (*position)(query.Get(posID))
			_ = pos
		}
	}
}

func runArcheQueryPointerHeap(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[position](&world)
	rotID := ecs.ComponentID[rotation](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := createQuery(&world, posID, rotID)
		b.StartTimer()
		for query.Next() {
			pos := (*position)(query.Get(posID))
			_ = pos
		}
	}
}

func createQuery(w *ecs.World, comps ...ecs.ID) *ecs.Query {
	q := w.Query(comps...)
	return &q
}

func runArcheQueryGeneric(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[position](&world)
	rotID := ecs.ComponentID[rotation](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := ecs.Query2[position, rotation](&world)
		b.StartTimer()
		for query.Next() {
			pos := query.Get1()
			_ = pos
		}
	}
}
func runArcheQuery5C(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	id0 := ecs.ComponentID[testStruct0](&world)
	id1 := ecs.ComponentID[testStruct1](&world)
	id2 := ecs.ComponentID[testStruct2](&world)
	id3 := ecs.ComponentID[testStruct3](&world)
	id4 := ecs.ComponentID[testStruct4](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, id0, id1, id2, id3, id4)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(id0, id1, id2, id3, id4)
		b.StartTimer()
		for query.Next() {
			t1 := (*testStruct0)(query.Get(id0))
			t2 := (*testStruct0)(query.Get(id1))
			t3 := (*testStruct0)(query.Get(id2))
			t4 := (*testStruct0)(query.Get(id3))
			t5 := (*testStruct0)(query.Get(id4))
			_, _, _, _, _ = t1, t2, t3, t4, t5
		}
	}
}

func runArcheQueryGeneric5C(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	id0 := ecs.ComponentID[testStruct0](&world)
	id1 := ecs.ComponentID[testStruct1](&world)
	id2 := ecs.ComponentID[testStruct2](&world)
	id3 := ecs.ComponentID[testStruct3](&world)
	id4 := ecs.ComponentID[testStruct4](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, id0, id1, id2, id3, id4)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := ecs.Query5[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4](&world)
		b.StartTimer()
		for query.Next() {
			t1, t2, t3, t4, t5 := query.GetAll()
			_, _, _, _, _ = t1, t2, t3, t4, t5
		}
	}
}

func runArcheQuery1kArch(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	registerAll(&world)

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
		query := world.Query(6)
		b.StartTimer()
		for query.Next() {
			pos := (*position)(query.Get(6))
			_ = pos
		}
	}
}

func runArcheWorld(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[position](&world)
	rotID := ecs.ComponentID[rotation](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(posID, rotID)
		b.StartTimer()
		for query.Next() {
			entity := query.Entity()
			pos := (*position)(world.Get(entity, posID))
			_ = pos
		}
	}
}

func runArcheWorldGeneric(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[position](&world)
	rotID := ecs.ComponentID[rotation](&world)

	get := ecs.NewMap[position](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(posID, rotID)
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

func BenchmarkArcheIterQueryIDPointer_1_000(b *testing.B) {
	runArcheQueryPointerHeap(b, 1000)
}

func BenchmarkArcheIterQueryIDPointer_10_000(b *testing.B) {
	runArcheQueryPointerHeap(b, 10000)
}

func BenchmarkArcheIterQueryIDPointer_100_000(b *testing.B) {
	runArcheQueryPointerHeap(b, 100000)
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
