//go:build tiny

package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitMaskTiny(t *testing.T) {
	mask := Mask{}
	mask.Set(id(100), true)

	assert.True(t, mask.IsZero())
}
