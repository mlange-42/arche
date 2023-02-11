package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogicFilters(t *testing.T) {

	hasA := All(0)
	hasB := All(1)
	hasAll := All(0, 1)
	hasNone := Mask{}

	var filter MaskFilter
	filter = hasA
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = &Or{hasA, hasB}
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = &XOr{hasA, hasB}
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = &hasAll
	assert.True(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = &And{hasA, hasB}
	assert.True(t, match(filter, hasAll))
	assert.False(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = Not(hasB)
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = &And{hasA, Not(hasB)}
	assert.False(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.False(t, match(filter, hasNone))

	filter = &Or{hasAll, Not(hasB)}
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.False(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	filter = &Or{OneOf(0, 1), Not(hasB)}
	assert.True(t, match(filter, hasAll))
	assert.True(t, match(filter, hasA))
	assert.True(t, match(filter, hasB))
	assert.True(t, match(filter, hasNone))

	assert.Equal(t, &Or{hasA, hasB}, OneOf(0, 1))
}

func match(f MaskFilter, m Mask) bool {
	return f.Matches(m.BitMask)
}
