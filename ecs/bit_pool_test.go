package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitPool(t *testing.T) {
	p := newBitPool()

	for i := 0; i < MaskTotalBits; i++ {
		assert.Equal(t, i, int(p.Get()))
	}

	assert.Panics(t, func() { p.Get() })

	for i := 0; i < 10; i++ {
		p.Recycle(uint8(i))
	}
	for i := 9; i >= 0; i-- {
		assert.Equal(t, i, int(p.Get()))
	}

	assert.Panics(t, func() { p.Get() })
}
