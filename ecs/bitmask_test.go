package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitMask(t *testing.T) {
	mask := NewBitMask(ID(1), ID(2), ID(13), ID(27))

	assert.Equal(t, 4, mask.TotalBitsSet())

	assert.True(t, mask.Get(1))
	assert.True(t, mask.Get(2))
	assert.True(t, mask.Get(13))
	assert.True(t, mask.Get(27))

	assert.False(t, mask.Get(0))
	assert.False(t, mask.Get(3))

	mask.Set(ID(0), true)
	mask.Set(ID(1), false)

	assert.True(t, mask.Get(0))
	assert.False(t, mask.Get(1))

	other1 := NewBitMask(ID(1), ID(2), ID(32))
	other2 := NewBitMask(ID(0), ID(2))

	assert.False(t, mask.Contains(other1))
	assert.True(t, mask.Contains(other2))

	mask.Reset()
	assert.Equal(t, 0, mask.TotalBitsSet())

	mask = NewBitMask(ID(1), ID(2), ID(13), ID(27))
	other1 = NewBitMask(ID(1), ID(32))
	other2 = NewBitMask(ID(0), ID(32))

	assert.True(t, mask.ContainsAny(other1))
	assert.False(t, mask.ContainsAny(other2))
}
