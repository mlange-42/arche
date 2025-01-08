package main

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/arche/benchmark"
	"github.com/mlange-42/arche/ecs"
)

func benchesWorld() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "World.Get", Desc: "random, 1000 entities", F: worldGet_1000, N: 1000},
		{Name: "World.GetUnchecked", Desc: "random, 1000 entities", F: worldGetUnchecked_1000, N: 1000},
		{Name: "World.Has", Desc: "random, 1000 entities", F: worldHas_1000, N: 1000},
		{Name: "World.HasUnchecked", Desc: "random, 1000 entities", F: worldHasUnchecked_1000, N: 1000},
		{Name: "World.Alive", Desc: "random, 1000 entities", F: worldAlive_1000, N: 1000},
		{Name: "World.Relations.Get", Desc: "random, 1000 entities", F: worldRelation_1000, N: 1000},
		{Name: "World.Relations.GetUnchecked", Desc: "random, 1000 entities", F: worldRelationUnchecked_1000, N: 1000},
	}
}

func worldGet_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewBuilder(&w, id1)
	query := builder.NewBatchQ(1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var comp *comp1
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			comp = (*comp1)(w.Get(e, id1))
		}
	}
	b.StopTimer()
	v := comp.V * comp.V
	_ = v
}

func worldGetUnchecked_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewBuilder(&w, id1)
	query := builder.NewBatchQ(1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var comp *comp1
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			comp = (*comp1)(w.GetUnchecked(e, id1))
		}
	}
	b.StopTimer()
	v := comp.V * comp.V
	_ = v
}

func worldHas_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewBuilder(&w, id1)
	query := builder.NewBatchQ(1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var has bool
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			has = w.Has(e, id1)
		}
	}
	b.StopTimer()
	v := !has
	_ = v
}

func worldHasUnchecked_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewBuilder(&w, id1)
	query := builder.NewBatchQ(1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var has bool
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			has = w.HasUnchecked(e, id1)
		}
	}
	b.StopTimer()
	v := !has
	_ = v
}

func worldAlive_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewBuilder(&w, id1)
	query := builder.NewBatchQ(1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var has bool
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			has = w.Alive(e)
		}
	}
	b.StopTimer()
	v := !has
	_ = v
}

func worldRelation_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[relComp1](&w)
	parent := w.NewEntity()

	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewBuilder(&w, id1).WithRelation(id1)
	query := builder.NewBatchQ(1000, parent)
	for query.Next() {
		entities = append(entities, query.Entity())
	}
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var par ecs.Entity
	rel := w.Relations()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			par = rel.Get(e, id1)
		}
	}
	b.StopTimer()
	v := par.IsZero()
	_ = v
}

func worldRelationUnchecked_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[relComp1](&w)
	parent := w.NewEntity()

	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewBuilder(&w, id1).WithRelation(id1)
	query := builder.NewBatchQ(1000, parent)
	for query.Next() {
		entities = append(entities, query.Entity())
	}
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var par ecs.Entity
	rel := w.Relations()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			par = rel.GetUnchecked(e, id1)
		}
	}
	b.StopTimer()
	v := par.IsZero()
	_ = v
}
