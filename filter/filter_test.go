package filter

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

type position struct {
	X int
	Y int
}

type rotation struct {
	Angle int
}

func TestLogicFilters(t *testing.T) {

	hasA := ecs.All(0)
	hasB := ecs.All(1)
	hasAll := ecs.All(0, 1)
	hasNone := ecs.All()

	var filter ecs.Filter
	filter = hasA
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = All(0, 1)
	assert.True(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = Not(All(0, 1))
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = Any(0, 1)
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = Or(hasA, hasB)
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = XOr(hasA, hasB)
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = &hasAll
	assert.True(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = And(hasA, hasB)
	assert.True(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = AnyNot(1)
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = AnyNot(0, 1)
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = NoneOf(0)
	assert.False(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = NoneOf(0, 1)
	assert.False(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = And(hasA, AnyNOT(hasB))
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = Or(hasAll, AnyNot(1))
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = Or(Any(0, 1), AnyNOT(hasB))
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	assert.Equal(t, AnyNot(0, 1), Any(0, 1).Not())
	assert.Equal(t, NoneOf(0, 1), All(0, 1).Not())
}

func TestFilter(t *testing.T) {
	w := ecs.NewWorld()

	posID := ecs.ComponentID[position](&w)
	rotID := ecs.ComponentID[rotation](&w)

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

	q := w.Query(All(posID, rotID))
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

	q = w.Query(&AND{L: All(rotID), R: AnyNOT(All(posID))})

	cnt = 0
	for q.Next() {
		_ = q.Entity()
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func match(f ecs.Filter, m Mask) bool {
	return f.Matches(m.BitMask)
}

func TestInterface(t *testing.T) {
	w := ecs.NewWorld()

	posID := ecs.ComponentID[position](&w)

	f := w.Query(All(posID))
	var f2 ecs.QueryIter = &f
	_ = f2
}

func BenchmarkFilterStackOr(b *testing.B) {
	b.StopTimer()
	mask := All(1, 2, 3, 4, 5)

	filter := OR{All(1), All(2)}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(mask.BitMask)
	}
}

func BenchmarkFilterStack5And(b *testing.B) {
	b.StopTimer()
	mask := All(1, 2, 3, 4, 5)

	a1 := AND{All(1), All(2)}
	a2 := AND{&a1, All(3)}
	a3 := AND{&a2, All(4)}
	filter := AND{&a3, All(5)}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(mask.BitMask)
	}
}

func BenchmarkFilterHeapOr(b *testing.B) {
	b.StopTimer()
	mask := All(1, 2, 3, 4, 5)

	filter := Or(All(1), All(2))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(mask.BitMask)
	}
}

func BenchmarkFilterHeap5And(b *testing.B) {
	b.StopTimer()
	mask := All(1, 2, 3, 4, 5)

	filter := And(All(1), And(All(2), And(All(3), And(All(4), All(5)))))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(mask.BitMask)
	}
}
