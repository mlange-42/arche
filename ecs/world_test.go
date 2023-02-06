package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldEntites(t *testing.T) {
	w := NewWorld()

	assert.Equal(t, Entity{0, 0}, w.NewEntity())
	assert.Equal(t, Entity{1, 0}, w.NewEntity())
	assert.Equal(t, Entity{2, 0}, w.NewEntity())

	assert.Equal(t, 0, int(w.entities[0].index))
	assert.Equal(t, 1, int(w.entities[1].index))
	assert.Equal(t, 2, int(w.entities[2].index))

	assert.True(t, w.RemEntity(Entity{1, 0}))
	assert.False(t, w.Alive(Entity{1, 0}))

	assert.Equal(t, 0, int(w.entities[0].index))
	assert.Equal(t, 1, int(w.entities[2].index))

	assert.Equal(t, Entity{1, 1}, w.NewEntity())
	assert.False(t, w.Alive(Entity{1, 0}))
	assert.True(t, w.Alive(Entity{1, 1}))

	assert.Equal(t, 2, int(w.entities[1].index))

	assert.False(t, w.RemEntity(Entity{1, 0}))

	assert.True(t, w.RemEntity(Entity{2, 0}))
	assert.True(t, w.RemEntity(Entity{1, 1}))
	assert.True(t, w.RemEntity(Entity{0, 0}))
}

func TestWorldComponents(t *testing.T) {
	w := NewWorld()

	posID := RegisterComponent[position](&w)
	rotID := RegisterComponent[rotation](&w)

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

	posID := RegisterComponent[position](&w)
	rotID := RegisterComponent[rotation](&w)

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

func TestRegisterComponents(t *testing.T) {
	world := NewWorld()

	RegisterComponent[position](&world)

	assert.Equal(t, ID(0), ComponentID[position](&world))
	assert.Equal(t, ID(1), ComponentID[rotation](&world))
}
