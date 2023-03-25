package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	world := NewWorld(NewConfig().WithCapacityIncrement(16))
	world.SetListener(func(e *EntityEvent) {})

	posID := ComponentID[position](&world)
	rotID := ComponentID[rotation](&world)

	comps := []Component{
		{ID: posID, Comp: &position{100, 200}},
		{ID: rotID, Comp: &rotation{300}},
	}

	b := world.Batch()

	b.NewEntities(100, posID)
	b.NewEntitiesWith(100, comps...)

	filter1 := All(posID).Exact()
	q := world.Query(&filter1)
	assert.Equal(t, 100, q.Count())
	q.Close()

	filter2 := All(posID, rotID).Exact()
	q = world.Query(&filter2)
	assert.Equal(t, 100, q.Count())
	q.Close()
}

func TestBatchQuery(t *testing.T) {
	world := NewWorld(NewConfig().WithCapacityIncrement(16))
	world.SetListener(func(e *EntityEvent) {})

	posID := ComponentID[position](&world)
	rotID := ComponentID[rotation](&world)

	comps := []Component{
		{ID: posID, Comp: &position{100, 200}},
		{ID: rotID, Comp: &rotation{300}},
	}

	b := world.Batch()

	q := b.NewEntitiesQuery(100, posID)
	q.Close()
	q = b.NewEntitiesWithQuery(100, comps...)
	q.Close()

	filter1 := All(posID).Exact()
	q = world.Query(&filter1)
	assert.Equal(t, 100, q.Count())
	q.Close()

	filter2 := All(posID, rotID).Exact()
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
