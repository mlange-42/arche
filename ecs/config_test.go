package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := newConfig(16)
	assert.Equal(t, 16, c.initialCapacity)
	assert.Equal(t, 16, c.initialCapacityRelations)

	c = newConfig(16, 8)
	assert.Equal(t, 16, c.initialCapacity)
	assert.Equal(t, 8, c.initialCapacityRelations)
}
