package ecs

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestPagedArr32(t *testing.T) {
	a := newPagedSlice[int](32)

	for i := 0; i < 66; i++ {
		a.Add(i)
		assert.Equal(t, i, *a.Get(i))
		assert.Equal(t, i+1, a.Len())
	}
}

func TestPagedArrPointerPersistence(t *testing.T) {
	a := newPagedSlice[int](32)

	a.Add(0)
	p1 := a.Get(0)

	for i := 1; i < 66; i++ {
		a.Add(i)
		assert.Equal(t, i, *a.Get(i))
		assert.Equal(t, i+1, a.Len())
	}

	p2 := a.Get(0)
	assert.Equal(t, unsafe.Pointer(p1), unsafe.Pointer(p2))
	*p1 = 100
	assert.Equal(t, 100, *p2)
}
