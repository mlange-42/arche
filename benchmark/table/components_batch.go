package main

import (
	"testing"

	"github.com/mlange-42/arche/benchmark"
	"github.com/mlange-42/arche/ecs"
)

func benchesComponentsBatch() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Batch.Add 1 Comp", Desc: "1000, memory already allocated", F: componentsBatchAdd1_1000, N: 1000},
		{Name: "Batch.Add 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchAdd5_1000, N: 1000},

		{Name: "Batch.Remove 1 Comp", Desc: "1000", F: componentsBatchRemove1_1000, N: 1000},
		{Name: "Batch.Remove 5 Comps", Desc: "1000", F: componentsBatchRemove5_1000, N: 1000},

		{Name: "Batch.Exchange 1 Comp", Desc: "1000, memory already allocated", F: componentsBatchExchange1_1000, N: 1000},
	}
}

func componentsBatchAdd1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	filter := ecs.All(id1)

	builder := ecs.NewBuilder(&w)
	builder.NewBatch(1000)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Add(ecs.All(), id1)
		b.StopTimer()
		w.Batch().Remove(filter, id1)
	}
}

func componentsBatchAdd5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	ids := []ecs.ID{id1, id2, id3, id4, id5}
	filter := ecs.All(ids...)

	builder := ecs.NewBuilder(&w)
	builder.NewBatch(1000)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Add(ecs.All(), ids...)
		b.StopTimer()
		w.Batch().Remove(filter, ids...)
	}
}

func componentsBatchRemove1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	filter := ecs.All(id1)

	builder := ecs.NewBuilder(&w, id1)
	builder.NewBatch(1000)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Remove(filter, id1)
		b.StopTimer()
		w.Batch().Add(ecs.All(), id1)
	}
}

func componentsBatchRemove5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	ids := []ecs.ID{id1, id2, id3, id4, id5}
	filter := ecs.All(ids...)

	builder := ecs.NewBuilder(&w, ids...)
	builder.NewBatch(1000)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Remove(filter, ids...)
		b.StopTimer()
		w.Batch().Add(ecs.All(), ids...)
	}
}

func componentsBatchExchange1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	filter1 := ecs.All(id1)
	filter2 := ecs.All(id2)
	ex1 := []ecs.ID{id1}
	ex2 := []ecs.ID{id2}

	builder := ecs.NewBuilder(&w, id1)
	builder.NewBatch(1000)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Exchange(filter1, ex2, ex1)
		b.StopTimer()
		w.Batch().Exchange(filter2, ex1, ex2)
	}
}
