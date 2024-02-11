package ecs

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMask(t *testing.T) {
	filter := All(id(0), id(2), id(4))
	other := All(id(0), id(1), id(2))

	assert.False(t, filter.Matches(&other))

	other = All(id(0), id(1), id(2), id(3), id(4))
	assert.True(t, filter.Matches(&other))
}

func TestQuery(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	rotID := ComponentID[rotation](&w)
	velID := ComponentID[Velocity](&w)
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
		pos := (*Position)(q.Get(posID))
		rot := (*rotation)(q.Get(rotID))
		assert.Equal(t, w.Mask(ent), q.Mask())
		_ = ent
		_ = pos
		_ = rot
		cnt++
	}
	assert.Equal(t, 2, cnt)

	q = w.Query(All(posID))
	assert.Equal(t, 3, q.Count())
	cnt = 0
	entities := []Entity{}
	for q.Next() {
		ent := q.Entity()
		pos := (*Position)(q.Get(posID))
		_ = ent
		_ = pos
		cnt++
		entities = append(entities, ent)
	}
	assert.Equal(t, 3, len(entities))

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

	if isDebug {
		assert.PanicsWithValue(t, "query iteration already finished", func() { q.Next() })
	} else {
		assert.PanicsWithError(t, "runtime error: index out of range [-1]", func() { q.Next() })
	}

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

func TestQueryCached(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	filterPos := w.Cache().Register(All(posID))
	filterPosVel := w.Cache().Register(All(posID, velID))

	q := w.Query(&filterPos)
	assert.Equal(t, 0, q.Count())
	q.Close()

	q = w.Query(&filterPosVel)
	assert.Equal(t, 0, q.Count())
	q.Close()

	NewBuilder(&w, posID).NewBatch(10)
	NewBuilder(&w, velID).NewBatch(10)
	NewBuilder(&w, posID, velID).NewBatch(10)

	q = w.Query(&filterPos)
	assert.Equal(t, 20, q.Count())
	q.Close()

	q = w.Query(&filterPosVel)
	assert.Equal(t, 10, q.Count())
	q.Close()

	NewBuilder(&w, posID).NewBatch(10)

	q = w.Query(&filterPos)
	assert.Equal(t, 30, q.Count())

	for q.Next() {
	}

	filterVel := w.Cache().Register(All(velID))
	q = w.Query(&filterVel)
	assert.Equal(t, 20, q.Count())
	q.Close()
}

func TestQueryCachedRelation(t *testing.T) {
	w := NewWorld()

	relID := ComponentID[testRelationA](&w)

	target1 := w.NewEntity()
	target2 := w.NewEntity()

	relFilter := NewRelationFilter(All(relID), target1)
	cf := w.Cache().Register(&relFilter)

	q := w.Query(&cf)
	assert.Equal(t, 0, q.Count())
	cnt := 0
	for q.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)

	NewBuilder(&w, relID).WithRelation(relID).NewBatch(10, target1)

	q = w.Query(&cf)
	assert.Equal(t, 10, q.Count())
	cnt = 0
	for q.Next() {
		cnt++
	}
	assert.Equal(t, 10, cnt)

	relFilter = NewRelationFilter(All(relID), target2)
	cf = w.Cache().Register(&relFilter)

	q = w.Query(&cf)
	assert.Equal(t, 0, q.Count())
	cnt = 0
	for q.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestQueryEmptyNode(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	relID := ComponentID[testRelationA](&w)

	target := w.NewEntity(posID)

	assert.False(t, w.nodes.Get(2).IsActive)

	builder := NewBuilder(&w, relID).WithRelation(relID)
	child := builder.New(target)

	w.RemoveEntity(child)
	w.RemoveEntity(target)

	assert.True(t, w.nodes.Get(2).HasRelation)
	assert.True(t, w.nodes.Get(2).IsActive)
	assert.Equal(t, 1, int(w.nodes.Get(2).archetypes.Len()))

	w.NewEntity(velID)

	q := w.Query(All())
	assert.Equal(t, 1, q.Count())
	q.Close()

	cf := w.Cache().Register(All())
	q = w.Query(&cf)
	assert.Equal(t, 1, q.Count())
	cnt := 0
	for q.Next() {
		cnt++
	}
	assert.Equal(t, 1, cnt)
}

func TestQueryCount(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
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
	q.Close()

	q = NewBuilder(&w, posID, rotID).NewBatchQ(25)
	assert.Equal(t, 25, q.Count())
	q.Close()
}

type testFilter struct{}

func (f testFilter) Matches(bits *Mask) bool {
	return true
}

func TestQueryInterface(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
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

	filter := testFilter{}
	q := w.Query(&filter)

	cnt := 0
	for q.Next() {
		_ = q.Entity()
		cnt++
	}

	assert.Equal(t, 5, cnt)
	assert.Equal(t, 5, q.Count())
}

func TestQueryStep(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	rotID := ComponentID[rotation](&w)

	_ = w.NewEntity(posID)
	_ = w.NewEntity(posID, rotID)
	_ = w.NewEntity(posID, rotID)
	_ = w.NewEntity(posID, rotID)
	_ = w.NewEntity(posID, rotID)
	_ = w.NewEntity(posID, velID)
	_ = w.NewEntity(posID, velID)
	_ = w.NewEntity(posID, velID)
	_ = w.NewEntity(posID, velID, rotID)
	_ = w.NewEntity(posID, velID, rotID)

	q := w.Query(All(posID))
	cnt := 0
	for q.Next() {
		assert.Equal(t, cnt+1, int(q.Entity().id))
		cnt++
	}
	assert.Equal(t, 10, cnt)

	q = w.Query(All(posID))
	assert.Equal(t, 10, q.Count())

	cnt = 0
	for q.Step(1) {
		cnt++
	}
	assert.Equal(t, 10, cnt)

	q = w.Query(All(posID))
	q.Next()
	assert.Equal(t, Entity{1, 0}, q.Entity())
	q.Step(1)
	assert.Equal(t, Entity{2, 0}, q.Entity())
	q.Step(2)
	assert.Equal(t, Entity{4, 0}, q.Entity())
	q.Step(3)
	assert.Equal(t, Entity{7, 0}, q.Entity())
	q.Step(3)
	assert.Equal(t, Entity{10, 0}, q.Entity())

	assert.True(t, w.IsLocked())

	assert.False(t, q.Step(3))
	assert.False(t, w.IsLocked())

	q = w.Query(All(posID))
	q.Step(1)
	assert.Equal(t, Entity{1, 0}, q.Entity())

	q = w.Query(All(posID))
	q.Step(2)
	assert.Equal(t, Entity{2, 0}, q.Entity())

	q = w.Query(All(posID))
	q.Step(10)
	assert.Equal(t, Entity{10, 0}, q.Entity())

	q = w.Query(All(posID))
	assert.PanicsWithValue(t, "step size must be positive", func() { q.Step(0) })
	q.Step(2)
	assert.PanicsWithValue(t, "step size must be positive", func() { q.Step(0) })

	q = w.Query(All(posID))
	cnt = 0
	for q.Step(2) {
		cnt++
	}
	assert.Equal(t, 5, cnt)

}

func TestQueryClosed(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Add(e0, posID)
	w.Add(e1, posID, rotID)
	w.Add(e2, posID, rotID)

	q := w.Query(All(posID, rotID))
	if isDebug {
		assert.PanicsWithValue(t, "query already iterated or iteration not started yet", func() { q.Entity() })
		assert.PanicsWithValue(t, "query already iterated or iteration not started yet", func() { q.Get(posID) })
	} else {
		assert.PanicsWithError(t, "runtime error: invalid memory address or nil pointer dereference", func() { q.Entity() })
		assert.PanicsWithError(t, "runtime error: invalid memory address or nil pointer dereference", func() { q.Get(posID) })
	}
	q.Close()
	if isDebug {
		assert.PanicsWithValue(t, "query already iterated or iteration not started yet", func() { q.Entity() })
		assert.PanicsWithValue(t, "query already iterated or iteration not started yet", func() { q.Get(posID) })
		assert.PanicsWithValue(t, "query iteration already finished", func() { q.Next() })
	} else {
		assert.PanicsWithError(t, "runtime error: invalid memory address or nil pointer dereference", func() { q.Entity() })
		assert.PanicsWithError(t, "runtime error: invalid memory address or nil pointer dereference", func() { q.Get(posID) })
		assert.PanicsWithError(t, "runtime error: index out of range [-1]", func() { q.Next() })
	}

	assert.True(t, w.locks.locks.IsZero())

	cnt := 0
	q = w.Query(All(posID, rotID))
	assert.False(t, w.locks.locks.IsZero())
	for q.Next() {
		cnt++
	}
	assert.True(t, w.locks.locks.IsZero())
	assert.Equal(t, 2, cnt)

	if isDebug {
		assert.PanicsWithValue(t, "query iteration already finished", func() { q.Next() })
	} else {
		assert.PanicsWithError(t, "runtime error: index out of range [-1]", func() { q.Next() })
	}

	cnt = 0
	excl := All(rotID).Exclusive()
	q = w.Query(&excl)
	assert.False(t, w.locks.locks.IsZero())
	for q.Next() {
		cnt++
	}
	assert.True(t, w.locks.locks.IsZero())
	assert.Equal(t, 0, cnt)

	if isDebug {
		assert.PanicsWithValue(t, "query iteration already finished", func() { q.Next() })
	} else {
		assert.PanicsWithError(t, "runtime error: index out of range [-1]", func() { q.Next() })
	}

	cnt = 0
	excl = All(posID).Exclusive()
	q = w.Query(&excl)
	assert.False(t, w.locks.locks.IsZero())
	for q.Next() {
		cnt++
	}
	assert.True(t, w.locks.locks.IsZero())
	assert.Equal(t, 1, cnt)

	if isDebug {
		assert.PanicsWithValue(t, "query iteration already finished", func() { q.Next() })
	} else {
		assert.PanicsWithError(t, "runtime error: index out of range [-1]", func() { q.Next() })
	}
}

func TestQueryNextArchetype(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[Position](&world)

	var entity Entity
	for i := 0; i < 10; i++ {
		entity = world.NewEntity()
		world.Add(entity, posID)
	}

	query := world.Query(All(posID))

	assert.True(t, query.nextArchetype())
	assert.False(t, query.nextArchetype())

	assert.PanicsWithError(t, "runtime error: index out of range [-1]", func() { query.nextArchetype() })

}

func TestQueryRelations(t *testing.T) {
	world := NewWorld()

	relID := ComponentID[testRelationA](&world)
	rel2ID := ComponentID[testRelationB](&world)
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	targ := world.NewEntity(posID)

	e1 := world.NewEntity(relID, velID)
	world.Relations().Set(e1, relID, targ)

	filter := All(relID)
	query := world.Query(filter)

	for query.Next() {
		targ2 := query.Relation(relID)

		assert.Equal(t, targ, targ2)
		assert.Equal(t, targ, query.relationUnchecked(relID))

		assert.PanicsWithValue(t, "entity has no component ecs.testRelationB, or it is not a relation component",
			func() { query.Relation(rel2ID) })
		assert.PanicsWithValue(t, "entity has no component ecs.Position, or it is not a relation component",
			func() { query.Relation(posID) })
		assert.PanicsWithValue(t, "entity has no component ecs.Velocity, or it is not a relation component",
			func() { query.Relation(velID) })
	}
}

func TestQueryIds(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	_ = world.NewEntity(velID)
	_ = world.NewEntity(velID, posID)

	filter := All()
	query := world.Query(filter)

	query.Next()
	assert.Equal(t, query.Ids(), []ID{id(1)})
	query.Next()
	assert.Equal(t, query.Ids(), []ID{id(0), id(1)})
}

func TestQueryEntityAt(t *testing.T) {
	world := NewWorld()

	rotID := ComponentID[rotation](&world)
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)
	relID := ComponentID[relationComp](&world)

	bPos := NewBuilder(&world, rotID, posID)
	bPosVelRel := NewBuilder(&world, rotID, posID, velID, relID).WithRelation(relID)
	bPosRel := NewBuilder(&world, rotID, posID, relID).WithRelation(relID)
	bVelRel := NewBuilder(&world, rotID, velID, relID).WithRelation(relID)
	bRel := NewBuilder(&world, rotID, relID).WithRelation(relID)

	parent1 := world.NewEntity()
	parent2 := world.NewEntity()

	bPos.NewBatch(10)
	bPosVelRel.NewBatch(10, parent2)
	bPosRel.NewBatch(10, parent1)
	bVelRel.NewBatch(10, parent1)
	bRel.NewBatch(10)

	query := world.Query(All())

	assert.PanicsWithValue(t, "can't get entity at negative index", func() { query.EntityAt(-1) })
	for i := 0; i < 52; i++ {
		assert.Equal(t, Entity{eid(i + 1), 0}, query.EntityAt(i))
	}
	assert.PanicsWithValue(t, "query index out of range: index 52, length 52", func() { query.EntityAt(52) })
	query.Close()

	relFilter := NewRelationFilter(All(relID), parent2)
	query = world.Query(&relFilter)
	for i := 0; i < 10; i++ {
		assert.Equal(t, Entity{eid(i + 13), 0}, query.EntityAt(i))
	}
	assert.PanicsWithValue(t, "query index out of range: index 10, length 10", func() { query.EntityAt(10) })
	query.Close()

	query = world.Query(All(relID))
	for i := 0; i < 40; i++ {
		assert.Equal(t, Entity{eid(i + 13), 0}, query.EntityAt(i))
	}
	assert.PanicsWithValue(t, "query index out of range: index 40, length 40", func() { query.EntityAt(40) })
	query.Close()

	query = world.Batch().ExchangeQ(All(posID), nil, []ID{posID})
	assert.Equal(t, 30, query.Count())

	for i := 0; i < 30; i++ {
		assert.Equal(t, Entity{eid(i + 3), 0}, query.EntityAt(i))
	}
	assert.PanicsWithValue(t, "query index out of range: index 30, length 30", func() { query.EntityAt(30) })
	query.Close()

	f := All(relID)
	filter := world.Cache().Register(f)
	query = world.Query(&filter)
	assert.Equal(t, 40, query.Count())
	for i := 0; i < 40; i++ {
		assert.Equal(t, (i+3)%10, int(query.EntityAt(i).id)%10)
	}
	assert.PanicsWithValue(t, "query index out of range: index 40, length 40", func() { query.EntityAt(40) })
	query.Close()
}

func BenchmarkQueryCreate(b *testing.B) {
	b.StopTimer()

	world := NewWorld()
	posID := ComponentID[Position](&world)
	bPos := NewBuilder(&world, posID)
	bPos.NewBatch(100)

	filter := All(posID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		query.Close()
	}
}

func BenchmarkQueryCreateCached(b *testing.B) {
	b.StopTimer()

	world := NewWorld()
	posID := ComponentID[Position](&world)
	bPos := NewBuilder(&world, posID)
	bPos.NewBatch(100)

	f := All(posID)
	filter := world.Cache().Register(f)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&filter)
		query.Close()
	}
}

func BenchmarkQueryEntityAt_1Arch_1000(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	posID := ComponentID[Position](&w)
	builder := NewBuilder(&w, posID)
	builder.NewBatch(1000)

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = rand.Intn(1000)
	}

	query := w.Query(All(posID))
	b.StartTimer()
	var e Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func BenchmarkQueryEntityAt_1Arch_1000_Registered(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	posID := ComponentID[Position](&w)
	builder := NewBuilder(&w, posID)
	builder.NewBatch(1000)

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = rand.Intn(1000)
	}

	f := All(posID)
	filter := w.Cache().Register(f)
	query := w.Query(&filter)
	b.StartTimer()
	var e Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func BenchmarkQueryEntityAt_10Arch_1000(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	id1 := ComponentID[testStruct0](&w)
	id2 := ComponentID[testStruct1](&w)
	id3 := ComponentID[testStruct2](&w)
	id4 := ComponentID[testStruct3](&w)

	comps := [][]ID{
		{id1}, {id2}, {id3}, {id4},
		{id1, id2}, {id1, id3}, {id1, id4},
		{id2, id3}, {id2, id4}, {id3, id4},
	}
	for _, c := range comps {
		b := NewBuilder(&w, c...)
		b.NewBatch(100)
	}

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = rand.Intn(1000)
	}

	query := w.Query(All())
	b.StartTimer()
	var e Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func BenchmarkQueryEntityAt_10Arch_1000_Registered(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	id1 := ComponentID[testStruct0](&w)
	id2 := ComponentID[testStruct1](&w)
	id3 := ComponentID[testStruct2](&w)
	id4 := ComponentID[testStruct3](&w)

	comps := [][]ID{
		{id1}, {id2}, {id3}, {id4},
		{id1, id2}, {id1, id3}, {id1, id4},
		{id2, id3}, {id2, id4}, {id3, id4},
	}
	for _, c := range comps {
		b := NewBuilder(&w, c...)
		b.NewBatch(100)
	}

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = rand.Intn(1000)
	}

	f := All()
	filter := w.Cache().Register(f)
	query := w.Query(&filter)
	b.StartTimer()
	var e Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}
