package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogicFilters(t *testing.T) {
	w := NewWorld()

	hasA := Mask1[testStruct0](&w)
	hasB := Mask1[testStruct1](&w)
	hasAll := Mask2[testStruct0, testStruct1](&w)
	hasNone := Mask{}

	var filter filter
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
}

func match(f filter, m Mask) bool {
	return f.Matches(m)
}
