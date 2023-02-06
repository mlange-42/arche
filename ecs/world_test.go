package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldEntites(t *testing.T) {
	w := newWorld()

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
	w := newWorld()

	posID := RegisterComponent[position](w)
	rotID := RegisterComponent[rotation](w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	assert.Equal(t, 1, len(w.archetypes))

	w.Add(e0, posID)
	assert.Equal(t, 2, len(w.archetypes))
	w.Add(e1, posID, rotID)
	assert.Equal(t, 3, len(w.archetypes))
	w.Add(e2, rotID)
	assert.Equal(t, 4, len(w.archetypes))

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
}

func TestRegisterComponents(t *testing.T) {
	world := NewWorld()

	RegisterComponent[position](world)

	assert.Equal(t, ID(0), ComponentID[position](world))
	assert.Equal(t, ID(1), ComponentID[rotation](world))
}
