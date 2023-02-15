package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeEvent(t *testing.T) {
	e := ChangeEvent{AddedRemoved: 0}

	assert.False(t, e.EntityAdded())
	assert.False(t, e.EntityRemoved())

	e = ChangeEvent{AddedRemoved: 1}

	assert.True(t, e.EntityAdded())
	assert.False(t, e.EntityRemoved())

	e = ChangeEvent{AddedRemoved: -1}

	assert.False(t, e.EntityAdded())
	assert.True(t, e.EntityRemoved())
}
