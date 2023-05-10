package ecs_test

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestCachedMaskFilter(t *testing.T) {
	f := ecs.All(1, 2, 3).Without(4)

	assert.True(t, f.Matches(ecs.All(1, 2, 3)))
	assert.True(t, f.Matches(ecs.All(1, 2, 3, 5)))

	assert.False(t, f.Matches(ecs.All(1, 2)))
	assert.False(t, f.Matches(ecs.All(1, 2, 3, 4)))
}

func TestCachedFilter(t *testing.T) {
	w := ecs.NewWorld()

	f := ecs.All(1, 2, 3)
	fc := w.Cache().Register(f)

	assert.Equal(t, f.Matches(ecs.All(1, 2, 3)), fc.Matches(ecs.All(1, 2, 3)))
	assert.Equal(t, f.Matches(ecs.All(1, 2)), fc.Matches(ecs.All(1, 2)))

	w.Cache().Unregister(&fc)
}

func ExampleMaskFilter() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	filter := ecs.All(posID).Without(velID)
	query := world.Query(&filter)

	for query.Next() {
		// ...
	}
	// Output:
}

func ExampleCachedFilter() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)

	filter := ecs.All(posID)
	cached := world.Cache().Register(filter)

	query := world.Query(&cached)

	for query.Next() {
		// ...
	}
	// Output:
}

func ExampleRelationFilter() {
	world := ecs.NewWorld()
	childID := ecs.ComponentID[ChildOf](&world)

	target := world.NewEntity()

	builder := ecs.NewBuilder(&world, childID).WithRelation(childID)
	builder.NewBatch(100, target)

	filter := ecs.RelationFilter(ecs.All(childID), target)

	query := world.Query(filter)
	for query.Next() {
		// ...
	}
	// Output:
}
