package main

import (
	"testing"

	"github.com/mlange-42/arche/benchmark"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchesEntities() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Entity.IsZero", Desc: "", F: entitiesIsZero_2, N: 2},

		{Name: "World.NewEntity", Desc: "memory already allocated", F: entitiesCreate_1000, N: 1000},
		{Name: "World.NewEntity w/ 1 Comp", Desc: "memory already allocated", F: entitiesCreate_1Comp_1000, N: 1000},
		{Name: "World.NewEntity w/ 5 Comps", Desc: "memory already allocated", F: entitiesCreate_5Comp_1000, N: 1000},

		{Name: "World.RemoveEntity", Desc: "", F: entitiesRemove_1000, N: 1000},
		{Name: "World.RemoveEntity w/ 1 Comp", Desc: "", F: entitiesRemove_1Comp_1000, N: 1000},
		{Name: "World.RemoveEntity w/ 5 Comps", Desc: "", F: entitiesRemove_5Comp_1000, N: 1000},

		{Name: "Map1.NewWith w/ 1 Comp", Desc: "memory already allocated", F: entitiesCreateWithGeneric_1Comp_1000, N: 1000},
		{Name: "Map5.NewWith w/ 5 Comps", Desc: "memory already allocated", F: entitiesCreateWithGeneric_5Comp_1000, N: 1000},

		{Name: "World.NewEntityWith w/ 1 Comp", Desc: "⚠️ deprecated, memory already allocated", F: entitiesCreateWith_1Comp_1000, N: 1000},
		{Name: "World.NewEntityWith w/ 5 Comps", Desc: "⚠️ deprecated, memory already allocated", F: entitiesCreateWith_5Comp_1000, N: 1000},
	}
}

func entitiesIsZero_2(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	e := w.NewEntity()
	z := ecs.Entity{}
	var zero1 bool
	var zero2 bool

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		zero1 = e.IsZero()
		zero2 = z.IsZero()
	}
	b.StopTimer()
	s := zero1 || zero2
	_ = s
}

func entitiesCreate_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = w.NewEntity()
		}
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesCreate_1Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = w.NewEntity(id1)
		}
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesCreate_5Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	ids := []ecs.ID{id1, id2, id3, id4, id5}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = w.NewEntity(ids...)
		}
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesCreateWith_1Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)

	comp1 := ecs.Component{ID: id1, Comp: &comp1{}}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = w.NewEntityWith(comp1)
		}
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesCreateWith_5Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)

	comps := []ecs.Component{
		{ID: id1, Comp: &comp1{}},
		{ID: id2, Comp: &comp2{}},
		{ID: id3, Comp: &comp3{}},
		{ID: id4, Comp: &comp4{}},
		{ID: id5, Comp: &comp5{}},
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = w.NewEntityWith(comps...)
		}
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesCreateWithGeneric_1Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	mapper := generic.NewMap1[comp1](&w)

	c1 := comp1{}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = mapper.NewWith(&c1)
		}
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesCreateWithGeneric_5Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	mapper := generic.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)

	c1 := comp1{}
	c2 := comp2{}
	c3 := comp3{}
	c4 := comp4{}
	c5 := comp5{}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = mapper.NewWith(&c1, &c2, &c3, &c4, &c5)
		}
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesRemove_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	builder := ecs.NewBuilder(&w)

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		query := builder.NewBatchQ(1000)
		for query.Next() {
			entities = append(entities, query.Entity())
		}
		b.StartTimer()
		for _, e := range entities {
			w.RemoveEntity(e)
		}
		b.StopTimer()
		entities = entities[:0]
	}
}

func entitiesRemove_1Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))
	id1 := ecs.ComponentID[comp1](&w)
	builder := ecs.NewBuilder(&w, id1)

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		query := builder.NewBatchQ(1000)
		for query.Next() {
			entities = append(entities, query.Entity())
		}
		b.StartTimer()
		for _, e := range entities {
			w.RemoveEntity(e)
		}
		b.StopTimer()
		entities = entities[:0]
	}
}

func entitiesRemove_5Comp_1000(b *testing.B) {
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
		query := builder.NewBatchQ(1000)
		for query.Next() {
			entities = append(entities, query.Entity())
		}
		b.StartTimer()
		for _, e := range entities {
			w.RemoveEntity(e)
		}
		b.StopTimer()
		entities = entities[:0]
	}
}
