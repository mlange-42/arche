package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCachedMaskFilter(t *testing.T) {
	f := All(id(1), id(2), id(3)).Without(id(4))

	assert.True(t, f.Matches(all(id(1), id(2), id(3))))
	assert.True(t, f.Matches(all(id(1), id(2), id(3), id(5))))

	assert.False(t, f.Matches(all(id(1), id(2))))
	assert.False(t, f.Matches(all(id(1), id(2), id(3), id(4))))
}

func TestCachedFilter(t *testing.T) {
	w := NewWorld()

	f := All(id(1), id(2), id(3))
	fc := w.Cache().Register(f)

	assert.Equal(t, f.Matches(all(id(1), id(2), id(3))), fc.Matches(all(id(1), id(2), id(3))))
	assert.Equal(t, f.Matches(all(id(1), id(2))), fc.Matches(all(id(1), id(2))))

	w.Cache().Unregister(&fc)
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

func ExampleRelationFilter() {
	world := NewWorld()
	childID := ComponentID[ChildOf](&world)

	target := world.NewEntity()

	builder := NewBuilder(&world, childID).WithRelation(childID)
	builder.NewBatch(100, target)

	filter := NewRelationFilter(All(childID), target)

	query := world.Query(&filter)
	for query.Next() {
		// ...
	}
	// Output:
}
