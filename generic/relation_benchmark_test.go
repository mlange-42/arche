package generic_test

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

type testRelationA struct {
	ecs.Relation
}

func benchmarkRelationGetQuery(b *testing.B, count int) {
	b.StopTimer()

	world := ecs.NewWorld(1024, 128)
	relID := ecs.ComponentID[testRelationA](&world)

	target := world.NewEntity()

	builder := ecs.NewBuilder(&world, relID).WithRelation(relID)
	builder.NewBatch(count, target)

	filter := generic.NewFilter1[testRelationA]().WithRelation(generic.T[testRelationA](), target)
	b.StartTimer()

	var tempTarget ecs.Entity
	for i := 0; i < b.N; i++ {
		query := filter.Query(&world)
		for query.Next() {
			tempTarget = query.Relation()
		}
	}

	_ = tempTarget
}

func benchmarkRelationGetWorld(b *testing.B, count int) {
	b.StopTimer()

	world := ecs.NewWorld(1024, 128)
	relID := ecs.ComponentID[testRelationA](&world)

	target := world.NewEntity()

	builder := ecs.NewBuilder(&world, relID).WithRelation(relID)
	q := builder.NewBatchQ(count, target)
	entities := make([]ecs.Entity, 0, count)
	for q.Next() {
		entities = append(entities, q.Entity())
	}

	mapper := generic.NewMap[testRelationA](&world)
	b.StartTimer()

	var tempTarget ecs.Entity
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			tempTarget = mapper.GetRelation(e)
		}
	}

	_ = tempTarget
}

func benchmarkRelationGetWorldUnchecked(b *testing.B, count int) {
	b.StopTimer()

	world := ecs.NewWorld(1024, 128)
	relID := ecs.ComponentID[testRelationA](&world)

	target := world.NewEntity()

	builder := ecs.NewBuilder(&world, relID).WithRelation(relID)
	q := builder.NewBatchQ(count, target)
	entities := make([]ecs.Entity, 0, count)
	for q.Next() {
		entities = append(entities, q.Entity())
	}
	mapper := generic.NewMap[testRelationA](&world)
	b.StartTimer()

	var tempTarget ecs.Entity
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			tempTarget = mapper.GetRelationUnchecked(e)
		}
	}

	_ = tempTarget
}

func benchmarkRelationSet(b *testing.B, count int) {
	b.StopTimer()

	world := ecs.NewWorld(1024, 128)
	relID := ecs.ComponentID[testRelationA](&world)

	target := world.NewEntity()

	builder := ecs.NewBuilder(&world, relID).WithRelation(relID)
	q := builder.NewBatchQ(count)
	entities := make([]ecs.Entity, 0, count)
	for q.Next() {
		entities = append(entities, q.Entity())
	}
	b.StartTimer()

	var tempTarget ecs.Entity
	for i := 0; i < b.N; i++ {
		trg := ecs.Entity{}
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
