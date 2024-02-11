package ecs_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestRelationsExchange(t *testing.T) {
	w := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&w)
	velID := ecs.ComponentID[Velocity](&w)
	childID := ecs.ComponentID[ChildOf](&w)

	parent1 := w.NewEntity(posID)
	parent2 := w.NewEntity(posID)

	builder := ecs.NewBuilder(&w, posID, childID).WithRelation(childID)
	builder.NewBatch(10, parent1)

	relFilter := ecs.NewRelationFilter(ecs.All(posID, childID), parent1)
	query := w.Query(&relFilter)
	assert.Equal(t, 10, query.Count())
	query.Close()

	filter := ecs.All(posID, childID)
	assert.PanicsWithValue(t, "can't add relation: resulting entity has no component Velocity", func() {
		w.Relations().ExchangeBatch(filter, nil, []ecs.ID{posID}, velID, parent2)
	})
	assert.PanicsWithValue(t, "can't add relation: Velocity is not a relation component", func() {
		w.Relations().ExchangeBatch(filter, []ecs.ID{velID}, []ecs.ID{posID}, velID, parent2)
	})
	assert.PanicsWithValue(t, "exchange operation has no effect, but a relation is specified. Use Batch.SetRelation instead", func() {
		w.Relations().ExchangeBatch(filter, nil, nil, childID, parent2)
	})
	assert.PanicsWithValue(t, "exchange operation has no effect, but a relation is specified. Use Batch.SetRelation instead", func() {
		_ = w.Relations().ExchangeBatchQ(filter, nil, nil, childID, parent2)
	})

	cnt := w.Relations().ExchangeBatch(filter, []ecs.ID{velID}, nil, childID, parent2)
	assert.Equal(t, 10, cnt)

	relFilter = ecs.NewRelationFilter(ecs.All(posID, velID, childID), parent2)
	query = w.Query(&relFilter)
	assert.Equal(t, 10, query.Count())
	query.Close()

	cnt = w.Relations().ExchangeBatch(filter, nil, []ecs.ID{velID}, childID, parent1)
	assert.Equal(t, 10, cnt)

	relFilter = ecs.NewRelationFilter(ecs.All(posID, childID), parent1)
	query = w.Query(&relFilter)
	assert.Equal(t, 10, query.Count())
	query.Close()

}

func ExampleRelations() {
	world := ecs.NewWorld()

	relID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()
	child := world.NewEntity(relID)

	world.Relations().Set(child, relID, parent)
	fmt.Println(world.Relations().Get(child, relID))
	// Output: {1 0}
}

func ExampleRelations_SetBatch() {
	world := ecs.NewWorld()

	relID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()
	ecs.NewBuilder(&world, relID).NewBatch(100)

	filter := ecs.All(relID)
	world.Relations().SetBatch(filter, relID, parent)
	// Output:
}

func ExampleRelations_SetBatchQ() {
	world := ecs.NewWorld()

	relID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()
	ecs.NewBuilder(&world, relID).NewBatch(100)

	filter := ecs.All(relID)
	query := world.Relations().SetBatchQ(filter, relID, parent)
	fmt.Println(query.Count())
	query.Close()
	// Output: 100
}
