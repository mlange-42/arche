package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	relID := ComponentID[testRelationA](&w)

	target := w.NewEntity()

	b1 := NewBuilder(&w, posID, velID, relID)

	e1 := b1.Build()
	assert.True(t, w.Has(e1, posID))
	assert.True(t, w.Has(e1, velID))

	assert.Panics(t, func() { b1.BuildRelation(target) })
	assert.Panics(t, func() { b1.BatchRelation(10, target) })
	assert.Panics(t, func() { b1.QueryRelation(10, target) })

	b1.Batch(10)
	q := b1.Query(10)
	assert.Equal(t, 10, q.Count())
	q.Close()

	b1 = NewBuilderWith(&w, Component{ID: posID, Comp: &Position{}})

	e1 = b1.Build()
	assert.True(t, w.Has(e1, posID))

	assert.Panics(t, func() { b1.BuildRelation(target) })
	assert.Panics(t, func() { b1.BatchRelation(10, target) })
	assert.Panics(t, func() { b1.QueryRelation(10, target) })

	b1.Batch(10)
	q = b1.Query(10)
	assert.Equal(t, 10, q.Count())
	q.Close()

	b1 = NewBuilder(&w, posID, velID, relID).WithRelation(relID)

	e1 = b1.Build()
	e2 := b1.BuildRelation(target)
	assert.Equal(t, target, w.GetRelation(e2, relID))

	b1.BatchRelation(10, target)
	q = b1.QueryRelation(10, target)
	assert.Equal(t, 10, q.Count())
	for q.Next() {
		assert.Equal(t, target, q.Relation(relID))
	}

	b1 = NewBuilderWith(&w,
		Component{ID: posID, Comp: &Position{}},
		Component{ID: relID, Comp: &testRelationA{}},
	).WithRelation(relID)

	e1 = b1.Build()
	e2 = b1.BuildRelation(target)
	assert.Equal(t, target, w.GetRelation(e2, relID))

	b1.BatchRelation(10, target)
	q = b1.QueryRelation(10, target)
	assert.Equal(t, 10, q.Count())
	for q.Next() {
		assert.Equal(t, target, q.Relation(relID))
	}
}
