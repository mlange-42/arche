package ecs

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"testing"

	"github.com/mlange-42/arche/ecs/stats"
	"github.com/stretchr/testify/assert"
)

func TestWorldConfig(t *testing.T) {
	_ = NewWorld(NewConfig())

	assert.Panics(t, func() { _ = NewWorld(Config{}) })
	assert.Panics(t, func() { _ = NewWorld(Config{}, Config{}) })
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

	assert.Equal(t, 1, w.archetypes.Len())

	w.Add(e0, posID)
	assert.Equal(t, 2, w.archetypes.Len())
	w.Add(e1, posID, rotID)
	assert.Equal(t, 3, w.archetypes.Len())
	w.Add(e2, posID, rotID)
	assert.Equal(t, 3, w.archetypes.Len())

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

	pos1 := (*Position)(w.Get(e1, posID))
	assert.Equal(t, &Position{}, pos1)

	pos1.X = 100
	pos1.Y = 101

	pos0 := (*Position)(w.Get(e0, posID))
	pos1 = (*Position)(w.Get(e1, posID))
	assert.Equal(t, &Position{}, pos0)
	assert.Equal(t, &Position{100, 101}, pos1)

	w.RemoveEntity(e0)
	assert.Panics(t, func() { w.Get(e0, posID) })
	_ = w.NewEntity(posID)
	assert.Panics(t, func() { w.Get(e0, posID) })

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

	for i := 0; i < 100; i++ {
		query := world.Query(All(posID, rotID))
		for query.Next() {
			pos := (*Position)(query.Get(posID))
			_ = pos
		}
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

	assert.Panics(t, func() { world.newEntitiesQuery(0, posID, rotID) })

	query := world.newEntitiesQuery(100, posID, rotID)
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

	query = world.newEntitiesQuery(100, posID, rotID)
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

	query = world.newEntitiesQuery(100, posID, rotID)
	assert.Equal(t, 301, len(events))
	query.Close()
	assert.Equal(t, 401, len(events))
	assert.Equal(t, 101, len(world.entities))

	world.newEntities(100, posID, rotID)
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

	assert.Panics(t, func() { world.newEntitiesWithQuery(0, comps...) })
	assert.Equal(t, 1, len(events))

	query := world.newEntitiesWithQuery(1)
	assert.Equal(t, 1, len(events))
	query.Close()
	assert.Equal(t, 2, len(events))

	query = world.newEntitiesWithQuery(100, comps...)
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

	query = world.newEntitiesWithQuery(100,
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

	world.newEntitiesWith(100, comps...)
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

	query := world.newEntitiesQuery(100, posID)
	assert.Equal(t, 100, query.Count())
	query.Close()
	assert.Equal(t, 100, len(events))

	query = world.newEntitiesQuery(100, posID, rotID)
	assert.Equal(t, 100, query.Count())
	query.Close()
	assert.Equal(t, 200, len(events))

	query = world.Query(All())
	assert.Equal(t, 200, query.Count())
	query.Close()

	filter := All(posID).Exclusive()
	cnt := world.removeEntities(&filter)
	assert.Equal(t, 100, cnt)
	assert.Equal(t, 300, len(events))

	query = world.Query(All())
	assert.Equal(t, 100, query.Count())
	query.Close()

	query = world.Query(All(posID, rotID))
	assert.Equal(t, 100, query.Count())
	query.Close()
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

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Add(e0, posID)
	w.Add(e1, posID, rotID)
	w.Add(e2, posID, rotID)

	stats := w.Stats()
	fmt.Println(stats.String())

	assert.Equal(t, 3, len(stats.Archetypes))
	assert.Equal(t, 3, stats.Entities.Used)

	w.NewEntity(velID)
	stats = w.Stats()
	assert.Equal(t, 4, len(stats.Archetypes))
	assert.Equal(t, 4, stats.Entities.Used)

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

func TestWorldReset(t *testing.T) {
	world := NewWorld()

	world.SetListener(func(e *EntityEvent) {})

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	world.NewEntity(posID)
	world.NewEntity(velID)
	world.NewEntity(posID, velID)
	world.NewEntity(posID, velID)

	AddResource(&world, &rotation{100})

	world.Reset()

	assert.Equal(t, 0, int(world.archetypes.Get(0).Len()))
	assert.Equal(t, 0, world.entityPool.Len())
	assert.Equal(t, 1, len(world.entities))

	query := world.Query(All())
	assert.Equal(t, 0, query.Count())
	query.Close()

	e1 := world.NewEntity(posID)
	e2 := world.NewEntity(velID)
	world.NewEntity(posID, velID)
	world.NewEntity(posID, velID)

	assert.Equal(t, Entity{1, 0}, e1)
	assert.Equal(t, Entity{2, 0}, e2)
}

func TestArchetypeGraph(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)
	rotID := ComponentID[rotation](&world)

	archEmpty := world.archetypes.Get(0)
	arch0 := world.findOrCreateArchetype(archEmpty, []ID{posID, velID}, []ID{})
	archEmpty2 := world.findOrCreateArchetype(arch0, []ID{}, []ID{velID, posID})
	assert.Equal(t, archEmpty, archEmpty2)
	assert.Equal(t, 2, world.archetypes.Len())
	assert.Equal(t, 3, world.graph.Len())

	archEmpty3 := world.findOrCreateArchetype(arch0, []ID{}, []ID{posID, velID})
	assert.Equal(t, archEmpty, archEmpty3)
	assert.Equal(t, 2, world.archetypes.Len())
	assert.Equal(t, 4, world.graph.Len())

	arch01 := world.findOrCreateArchetype(arch0, []ID{velID}, []ID{})
	arch012 := world.findOrCreateArchetype(arch01, []ID{rotID}, []ID{})

	assert.Equal(t, []ID{0, 1, 2}, arch012.Ids)

	archEmpty4 := world.findOrCreateArchetype(arch012, []ID{}, []ID{posID, rotID, velID})
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
	}, events[len(events)-1])

	w.RemoveEntity(e0)
	assert.Equal(t, 2, len(events))
	assert.Equal(t, EntityEvent{
		Entity: e0, AddedRemoved: -1,
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
	assert.Equal(t, 1024, w.archetypes.Len())

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
	printTypeSize[archetypeNode]()
	printTypeSize[layout]()
	printTypeSize[entityPool]()
	printTypeSizeName[componentRegistry[ID]]("componentRegistry")
	printTypeSize[bitPool]()
	printTypeSize[Query]()
	printTypeSize[Resources]()
	printTypeSizeName[reflect.Value]("reflect.Value")
	printTypeSize[EntityEvent]()
	printTypeSize[Cache]()
	printTypeSizeName[idMap[archetype]]("idMap")
}

func printTypeSize[T any]() {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%18s: %5d B\n", tp.Name(), tp.Size())
}

func printTypeSizeName[T any](name string) {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%18s: %5d B\n", name, tp.Size())
}

func BenchmarkEntityAlive(b *testing.B) {
	b.StopTimer()

	world := NewWorld(NewConfig().WithCapacityIncrement(1024))
	posID := ComponentID[Position](&world)

	entities := make([]Entity, 0, 1000)
	q := world.newEntitiesQuery(1000, posID)
	for q.Next() {
		entities = append(entities, q.Entity())
	}

	b.StartTimer()

	var alive bool
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			alive = world.Alive(e)
		}
	}

	_ = alive
}

func BenchmarkGetResource(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	AddResource(&w, &Position{1, 2})
	posID := ResourceID[Position](&w)

	b.StartTimer()

	var res *Position
	for i := 0; i < b.N; i++ {
		res = w.Resources().Get(posID).(*Position)
	}

	_ = res
}

func BenchmarkGetResourceShortcut(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	AddResource(&w, &Position{1, 2})

	b.StartTimer()

	var res *Position
	for i := 0; i < b.N; i++ {
		res = GetResource[Position](&w)
	}

	_ = res
}

func BenchmarkNewEntities_10_000_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := NewWorld(NewConfig().WithCapacityIncrement(1024))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		for i := 0; i < 10000; i++ {
			_ = world.NewEntity(posID, velID)
		}
	}
}

func BenchmarkNewEntitiesBatch_10_000_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := NewWorld(NewConfig().WithCapacityIncrement(1024))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		world.newEntities(10000, posID, velID)
	}
}

func BenchmarkNewEntities_10_000_Reset(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(1024))

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	for i := 0; i < 10000; i++ {
		_ = world.NewEntity(posID, velID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		world.Reset()
		for i := 0; i < 10000; i++ {
			_ = world.NewEntity(posID, velID)
		}
	}
}

func BenchmarkNewEntitiesBatch_10_000_Reset(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(1024))

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	for i := 0; i < 10000; i++ {
		_ = world.NewEntity(posID, velID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		world.Reset()
		world.newEntities(10000, posID, velID)
	}
}

func BenchmarkRemoveEntities_10_000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := NewWorld(NewConfig().WithCapacityIncrement(10000))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		entities := make([]Entity, 10000)
		q := world.newEntitiesQuery(10000, posID, velID)

		cnt := 0
		for q.Next() {
			entities[cnt] = q.Entity()
			cnt++
		}

		b.StartTimer()

		for _, e := range entities {
			world.RemoveEntity(e)
		}
	}
}

func BenchmarkRemoveEntitiesBatch_10_000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := NewWorld(NewConfig().WithCapacityIncrement(10000))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		q := world.newEntitiesQuery(10000, posID, velID)
		q.Close()
		b.StartTimer()
		world.removeEntities(All(posID, velID))
	}
}

func BenchmarkWorldStats_1Arch(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	w.NewEntity()

	b.StartTimer()

	var st *stats.WorldStats
	for i := 0; i < b.N; i++ {
		st = w.Stats()
	}
	_ = st
}

func BenchmarkWorldStats_10Arch(b *testing.B) {
	b.StopTimer()

	w := NewWorld()

	ids := []ID{
		ComponentID[testStruct0](&w),
		ComponentID[testStruct1](&w),
		ComponentID[testStruct2](&w),
		ComponentID[testStruct3](&w),
		ComponentID[testStruct4](&w),
		ComponentID[testStruct5](&w),
		ComponentID[testStruct6](&w),
		ComponentID[testStruct7](&w),
		ComponentID[testStruct8](&w),
		ComponentID[testStruct9](&w),
	}

	for _, id := range ids {
		w.NewEntity(id)
	}

	b.StartTimer()

	var st *stats.WorldStats
	for i := 0; i < b.N; i++ {
		st = w.Stats()
	}
	_ = st
}

func ExampleComponentID() {
	world := NewWorld()
	posID := ComponentID[Position](&world)

	world.NewEntity(posID)
	// Output:
}

func ExampleTypeID() {
	world := NewWorld()
	posID := TypeID(&world, reflect.TypeOf(Position{}))

	world.NewEntity(posID)
	// Output:
}

func ExampleResourceID() {
	world := NewWorld()
	resID := ResourceID[Position](&world)

	world.Resources().Add(resID, &Position{100, 100})
	// Output:
}

func ExampleGetResource() {
	world := NewWorld()

	myRes := Position{100, 100}

	AddResource(&world, &myRes)
	res := GetResource[Position](&world)
	fmt.Println(res)
	// Output: &{100 100}
}

func ExampleAddResource() {
	world := NewWorld()

	myRes := Position{100, 100}
	AddResource(&world, &myRes)

	res := GetResource[Position](&world)
	fmt.Println(res)
	// Output: &{100 100}
}

func ExampleWorld() {
	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	_ = world.NewEntity(posID, velID)
	// Output:
}

func ExampleWorld_NewEntity() {
	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	_ = world.NewEntity(posID, velID)
	// Output:
}

func ExampleWorld_NewEntityWith() {
	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	_ = world.NewEntityWith(
		Component{ID: posID, Comp: &Position{X: 0, Y: 0}},
		Component{ID: velID, Comp: &Velocity{X: 10, Y: 2}},
	)
	// Output:
}

func ExampleWorld_RemoveEntity() {
	world := NewWorld()
	e := world.NewEntity()
	world.RemoveEntity(e)
	// Output:
}

func ExampleWorld_Get() {
	world := NewWorld()
	posID := ComponentID[Position](&world)

	e := world.NewEntity(posID)

	pos := (*Position)(world.Get(e, posID))
	pos.X, pos.Y = 10, 5
	// Output:
}

func ExampleWorld_Has() {
	world := NewWorld()
	posID := ComponentID[Position](&world)

	e := world.NewEntity(posID)

	if world.Has(e, posID) {
		world.Remove(e, posID)
	}
	// Output:
}

func ExampleWorld_Add() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	e := world.NewEntity()

	world.Add(e, posID, velID)
	// Output:
}

func ExampleWorld_Assign() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	e := world.NewEntity()

	world.Assign(e,
		Component{ID: posID, Comp: &Position{X: 0, Y: 0}},
		Component{ID: velID, Comp: &Velocity{X: 10, Y: 2}},
	)
	// Output:
}

func ExampleWorld_Set() {
	world := NewWorld()
	posID := ComponentID[Position](&world)

	e := world.NewEntity(posID)

	world.Set(e, posID, &Position{X: 0, Y: 0})
	// Output:
}

func ExampleWorld_Remove() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	e := world.NewEntity(posID, velID)

	world.Remove(e, posID, velID)
	// Output:
}

func ExampleWorld_Exchange() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	e := world.NewEntity(posID)

	world.Exchange(e, []ID{velID}, []ID{posID})
	// Output:
}

func ExampleWorld_Reset() {
	world := NewWorld()
	_ = world.NewEntity()

	world.Reset()
	// Output:
}

func ExampleWorld_Query() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	filter := All(posID, velID)
	query := world.Query(filter)
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		vel := (*Velocity)(query.Get(velID))
		pos.X += vel.X
		pos.Y += vel.Y
	}
	// Output:
}
func ExampleWorld_Resources() {
	world := NewWorld()

	resID := ResourceID[Position](&world)

	myRes := Position{}
	world.Resources().Add(resID, &myRes)

	res := (world.Resources().Get(resID)).(*Position)
	res.X, res.Y = 10, 5
	// Output:
}

func ExampleWorld_Batch() {
	world := NewWorld()
	world.Batch().NewEntities(10_000)
	// Output:
}

func ExampleWorld_Cache() {
	world := NewWorld()
	posID := ComponentID[Position](&world)

	filter := All(posID)
	cached := world.Cache().Register(filter)
	query := world.Query(&cached)

	for query.Next() {
		// handle entities...
	}
	// Output:
}

func ExampleWorld_SetListener() {
	world := NewWorld()

	listener := func(evt *EntityEvent) {
		fmt.Println(evt)
	}
	world.SetListener(listener)

	world.NewEntity()
	// Output: &{{1 0} {0 0} {0 0} [] [] [] 1}
}

func ExampleWorld_Stats() {
	world := NewWorld()
	stats := world.Stats()
	fmt.Println(stats.Entities.String())
	// Output: Entities -- Used: 0, Recycled: 0, Total: 0, Capacity: 128
}
