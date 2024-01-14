package ecs

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestCapacity(t *testing.T) {
	assert.Equal(t, 0, capacity(0, 8))
	assert.Equal(t, 8, capacity(1, 8))
	assert.Equal(t, 8, capacity(8, 8))
	assert.Equal(t, 16, capacity(9, 8))
}

func TestCapacityNonZero(t *testing.T) {
	assert.Equal(t, 8, capacityNonZero(0, 8))
	assert.Equal(t, 8, capacityNonZero(1, 8))
	assert.Equal(t, 8, capacityNonZero(8, 8))
	assert.Equal(t, 16, capacityNonZero(9, 8))
}

func TestCapacityU32(t *testing.T) {
	assert.Equal(t, 0, int(capacityU32(0, 8)))
	assert.Equal(t, 8, int(capacityU32(1, 8)))
	assert.Equal(t, 8, int(capacityU32(8, 8)))
	assert.Equal(t, 16, int(capacityU32(9, 8)))
}

func TestLockMask(t *testing.T) {
	locks := lockMask{}

	assert.False(t, locks.IsLocked())

	l1 := locks.Lock()
	assert.True(t, locks.IsLocked())
	assert.Equal(t, 0, int(l1))

	l2 := locks.Lock()
	assert.True(t, locks.IsLocked())
	assert.Equal(t, 1, int(l2))

	locks.Unlock(l1)
	assert.True(t, locks.IsLocked())

	assert.Panics(t, func() { locks.Unlock(l1) })

	locks.Unlock(l2)
	assert.False(t, locks.IsLocked())
}

func TestPagedSlice(t *testing.T) {
	a := pagedSlice[int32]{}

	var i int32
	for i = 0; i < 66; i++ {
		a.Add(i)
		assert.Equal(t, i, *a.Get(i))
		assert.Equal(t, i+1, a.Len())
	}

	a.Set(3, 100)
	assert.Equal(t, int32(100), *a.Get(3))
}

func TestPagedSlicePointerPersistence(t *testing.T) {
	a := pagedSlice[int32]{}

	a.Add(0)
	p1 := a.Get(0)

	var i int32
	for i = 1; i < 66; i++ {
		a.Add(i)
		assert.Equal(t, i, *a.Get(i))
		assert.Equal(t, i+1, a.Len())
	}

	p2 := a.Get(0)
	assert.Equal(t, unsafe.Pointer(p1), unsafe.Pointer(p2))
	*p1 = 100
	assert.Equal(t, int32(100), *p2)
}

func BenchmarkPagedSlice_Get(b *testing.B) {
	b.StopTimer()

	count := 128
	s := pagedSlice[int]{}

	for i := 0; i < count; i++ {
		s.Add(1)
	}

	b.StartTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		sum += *s.Get(int32(i % count))
	}
}
