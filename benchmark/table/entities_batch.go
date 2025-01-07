package main

import (
	"testing"

	"github.com/mlange-42/arche/benchmark"
	"github.com/mlange-42/arche/ecs"
)

func benchesEntitiesBatch() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Builder.NewBatch", Desc: "1000, memory already allocated", F: entitiesBatchCreate_1000, N: 1000},
		{Name: "Builder.NewBatch w/ 1 Comp", Desc: "1000, memory already allocated", F: entitiesBatchCreate_1Comp_1000, N: 1000},
		{Name: "Builder.NewBatch w/ 5 Comps", Desc: "1000, memory already allocated", F: entitiesBatchCreate_5Comp_1000, N: 1000},

		{Name: "Batch.RemoveEntities", Desc: "1000", F: entitiesBatchRemove_1000, N: 1000},
		{Name: "Batch.RemoveEntities w/ 1 Comp", Desc: "1000", F: entitiesBatchRemove_1Comp_1000, N: 1000},
		{Name: "Batch.RemoveEntities w/ 5 Comps", Desc: "1000", F: entitiesBatchRemove_5Comp_1000, N: 1000},
	}
}

func entitiesBatchCreate_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	builder := ecs.NewBuilder(&w)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		builder.NewBatch(1000)
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesBatchCreate_1Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	builder := ecs.NewBuilder(&w, id1)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		builder.NewBatch(1000)
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesBatchCreate_5Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	ids := []ecs.ID{id1, id2, id3, id4, id5}
	builder := ecs.NewBuilder(&w, ids...)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		builder.NewBatch(1000)
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesBatchRemove_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	builder := ecs.NewBuilder(&w)

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		builder.NewBatch(1000)
		b.StartTimer()
		w.Batch().RemoveEntities(ecs.All())
		b.StopTimer()
		entities = entities[:0]
	}
}

func entitiesBatchRemove_1Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	builder := ecs.NewBuilder(&w, id1)

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		builder.NewBatch(1000)
		b.StartTimer()
		w.Batch().RemoveEntities(ecs.All())
		b.StopTimer()
		entities = entities[:0]
	}
}

func entitiesBatchRemove_5Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	ids := []ecs.ID{id1, id2, id3, id4, id5}
	builder := ecs.NewBuilder(&w, ids...)

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		builder.NewBatch(1000)
		b.StartTimer()
		w.Batch().RemoveEntities(ecs.All())
		b.StopTimer()
		entities = entities[:0]
	}
}
