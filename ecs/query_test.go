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

	relFilter := RelationFilter(All(relID), target1)
	cf := w.Cache().Register(relFilter)

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

	relFilter = RelationFilter(All(relID), target2)
	cf = w.Cache().Register(relFilter)

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

	q = NewBuilder(&w, posID, rotID).NewQuery(25)
	assert.Equal(t, 25, q.Count())
	q.Close()
}

type testFilter struct{}

func (f testFilter) Matches(bits Mask) bool {
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

	q := w.Query(testFilter{})

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
	assert.Panics(t, func() { q.Step(0) })
	q.Step(2)
	assert.Panics(t, func() { q.Step(0) })

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
	assert.Panics(t, func() { q.Entity() })
	assert.Panics(t, func() { q.Get(posID) })

	q.Close()
	assert.Panics(t, func() { q.Entity() })
	assert.Panics(t, func() { q.Get(posID) })
	assert.Panics(t, func() { q.Next() })
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
	assert.Panics(t, func() { query.nextArchetype() })
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

		assert.Panics(t, func() { query.Relation(rel2ID) })
		assert.Panics(t, func() { query.Relation(posID) })
		assert.Panics(t, func() { query.Relation(velID) })
	}
}
