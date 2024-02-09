package main

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func benchesEntities() []bench {
	return []bench{
		{Name: "World.NewEntity", Desc: "memory already allocated", F: entitiesCreate_1000, N: 1000},
		{Name: "World.NewEntity w/ 1 Comp", Desc: "memory already allocated", F: entitiesCreate_1Comp_1000, N: 1000},
		{Name: "World.NewEntity w/ 5 Comps", Desc: "memory already allocated", F: entitiesCreate_5Comp_1000, N: 1000},

		{Name: "World.RemoveEntity", Desc: "", F: entitiesRemove_1000, N: 1000},
		{Name: "World.RemoveEntity w/ 1 Comp", Desc: "", F: entitiesRemove_1Comp_1000, N: 1000},
		{Name: "World.RemoveEntity w/ 5 Comps", Desc: "", F: entitiesRemove_5Comp_1000, N: 1000},
	}
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
