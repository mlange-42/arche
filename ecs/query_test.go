package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMask(t *testing.T) {
	filter := All(0, 2, 4)
	other := All(0, 1, 2)

	assert.False(t, filter.Matches(other))

	other = All(0, 1, 2, 3, 4)
	assert.True(t, filter.Matches(other))
}

func TestQuery(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)
	velID := ComponentID[velocity](&w)
	s0ID := ComponentID[testStruct0](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	e4 := w.NewEntity()

	w.Add(e0, posID)
	w.Add(e1, posID, rotID)
	w.Add(e2, posID, rotID)
	w.Add(e3, rotID, velID)
	w.Add(e4, rotID)

	q := w.Query(All(posID, rotID))
	cnt := 0
	for q.Next() {
		ent := q.Entity()
		pos := (*position)(q.Get(posID))
		rot := (*rotation)(q.Get(rotID))
		assert.Equal(t, w.Mask(ent), q.Mask())
		_ = ent
		_ = pos
		_ = rot
		cnt++
	}
	assert.Equal(t, 2, cnt)

	q = w.Query(All(posID))
	cnt = 0
	for q.Next() {
		ent := q.Entity()
		pos := (*position)(q.Get(posID))
		_ = ent
		_ = pos
		cnt++
	}
	assert.Equal(t, 3, cnt)

	q = w.Query(All(rotID))
	cnt = 0
	for q.Next() {
		ent := q.Entity()
		rot := (*rotation)(q.Get(rotID))
		_ = ent
		_ = rot
		hasPos := q.Has(posID)
		_ = hasPos
		cnt++
	}
	assert.Equal(t, 4, cnt)

	assert.Panics(t, func() { q.Next() })

	filter := All(rotID).Without(posID)
	q = w.Query(&filter)

	cnt = 0
	for q.Next() {
		_ = q.Entity()
		cnt++
	}
	assert.Equal(t, 2, cnt)

	filter = All(rotID).Without(posID, velID)
	q = w.Query(&filter)

	cnt = 0
	for q.Next() {
		_ = q.Entity()
		cnt++
	}
	assert.Equal(t, 1, cnt)

	filter = All(rotID, s0ID).Without()
	q = w.Query(&filter)

	cnt = 0
	for q.Next() {
		_ = q.Entity()
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestQueryCount(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	e4 := w.NewEntity()

	w.Add(e0, posID)
	w.Add(e1, posID, rotID)
	w.Add(e2, posID, rotID)
	w.Add(e3, posID, rotID)
	w.Add(e4, rotID)

	q := w.Query(All(posID))
	assert.Equal(t, 4, q.Count())
	assert.Equal(t, 4, q.Count())
}

func TestQueryClosed(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Add(e0, posID)
	w.Add(e1, posID, rotID)
	w.Add(e2, posID, rotID)

	q := w.Query(All(posID, rotID))
	q.Close()
	assert.Panics(t, func() { q.Next() })
	assert.Panics(t, func() { q.Get(posID) })
}
