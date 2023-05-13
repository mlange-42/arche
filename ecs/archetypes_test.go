package ecs

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestArchetypePointers(t *testing.T) {
	pt := pointers[archetype]{}

	a1 := archetype{}
	a2 := archetype{}
	a3 := archetype{}

	pt.Add(&a1)
	pt.Add(&a2)
	pt.Add(&a3)

	assert.Equal(t, int32(3), pt.Len())

	var last archetype
	for i := 0; i < 15; i++ {
		last = archetype{}
		pt.Add(&last)
	}

	assert.Equal(t, unsafe.Pointer(&a1), unsafe.Pointer(pt.Get(0)))
	assert.Equal(t, unsafe.Pointer(&a2), unsafe.Pointer(pt.Get(1)))
	assert.Equal(t, unsafe.Pointer(&a3), unsafe.Pointer(pt.Get(2)))

	assert.Equal(t, int32(18), pt.Len())

	pt.RemoveAt(1)
	assert.Equal(t, int32(17), pt.Len())
	assert.Equal(t, unsafe.Pointer(&last), unsafe.Pointer(pt.Get(1)))
	assert.Equal(t, unsafe.Pointer(&a3), unsafe.Pointer(pt.Get(2)))
}

func TestBatchArchetype(t *testing.T) {
	arch := archetype{}
	batch := batchArchetypes{}
	batch.Add(&arch, nil, 0, 1)

	assert.Equal(t, &arch, batch.Get(0))
	assert.Equal(t, int32(1), batch.Len())
}
