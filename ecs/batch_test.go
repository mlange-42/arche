package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	world := NewWorld(NewConfig().WithCapacityIncrement(16))
	world.SetListener(func(e *EntityEvent) {})

	posID := ComponentID[Position](&world)
	rotID := ComponentID[rotation](&world)

	comps := []Component{
		{ID: posID, Comp: &Position{100, 200}},
		{ID: rotID, Comp: &rotation{300}},
	}

	b := world.Batch()

	b.NewEntities(100, posID)
	b.NewEntitiesWith(100, comps...)

	filter1 := All(posID).Exclusive()
	q := world.Query(&filter1)
	assert.Equal(t, 100, q.Count())
	q.Close()

	filter2 := All(posID, rotID).Exclusive()
	q = world.Query(&filter2)
	assert.Equal(t, 100, q.Count())
	q.Close()
}

func TestBatchQuery(t *testing.T) {
	world := NewWorld(NewConfig().WithCapacityIncrement(16))
	world.SetListener(func(e *EntityEvent) {})

	posID := ComponentID[Position](&world)
	rotID := ComponentID[rotation](&world)

	comps := []Component{
		{ID: posID, Comp: &Position{100, 200}},
		{ID: rotID, Comp: &rotation{300}},
	}

	b := world.Batch()

	q := b.NewEntitiesQuery(100, posID)
	q.Close()
	q = b.NewEntitiesWithQuery(100, comps...)
	q.Close()

	filter1 := All(posID).Exclusive()
	q = world.Query(&filter1)
	assert.Equal(t, 100, q.Count())
	q.Close()

	filter2 := All(posID, rotID).Exclusive()
	q = world.Query(&filter2)
	assert.Equal(t, 100, q.Count())
	q.Close()

	b.RemoveEntities(&filter2)

	q = world.Query(&filter1)
	assert.Equal(t, 100, q.Count())
	q.Close()

	q = world.Query(&filter2)
	assert.Equal(t, 0, q.Count())
	q.Close()
}

func ExampleBatch() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	world.Batch().NewEntities(10_000, posID, velID)
	// Output:
}

func ExampleBatch_NewEntities() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	world.Batch().NewEntities(10_000, posID, velID)
	// Output:
}

func ExampleBatch_NewEntitiesQuery() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	query := world.Batch().NewEntitiesQuery(10_000, posID, velID)

	for query.Next() {
		// initialize components of the newly created entities
	}
	// Output:
}

func ExampleBatch_NewEntitiesWith() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	components := []Component{
		{ID: posID, Comp: &Position{X: 0, Y: 0}},
		{ID: velID, Comp: &Velocity{X: 10, Y: 2}},
	}

	world.Batch().NewEntitiesWith(10_000, components...)
	// Output:
}

func ExampleBatch_NewEntitiesWithQuery() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	components := []Component{
		{ID: posID, Comp: &Position{X: 0, Y: 0}},
		{ID: velID, Comp: &Velocity{X: 10, Y: 2}},
	}

	query := world.Batch().NewEntitiesWithQuery(10_000, components...)

	for query.Next() {
		// initialize components of the newly created entities
	}
	// Output:
}

func ExampleBatch_RemoveEntities() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	world.Batch().NewEntities(10_000, posID, velID)

	filter := All(posID, velID).Exclusive()
	world.Batch().RemoveEntities(&filter)
	// Output:
}
