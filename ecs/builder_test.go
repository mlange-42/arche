package ecs_test

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	w := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&w)
	velID := ecs.ComponentID[Velocity](&w)
	relID := ecs.ComponentID[ChildOf](&w)

	target := w.NewEntity()

	b1 := ecs.NewBuilder(&w, posID, velID, relID)

	e1 := b1.New()
	assert.True(t, w.Has(e1, posID))
	assert.True(t, w.Has(e1, velID))

	e2 := w.NewEntity()
	b1.Add(e2)

	assert.PanicsWithValue(t, "can't set target entity: builder has no relation", func() { b1.New(target) })
	assert.PanicsWithValue(t, "can't set target entity: builder has no relation", func() { b1.NewBatch(10, target) })
	assert.PanicsWithValue(t, "can't set target entity: builder has no relation", func() { b1.NewBatchQ(10, target) })
	assert.PanicsWithValue(t, "can't set target entity: builder has no relation", func() { b1.Add(e1, target) })

	b1.NewBatch(10)
	q := b1.NewBatchQ(10)
	assert.Equal(t, 10, q.Count())
	q.Close()

	b1 = ecs.NewBuilderWith(&w, ecs.Component{ID: posID, Comp: &Position{}})

	e1 = b1.New()
	assert.True(t, w.Has(e1, posID))

	e2 = w.NewEntity()
	b1.Add(e2)
	e2 = w.NewEntity()

	assert.PanicsWithValue(t, "can't set target entity: builder has no relation", func() { b1.New(target) })
	assert.PanicsWithValue(t, "can't set target entity: builder has no relation", func() { b1.NewBatch(10, target) })
	assert.PanicsWithValue(t, "can't set target entity: builder has no relation", func() { b1.NewBatchQ(10, target) })
	assert.PanicsWithValue(t, "can't set target entity: builder has no relation", func() { b1.Add(e2, target) })

	b1.NewBatch(10)
	q = b1.NewBatchQ(10)
	assert.Equal(t, 10, q.Count())
	q.Close()

	b1 = ecs.NewBuilder(&w, posID, velID, relID).WithRelation(relID)

	b1.New()
	e2 = b1.New(target)
	assert.Equal(t, target, w.Relations().Get(e2, relID))

	e2 = w.NewEntity()
	b1.Add(e2, target)
	assert.Equal(t, target, w.Relations().Get(e2, relID))

	b1.NewBatch(10, target)
	q = b1.NewBatchQ(10, target)
	assert.Equal(t, 10, q.Count())
	for q.Next() {
		assert.Equal(t, target, q.Relation(relID))
	}

	b1 = ecs.NewBuilderWith(&w,
		ecs.Component{ID: posID, Comp: &Position{}},
		ecs.Component{ID: relID, Comp: &ChildOf{}},
	).WithRelation(relID)

	b1.New()
	e2 = b1.New(target)
	assert.Equal(t, target, w.Relations().Get(e2, relID))

	e2 = w.NewEntity()
	b1.Add(e2, target)
	assert.Equal(t, target, w.Relations().Get(e2, relID))

	b1.NewBatch(10, target)
	q = b1.NewBatchQ(10, target)
	assert.Equal(t, 10, q.Count())
	for q.Next() {
		assert.Equal(t, target, q.Relation(relID))
	}
}
func ExampleBuilder() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	builder := ecs.NewBuilder(&world, posID, velID)

	_ = builder.New()
	// Output:
}

func ExampleNewBuilder() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	builder := ecs.NewBuilder(&world, posID, velID)

	_ = builder.New()
	// Output:
}

func ExampleNewBuilderWith() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	components := []ecs.Component{
		{ID: posID, Comp: &Position{X: 0, Y: 0}},
		{ID: velID, Comp: &Velocity{X: 10, Y: 2}},
	}

	builder := ecs.NewBuilderWith(&world, components...)

	_ = builder.New()
	// Output:
}

func ExampleBuilder_New() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	builder := ecs.NewBuilder(&world, posID, velID)

	_ = builder.New()
	// Output:
}

func ExampleBuilder_NewBatch() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	builder := ecs.NewBuilder(&world, posID, velID)

	builder.NewBatch(1000)
	// Output:
}

func ExampleBuilder_NewBatchQ() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	builder := ecs.NewBuilder(&world, posID, velID)

	query := builder.NewBatchQ(1000)

	for query.Next() {
		// initialize components of the newly created entities
	}
	// Output:
}

func ExampleBuilder_Add() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	builder := ecs.NewBuilder(&world, posID, velID)

	entity := world.NewEntity()
	builder.Add(entity)
	// Output:
}

func ExampleBuilder_WithRelation() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	childID := ecs.ComponentID[ChildOf](&world)

	target := world.NewEntity()

	builder := ecs.NewBuilder(&world, posID, childID).
		WithRelation(childID)

	builder.New(target)
	// Output:
}
