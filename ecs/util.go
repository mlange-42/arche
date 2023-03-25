package ecs

// Calculates the capacity required for size, given an increment.
func capacity(size, increment int) int {
	cap := increment * (size / increment)
	if size%increment != 0 {
		cap += increment
	}
	return cap
}

// Calculates the capacity required for size, given an increment.
func capacityU32(size, increment uint32) uint32 {
	cap := increment * (size / increment)
	if size%increment != 0 {
		cap += increment
	}
	return cap
}

// Manages locks by mask bits.
type lockMask struct {
	locks   Mask    // The actual locks.
	bitPool bitPool // The bit pool for getting and recycling bits.
}

// Lock the world and get the Lock bit for later unlocking.
func (m *lockMask) Lock() uint8 {
	lock := m.bitPool.Get()
	m.locks.Set(ID(lock), true)
	return lock
}

// Unlock unlocks the given lock bit.
func (m *lockMask) Unlock(l uint8) {
	if !m.locks.Get(ID(l)) {
		panic("unbalanced unlock")
	}
	m.locks.Set(ID(l), false)
	m.bitPool.Recycle(l)
}

// IsLocked returns whether the world is locked by any queries.
func (m *lockMask) IsLocked() bool {
	return !m.locks.IsZero()
}

// Reset the locks and the pool.
func (m *lockMask) Reset() {
	m.locks = Mask{}
	m.bitPool.Reset()
}
