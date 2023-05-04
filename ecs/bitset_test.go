package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitSet(t *testing.T) {
	b := bitSet{}

	b.ExtendTo(64)
	assert.Equal(t, 1, len(b.data))
	b.ExtendTo(65)
	assert.Equal(t, 2, len(b.data))
	b.ExtendTo(120)
	assert.Equal(t, 2, len(b.data))

	assert.False(t, b.Get(127))
	b.Set(127, true)
	assert.True(t, b.Get(127))
	b.Set(127, false)
	assert.False(t, b.Get(127))

	b.Set(63, true)
	b.Reset()
	assert.False(t, b.Get(63))
}
