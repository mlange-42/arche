package ecs

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldConfig(t *testing.T) {
	_ = NewWorld(NewConfig())

	assert.Panics(t, func() { _ = NewWorld(Config{}) })
	assert.Panics(t, func() { _ = NewWorld(Config{}, Config{}) })

	world := NewWorld(
		NewConfig().WithCapacityIncrement(32).WithRelationCapacityIncrement(8),
	)

	relID := ComponentID[testRelationA](&world)

	world.NewEntity()
	world.NewEntity(relID)

	assert.Equal(t, uint32(32), world.nodes.Get(0).capacityIncrement)
	assert.Equal(t, uint32(8), world.nodes.Get(1).capacityIncrement)
}

func TestWorldEntites(t *testing.T) {
	w := NewWorld()

	assert.Equal(t, newEntityGen(1, 0), w.NewEntity())
	assert.Equal(t, newEntityGen(2, 0), w.NewEntity())
	assert.Equal(t, newEntityGen(3, 0), w.NewEntity())

	assert.Equal(t, 0, int(w.entities[0].index))
	assert.Equal(t, 0, int(w.entities[1].index))
	assert.Equal(t, 1, int(w.entities[2].index))
	assert.Equal(t, 2, int(w.entities[3].index))
	w.RemoveEntity(newEntityGen(2, 0))
	assert.False(t, w.Alive(newEntityGen(2, 0)))

	assert.Equal(t, 0, int(w.entities[1].index))
	assert.Equal(t, 1, int(w.entities[3].index))

	assert.Equal(t, newEntityGen(2, 1), w.NewEntity())
	assert.False(t, w.Alive(newEntityGen(2, 0)))
	assert.True(t, w.Alive(newEntityGen(2, 1)))

	assert.Equal(t, 2, int(w.entities[2].index))

	w.RemoveEntity(newEntityGen(3, 0))
	w.RemoveEntity(newEntityGen(2, 1))
	w.RemoveEntity(newEntityGen(1, 0))

	assert.Panics(t, func() { w.RemoveEntity(newEntityGen(3, 0)) })
	assert.Panics(t, func() { w.RemoveEntity(newEntityGen(2, 1)) })
	assert.Panics(t, func() { w.RemoveEntity(newEntityGen(1, 0)) })
}

func TestWorldNewEntites(t *testing.T) {
	w := NewWorld(NewConfig().WithCapacityIncrement(32))

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity(posID, velID, rotID)
	e2 := w.NewEntityWith(
		Component{posID, &Position{1, 2}},
		Component{velID, &Velocity{3, 4}},
		Component{rotID, &rotation{5}},
	)
	e3 := w.NewEntityWith()

	assert.Equal(t, All(), w.Mask(e0))
	assert.Equal(t, All(posID, velID, rotID), w.Mask(e1))
	assert.Equal(t, All(posID, velID, rotID), w.Mask(e2))
	assert.Equal(t, All(), w.Mask(e3))

	pos := (*Position)(w.Get(e2, posID))
	vel := (*Velocity)(w.Get(e2, velID))
	rot := (*rotation)(w.Get(e2, rotID))

	assert.Equal(t, &Position{1, 2}, pos)
	assert.Equal(t, &Velocity{3, 4}, vel)
	assert.Equal(t, &rotation{5}, rot)

	w.RemoveEntity(e0)
	w.RemoveEntity(e1)
	w.RemoveEntity(e2)
	w.RemoveEntity(e3)

	for i := 0; i < 35; i++ {
		e := w.NewEntityWith(
			Component{posID, &Position{i + 1, i + 2}},
			Component{velID, &Velocity{i + 3, i + 4}},
			Component{rotID, &rotation{i + 5}},
		)

		pos := (*Position)(w.Get(e, posID))
		vel := (*Velocity)(w.Get(e, velID))
		rot := (*rotation)(w.Get(e, rotID))

		assert.Equal(t, &Position{i + 1, i + 2}, pos)
		assert.Equal(t, &Velocity{i + 3, i + 4}, vel)
		assert.Equal(t, &rotation{i + 5}, rot)
	}
}

func TestWorldComponents(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	rotID := ComponentID[rotation](&w)

	tPosID := TypeID(&w, reflect.TypeOf(Position{}))
	tRotID := TypeID(&w, reflect.TypeOf(rotation{}))

	assert.Equal(t, posID, tPosID)
	assert.Equal(t, rotID, tRotID)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	assert.Equal(t, int32(1), w.archetypes.Len())

	w.Add(e0, posID)
	assert.Equal(t, int32(2), w.archetypes.Len())
	w.Add(e1, posID, rotID)
	assert.Equal(t, int32(3), w.archetypes.Len())
	w.Add(e2, posID, rotID)
	assert.Equal(t, int32(3), w.archetypes.Len())

	assert.Equal(t, All(posID), w.Mask(e0))
	assert.Equal(t, All(posID, rotID), w.Mask(e1))

	w.Remove(e2, posID)

	maskNone := All()
	maskPos := All(posID)
	maskRot := All(rotID)
	maskPosRot := All(posID, rotID)

	archNone, ok := w.findArchetypeSlow(maskNone)
	assert.True(t, ok)
	archPos, ok := w.findArchetypeSlow(maskPos)
	assert.True(t, ok)
	archRot, ok := w.findArchetypeSlow(maskRot)
	assert.True(t, ok)
	archPosRot, ok := w.findArchetypeSlow(maskPosRot)
	assert.True(t, ok)

	assert.Equal(t, 0, int(archNone.archetype.Len()))
	assert.Equal(t, 1, int(archPos.archetype.Len()))
	assert.Equal(t, 1, int(archRot.archetype.Len()))
	assert.Equal(t, 1, int(archPosRot.archetype.Len()))

	w.Remove(e1, posID)

	assert.Equal(t, 0, int(archNone.archetype.Len()))
	assert.Equal(t, 1, int(archPos.archetype.Len()))
	assert.Equal(t, 2, int(archRot.archetype.Len()))
	assert.Equal(t, 0, int(archPosRot.archetype.Len()))

	w.Add(e0, rotID)
	assert.Equal(t, 0, int(archPos.archetype.Len()))
	assert.Equal(t, 1, int(archPosRot.archetype.Len()))

	w.Remove(e2, rotID)
	// No-op add/remove
	w.Add(e0)
	w.Remove(e0)

	w.RemoveEntity(e0)
	assert.Panics(t, func() { w.Has(newEntityGen(1, 0), posID) })
	assert.Panics(t, func() { w.Get(newEntityGen(1, 0), posID) })
}

func TestWorldLabels(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	labID := ComponentID[label](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Add(e0, posID, labID)
	w.Add(e1, labID)
	w.Add(e1, posID)

	lab0 := (*label)(w.Get(e0, labID))
	assert.NotNil(t, lab0)

	lab1 := (*label)(w.Get(e1, labID))
	assert.NotNil(t, lab1)

	assert.True(t, w.Has(e0, labID))
	assert.True(t, w.Has(e1, labID))
	assert.False(t, w.Has(e2, labID))

	assert.Equal(t, lab0, lab1)
}

func TestWorldExchange(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	rotID := ComponentID[rotation](&w)
	rel1ID := ComponentID[testRelationA](&w)
	rel2ID := ComponentID[testRelationB](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Exchange(e0, []ID{posID}, []ID{})
	w.Exchange(e1, []ID{posID, rotID}, []ID{})
	w.Exchange(e2, []ID{rotID}, []ID{})

	assert.True(t, w.Has(e0, posID))
	assert.False(t, w.Has(e0, rotID))

	assert.True(t, w.Has(e1, posID))
	assert.True(t, w.Has(e1, rotID))

	assert.False(t, w.Has(e2, posID))
	assert.True(t, w.Has(e2, rotID))

	w.Exchange(e2, []ID{posID}, []ID{})
	assert.True(t, w.Has(e2, posID))
	assert.True(t, w.Has(e2, rotID))

	w.Exchange(e0, []ID{rotID}, []ID{posID})
	assert.False(t, w.Has(e0, posID))
	assert.True(t, w.Has(e0, rotID))

	w.Exchange(e1, []ID{velID}, []ID{posID})
	assert.False(t, w.Has(e1, posID))
	assert.True(t, w.Has(e1, rotID))
	assert.True(t, w.Has(e1, velID))

	assert.Panics(t, func() { w.Exchange(e1, []ID{velID}, []ID{}) })
	assert.Panics(t, func() { w.Exchange(e1, []ID{}, []ID{posID}) })

	w.RemoveEntity(e0)
	_ = w.NewEntity()
	assert.Panics(t, func() { w.Exchange(e0, []ID{posID}, []ID{}) })

	target := w.NewEntity()
	e0 = w.NewEntity(rel1ID)

	assert.Panics(t, func() { w.exchange(e0, []ID{rel2ID}, nil, int8(rel2ID), target) })
	assert.Panics(t, func() { w.exchange(e0, []ID{posID}, nil, int8(posID), target) })

	w.Remove(e0, rel1ID)
	assert.Panics(t, func() { w.exchange(e0, []ID{posID}, nil, int8(rel1ID), target) })
}

func TestWorldExchangeBatch(t *testing.T) {
	w := NewWorld()

	events := []EntityEvent{}
	w.SetListener(func(e *EntityEvent) {
		events = append(events, *e)
	})

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	relID := ComponentID[testRelationA](&w)

	target1 := w.NewEntity(velID)
	target2 := w.NewEntity(velID)

	builder := NewBuilder(&w, posID, relID).WithRelation(relID)
	builder.NewBatch(100, target1)
	builder.NewBatch(100, target2)

	cnt := 0
	filter := All(posID, relID)
	w.Batch().Exchange(filter, []ID{velID}, []ID{posID}, func(q Query) {
		assert.Equal(t, 100, q.Count())
		for q.Next() {
			assert.True(t, q.Has(velID))
			assert.True(t, q.Has(relID))
			assert.False(t, q.Has(posID))
		}
		cnt++
	})
	assert.Equal(t, 2, cnt)

	query := w.Query(All(posID))
	assert.Equal(t, 0, query.Count())
	query.Close()

	cnt = 0
	filter2 := NewRelationFilter(All(relID), target1)
	w.Batch().Exchange(&filter2, []ID{posID}, []ID{velID}, func(q Query) {
		assert.Equal(t, 100, q.Count())
		for q.Next() {
			assert.True(t, q.Has(posID))
			assert.True(t, q.Has(relID))
			assert.False(t, q.Has(velID))
			assert.Equal(t, target1, q.Relation(relID))
		}
		cnt++
	})
	assert.Equal(t, 1, cnt)

	query = w.Query(All(posID))
	assert.Equal(t, 100, query.Count())
	query.Close()

	w.Batch().Exchange(All(posID), nil, nil, nil)

	relFilter := NewRelationFilter(All(relID), target2)
	w.Batch().Exchange(&relFilter, nil, []ID{relID}, nil)
	w.Batch().Exchange(All(relID), nil, []ID{relID}, nil)

	w.Batch().RemoveEntities(All(posID))

	assert.Equal(t, 802, len(events))
	assert.Equal(t, 1, events[0].AddedRemoved)
	assert.Equal(t, 1, events[1].AddedRemoved)
	assert.Equal(t, 1, events[2].AddedRemoved)
	assert.Equal(t, 1, events[201].AddedRemoved)

	assert.Equal(t, 0, events[202].AddedRemoved)
	assert.Equal(t, []ID{velID}, events[202].Added)
	assert.Equal(t, []ID{posID}, events[202].Removed)

	filter = All(velID)
	w.Batch().Remove(filter, nil, velID)

	filter = All(velID)
	q := w.Query(filter)
	assert.Equal(t, 0, q.Count())
	q.Close()

	filter = All()
	w.Batch().Add(filter, nil, velID)

	assert.Panics(t, func() { w.Batch().Remove(All(velID), func(q Query) {}, velID) })
}

func TestWorldAssignSet(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()

	assert.Panics(t, func() { w.Assign(e0) })

	w.Assign(e0, Component{posID, &Position{2, 3}})
	pos := (*Position)(w.Get(e0, posID))
	assert.Equal(t, 2, pos.X)
	pos.X = 5

	pos = (*Position)(w.Get(e0, posID))
	assert.Equal(t, 5, pos.X)

	assert.Panics(t, func() { w.Assign(e0, Component{posID, &Position{2, 3}}) })
	assert.Panics(t, func() { _ = (*Position)(w.copyTo(e1, posID, &Position{2, 3})) })

	e2 := w.NewEntity()
	w.Assign(e2,
		Component{posID, &Position{4, 5}},
		Component{velID, &Velocity{1, 2}},
		Component{rotID, &rotation{3}},
	)
	assert.True(t, w.Has(e2, velID))
	assert.True(t, w.Has(e2, rotID))
	assert.True(t, w.Has(e2, posID))

	pos = (*Position)(w.Get(e2, posID))
	rot := (*rotation)(w.Get(e2, rotID))
	vel := (*Velocity)(w.Get(e2, velID))
	assert.Equal(t, &Position{4, 5}, pos)
	assert.Equal(t, &rotation{3}, rot)
	assert.Equal(t, &Velocity{1, 2}, vel)

	_ = (*Position)(w.Set(e2, posID, &Position{7, 8}))
	pos = (*Position)(w.Get(e2, posID))
	assert.Equal(t, 7, pos.X)

	*pos = Position{8, 9}
	pos = (*Position)(w.Get(e2, posID))
	assert.Equal(t, 8, pos.X)

	w.RemoveEntity(e0)
	_ = w.NewEntity()
	assert.Panics(t, func() { w.Assign(e0, Component{posID, &Position{2, 3}}) })
}

func TestWorldGetComponents(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Add(e0, posID, rotID)
	w.Add(e1, posID, rotID)
	w.Add(e2, rotID)

	assert.False(t, w.Has(e2, posID))
	assert.True(t, w.Has(e2, rotID))
	assert.False(t, w.HasUnchecked(e2, posID))
	assert.True(t, w.HasUnchecked(e2, rotID))

	pos1 := (*Position)(w.Get(e1, posID))
	assert.Equal(t, &Position{}, pos1)

	pos1.X = 100
	pos1.Y = 101

	pos0 := (*Position)(w.Get(e0, posID))
	pos1 = (*Position)(w.Get(e1, posID))
	assert.Equal(t, &Position{}, pos0)
	assert.Equal(t, &Position{100, 101}, pos1)

	pos0 = (*Position)(w.GetUnchecked(e0, posID))
	pos1 = (*Position)(w.GetUnchecked(e1, posID))
	assert.Equal(t, &Position{}, pos0)
	assert.Equal(t, &Position{100, 101}, pos1)

	w.RemoveEntity(e0)
	assert.Panics(t, func() { w.Get(e0, posID) })
	assert.Panics(t, func() { w.Mask(e0) })

	_ = w.NewEntity(posID)
	assert.Panics(t, func() { w.Get(e0, posID) })
	assert.Panics(t, func() { w.Mask(e0) })

	pos1 = (*Position)(w.Get(e1, posID))
	assert.Equal(t, &Position{100, 101}, pos1)

	pos2 := (*Position)(w.Get(e2, posID))
	assert.True(t, pos2 == nil)
}

func TestWorldIter(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[Position](&world)
	rotID := ComponentID[rotation](&world)

	for i := 0; i < 1000; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}

	world.NewEntity(rotID)

	for i := 0; i < 10; i++ {
		query := world.Query(All(posID, rotID))
		cnt := 0
		for query.Next() {
			pos := (*Position)(query.Get(posID))
			_ = pos
			cnt++
		}
		assert.Equal(t, 1000, cnt)
		assert.Panics(t, func() { query.Next() })
	}

	for i := 0; i < MaskTotalBits-1; i++ {
		query := world.Query(All(posID, rotID))
		for query.Next() {
			pos := (*Position)(query.Get(posID))
			_ = pos
			break
		}
	}
	query := world.Query(All(posID, rotID))

	assert.Panics(t, func() { world.Query(All(posID, rotID)) })

	query.Close()
	assert.Panics(t, func() { query.Close() })
}

func TestWorldNewEntities(t *testing.T) {
	world := NewWorld(NewConfig().WithCapacityIncrement(16))

	events := []EntityEvent{}
	world.SetListener(func(e *EntityEvent) {
		assert.Equal(t, world.IsLocked(), e.EntityRemoved())
		events = append(events, *e)
	})

	posID := ComponentID[Position](&world)
	rotID := ComponentID[rotation](&world)

	world.NewEntity(posID, rotID)
	assert.Equal(t, 2, len(world.entities))

	assert.Panics(t, func() { world.newEntitiesQuery(0, -1, Entity{}, posID, rotID) })

	query := world.newEntitiesQuery(100, -1, Entity{}, posID, rotID)
	assert.Equal(t, 100, query.Count())
	assert.Equal(t, 102, len(world.entities))
	assert.Equal(t, 1, len(events))

	cnt := 0
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		pos.X = cnt + 1
		pos.Y = cnt + 1
		cnt++
	}
	assert.Equal(t, 100, cnt)
	assert.Equal(t, 101, len(events))

	query = world.Query(All(posID, rotID))
	assert.Equal(t, 101, query.Count())

	cnt = 0
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		assert.Equal(t, cnt, pos.X)
		cnt++
	}

	world.Reset()
	assert.Equal(t, 1, len(world.entities))

	query = world.newEntitiesQuery(100, -1, Entity{}, posID, rotID)
	assert.Equal(t, 100, query.Count())
	assert.Equal(t, 101, len(events))

	entities := make([]Entity, query.Count())
	cnt = 0
	for query.Next() {
		entities[cnt] = query.Entity()
		cnt++
	}
	assert.Equal(t, 100, cnt)
	assert.Equal(t, 201, len(events))

	for _, e := range entities {
		world.RemoveEntity(e)
	}
	assert.Equal(t, 301, len(events))
	assert.Equal(t, 101, len(world.entities))

	query = world.newEntitiesQuery(100, -1, Entity{}, posID, rotID)
	assert.Equal(t, 301, len(events))
	query.Close()
	assert.Equal(t, 401, len(events))
	assert.Equal(t, 101, len(world.entities))

	world.newEntities(100, -1, Entity{}, posID, rotID)
	assert.Equal(t, 501, len(events))
	assert.Equal(t, 201, len(world.entities))
}

func TestWorldNewEntitiesWith(t *testing.T) {
	world := NewWorld(NewConfig().WithCapacityIncrement(16))

	events := []EntityEvent{}
	world.SetListener(func(e *EntityEvent) {
		assert.Equal(t, world.IsLocked(), e.EntityRemoved())
		events = append(events, *e)
	})

	posID := ComponentID[Position](&world)
	rotID := ComponentID[rotation](&world)

	comps := []Component{
		{ID: posID, Comp: &Position{100, 200}},
		{ID: rotID, Comp: &rotation{300}},
	}

	world.NewEntity(posID, rotID)
	assert.Equal(t, 1, len(events))

	assert.Panics(t, func() { world.newEntitiesWithQuery(0, -1, Entity{}, comps...) })
	assert.Equal(t, 1, len(events))

	query := world.newEntitiesWithQuery(1, -1, Entity{})
	assert.Equal(t, 1, len(events))
	query.Close()
	assert.Equal(t, 2, len(events))

	query = world.newEntitiesWithQuery(100, -1, Entity{}, comps...)
	assert.Equal(t, 100, query.Count())
	assert.Equal(t, 2, len(events))

	cnt := 0
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		assert.Equal(t, 100, pos.X)
		assert.Equal(t, 200, pos.Y)
		pos.X = cnt + 1
		pos.Y = cnt + 1
		cnt++
	}
	assert.Equal(t, 100, cnt)
	assert.Equal(t, 102, len(events))

	query = world.Query(All(posID, rotID))
	assert.Equal(t, 101, query.Count())

	cnt = 0
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		assert.Equal(t, cnt, pos.X)
		cnt++
	}

	world.Reset()

	query = world.newEntitiesWithQuery(100, -1, Entity{},
		Component{ID: posID, Comp: &Position{100, 200}},
		Component{ID: rotID, Comp: &rotation{300}},
	)
	assert.Equal(t, 100, query.Count())
	assert.Equal(t, 102, len(events))

	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 100, cnt)
	assert.Equal(t, 202, len(events))

	world.newEntitiesWith(100, -1, Entity{}, comps...)
	assert.Equal(t, 302, len(events))
}

func TestWorldRemoveEntities(t *testing.T) {
	world := NewWorld(NewConfig().WithCapacityIncrement(16))

	events := []EntityEvent{}
	world.SetListener(func(e *EntityEvent) {
		assert.Equal(t, world.IsLocked(), e.EntityRemoved())
		events = append(events, *e)
	})

	posID := ComponentID[Position](&world)
	rotID := ComponentID[rotation](&world)

	query := world.newEntitiesQuery(100, -1, Entity{}, posID)
	assert.Equal(t, 100, query.Count())
	query.Close()
	assert.Equal(t, 100, len(events))

	query = world.newEntitiesQuery(100, -1, Entity{}, posID, rotID)
	assert.Equal(t, 100, query.Count())
	query.Close()
	assert.Equal(t, 200, len(events))

	query = world.Query(All())
	assert.Equal(t, 200, query.Count())
	query.Close()

	filter := All(posID).Exclusive()
	cnt := world.Batch().RemoveEntities(&filter)
	assert.Equal(t, 100, cnt)
	assert.Equal(t, 300, len(events))

	query = world.Query(All())
	assert.Equal(t, 100, query.Count())
	query.Close()

	query = world.Query(All(posID, rotID))
	assert.Equal(t, 100, query.Count())
	query.Close()
}

func TestWorldRelationSet(t *testing.T) {
	world := NewWorld()

	events := []EntityEvent{}
	world.SetListener(func(e *EntityEvent) {
		events = append(events, *e)
	})

	rotID := ComponentID[rotation](&world)
	relID := ComponentID[testRelationA](&world)
	rel2ID := ComponentID[testRelationB](&world)

	targ := world.NewEntity()
	e1 := world.NewEntity(relID, rotID)
	e2 := world.NewEntity(relID, rotID)

	assert.Equal(t, int32(3), world.nodes.Len())
	assert.Equal(t, int32(1), world.nodes.Get(2).archetypes.Len())
	assert.Equal(t, int32(1), world.archetypes.Len())

	assert.Equal(t, Entity{}, world.Relations().Get(e1, relID))
	assert.Equal(t, Entity{}, world.Relations().GetUnchecked(e1, relID))
	world.Relations().Set(e1, relID, targ)

	assert.Equal(t, targ, world.Relations().Get(e1, relID))
	assert.Equal(t, targ, world.Relations().GetUnchecked(e1, relID))
	assert.Equal(t, int32(3), world.nodes.Len())
	assert.Equal(t, int32(2), world.nodes.Get(2).archetypes.Len())
	assert.Equal(t, int32(1), world.archetypes.Len())

	world.Relations().Set(e1, relID, Entity{})

	assert.Panics(t, func() { world.Relations().Get(e1, rotID) })
	assert.Panics(t, func() { world.Relations().Get(e1, rel2ID) })
	assert.Panics(t, func() { world.Relations().Set(e1, rotID, Entity{}) })
	assert.Panics(t, func() { world.Relations().Set(e1, rel2ID, Entity{}) })

	// Should do nothing
	world.Relations().Set(e1, relID, Entity{})

	assert.Equal(t, Entity{}, world.Relations().Get(e1, relID))
	assert.Equal(t, int32(3), world.nodes.Len())
	assert.Equal(t, int32(1), world.archetypes.Len())

	world.Remove(e2, relID)

	assert.Panics(t, func() { world.Relations().Get(e2, relID) })
	assert.Panics(t, func() { world.Relations().Set(e2, relID, Entity{}) })

	assert.Panics(t, func() { world.NewEntity(relID, rel2ID) })
	assert.Panics(t, func() { world.Add(e1, rel2ID) })

	world.RemoveEntity(e1)
	assert.Panics(t, func() { world.Relations().Get(e1, relID) })
	assert.Panics(t, func() { world.Relations().Set(e1, relID, targ) })

	e3 := world.NewEntity(relID, rotID)
	world.RemoveEntity(targ)
	assert.Panics(t, func() { world.Relations().Set(e3, relID, targ) })

	assert.Equal(t, int32(2), world.nodes.Get(2).archetypes.Len())
	assert.True(t, world.nodes.Get(2).archetypes.Get(0).IsActive())
	assert.False(t, world.nodes.Get(2).archetypes.Get(1).IsActive())

	assert.Equal(t, 9, len(events))
}

func TestWorldRelationSetBatch(t *testing.T) {
	world := NewWorld()

	events := []EntityEvent{}
	world.SetListener(func(e *EntityEvent) {
		events = append(events, *e)
	})

	posID := ComponentID[Position](&world)
	rotID := ComponentID[rotation](&world)
	relID := ComponentID[testRelationA](&world)

	targ1 := world.NewEntity(posID)
	targ2 := world.NewEntity(posID)
	targ3 := world.NewEntity(posID)

	builder := NewBuilder(&world, rotID, relID).WithRelation(relID)
	builder.NewBatch(100, targ1)
	builder.NewBatch(100, targ2)
	builder.NewBatch(100, targ3)

	relFilter := NewRelationFilter(All(relID), targ2)
	world.Batch().SetRelation(&relFilter, relID, targ1, func(q Query) {
		assert.Equal(t, 100, q.Count())
		for q.Next() {
			assert.Equal(t, targ1, q.Relation(relID))
		}
	})

	total := 0
	world.Batch().SetRelation(All(relID), relID, targ3, func(q Query) {
		total += q.Count()
		for q.Next() {
			assert.Equal(t, targ3, q.Relation(relID))
		}
	})
	assert.Equal(t, 300, total)

	relFilter = NewRelationFilter(All(relID), targ3)
	world.Batch().SetRelation(&relFilter, relID, Entity{}, func(q Query) {
		assert.Equal(t, 300, q.Count())
		for q.Next() {
			assert.True(t, q.Relation(relID).IsZero())
		}
	})

	relFilter = NewRelationFilter(All(relID), Entity{})
	world.Batch().SetRelation(&relFilter, relID, targ1, nil)

	world.RemoveEntity(targ3)
	assert.Panics(t, func() {
		world.Batch().SetRelation(All(relID), relID, targ3, nil)
	})

	assert.Equal(t, 1304, len(events))

	world.Relations().SetBatch(All(relID), relID, targ1, nil)

	assert.Panics(t, func() {
		world.Batch().SetRelation(All(relID), relID, targ2, func(q Query) {})
	})

	fmt.Println(debugPrintWorld(&world))
}

func TestWorldRelationRemove(t *testing.T) {
	world := NewWorld()

	events := []EntityEvent{}
	world.SetListener(func(e *EntityEvent) { events = append(events, *e) })

	rotID := ComponentID[rotation](&world)
	relID := ComponentID[testRelationA](&world)

	targ := world.NewEntity()
	targ2 := world.NewEntity()
	targ3 := world.NewEntity()

	e1 := world.NewEntity(relID, rotID)
	e2 := world.NewEntity(relID, rotID)

	filter := NewRelationFilter(All(relID), targ)
	world.Cache().Register(&filter)

	assert.Equal(t, int32(3), world.nodes.Len())
	assert.Equal(t, int32(1), world.nodes.Get(2).archetypes.Len())
	assert.Equal(t, int32(1), world.archetypes.Len())

	world.Relations().Set(e1, relID, targ)
	world.Relations().Set(e2, relID, targ)

	assert.Equal(t, int32(2), world.nodes.Get(2).archetypes.Len())
	assert.Equal(t, int32(1), world.archetypes.Len())

	world.RemoveEntity(targ)
	assert.Equal(t, int32(1), world.archetypes.Len())

	world.Relations().Set(e1, relID, Entity{})
	world.Relations().Set(e2, relID, Entity{})

	assert.Equal(t, int32(2), world.nodes.Get(2).archetypes.Len())
	assert.Equal(t, int32(1), world.archetypes.Len())

	world.Relations().Set(e1, relID, targ2)
	world.Relations().Set(e2, relID, targ2)

	assert.Equal(t, int32(2), world.nodes.Get(2).archetypes.Len())
	assert.Equal(t, int32(1), world.archetypes.Len())

	world.Relations().Set(e1, relID, Entity{})
	world.Relations().Set(e2, relID, Entity{})

	_ = world.Stats()

	world.RemoveEntity(targ2)
	assert.Equal(t, int32(1), world.archetypes.Len())

	world.Relations().Set(e1, relID, targ3)
	world.Relations().Set(e2, relID, targ3)

	assert.Equal(t, int32(2), world.nodes.Get(2).archetypes.Len())
	assert.Equal(t, targ3, world.nodes.Get(2).archetypes.Get(1).RelationTarget)
	assert.Equal(t, int32(1), world.archetypes.Len())

	world.Batch().RemoveEntities(All())
	world.Batch().RemoveEntities(All())

	assert.Equal(t, int32(2), world.nodes.Get(2).archetypes.Len())
	assert.True(t, world.nodes.Get(2).archetypes.Get(0).IsActive())
	assert.False(t, world.nodes.Get(2).archetypes.Get(1).IsActive())
}

func TestWorldRelationQuery(t *testing.T) {
	world := NewWorld()

	rotID := ComponentID[rotation](&world)
	relID := ComponentID[testRelationA](&world)

	targ0 := world.NewEntityWith(Component{ID: rotID, Comp: &rotation{Angle: 0}})

	targ1 := world.NewEntityWith(Component{ID: rotID, Comp: &rotation{Angle: 1}})
	targ2 := world.NewEntityWith(Component{ID: rotID, Comp: &rotation{Angle: 2}})
	targ3 := world.NewEntityWith(Component{ID: rotID, Comp: &rotation{Angle: 3}})

	e1 := world.NewEntity(relID)
	world.Relations().Set(e1, relID, targ0)

	for i := 0; i < 4; i++ {
		e1 := world.NewEntity(relID)
		world.Relations().Set(e1, relID, targ1)

		e2 := world.NewEntity(relID)
		world.Relations().Set(e2, relID, targ2)
	}

	world.RemoveEntity(e1)
	world.RemoveEntity(targ0)

	filter := All(relID)
	query := world.Query(filter)
	assert.Equal(t, 8, query.Count())
	cnt := 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 8, cnt)

	filter2 := NewRelationFilter(All(relID), targ1)
	query = world.Query(&filter2)
	assert.Equal(t, 4, query.Count())
	query.Close()

	filter2 = NewRelationFilter(All(relID), targ2)
	query = world.Query(&filter2)
	assert.Equal(t, 4, query.Count())
	query.Close()

	filter2 = NewRelationFilter(All(relID), targ3)
	query = world.Query(&filter2)
	assert.Equal(t, 0, query.Count())
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestWorldRelationQueryCached(t *testing.T) {
	world := NewWorld()

	rotID := ComponentID[rotation](&world)
	relID := ComponentID[testRelationA](&world)

	targ0 := world.NewEntityWith(Component{ID: rotID, Comp: &rotation{Angle: 0}})

	targ1 := world.NewEntityWith(Component{ID: rotID, Comp: &rotation{Angle: 1}})
	targ2 := world.NewEntityWith(Component{ID: rotID, Comp: &rotation{Angle: 2}})
	targ3 := world.NewEntityWith(Component{ID: rotID, Comp: &rotation{Angle: 3}})

	e1 := world.NewEntity(relID)
	world.Relations().Set(e1, relID, targ0)

	for i := 0; i < 4; i++ {
		e1 := world.NewEntity(relID)
		world.Relations().Set(e1, relID, targ1)

		e2 := world.NewEntity(relID)
		world.Relations().Set(e2, relID, targ2)
	}

	world.RemoveEntity(e1)
	world.RemoveEntity(targ0)

	filter := All(relID)
	regFilter := world.Cache().Register(filter)
	query := world.Query(&regFilter)
	assert.Equal(t, 8, query.Count())
	cnt := 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 8, cnt)
	world.Cache().Unregister(&regFilter)

	filter2 := NewRelationFilter(All(relID), targ1)
	regFilter2 := world.Cache().Register(&filter2)
	query = world.Query(&regFilter2)
	assert.Equal(t, 4, query.Count())
	query.Close()
	world.Cache().Unregister(&regFilter2)

	filter2 = NewRelationFilter(All(relID), targ2)
	regFilter2 = world.Cache().Register(&filter2)
	query = world.Query(&regFilter2)
	assert.Equal(t, 4, query.Count())
	query.Close()
	world.Cache().Unregister(&regFilter2)

	filter2 = NewRelationFilter(All(relID), targ3)
	regFilter2 = world.Cache().Register(&filter2)
	query = world.Query(&regFilter2)
	assert.Equal(t, 0, query.Count())
	query.Close()
	world.Cache().Unregister(&regFilter2)
}

func TestWorldRelation(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[Position](&world)
	relID := ComponentID[testRelationA](&world)

	parents := make([]Entity, 25)
	for i := 0; i < 25; i++ {
		parents[i] = world.NewEntityWith(Component{ID: posID, Comp: &Position{X: i, Y: 0}})
	}

	for i := 0; i < 2500; i++ {
		par := parents[i/100]
		e := world.NewEntity(relID)
		world.Relations().Set(e, relID, par)
	}

	parFilter := All(posID)
	parQuery := world.Query(parFilter)
	assert.Equal(t, 25, parQuery.Count())
	for parQuery.Next() {
		targ := (*Position)(parQuery.Get(posID))
		filter := NewRelationFilter(All(relID), parQuery.Entity())
		query := world.Query(&filter)
		assert.Equal(t, 100, query.Count())
		for query.Next() {
			targ.Y++
		}
	}

	parQuery = world.Query(parFilter)
	for parQuery.Next() {
		targ := (*Position)(parQuery.Get(posID))
		assert.Equal(t, 100, targ.Y)
	}
}

func TestWorldRelationCreate(t *testing.T) {
	world := NewWorld()
	world.SetListener(func(e *EntityEvent) {})

	posID := ComponentID[Position](&world)
	relID := ComponentID[testRelationA](&world)

	alive := world.NewEntity()
	dead := world.NewEntity()
	world.RemoveEntity(dead)

	world.newEntities(5, int8(relID), alive, posID, relID)
	assert.Panics(t, func() { world.newEntitiesNoNotify(5, int8(relID), dead, posID, relID) })

	world.newEntityTarget(relID, alive, posID, relID)
	assert.Panics(t, func() { world.newEntityTarget(relID, dead, posID, relID) })

	world.newEntitiesWith(5, int8(relID), alive,
		Component{ID: posID, Comp: &Position{}},
		Component{ID: relID, Comp: &testRelationA{}},
	)
	assert.Panics(t, func() {
		world.newEntitiesWith(5, int8(relID), dead,
			Component{ID: posID, Comp: &Position{}},
			Component{ID: relID, Comp: &testRelationA{}},
		)
	})

	world.newEntityTargetWith(relID, alive,
		Component{ID: posID, Comp: &Position{}},
		Component{ID: relID, Comp: &testRelationA{}},
	)
	assert.Panics(t, func() {
		world.newEntityTargetWith(relID, dead,
			Component{ID: posID, Comp: &Position{}},
			Component{ID: relID, Comp: &testRelationA{}},
		)
	})
}

func TestWorldRelationMove(t *testing.T) {
	world := NewWorld()
	world.SetListener(func(e *EntityEvent) {})

	posID := ComponentID[Position](&world)
	relID := ComponentID[testRelationA](&world)

	target1 := world.NewEntity()
	target2 := world.NewEntity()

	entities := []Entity{}
	for _, trg := range [...]Entity{target1, target2} {
		query := NewBuilder(&world, relID).WithRelation(relID).NewQuery(100, trg)
		for query.Next() {
			entities = append(entities, query.Entity())
		}
	}

	for _, e := range entities {
		world.Add(e, posID)
	}

	for i, e := range entities {
		trg := world.Relations().Get(e, relID)

		if i < 100 {
			assert.Equal(t, target1, trg)
		} else {
			assert.Equal(t, target2, trg)
		}
	}

	for _, e := range entities {
		world.Remove(e, posID)
	}

	for i, e := range entities {
		trg := world.Relations().Get(e, relID)

		if i < 100 {
			assert.Equal(t, target1, trg)
		} else {
			assert.Equal(t, target2, trg)
		}
	}

	for _, e := range entities {
		world.Remove(e, relID)
	}
}

func TestWorldLock(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[Position](&world)
	rotID := ComponentID[rotation](&world)

	var entity Entity
	for i := 0; i < 100; i++ {
		entity = world.NewEntity()
		world.Add(entity, posID)
	}

	query1 := world.Query(All(posID))
	query2 := world.Query(All(posID))
	assert.True(t, world.IsLocked())
	query1.Close()
	assert.True(t, world.IsLocked())
	query2.Close()
	assert.False(t, world.IsLocked())

	query1 = world.Query(All(posID))

	assert.Panics(t, func() { world.NewEntity() })
	assert.Panics(t, func() { world.RemoveEntity(entity) })
	assert.Panics(t, func() { world.Add(entity, rotID) })
	assert.Panics(t, func() { world.Remove(entity, posID) })
}

func TestWorldStats(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	rotID := ComponentID[rotation](&w)
	velID := ComponentID[Velocity](&w)
	relID := ComponentID[testRelationA](&w)

	_ = w.Stats()

	e0 := w.NewEntity()
	e1 := w.NewEntity(posID, rotID)
	w.NewEntity(posID, rotID)

	stats := w.Stats()
	_ = stats.Nodes[1].String()
	s := stats.Nodes[2].String()
	fmt.Println(s)

	assert.Equal(t, 3, len(stats.Nodes))
	assert.Equal(t, 3, stats.Entities.Used)
	_ = w.Stats()

	w.Add(e0, posID)

	w.NewEntity(velID)
	stats = w.Stats()
	assert.Equal(t, 4, len(stats.Nodes))
	assert.Equal(t, 4, stats.Entities.Used)

	stats = w.Stats()
	assert.Equal(t, 4, len(stats.Nodes))

	builder := NewBuilder(&w, relID).WithRelation(relID)

	builder.NewBatch(10)
	builder.NewBatch(10, e0)
	_ = w.Stats()

	builder.NewBatch(5, e1)

	stats = w.Stats()
	assert.Equal(t, 29, stats.Entities.Used)
	assert.Equal(t, 5, len(stats.Nodes))

	node := &stats.Nodes[4]
	assert.Equal(t, 3, len(node.Archetypes))
	assert.Equal(t, 10, node.Archetypes[0].Size)
	assert.Equal(t, 10, node.Archetypes[1].Size)
	assert.Equal(t, 5, node.Archetypes[2].Size)

	f := All(relID).Exclusive()
	w.Batch().RemoveEntities(&f)
	w.RemoveEntity(e0)
	stats = w.Stats()

	s = stats.String()
	fmt.Println(s)
}

func TestWorldResources(t *testing.T) {
	w := NewWorld()

	posID := ResourceID[Position](&w)
	rotID := ResourceID[rotation](&w)

	assert.False(t, w.Resources().Has(posID))
	assert.Nil(t, w.Resources().Get(posID))

	AddResource(&w, &Position{1, 2})

	assert.True(t, w.Resources().Has(posID))
	pos, ok := w.Resources().Get(posID).(*Position)

	assert.True(t, ok)
	assert.Equal(t, Position{1, 2}, *pos)

	assert.Panics(t, func() { w.Resources().Add(posID, &Position{1, 2}) })

	pos = GetResource[Position](&w)
	assert.Equal(t, Position{1, 2}, *pos)

	w.Resources().Add(rotID, &rotation{5})
	assert.True(t, w.Resources().Has(rotID))
	w.Resources().Remove(rotID)
	assert.False(t, w.Resources().Has(rotID))
	assert.Panics(t, func() { w.Resources().Remove(rotID) })
}

func TestWorldComponentType(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	rotID := ComponentID[rotation](&w)

	tp, ok := w.ComponentType(posID)
	assert.True(t, ok)
	assert.Equal(t, reflect.TypeOf(Position{}), tp)

	tp, ok = w.ComponentType(rotID)
	assert.True(t, ok)
	assert.Equal(t, reflect.TypeOf(rotation{}), tp)

	_, ok = w.ComponentType(2)
	assert.False(t, ok)
}

func TestRegisterComponents(t *testing.T) {
	world := NewWorld()

	ComponentID[Position](&world)

	assert.Equal(t, ID(0), ComponentID[Position](&world))
	assert.Equal(t, ID(1), ComponentID[rotation](&world))
}

func TestWorldBatchRemove(t *testing.T) {
	world := NewWorld()

	rotID := ComponentID[rotation](&world)
	relID := ComponentID[testRelationA](&world)

	target1 := world.NewEntity()
	target2 := world.NewEntity()
	target3 := world.NewEntity()

	builder := NewBuilder(&world, rotID, relID).WithRelation(relID)

	builder.NewBatch(10, target1)
	builder.NewBatch(10, target2)
	builder.NewBatch(10, target3)

	filter := All(rotID).Exclusive()
	filter2 := world.Cache().Register(&filter)
	world.Batch().RemoveEntities(&filter)
	world.Cache().Unregister(&filter2)

	relFilter := NewRelationFilter(All(rotID, relID), target1)
	world.Batch().RemoveEntities(&relFilter)

	relFilter = NewRelationFilter(All(rotID, relID), target2)
	world.Batch().RemoveEntities(&relFilter)

	filter = All().Exclusive()
	world.Batch().RemoveEntities(&filter)

	relFilter = NewRelationFilter(All(rotID, relID), target3)
	world.Batch().RemoveEntities(&relFilter)

	query := world.Query(All())
	assert.Equal(t, 0, query.Count())
	query.Close()
}

func TestWorldReset(t *testing.T) {
	world := NewWorld()

	world.SetListener(func(e *EntityEvent) {})
	AddResource(&world, &rotation{100})

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)
	relID := ComponentID[testRelationA](&world)

	target1 := world.NewEntity()
	target2 := world.NewEntity()

	world.NewEntity(velID)
	world.NewEntity(posID, velID)
	world.NewEntity(posID, velID)
	e1 := world.NewEntity(posID, relID)
	e2 := world.NewEntity(posID, relID)

	world.Relations().Set(e1, relID, target1)
	world.Relations().Set(e2, relID, target2)

	world.RemoveEntity(e1)
	world.RemoveEntity(target1)

	world.Reset()

	assert.Equal(t, 0, int(world.archetypes.Get(0).Len()))
	assert.Equal(t, 0, int(world.archetypes.Get(1).Len()))
	assert.Equal(t, 0, world.entityPool.Len())
	assert.Equal(t, 1, len(world.entities))

	query := world.Query(All())
	assert.Equal(t, 0, query.Count())
	query.Close()

	e1 = world.NewEntity(posID)
	e2 = world.NewEntity(velID)
	world.NewEntity(posID, velID)
	world.NewEntity(posID, velID)

	assert.Equal(t, Entity{1, 0}, e1)
	assert.Equal(t, Entity{2, 0}, e2)

	query = world.Query(All())
	assert.Equal(t, 4, query.Count())
	query.Close()
}

func TestArchetypeGraph(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)
	rotID := ComponentID[rotation](&world)

	archEmpty := world.archetypes.Get(0)
	arch0 := world.findOrCreateArchetype(archEmpty, []ID{posID, velID}, []ID{}, Entity{})
	archEmpty2 := world.findOrCreateArchetype(arch0, []ID{}, []ID{velID, posID}, Entity{})
	assert.Equal(t, archEmpty, archEmpty2)
	assert.Equal(t, int32(2), world.archetypes.Len())
	assert.Equal(t, int32(3), world.nodes.Len())

	archEmpty3 := world.findOrCreateArchetype(arch0, []ID{}, []ID{posID, velID}, Entity{})
	assert.Equal(t, archEmpty, archEmpty3)
	assert.Equal(t, int32(2), world.archetypes.Len())
	assert.Equal(t, int32(4), world.nodes.Len())

	arch01 := world.findOrCreateArchetype(arch0, []ID{velID}, []ID{}, Entity{})
	arch012 := world.findOrCreateArchetype(arch01, []ID{rotID}, []ID{}, Entity{})

	assert.Equal(t, []ID{0, 1, 2}, arch012.node.Ids)

	archEmpty4 := world.findOrCreateArchetype(arch012, []ID{}, []ID{posID, rotID, velID}, Entity{})
	assert.Equal(t, archEmpty, archEmpty4)
}

func TestWorldListener(t *testing.T) {
	events := []EntityEvent{}
	listen := func(e *EntityEvent) {
		events = append(events, *e)
	}

	w := NewWorld()

	w.SetListener(listen)

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	assert.Equal(t, 1, len(events))
	assert.Equal(t, EntityEvent{
		Entity: e0, AddedRemoved: 1,
		Current: []ID{},
	}, events[len(events)-1])

	w.RemoveEntity(e0)
	assert.Equal(t, 2, len(events))
	assert.Equal(t, EntityEvent{
		Entity: e0, AddedRemoved: -1,
		Removed: []ID{},
	}, events[len(events)-1])

	e0 = w.NewEntity(posID, velID)
	assert.Equal(t, 3, len(events))
	assert.Equal(t, EntityEvent{
		Entity:       e0,
		NewMask:      All(posID, velID),
		Added:        []ID{posID, velID},
		Current:      []ID{posID, velID},
		AddedRemoved: 1,
	}, events[len(events)-1])

	w.RemoveEntity(e0)
	assert.Equal(t, 4, len(events))
	assert.Equal(t, EntityEvent{
		Entity:       e0,
		OldMask:      All(posID, velID),
		NewMask:      Mask{},
		Removed:      []ID{posID, velID},
		Current:      nil,
		AddedRemoved: -1,
	}, events[len(events)-1])

	e0 = w.NewEntityWith(Component{posID, &Position{}}, Component{velID, &Velocity{}})
	assert.Equal(t, 5, len(events))
	assert.Equal(t, EntityEvent{
		Entity:       e0,
		NewMask:      All(posID, velID),
		Added:        []ID{posID, velID},
		Current:      []ID{posID, velID},
		AddedRemoved: 1,
	}, events[len(events)-1])

	w.Add(e0, rotID)
	assert.Equal(t, 6, len(events))
	assert.Equal(t, EntityEvent{
		Entity:       e0,
		OldMask:      All(posID, velID),
		NewMask:      All(posID, velID, rotID),
		Added:        []ID{rotID},
		Current:      []ID{posID, velID, rotID},
		AddedRemoved: 0,
	}, events[len(events)-1])

	w.Remove(e0, posID)
	assert.Equal(t, 7, len(events))
	assert.Equal(t, EntityEvent{
		Entity:       e0,
		OldMask:      All(posID, velID, rotID),
		NewMask:      All(velID, rotID),
		Removed:      []ID{posID},
		Current:      []ID{velID, rotID},
		AddedRemoved: 0,
	}, events[len(events)-1])

}

type withSlice struct {
	Slice []int
}

func TestWorldRemoveGC(t *testing.T) {
	w := NewWorld()
	compID := ComponentID[withSlice](&w)

	runtime.GC()
	mem1 := runtime.MemStats{}
	mem2 := runtime.MemStats{}
	runtime.ReadMemStats(&mem1)

	entities := []Entity{}
	for i := 0; i < 100; i++ {
		e := w.NewEntity(compID)
		ws := (*withSlice)(w.Get(e, compID))
		ws.Slice = make([]int, 10000)
		entities = append(entities, e)
	}

	runtime.GC()
	runtime.ReadMemStats(&mem2)
	heap := int(mem2.HeapInuse - mem1.HeapInuse)
	assert.Greater(t, heap, 8000000)
	assert.Less(t, heap, 10000000)

	rand.Shuffle(len(entities), func(i, j int) {
		entities[i], entities[j] = entities[j], entities[i]
	})

	for _, e := range entities {
		w.RemoveEntity(e)
	}

	runtime.GC()
	runtime.ReadMemStats(&mem2)
	heap = int(mem2.HeapInuse - mem1.HeapInuse)
	assert.Less(t, heap, 800000)

	w.NewEntity(compID)
}

func TestWorldResetGC(t *testing.T) {
	w := NewWorld()
	compID := ComponentID[withSlice](&w)

	runtime.GC()
	mem1 := runtime.MemStats{}
	mem2 := runtime.MemStats{}
	runtime.ReadMemStats(&mem1)

	for i := 0; i < 100; i++ {
		e := w.NewEntity(compID)
		ws := (*withSlice)(w.Get(e, compID))
		ws.Slice = make([]int, 10000)
	}

	runtime.ReadMemStats(&mem2)
	heap := int(mem2.HeapInuse - mem1.HeapInuse)
	assert.Greater(t, heap, 8000000)
	assert.Less(t, heap, 10000000)

	runtime.GC()
	runtime.ReadMemStats(&mem2)
	heap = int(mem2.HeapInuse - mem1.HeapInuse)
	assert.Greater(t, heap, 8000000)
	assert.Less(t, heap, 10000000)

	w.Reset()

	runtime.GC()
	runtime.ReadMemStats(&mem2)
	heap = int(mem2.HeapInuse - mem1.HeapInuse)
	assert.Less(t, heap, 800000)

	w.NewEntity(compID)
}

func Test1000Archetypes(t *testing.T) {
	_ = testStruct0{1}
	_ = testStruct1{1}
	_ = testStruct2{1}
	_ = testStruct3{1}
	_ = testStruct4{1}
	_ = testStruct5{1}
	_ = testStruct6{1}
	_ = testStruct7{1}
	_ = testStruct8{1}
	_ = testStruct9{1}
	_ = testStruct10{1}

	w := NewWorld()

	ids := [10]ID{}
	ids[0] = ComponentID[testStruct0](&w)
	ids[1] = ComponentID[testStruct1](&w)
	ids[2] = ComponentID[testStruct2](&w)
	ids[3] = ComponentID[testStruct3](&w)
	ids[4] = ComponentID[testStruct4](&w)
	ids[5] = ComponentID[testStruct5](&w)
	ids[6] = ComponentID[testStruct6](&w)
	ids[7] = ComponentID[testStruct7](&w)
	ids[8] = ComponentID[testStruct8](&w)
	ids[9] = ComponentID[testStruct9](&w)

	for i := 0; i < 1024; i++ {
		mask := Mask{uint64(i), 0}
		add := make([]ID, 0, 10)
		for j := 0; j < 10; j++ {
			id := ID(j)
			if mask.Get(id) {
				add = append(add, id)
			}
		}
		entity := w.NewEntity()
		w.Add(entity, add...)
	}
	assert.Equal(t, int32(1024), w.archetypes.Len())

	cnt := 0
	query := w.Query(All(0, 7))
	for query.Next() {
		cnt++
	}

	assert.Equal(t, 256, cnt)
}

func TestTypeSizes(t *testing.T) {
	printTypeSize[Entity]()
	printTypeSize[entityIndex]()
	printTypeSize[Mask]()
	printTypeSize[World]()
	printTypeSizeName[pagedSlice[archetype]]("pagedArr32")
	printTypeSize[archetype]()
	printTypeSize[archetypeAccess]()
	printTypeSize[archNode]()
	printTypeSize[nodeData]()
	printTypeSize[layout]()
	printTypeSize[entityPool]()
	printTypeSizeName[componentRegistry[ID]]("componentRegistry")
	printTypeSize[bitPool]()
	printTypeSize[Query]()
	printTypeSize[Resources]()
	printTypeSizeName[reflect.Value]("reflect.Value")
	printTypeSize[EntityEvent]()
	printTypeSize[Cache]()
	printTypeSizeName[idMap[uint32]]("idMap")
}

func printTypeSize[T any]() {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%18s: %5d B\n", tp.Name(), tp.Size())
}

func printTypeSizeName[T any](name string) {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%18s: %5d B\n", name, tp.Size())
}
