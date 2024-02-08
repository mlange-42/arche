package main

import (
	"math/rand"
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

func querEntityAt_1Arch_1000(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	builder := ecs.NewBuilder(&w, id1)
	builder.NewBatch(1000)

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = rand.Intn(1000)
	}

	query := w.Query(ecs.All(id1))
	b.StartTimer()
	var e ecs.Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func querEntityAtRegistered_1Arch_1000(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	builder := ecs.NewBuilder(&w, id1)
	builder.NewBatch(1000)

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = rand.Intn(1000)
	}

	f := ecs.All(id1)
	filter := w.Cache().Register(f)

	query := w.Query(&filter)
	b.StartTimer()
	var e ecs.Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func querEntityAt_5Arch_1000(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	builder := ecs.NewBuilder(&w, id1, id2, id3, id4, id5)
	builder.NewBatch(1000)

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = rand.Intn(1000)
	}

	query := w.Query(ecs.All(id1))
	b.StartTimer()
	var e ecs.Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func querEntityAtRegistered_5Arch_1000(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	builder := ecs.NewBuilder(&w, id1, id2, id3, id4, id5)
	builder.NewBatch(1000)

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = rand.Intn(1000)
	}

	f := ecs.All(id1)
	filter := w.Cache().Register(f)

	query := w.Query(&filter)
	b.StartTimer()
	var e ecs.Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func queryCreate(b *testing.B) {
	b.StopTimer()

	world := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&world)
	builder := ecs.NewBuilder(&world, id1)
	builder.NewBatch(100)

	filter := ecs.All(id1)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		query.Close()
	}
}

func queryCreateCached(b *testing.B) {
	b.StopTimer()

	world := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&world)
	builder := ecs.NewBuilder(&world, id1)
	builder.NewBatch(100)

	f := ecs.All(id1)
	filter := world.Cache().Register(f)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&filter)
		query.Close()
	}
}
