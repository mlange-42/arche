package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityEvent(t *testing.T) {
	e := EntityEvent{AddedRemoved: 0}

	assert.False(t, e.EntityAdded())
	assert.False(t, e.EntityRemoved())

	e = EntityEvent{AddedRemoved: 1}

	assert.True(t, e.EntityAdded())
	assert.False(t, e.EntityRemoved())

	e = EntityEvent{AddedRemoved: -1}

	assert.False(t, e.EntityAdded())
	assert.True(t, e.EntityRemoved())
}
