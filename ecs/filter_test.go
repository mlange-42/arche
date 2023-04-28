package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCachedMaskFilter(t *testing.T) {
	f := All(1, 2, 3).Without(4)

	assert.True(t, f.Matches(All(1, 2, 3), nil))
	assert.True(t, f.Matches(All(1, 2, 3, 5), nil))

	assert.False(t, f.Matches(All(1, 2), nil))
	assert.False(t, f.Matches(All(1, 2, 3, 4), nil))
}

func TestCachedFilter(t *testing.T) {
	f := All(1, 2, 3)
	fc := CachedFilter{filter: f, id: 0}

	assert.Equal(t, f.Matches(All(1, 2, 3), nil), fc.Matches(All(1, 2, 3), nil))
	assert.Equal(t, f.Matches(All(1, 2), nil), fc.Matches(All(1, 2), nil))
}

func ExampleMaskFilter() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	filter := All(posID).Without(velID)
	query := world.Query(&filter)

	for query.Next() {
		// ...
	}
	// Output:
}

func ExampleCachedFilter() {
	world := NewWorld()
	posID := ComponentID[Position](&world)

	filter := All(posID)
	cached := world.Cache().Register(filter)

	query := world.Query(&cached)

	for query.Next() {
		// ...
	}
	// Output:
}
