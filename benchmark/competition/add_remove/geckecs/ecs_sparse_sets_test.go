package geckecs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSparseSet(t *testing.T) {
	// from https://skypjack.github.io/2020-08-02-ecs-baf-part-9/
	w := NewWorld()
	e := w.EntityFromU32(2)

	ss := NewSparseSet[int](nil)
	assert.False(t, ss.Has(e))

	ss.Set(1, e)
	assert.True(t, ss.Has(e))

	ss.Set(2, w.EntityFromU32(7))
	ss.Set(3, w.EntityFromU32(0))

	assert.ElementsMatch(t, ss.sparsePages, []int{2, DeadEntityID, 0, DeadEntityID, DeadEntityID, DeadEntityID, DeadEntityID, 1})
	assert.ElementsMatch(t, ss.dense, []Entity{e, w.EntityFromU32(7), w.EntityFromU32(0)})
	assert.ElementsMatch(t, ss.components, []int{1, 2, 3})

	ss.Set(4, w.EntityFromU32(6))
	assert.ElementsMatch(t, ss.sparsePages, []int{2, DeadEntityID, 0, DeadEntityID, DeadEntityID, DeadEntityID, 3, 1})
	assert.ElementsMatch(t, ss.dense, []Entity{e, w.EntityFromU32(7), w.EntityFromU32(0), w.EntityFromU32(6)})
	assert.ElementsMatch(t, ss.components, []int{1, 2, 3, 4})

	ss.Remove(w.EntityFromU32(7))
	assert.ElementsMatch(t, ss.sparsePages, []int{2, DeadEntityID, 0, DeadEntityID, DeadEntityID, DeadEntityID, 1, DeadEntityID}, "removed sparse")
	assert.ElementsMatch(t, ss.dense, []Entity{e, w.EntityFromU32(6), w.EntityFromU32(0)}, "removed dense")
	assert.ElementsMatch(t, ss.components, []int{1, 4, 3}, "removed components")

	// No sort set, should be same
	ss.Sort()
	assert.ElementsMatch(t, ss.sparsePages, []int{2, DeadEntityID, 0, DeadEntityID, DeadEntityID, DeadEntityID, 1, DeadEntityID}, "removed sparse")
	assert.ElementsMatch(t, ss.dense, []Entity{e, w.EntityFromU32(6), w.EntityFromU32(0)}, "removed dense")
	assert.ElementsMatch(t, ss.components, []int{1, 4, 3}, "removed components")

	// Set sort function
	ss.LessThan = func(a, b Entity) bool {
		aT, _ := ss.Read(a)
		bT, _ := ss.Read(b)

		if aT < bT {
			return aT < bT
		}
		return false
	}
	ss.Sort()
	assert.Equal(t, 0, ss.sparsePages[2])
	assert.Equal(t, 2, ss.sparsePages[0])
	assert.Equal(t, 1, ss.sparsePages[6])
	assert.Equal(t, uint32(2), ss.dense[0].val)
	assert.Equal(t, uint32(0), ss.dense[1].val)
	assert.Equal(t, uint32(6), ss.dense[2].val)
	assert.Equal(t, 1, ss.components[0])
	assert.Equal(t, 3, ss.components[1])
	assert.Equal(t, 4, ss.components[2])
}

func BenchmarkSparseSet(b *testing.B) {
	for entityCount := uint32(10); entityCount <= 1_000_000; entityCount *= 10 {
		name := fmt.Sprintf("N=%d", entityCount)
		b.Run(name, func(b *testing.B) {
			b.Run("Set", func(b *testing.B) {
				w := NewWorld()
				s := NewSparseSet[int](nil)
				for i := 0; i < b.N; i++ {
					ee := make([]Entity, entityCount)
					for j := uint32(0); j < entityCount; j++ {
						ee[j] = w.EntityFromU32(j)
					}
					s.Set(1, ee...)
				}
			})

			b.Run("Contains", func(b *testing.B) {
				w := NewWorld()
				s := NewSparseSet[int](nil)
				for i := 0; i < b.N; i++ {
					for j := uint32(0); j < entityCount; j++ {
						s.Has(w.EntityFromU32(j))
					}
				}
			})

			b.Run("Remove", func(b *testing.B) {
				w := NewWorld()
				s := NewSparseSet[int](nil)
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					ee := make([]Entity, entityCount)
					for j := uint32(0); j < entityCount; j++ {
						ee[j] = w.EntityFromU32(j)
					}
					s.Set(i, ee...)
					b.StartTimer()

					s.Remove(ee...)
				}
			})

			b.Run("Sort", func(b *testing.B) {
				w := NewWorld()
				s := NewSparseSet[uint32](nil)
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					for j := uint32(0); j < entityCount; j++ {
						s.Set(j, w.EntityFromU32(j))
					}
					b.StartTimer()

					s.Sort()
				}
			})
		})
	}
}
