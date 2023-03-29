package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCachedMaskFilter(t *testing.T) {
	f := All(1, 2, 3).Without(4)

	assert.True(t, f.Matches(All(1, 2, 3)))
	assert.True(t, f.Matches(All(1, 2, 3, 5)))

	assert.False(t, f.Matches(All(1, 2)))
	assert.False(t, f.Matches(All(1, 2, 3, 4)))
}

func TestCachedFilter(t *testing.T) {
	f := All(1, 2, 3)
	fc := CachedFilter{filter: f, id: 0}

	assert.Equal(t, f.Matches(All(1, 2, 3)), fc.Matches(All(1, 2, 3)))
	assert.Equal(t, f.Matches(All(1, 2)), fc.Matches(All(1, 2)))
}
