package ecs

import (
	"testing"

	"github.com/mlange-42/arche/internal/base"
	"github.com/stretchr/testify/assert"
)

func TestBitMask(t *testing.T) {
	mask := base.NewBitMask(ID(1), ID(2), ID(13), ID(27))

	assert.Equal(t, uint(4), mask.TotalBitsSet())

	assert.True(t, mask.Get(1))
	assert.True(t, mask.Get(2))
	assert.True(t, mask.Get(13))
	assert.True(t, mask.Get(27))

	assert.False(t, mask.Get(0))
	assert.False(t, mask.Get(3))
}
