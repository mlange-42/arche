package main

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func queryIter_100_000(b *testing.B) {
	w := ecs.NewWorld()
	builder := ecs.NewBuilder(&w)
	builder.NewBatch(100_000)
	filter := ecs.All()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := w.Query(filter)
		b.StartTimer()
		for query.Next() {
		}
	}
}

func queryIterGet_1_100_000(b *testing.B) {
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)

	builder := ecs.NewBuilder(&w, id1)
	builder.NewBatch(100_000)
	filter := ecs.All(id1)

	var c1 *comp1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := w.Query(filter)
		b.StartTimer()
		for query.Next() {
			c1 = (*comp1)(query.Get(id1))
		}
	}
	b.StopTimer()
	sum := c1.V
	_ = sum
}

func queryIterGet_2_100_000(b *testing.B) {
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)

	builder := ecs.NewBuilder(&w, id1, id2)
	builder.NewBatch(100_000)
	filter := ecs.All(id1, id2)

	var c1 *comp1
	var c2 *comp2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := w.Query(filter)
		b.StartTimer()
		for query.Next() {
			c1 = (*comp1)(query.Get(id1))
			c2 = (*comp2)(query.Get(id2))
		}
	}
	b.StopTimer()
	sum := c1.V + c2.V
	_ = sum
}

func queryIterGet_5_100_000(b *testing.B) {
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)

	builder := ecs.NewBuilder(&w, id1, id2, id3, id4, id5)
	builder.NewBatch(100_000)
	filter := ecs.All(id1, id2, id3, id4, id5)

	var c1 *comp1
	var c2 *comp2
	var c3 *comp3
	var c4 *comp4
	var c5 *comp5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := w.Query(filter)
		b.StartTimer()
		for query.Next() {
			c1 = (*comp1)(query.Get(id1))
			c2 = (*comp2)(query.Get(id2))
			c3 = (*comp3)(query.Get(id3))
			c4 = (*comp4)(query.Get(id4))
			c5 = (*comp5)(query.Get(id5))
		}
	}
	b.StopTimer()
	sum := c1.V + c2.V + c3.V + c4.V + c5.V
	_ = sum
}
