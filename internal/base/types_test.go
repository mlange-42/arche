package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMask(t *testing.T) {
	filter := NewMask(0, 2, 4)
	other := NewBitMask(0, 1, 2)

	assert.False(t, filter.Matches(other))

	other = NewBitMask(0, 1, 2, 3, 4)
	assert.True(t, filter.Matches(other))
}
