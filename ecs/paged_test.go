package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagedArr32(t *testing.T) {
	a := PagedArr32[int]{}

	for i := 0; i < 66; i++ {
		a.Add(i)
		assert.Equal(t, i, *a.Get(i))
		assert.Equal(t, i+1, a.Len())
	}
}
