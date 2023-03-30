package ecs

// Page size of pagedSlice type
const pageSize = 32

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
//
// The number of simultaneous locks at a given time is limited to [MaskTotalBits].
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

// pagedSlice is a paged collection working with pages of length 32 slices.
// It's primary purpose is pointer persistence, which is not given using simple slices.
//
// Implements [archetypes].
type pagedSlice[T any] struct {
	pages   [][]T
	len     int
	lenLast int
}

// newPagedSlice creates a new pagedSlice with the given page size/capacity increment.
func newPagedSlice[T any]() pagedSlice[T] {
	return pagedSlice[T]{}
}

// Add adds a value to the paged slice.
func (p *pagedSlice[T]) Add(value T) {
	if p.len == 0 || p.lenLast == pageSize {
		p.pages = append(p.pages, make([]T, pageSize))
		p.lenLast = 0
	}
	p.pages[len(p.pages)-1][p.lenLast] = value
	p.len++
	p.lenLast++
}

// Get returns the value at the given index.
func (p *pagedSlice[T]) Get(index int) *T {
	return &p.pages[index/pageSize][index%pageSize]
}

// Len returns the current number of items in the paged slice.
func (p *pagedSlice[T]) Len() int {
	return p.len
}
