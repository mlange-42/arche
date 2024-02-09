package main

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func benchesComponents() []bench {
	return []bench{
		{Name: "World.Add 1 Comp", Desc: "memory already allocated", F: componentsAdd1_1000, N: 1000},
		{Name: "World.Add 5 Comps", Desc: "memory already allocated", F: componentsAdd5_1000, N: 1000},

		{Name: "World.Remove 1 Comp", Desc: "", F: componentsRemove1_1000, N: 1000},
		{Name: "World.Remove 5 Comps", Desc: "", F: componentsRemove5_1000, N: 1000},

		{Name: "World.Exchange 1 Comp", Desc: "memory already allocated", F: componentsExchange1_1000, N: 1000},
	}
}

func componentsAdd1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	filter := ecs.All(id1)

	builder := ecs.NewBuilder(&w)
	query := builder.NewBatchQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Add(e, id1)
		}
		b.StopTimer()
		w.Batch().Remove(filter, id1)
	}
}

func componentsAdd5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	filter := ecs.All(id1, id2, id3, id4, id5)

	builder := ecs.NewBuilder(&w)
	query := builder.NewBatchQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Add(e, id1, id2, id3, id4, id5)
		}
		b.StopTimer()
		w.Batch().Remove(filter, id1, id2, id3, id4, id5)
	}
}

func componentsRemove1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	filter := ecs.All()

	builder := ecs.NewBuilder(&w, id1)
	query := builder.NewBatchQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Remove(e, id1)
		}
		b.StopTimer()
		w.Batch().Add(filter, id1)
	}
}

func componentsRemove5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	filter := ecs.All()

	builder := ecs.NewBuilder(&w, id1, id2, id3, id4, id5)
	query := builder.NewBatchQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Remove(e, id1, id2, id3, id4, id5)
		}
		b.StopTimer()
		w.Batch().Add(filter, id1, id2, id3, id4, id5)
	}
}

func componentsExchange1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	filter := ecs.All(id2)
	ex1 := []ecs.ID{id1}
	ex2 := []ecs.ID{id2}

	builder := ecs.NewBuilder(&w, id1)
	query := builder.NewBatchQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Exchange(e, ex2, ex1)
		}
		b.StopTimer()
		w.Batch().Exchange(filter, ex1, ex2)
	}
}
