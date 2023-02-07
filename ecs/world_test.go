package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldEntites(t *testing.T) {
	w := NewWorld()

	assert.Equal(t, Entity{1, 0}, w.NewEntity())
	assert.Equal(t, Entity{2, 0}, w.NewEntity())
	assert.Equal(t, Entity{3, 0}, w.NewEntity())

	assert.Equal(t, 0, int(w.entities[0].index))
	assert.Equal(t, 0, int(w.entities[1].index))
	assert.Equal(t, 1, int(w.entities[2].index))
	assert.Equal(t, 2, int(w.entities[3].index))
	w.RemEntity(Entity{2, 0})
	assert.False(t, w.Alive(Entity{2, 0}))

	assert.Equal(t, 0, int(w.entities[1].index))
	assert.Equal(t, 1, int(w.entities[3].index))

	assert.Equal(t, Entity{2, 1}, w.NewEntity())
	assert.False(t, w.Alive(Entity{2, 0}))
	assert.True(t, w.Alive(Entity{2, 1}))

	assert.Equal(t, 2, int(w.entities[2].index))

	w.RemEntity(Entity{3, 0})
	w.RemEntity(Entity{2, 1})
	w.RemEntity(Entity{1, 0})

	assert.Panics(t, func() { w.RemEntity(Entity{3, 0}) })
	assert.Panics(t, func() { w.RemEntity(Entity{2, 1}) })
	assert.Panics(t, func() { w.RemEntity(Entity{1, 0}) })
}

func TestWorldComponents(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	assert.Equal(t, 1, len(w.archetypes))

	w.Add(e0, posID)
	assert.Equal(t, 2, len(w.archetypes))
	w.Add(e1, posID, rotID)
	assert.Equal(t, 3, len(w.archetypes))
	w.Add(e2, posID, rotID)
	assert.Equal(t, 3, len(w.archetypes))

	w.Remove(e2, posID)

	maskNone := NewMask()
	maskPos := NewMask(posID)
	maskRot := NewMask(rotID)
	maskPosRot := NewMask(posID, rotID)

	archNone, ok := w.findArchetype(maskNone)
	assert.True(t, ok)
	archPos, ok := w.findArchetype(maskPos)
	assert.True(t, ok)
	archRot, ok := w.findArchetype(maskRot)
	assert.True(t, ok)
	archPosRot, ok := w.findArchetype(maskPosRot)
	assert.True(t, ok)

	assert.Equal(t, 0, int(w.archetypes[archNone].Len()))
	assert.Equal(t, 1, int(w.archetypes[archPos].Len()))
	assert.Equal(t, 1, int(w.archetypes[archRot].Len()))
	assert.Equal(t, 1, int(w.archetypes[archPosRot].Len()))

	w.Remove(e1, posID)

	assert.Equal(t, 0, int(w.archetypes[archNone].Len()))
	assert.Equal(t, 1, int(w.archetypes[archPos].Len()))
	assert.Equal(t, 2, int(w.archetypes[archRot].Len()))
	assert.Equal(t, 0, int(w.archetypes[archPosRot].Len()))

	w.Add(e0, rotID)
	assert.Equal(t, 0, int(w.archetypes[archPos].Len()))
	assert.Equal(t, 1, int(w.archetypes[archPosRot].Len()))

	w.Remove(e2, rotID)
	// No-op add/remove
	w.Add(e0)
	w.Remove(e0)
}

func TestWorldGetComponents(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Add(e0, posID, rotID)
	w.Add(e1, posID, rotID)
	w.Add(e2, rotID)

	assert.False(t, w.Has(e2, posID))
	assert.True(t, w.Has(e2, rotID))

	pos1 := (*position)(w.Get(e1, posID))
	assert.Equal(t, &position{}, pos1)

	pos1.X = 100
	pos1.Y = 101

	pos0 := (*position)(w.Get(e0, posID))
	pos1 = (*position)(w.Get(e1, posID))
	assert.Equal(t, &position{}, pos0)
	assert.Equal(t, &position{100, 101}, pos1)

	w.RemEntity(e0)

	pos1 = (*position)(w.Get(e1, posID))
	assert.Equal(t, &position{100, 101}, pos1)

	pos2 := (*position)(w.Get(e2, posID))
	assert.True(t, pos2 == nil)
}

func TestWorldIter(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[position](&world)
	rotID := ComponentID[rotation](&world)

	for i := 0; i < 1000; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}

	for i := 0; i < 100; i++ {
		query := world.Query(posID, rotID)
		for query.Next() {
			pos := (*position)(query.Get(posID))
			_ = pos
		}
		assert.Panics(t, func() { query.Next() })
	}

	for i := 0; i < MaskTotalBits-1; i++ {
		query := world.Query(posID, rotID)
		for query.Next() {
			pos := (*position)(query.Get(posID))
			_ = pos
			break
		}
	}
	query := world.Query(posID, rotID)

	assert.Panics(t, func() { world.Query(posID, rotID) })

	query.Close()
	assert.Panics(t, func() { query.Close() })
}

func TestWorldLock(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[position](&world)
	rotID := ComponentID[rotation](&world)

	var entity Entity
	for i := 0; i < 100; i++ {
		entity = world.NewEntity()
		world.Add(entity, posID)
	}

	query1 := world.Query(posID)
	query2 := world.Query(posID)
	assert.True(t, world.IsLocked())
	query1.Close()
	assert.True(t, world.IsLocked())
	query2.Close()
	assert.False(t, world.IsLocked())

	query1 = world.Query(posID)

	assert.Panics(t, func() { world.NewEntity() })
	assert.Panics(t, func() { world.RemEntity(entity) })
	assert.Panics(t, func() { world.Add(entity, rotID) })
	assert.Panics(t, func() { world.Remove(entity, posID) })
}

func TestRegisterComponents(t *testing.T) {
	world := NewWorld()

	ComponentID[position](&world)

	assert.Equal(t, ID(0), ComponentID[position](&world))
	assert.Equal(t, ID(1), ComponentID[rotation](&world))
}
