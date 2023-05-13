package ecs

import (
	"testing"
)

func benchmarkRelationGetQuery(b *testing.B, count int) {
	b.StopTimer()

	world := NewWorld(NewConfig().WithCapacityIncrement(1024).WithRelationCapacityIncrement(128))
	relID := ComponentID[testRelationA](&world)

	target := world.NewEntity()

	builder := NewBuilder(&world, relID).WithRelation(relID)
	builder.NewBatch(count, target)

	filter := All(relID)
	b.StartTimer()

	var tempTarget Entity
	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		for query.Next() {
			tempTarget = query.Relation(relID)
		}
	}

	_ = tempTarget
}

func benchmarkRelationGetQueryUnchecked(b *testing.B, count int) {
	b.StopTimer()

	world := NewWorld(NewConfig().WithCapacityIncrement(1024).WithRelationCapacityIncrement(128))
	relID := ComponentID[testRelationA](&world)

	target := world.NewEntity()

	builder := NewBuilder(&world, relID).WithRelation(relID)
	builder.NewBatch(count, target)

	filter := All(relID)
	b.StartTimer()

	var tempTarget Entity
	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		for query.Next() {
			tempTarget = query.relationUnchecked(relID)
		}
	}

	_ = tempTarget
}

func benchmarkRelationGetWorld(b *testing.B, count int) {
	b.StopTimer()

	world := NewWorld(NewConfig().WithCapacityIncrement(1024).WithRelationCapacityIncrement(128))
	relID := ComponentID[testRelationA](&world)

	target := world.NewEntity()

	builder := NewBuilder(&world, relID).WithRelation(relID)
	q := builder.NewBatchQ(count, target)
	entities := make([]Entity, 0, count)
	for q.Next() {
		entities = append(entities, q.Entity())
	}
	b.StartTimer()

	var tempTarget Entity
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			tempTarget = world.Relations().Get(e, relID)
		}
	}

	_ = tempTarget
}

func benchmarkRelationGetWorldUnchecked(b *testing.B, count int) {
	b.StopTimer()

	world := NewWorld(NewConfig().WithCapacityIncrement(1024).WithRelationCapacityIncrement(128))
	relID := ComponentID[testRelationA](&world)

	target := world.NewEntity()

	builder := NewBuilder(&world, relID).WithRelation(relID)
	q := builder.NewBatchQ(count, target)
	entities := make([]Entity, 0, count)
	for q.Next() {
		entities = append(entities, q.Entity())
	}
	b.StartTimer()

	var tempTarget Entity
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			tempTarget = world.Relations().GetUnchecked(e, relID)
		}
	}

	_ = tempTarget
}

func benchmarkRelationSet(b *testing.B, count int) {
	b.StopTimer()

	world := NewWorld(NewConfig().WithCapacityIncrement(1024).WithRelationCapacityIncrement(128))
	relID := ComponentID[testRelationA](&world)

	target := world.NewEntity()

	builder := NewBuilder(&world, relID).WithRelation(relID)
	q := builder.NewBatchQ(count)
	entities := make([]Entity, 0, count)
	for q.Next() {
		entities = append(entities, q.Entity())
	}
	b.StartTimer()

	var tempTarget Entity
	for i := 0; i < b.N; i++ {
		trg := Entity{}
		if i%2 == 0 {
			trg = target
		}
		for _, e := range entities {
			world.Relations().Set(e, relID, trg)
		}
	}

	_ = tempTarget
}

func BenchmarkRelationGetQuery_1000(b *testing.B) {
	benchmarkRelationGetQuery(b, 1000)
}

func BenchmarkRelationGetQuery_10000(b *testing.B) {
	benchmarkRelationGetQuery(b, 10000)
}

func BenchmarkRelationGetQuery_100000(b *testing.B) {
	benchmarkRelationGetQuery(b, 100000)
}

func BenchmarkRelationGetQueryUnchecked_1000(b *testing.B) {
	benchmarkRelationGetQueryUnchecked(b, 1000)
}

func BenchmarkRelationGetQueryUnchecked_10000(b *testing.B) {
	benchmarkRelationGetQueryUnchecked(b, 10000)
}

func BenchmarkRelationGetQueryUnchecked_100000(b *testing.B) {
	benchmarkRelationGetQueryUnchecked(b, 100000)
}

func BenchmarkRelationGetWorld_1000(b *testing.B) {
	benchmarkRelationGetWorld(b, 1000)
}

func BenchmarkRelationGetWorld_10000(b *testing.B) {
	benchmarkRelationGetWorld(b, 10000)
}

func BenchmarkRelationGetWorld_100000(b *testing.B) {
	benchmarkRelationGetWorld(b, 100000)
}

func BenchmarkRelationGetWorldUnchecked_1000(b *testing.B) {
	benchmarkRelationGetWorldUnchecked(b, 1000)
}

func BenchmarkRelationGetWorldUnchecked_10000(b *testing.B) {
	benchmarkRelationGetWorldUnchecked(b, 10000)
}

func BenchmarkRelationGetWorldUnchecked_100000(b *testing.B) {
	benchmarkRelationGetWorldUnchecked(b, 100000)
}

func BenchmarkRelationSet_1000(b *testing.B) {
	benchmarkRelationSet(b, 1000)
}

func BenchmarkRelationSet_10000(b *testing.B) {
	benchmarkRelationSet(b, 10000)
}

func BenchmarkRelationSet_100000(b *testing.B) {
	benchmarkRelationSet(b, 100000)
}
