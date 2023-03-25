package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapacity(t *testing.T) {
	assert.Equal(t, 0, capacity(0, 8))
	assert.Equal(t, 8, capacity(1, 8))
	assert.Equal(t, 8, capacity(8, 8))
	assert.Equal(t, 16, capacity(9, 8))
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
