package ecs

import (
	"testing"

	"github.com/mlange-42/arche/internal/base"
	"github.com/stretchr/testify/assert"
)

func TestMask(t *testing.T) {
	filter := NewMask(0, 2, 4)
	other := base.NewBitMask(0, 1, 2)

	assert.False(t, filter.Matches(other))

	other = base.NewBitMask(0, 1, 2, 3, 4)
	assert.True(t, filter.Matches(other))
}

func TestQuery(t *testing.T) {
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
	w.Add(e3, rotID)
	w.Add(e4, rotID)

	q := w.Filter(All(posID, rotID))
	cnt := 0
	for q.Next() {
		ent := q.Entity()
		pos := (*position)(q.Get(posID))
		rot := (*rotation)(q.Get(rotID))
		_ = ent
		_ = pos
		_ = rot
		cnt++
	}
	assert.Equal(t, 2, cnt)

	q = w.Filter(All(posID))
	cnt = 0
	for q.Next() {
		ent := q.Entity()
		pos := (*position)(q.Get(posID))
		_ = ent
		_ = pos
		cnt++
	}
	assert.Equal(t, 3, cnt)

	q = w.Filter(All(rotID))
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

	q = w.Filter(All(rotID).Not(posID))

	cnt = 0
	for q.Next() {
		_ = q.Entity()
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func TestInterface(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)

	q := w.Filter(All(posID, rotID))
	var q2 EntityIter = &q
	_ = q2
}
