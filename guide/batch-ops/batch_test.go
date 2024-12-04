package batchops

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// Position component.
type Position struct {
	X float64
	Y float64
}

// Heading component
type Heading struct {
	Angle float64
}

// ChildOf relation component
type ChildOf struct {
	ecs.Relation
}

func TestBatchCreate(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	builder := ecs.NewBuilder(&world, posID, headID)

	builder.NewBatch(100)
}

func TestBatchCreateGeneric(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap2[Position, Heading](&world)

	builder.NewBatch(100)
}

func TestBatchCreateQuery(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	builder := ecs.NewBuilder(&world, posID, headID)

	query := builder.NewBatchQ(100)
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		head := (*Heading)(query.Get(headID))

		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100
		head.Angle = rand.Float64() * 360
	}
}

func TestBatchCreateQueryGeneric(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap2[Position, Heading](&world)

	query := builder.NewBatchQ(100)
	for query.Next() {
		pos, head := query.Get()

		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100
		head.Angle = rand.Float64() * 360
	}
}

func TestBatchAddQuery(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	// Create 100 entities with Position.
	builder := ecs.NewBuilder(&world, posID)
	builder.NewBatch(100)

	// Filter for entities with Position.
	filter := ecs.All(posID)
	// Batch-add Heading to them, using the query method for initialization.
	query := world.Batch().AddQ(filter, headID)
	for query.Next() {
		head := (*Heading)(query.Get(headID))
		head.Angle = rand.Float64() * 360
	}
}

func TestBatchAddQueryGeneric(t *testing.T) {
	world := ecs.NewWorld()

	// Create 100 entities with Position.
	builder := generic.NewMap1[Position](&world)
	builder.NewBatch(100)

	// Create a generic map to perform the batch operation
	adder := generic.NewMap1[Heading](&world)
	// Filter for entities with Position.
	filter := generic.NewFilter1[Position]()
	// Batch-add Heading to them, using the query method for initialization.
	query := adder.AddBatchQ(filter.Filter(&world))
	for query.Next() {
		head := query.Get()
		head.Angle = rand.Float64() * 360
	}
}

func TestBatchRelations(t *testing.T) {
	world := ecs.NewWorld()

	childID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()

	// Create 100 entities with ChildOf relation.
	builder := ecs.NewBuilder(&world, childID)
	builder.NewBatch(100)

	// Filter for entities with ChildOf.
	filter := ecs.All(childID)
	// Batch-set their relation target to parent.
	world.Batch().SetRelation(filter, childID, parent)
}

func TestBatchRelationsGeneric(t *testing.T) {
	world := ecs.NewWorld()

	childID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()

	// Create 100 entities with ChildOf relation.
	builder := generic.NewMap1[ChildOf](&world)
	builder.NewBatch(100)

	// Create a generic map to perform the batch operation
	mapper := generic.NewMap[ChildOf](&world)
	// Filter for entities with ChildOf.
	filter := ecs.All(childID)
	// Batch-set their relation target to parent.
	mapper.SetRelationBatch(filter, parent)
}

func TestBatchRemoveEntities(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)

	// Create 100 entities with Position.
	builder := ecs.NewBuilder(&world, posID)
	builder.NewBatch(100)

	// Filter for entities with Position.
	filter := ecs.All(posID)
	// Batch-remove matching entities.
	world.Batch().RemoveEntities(filter)
}

func TestBatchRemoveEntitiesGeneric(t *testing.T) {
	world := ecs.NewWorld()

	// Create 100 entities with Position.
	builder := generic.NewMap1[Position](&world)
	builder.NewBatch(100)

	// Filter for entities with Position.
	filter := generic.NewFilter1[Position]()
	// Batch-remove matching entities.
	world.Batch().RemoveEntities(filter.Filter(&world))
}
