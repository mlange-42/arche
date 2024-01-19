package filter_test

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	f "github.com/mlange-42/arche/filter"
	"github.com/stretchr/testify/assert"
)

type position struct {
	X int
	Y int
}

type rotation struct {
	Angle int
}

type relation struct {
	ecs.Relation
}

type TestStruct0 struct{ Val int32 }
type TestStruct1 struct{ Val int32 }
type TestStruct2 struct{ Val int32 }
type TestStruct3 struct{ Val int32 }
type TestStruct4 struct{ Val int32 }
type TestStruct5 struct{ Val int32 }
type TestStruct6 struct{ Val int32 }
type TestStruct7 struct{ Val int32 }
type TestStruct8 struct{ Val int32 }
type TestStruct9 struct{ Val int32 }
type TestStruct10 struct{ Val int32 }

func RegisterAll(w *ecs.World) []ecs.ID {
	_ = TestStruct0{1}
	_ = TestStruct1{1}
	_ = TestStruct2{1}
	_ = TestStruct3{1}
	_ = TestStruct4{1}
	_ = TestStruct5{1}
	_ = TestStruct6{1}
	_ = TestStruct7{1}
	_ = TestStruct8{1}
	_ = TestStruct9{1}
	_ = TestStruct10{1}

	ids := make([]ecs.ID, 11)
	ids[0] = ecs.ComponentID[TestStruct0](w)
	ids[1] = ecs.ComponentID[TestStruct1](w)
	ids[2] = ecs.ComponentID[TestStruct2](w)
	ids[3] = ecs.ComponentID[TestStruct3](w)
	ids[4] = ecs.ComponentID[TestStruct4](w)
	ids[5] = ecs.ComponentID[TestStruct5](w)
	ids[6] = ecs.ComponentID[TestStruct6](w)
	ids[7] = ecs.ComponentID[TestStruct7](w)
	ids[8] = ecs.ComponentID[TestStruct8](w)
	ids[9] = ecs.ComponentID[TestStruct9](w)
	ids[10] = ecs.ComponentID[TestStruct10](w)

	return ids
}

func TestLogicFilters(t *testing.T) {
	w := ecs.NewWorld()
	ids := RegisterAll(&w)

	hasA := ecs.All(ids[0])
	hasB := ecs.All(ids[1])
	hasAll := ecs.All(ids[0], ids[1])
	hasNone := ecs.All()

	var filter ecs.Filter
	filter = hasA
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = f.All(ids[0], ids[1])
	assert.True(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = f.Not(f.All(ids[0], ids[1]))
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = f.Any(ids[0], ids[1])
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = f.Or(hasA, hasB)
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = f.XOr(hasA, hasB)
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = &hasAll
	assert.True(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = f.And(hasA, hasB)
	assert.True(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = f.AnyNot(ids[1])
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = f.AnyNot(ids[0], ids[1])
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = f.NoneOf(ids[0])
	assert.False(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = f.NoneOf(ids[0], ids[1])
	assert.False(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = f.And(hasA, f.AnyNOT(hasB))
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = f.Or(hasAll, f.AnyNot(ids[1]))
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = f.Or(f.Any(ids[0], ids[1]), f.AnyNOT(hasB))
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))
}

func TestLogicFiltersRelation(t *testing.T) {
	w := ecs.NewWorld()

	posID := ecs.ComponentID[position](&w)
	relID := ecs.ComponentID[relation](&w)

	parent := w.NewEntity()

	w.NewEntity(posID)
	w.NewEntity(posID, relID)
	ecs.NewBuilder(&w, posID, relID).WithRelation(relID).New(parent)

	relFilter := ecs.NewRelationFilter(ecs.All(relID), parent)
	not := f.NOT{&relFilter}

	query := w.Query(&not)
	assert.Equal(t, 2, query.Count())
	query.Close()

	posFilter := ecs.All(posID)
	and := f.AND{L: posFilter, R: &relFilter}

	query = w.Query(&and)
	assert.Equal(t, 2, query.Count())
	query.Close()

	or := f.OR{L: posFilter, R: &relFilter}

	query = w.Query(&or)
	assert.Equal(t, 3, query.Count())
	query.Close()

	xor := f.XOR{L: posFilter, R: &relFilter}

	query = w.Query(&xor)
	assert.Equal(t, 1, query.Count())
	query.Close()
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

	q := w.Query(f.All(posID, rotID))
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

	q = w.Query(f.All(posID))
	cnt = 0
	for q.Next() {
		ent := q.Entity()
		pos := (*position)(q.Get(posID))
		_ = ent
		_ = pos
		cnt++
	}
	assert.Equal(t, 3, cnt)

	q = w.Query(f.All(rotID))
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

	q = w.Query(&f.AND{L: f.All(rotID), R: f.AnyNOT(f.All(posID))})

	cnt = 0
	for q.Next() {
		_ = q.Entity()
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func match(f ecs.Filter, m ecs.Mask) bool {
	return f.Matches(&m)
}

func TestInterface(t *testing.T) {
	w := ecs.NewWorld()

	posID := ecs.ComponentID[position](&w)

	f := w.Query(f.All(posID))
	var f2 *ecs.Query = &f
	_ = f2
}

func BenchmarkFilterStackOr(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	ids := RegisterAll(&w)
	mask := f.All(ids[1], ids[2], ids[3], ids[4], ids[5])

	filter := f.OR{f.All(ids[1]), f.All(ids[2])}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(&mask)
	}
}

func BenchmarkFilterHeapOr(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	ids := RegisterAll(&w)
	mask := f.All(ids[1], ids[2], ids[3], ids[4], ids[5])

	filter := f.Or(f.All(ids[1]), f.All(ids[2]))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(&mask)
	}
}

func BenchmarkFilterStack5And(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	ids := RegisterAll(&w)
	mask := f.All(ids[1], ids[2], ids[3], ids[4], ids[5])

	a1 := f.AND{f.All(ids[1]), f.All(ids[2])}
	a2 := f.AND{&a1, f.All(ids[3])}
	a3 := f.AND{&a2, f.All(ids[4])}
	filter := f.AND{&a3, f.All(ids[5])}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(&mask)
	}
}

func BenchmarkFilterHeap5And(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	ids := RegisterAll(&w)
	mask := f.All(ids[1], ids[2], ids[3], ids[4], ids[5])

	filter := f.And(f.All(ids[1]), f.And(f.All(ids[2]), f.And(f.All(ids[3]), f.And(f.All(ids[4]), f.All(ids[5])))))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(&mask)
	}
}
