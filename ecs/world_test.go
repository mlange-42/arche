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

	w.RemEntity(Entity{1, 0})
	assert.Equal(t, Entity{1, 1}, w.NewEntity())
}
