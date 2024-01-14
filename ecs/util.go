package ecs

import (
	"fmt"
	"strings"
)

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
// Always returns a value greater than zero.
func capacityNonZero(size, increment int) int {
	if size == 0 {
		return increment
	}
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

func maskToTypes(mask Mask, reg *componentRegistry[ID]) []componentType {
	count := int(mask.TotalBitsSet())
	types := make([]componentType, count)

	idx := 0
	for i := range mask.bits {
		if mask.bits[i] == 0 {
			continue
		}
		for j := 0; j < wordSize; j++ {
			id := ID(i*wordSize + j)
			if mask.Get(id) {
				types[idx] = componentType{ID: id, Type: reg.Types[id]}
				idx++
			}
		}
	}
	return types
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
	len     int32
	lenLast int32
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
func (p *pagedSlice[T]) Get(index int32) *T {
	return &p.pages[index/pageSize][index%pageSize]
}

// Set sets the value at the given index.
func (p *pagedSlice[T]) Set(index int32, value T) {
	p.pages[index/pageSize][index%pageSize] = value
}

// Len returns the current number of items in the paged slice.
func (p *pagedSlice[T]) Len() int32 {
	return p.len
}

// Prints world nodes and archetypes.
func debugPrintWorld(w *World) string {
	sb := strings.Builder{}

	ln := w.nodes.Len()
	var i int32
	for i = 0; i < ln; i++ {
		nd := w.nodes.Get(i)
		if !nd.IsActive {
			fmt.Fprintf(&sb, "Node %v (inactive)\n", nd.Ids)
			continue
		}
		nodeArches := nd.Archetypes()
		ln2 := int32(nodeArches.Len())
		fmt.Fprintf(&sb, "Node %v (%d arch), relation: %t\n", nd.Ids, ln2, nd.HasRelation)
		var j int32
		for j = 0; j < ln2; j++ {
			a := nodeArches.Get(j)
			if a.IsActive() {
				fmt.Fprintf(&sb, "   Arch %v (%d entities)\n", a.RelationTarget, a.Len())
			} else {
				fmt.Fprintf(&sb, "   Arch %v (inactive)\n", a.RelationTarget)
			}
		}
	}

	return sb.String()
}
