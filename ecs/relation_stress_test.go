package ecs_test

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestRelationStress(t *testing.T) {
	numRuns := 25
	numParents := 100
	numChildren := 100

	world := ecs.NewWorld(
		ecs.NewConfig().WithCapacityIncrement(1024).WithRelationCapacityIncrement(128),
	)

	posID := ecs.ComponentID[Position](&world)
	relID := ecs.ComponentID[ChildOf](&world)

	for run := 0; run < numRuns; run++ {
		world.Reset()

		parents := make([]ecs.Entity, 0, numParents)
		parBuilder := ecs.NewBuilder(&world, posID)
		childBuilder := ecs.NewBuilder(&world, relID).WithRelation(relID)

		query := parBuilder.NewBatchQ(numParents)
		for query.Next() {
			parents = append(parents, query.Entity())
		}

		for _, parent := range parents {
			childBuilder.NewBatch(numChildren, parent)
		}

		for i := 0; i < 1000; i++ {
			parIdx := rand.Intn(numParents)
			parent := parents[parIdx]

			childFilter := ecs.NewRelationFilter(ecs.All(relID), parent)
			removed := world.Batch().RemoveEntities(&childFilter)
			world.RemoveEntity(parent)
			assert.Equal(t, numChildren, removed)

			stats := world.Stats()
			assert.Equal(t, numParents, stats.Nodes[2].ArchetypeCount)
			assert.Equal(t, numParents-1, stats.Nodes[2].ActiveArchetypeCount)

			parent = parBuilder.New()
			parents[parIdx] = parent
			childBuilder.NewBatch(numChildren, parent)

			stats = world.Stats()
			assert.Equal(t, numParents, stats.Nodes[2].ArchetypeCount)
		}
	}
}
