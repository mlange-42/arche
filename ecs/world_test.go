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

	assert.Equal(t, &w.archetypes[0], w.entities[0].arch)

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

func TestRegisterComponents(t *testing.T) {
	world := NewWorld()

	RegisterComponent[position](world)

	assert.Equal(t, ID(0), ComponentID[position](world))
	assert.Equal(t, ID(1), ComponentID[rotation](world))
}
