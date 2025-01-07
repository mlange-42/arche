package main

import (
	"testing"

	"github.com/mlange-42/arche/benchmark"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchesComponents() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "World.Add 1 Comp", Desc: "memory already allocated", F: componentsAdd1_1000, N: 1000},
		{Name: "World.Add 5 Comps", Desc: "memory already allocated", F: componentsAdd5_1000, N: 1000},

		{Name: "World.Remove 1 Comp", Desc: "", F: componentsRemove1_1000, N: 1000},
		{Name: "World.Remove 5 Comps", Desc: "", F: componentsRemove5_1000, N: 1000},

		{Name: "World.Exchange 1 Comp", Desc: "memory already allocated", F: componentsExchange1_1000, N: 1000},

		{Name: "Map1.Assign 1 Comps", Desc: "memory already allocated", F: componentsAssignGeneric1_1000, N: 1000},
		{Name: "Map5.Assign 5 Comps", Desc: "memory already allocated", F: componentsAssignGeneric5_1000, N: 1000},

		{Name: "World.Assign 1 Comp", Desc: "⚠️ deprecated, memory already allocated", F: componentsAssign1_1000, N: 1000},
		{Name: "World.Assign 5 Comps", Desc: "⚠️ deprecated, memory already allocated", F: componentsAssign5_1000, N: 1000},
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

	ids := []ecs.ID{id1, id2, id3, id4, id5}
	filter := ecs.All(ids...)

	builder := ecs.NewBuilder(&w)
	query := builder.NewBatchQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Add(e, ids...)
		}
		b.StopTimer()
		w.Batch().Remove(filter, ids...)
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
	ids := []ecs.ID{id1, id2, id3, id4, id5}
	filter := ecs.All()

	builder := ecs.NewBuilder(&w, ids...)
	query := builder.NewBatchQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Remove(e, ids...)
		}
		b.StopTimer()
		w.Batch().Add(filter, ids...)
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

func componentsAssign1_1000(b *testing.B) {
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

	comp1 := ecs.Component{ID: id1, Comp: &comp1{}}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Assign(e, comp1)
		}
		b.StopTimer()
		w.Batch().Remove(filter, id1)
	}
}

func componentsAssign5_1000(b *testing.B) {
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
	query := builder.NewBatchQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	comps := []ecs.Component{
		{ID: id1, Comp: &comp1{}},
		{ID: id2, Comp: &comp2{}},
		{ID: id3, Comp: &comp3{}},
		{ID: id4, Comp: &comp4{}},
		{ID: id5, Comp: &comp5{}},
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Assign(e, comps...)
		}
		b.StopTimer()
		w.Batch().Remove(filter, ids...)
	}
}

func componentsAssignGeneric1_1000(b *testing.B) {
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

	c1 := comp1{}

	mapper := generic.NewMap1[comp1](&w)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			mapper.Assign(e, &c1)
		}
		b.StopTimer()
		w.Batch().Remove(filter, id1)
	}
}

func componentsAssignGeneric5_1000(b *testing.B) {
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
	query := builder.NewBatchQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	c1 := comp1{}
	c2 := comp2{}
	c3 := comp3{}
	c4 := comp4{}
	c5 := comp5{}

	mapper := generic.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			mapper.Assign(e, &c1, &c2, &c3, &c4, &c5)
		}
		b.StopTimer()
		w.Batch().Remove(filter, ids...)
	}
}
