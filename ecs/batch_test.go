package ecs_test

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestBatchExchangeRelation(t *testing.T) {
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
	assert.Panics(t, func() {
		w.Batch().ExchangeRelation(filter, nil, []ecs.ID{posID}, velID, parent2)
	})
	assert.Panics(t, func() {
		w.Batch().ExchangeRelation(filter, []ecs.ID{velID}, []ecs.ID{posID}, velID, parent2)
	})
	assert.Panics(t, func() {
		w.Batch().ExchangeRelation(filter, nil, nil, childID, parent2)
	})
	assert.Panics(t, func() {
		_ = w.Batch().ExchangeRelationQ(filter, nil, nil, childID, parent2)
	})

	w.Batch().ExchangeRelation(filter, []ecs.ID{velID}, nil, childID, parent2)

	relFilter = ecs.NewRelationFilter(ecs.All(posID, velID, childID), parent2)
	query = w.Query(&relFilter)
	assert.Equal(t, 10, query.Count())
	query.Close()

	w.Batch().ExchangeRelation(filter, nil, []ecs.ID{velID}, childID, parent1)

	relFilter = ecs.NewRelationFilter(ecs.All(posID, childID), parent1)
	query = w.Query(&relFilter)
	assert.Equal(t, 10, query.Count())
	query.Close()

}
