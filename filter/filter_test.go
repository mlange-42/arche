package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogicFilters(t *testing.T) {

	hasA := All(0)
	hasB := All(1)
	hasAll := All(0, 1)
	hasNone := All()

	var filter MaskFilter
	filter = hasA
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = Or(hasA, hasB)
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = XOr(hasA, hasB)
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = &hasAll
	assert.True(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = And(hasA, hasB)
	assert.True(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = NotAny(1)
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = NotAny(0, 1)
	assert.False(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = Not(0)
	assert.False(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = Not(0, 1)
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = And(hasA, NotANY(hasB))
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = Or(hasAll, NotAny(1))
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = Or(OneOf(0, 1), NotANY(hasB))
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	assert.Equal(t, Or(hasA, hasB), OneOf(0, 1))
}

func match(f MaskFilter, m Mask) bool {
	return f.Matches(m.BitMask)
}

func BenchmarkFilterStackOr(b *testing.B) {
	b.StopTimer()
	mask := All(1, 2, 3, 4, 5)

	filter := OR{All(1), All(2)}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(mask.BitMask)
	}
}

func BenchmarkFilterStack5And(b *testing.B) {
	b.StopTimer()
	mask := All(1, 2, 3, 4, 5)

	a1 := AND{All(1), All(2)}
	a2 := AND{&a1, All(3)}
	a3 := AND{&a2, All(4)}
	filter := AND{&a3, All(5)}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(mask.BitMask)
	}
}

func BenchmarkFilterHeapOr(b *testing.B) {
	b.StopTimer()
	mask := All(1, 2, 3, 4, 5)

	filter := Or(All(1), All(2))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(mask.BitMask)
	}
}

func BenchmarkFilterHeap5And(b *testing.B) {
	b.StopTimer()
	mask := All(1, 2, 3, 4, 5)

	filter := And(All(1), And(All(2), And(All(3), And(All(4), All(5)))))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = filter.Matches(mask.BitMask)
	}
}
