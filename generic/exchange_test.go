package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestExchange(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewExchange(&w).Adds(T2[testStruct0, testStruct1]()...)
	mutR := NewExchange(&w).Removes(T2[testStruct0, testStruct1]()...)
	mutEx := NewExchange(&w).
		Removes(T2[testStruct0, testStruct1]()...).
		Adds(T2[testStruct2, testStruct3]()...)

	mapper := NewMap2[testStruct0, testStruct1](&w)
	mapper2 := NewMap2[testStruct2, testStruct3](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.NewEntity()
	s0, s1 := mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutR.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1 = mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutEx.Exchange(e)
	s0, s1 = mapper.Get(e)
	assert.Nil(t, s0)
	assert.Nil(t, s1)
	s2, s3 := mapper2.Get(e)
	assert.NotNil(t, s2)
	assert.NotNil(t, s3)
}

func TestExchangeRelation(t *testing.T) {
	w := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&w)
	velID := ecs.ComponentID[Velocity](&w)
	rel1ID := ecs.ComponentID[testRelationA](&w)

	ex1 := NewExchange(&w).
		Adds(T3[Position, Velocity, testRelationA]()...).
		WithRelation(T[testRelationA]())

	ex2 := NewExchange(&w).
		Removes(T[Position]()).
		WithRelation(T[testRelationA]())

	ex3 := NewExchange(&w).
		WithRelation(T[testRelationA]()).
		Adds(T[Position]()).
		Removes(T[Velocity]())

	relMap := NewMap[testRelationA](&w)

	e0 := ex1.NewEntity()
	e1 := ex1.NewEntity()
	e2 := ex1.NewEntity(e0)

	assert.Equal(t, []ecs.ID{posID, velID, rel1ID}, w.Ids(e0))
	assert.Equal(t, []ecs.ID{posID, velID, rel1ID}, w.Ids(e1))
	assert.Equal(t, []ecs.ID{posID, velID, rel1ID}, w.Ids(e2))

	assert.Equal(t, ecs.Entity{}, relMap.GetRelation(e0))
	assert.Equal(t, ecs.Entity{}, relMap.GetRelation(e1))
	assert.Equal(t, e0, relMap.GetRelation(e2))

	e3 := w.NewEntity()
	e4 := w.NewEntity()
	ex1.Add(e3)
	ex1.Add(e4, e0)

	assert.Equal(t, []ecs.ID{posID, velID, rel1ID}, w.Ids(e3))
	assert.Equal(t, []ecs.ID{posID, velID, rel1ID}, w.Ids(e4))
	assert.Equal(t, ecs.Entity{}, relMap.GetRelation(e3))
	assert.Equal(t, e0, relMap.GetRelation(e4))

	ex2.Remove(e3, e1)
	ex2.Remove(e4)

	assert.Equal(t, []ecs.ID{velID, rel1ID}, w.Ids(e3))
	assert.Equal(t, []ecs.ID{velID, rel1ID}, w.Ids(e4))
	assert.Equal(t, e1, relMap.GetRelation(e3))
	assert.Equal(t, e0, relMap.GetRelation(e4))

	ex3.Exchange(e3)
	ex3.Exchange(e4, e2)

	assert.Equal(t, []ecs.ID{posID, rel1ID}, w.Ids(e3))
	assert.Equal(t, []ecs.ID{posID, rel1ID}, w.Ids(e4))
	assert.Equal(t, e1, relMap.GetRelation(e3))
	assert.Equal(t, e2, relMap.GetRelation(e4))
}

func TestExchangeRelationPanics(t *testing.T) {
	w := ecs.NewWorld()

	e0 := w.NewEntity()

	ex1 := NewExchange(&w).
		Adds(T2[Position, testRelationA]()...).
		Removes(T[Velocity]())

	assert.Panics(t, func() {
		_ = ex1.NewEntity(e0)
	})

	assert.Panics(t, func() {
		ex1.Add(e0, e0)
	})

	assert.Panics(t, func() {
		ex1.Remove(e0, e0)
	})

	assert.Panics(t, func() {
		ex1.Exchange(e0, e0)
	})

	assert.Panics(t, func() {
		ex1.ExchangeBatch(ecs.All(), e0)
	})
}

func TestExchangeBatch(t *testing.T) {
	w := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&w)

	ex1 := NewExchange(&w).
		Adds(T3[Position, Velocity, testRelationA]()...).
		WithRelation(T[testRelationA]())

	parent1 := w.NewEntity(posID)

	b := ecs.NewBuilder(&w)
	b.NewBatch(10)

	filter := ecs.All().Exclusive()
	ex1.ExchangeBatch(&filter, parent1)

	b.NewBatch(10)
	ex1.ExchangeBatch(&filter)
}
