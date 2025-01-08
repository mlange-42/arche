package iterate

import (
	"testing"

	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/filter"
	"github.com/mlange-42/arche/generic"
)

func runIter(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(count)

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	ecs.NewBuilder(&world, posID, rotID).NewBatch(count)
	filter := ecs.All(posID, rotID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&filter)
		cnt := 0
		for query.Next() {
			cnt++
		}
		_ = cnt
	}
}

func runGet(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(count)

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	ecs.NewBuilder(&world, posID, rotID).NewBatch(count)

	filter := ecs.All(posID, rotID)
	query := world.Query(&filter)
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

func runGetEntity(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(count)

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	ecs.NewBuilder(&world, posID, rotID).NewBatch(count)

	filter := ecs.All(posID, rotID)
	query := world.Query(&filter)
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

func runQuery(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(count)

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	ecs.NewBuilder(&world, posID, rotID).NewBatch(count)
	filter := ecs.All(posID, rotID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			pos := (*c.Position)(query.Get(posID))
			pos.X = 1.0
		}
	}
}

func runQueryCached(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(count)

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	ecs.NewBuilder(&world, posID, rotID).NewBatch(count)

	cf := world.Cache().Register(ecs.All(posID, rotID))
	var filter ecs.Filter = &cf
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		for query.Next() {
			pos := (*c.Position)(query.Get(posID))
			pos.X = 1.0
		}
	}
}

func runFilter(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(count)

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	ecs.NewBuilder(&world, posID, rotID).NewBatch(count)
	filter := filter.All(posID, rotID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			pos := (*c.Position)(query.Get(posID))
			pos.X = 1.0
		}
	}
}

func runQueryGeneric(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(count)

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	ecs.NewBuilder(&world, posID, rotID).NewBatch(count)

	query := generic.NewFilter1[c.Position]()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q := query.Query(&world)
		for q.Next() {
			pos := q.Get()
			pos.X = 1.0
		}
	}
}

func runQuery5C(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(count)

	id0 := ecs.ComponentID[c.TestStruct0](&world)
	id1 := ecs.ComponentID[c.TestStruct1](&world)
	id2 := ecs.ComponentID[c.TestStruct2](&world)
	id3 := ecs.ComponentID[c.TestStruct3](&world)
	id4 := ecs.ComponentID[c.TestStruct4](&world)

	ecs.NewBuilder(&world, id0, id1, id2, id3, id4).NewBatch(count)
	filter := ecs.All(id0, id1, id2, id3, id4)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&filter)
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

func runQueryGeneric5C(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(count)

	id0 := ecs.ComponentID[c.TestStruct0](&world)
	id1 := ecs.ComponentID[c.TestStruct1](&world)
	id2 := ecs.ComponentID[c.TestStruct2](&world)
	id3 := ecs.ComponentID[c.TestStruct3](&world)
	id4 := ecs.ComponentID[c.TestStruct4](&world)

	ecs.NewBuilder(&world, id0, id1, id2, id3, id4).NewBatch(count)

	query := generic.NewFilter5[c.TestStruct0, c.TestStruct1, c.TestStruct2, c.TestStruct3, c.TestStruct4]()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q := query.Query(&world)
		for q.Next() {
			t1, t2, t3, t4, t5 := q.Get()
			t1.Val, t2.Val, t3.Val, t4.Val, t5.Val = 1, 1, 1, 1, 1
		}
	}
}

func BenchmarkIter_1_000(b *testing.B) {
	runIter(b, 1000)
}

func BenchmarkIter_10_000(b *testing.B) {
	runIter(b, 10000)
}

func BenchmarkIter_100_000(b *testing.B) {
	runIter(b, 100000)
}

func BenchmarkGet_1_000(b *testing.B) {
	runGet(b, 1000)
}

func BenchmarkGet_10_000(b *testing.B) {
	runGet(b, 10000)
}

func BenchmarkGet_100_000(b *testing.B) {
	runGet(b, 100000)
}

func BenchmarkGetEntity_1_000(b *testing.B) {
	runGetEntity(b, 1000)
}

func BenchmarkGetEntity_10_000(b *testing.B) {
	runGetEntity(b, 10000)
}

func BenchmarkGetEntity_100_000(b *testing.B) {
	runGetEntity(b, 100000)
}

func BenchmarkIterQueryID_1_000(b *testing.B) {
	runQuery(b, 1000)
}

func BenchmarkIterQueryID_10_000(b *testing.B) {
	runQuery(b, 10000)
}

func BenchmarkIterQueryID_100_000(b *testing.B) {
	runQuery(b, 100000)
}

func BenchmarkIterQueryIDCached_1_000(b *testing.B) {
	runQueryCached(b, 1000)
}

func BenchmarkIterQueryIDCached_10_000(b *testing.B) {
	runQueryCached(b, 10000)
}

func BenchmarkIterQueryIDCached_100_000(b *testing.B) {
	runQueryCached(b, 100000)
}

func BenchmarkIterFilter_1_000(b *testing.B) {
	runFilter(b, 1000)
}

func BenchmarkIterFilter_10_000(b *testing.B) {
	runFilter(b, 10000)
}

func BenchmarkIterFilter_100_000(b *testing.B) {
	runFilter(b, 100000)
}

func BenchmarkIterQueryGeneric_1_000(b *testing.B) {
	runQueryGeneric(b, 1000)
}

func BenchmarkIterQueryGeneric_10_000(b *testing.B) {
	runQueryGeneric(b, 10000)
}

func BenchmarkIterQueryGeneric_100_000(b *testing.B) {
	runQueryGeneric(b, 100000)
}

func BenchmarkIterQueryID_5C_1_000(b *testing.B) {
	runQuery5C(b, 1000)
}

func BenchmarkIterQueryID_5C_10_000(b *testing.B) {
	runQuery5C(b, 10000)
}

func BenchmarkIterQueryID_5C_100_000(b *testing.B) {
	runQuery5C(b, 100000)
}

func BenchmarkIterQueryGeneric_5C_1_000(b *testing.B) {
	runQueryGeneric5C(b, 1000)
}

func BenchmarkIterQueryGeneric_5C_10_000(b *testing.B) {
	runQueryGeneric5C(b, 10000)
}

func BenchmarkIterQueryGeneric_5C_100_000(b *testing.B) {
	runQueryGeneric5C(b, 100000)
}
