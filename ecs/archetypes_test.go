package ecs

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestArchetypePointers(t *testing.T) {
	pt := archetypePointers{}

	a1 := archetype{}
	a2 := archetype{}
	a3 := archetype{}

	pt.Add(&a1)
	pt.Add(&a2)
	pt.Add(&a3)

	assert.Equal(t, 3, pt.Len())

	for i := 0; i < 45; i++ {
		pt.Add(&a3)
	}

	assert.Equal(t, unsafe.Pointer(&a1), unsafe.Pointer(pt.Get(0)))
	assert.Equal(t, unsafe.Pointer(&a2), unsafe.Pointer(pt.Get(1)))
	assert.Equal(t, unsafe.Pointer(&a3), unsafe.Pointer(pt.Get(2)))
}
