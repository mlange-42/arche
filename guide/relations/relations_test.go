package main

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// Position component.
type Position struct {
	X float64
	Y float64
}

// ChildOf relation component
type ChildOf struct {
	ecs.Relation
}

func TestCreateEntity(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	childID := ecs.ComponentID[ChildOf](&world)

	// We set the relation component with WithRelation.
	builder := ecs.NewBuilder(&world, posID, childID).WithRelation(childID)

	_ = builder.New() // An entity with a zero target.

	parent := world.NewEntity()
	_ = builder.New(parent) // An entity with parent as target.
}

func TestCreateEntityGeneric(t *testing.T) {
	world := ecs.NewWorld()

	// The second argument specifies the relation component.
	builder := generic.NewMap2[Position, ChildOf](&world, generic.T[ChildOf]())

	_ = builder.New() // An entity with a zero target.

	parent := world.NewEntity()
	_ = builder.New(parent) // An entity with parent as target.
}

func TestAddRelation(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	childID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()
	child := world.NewEntity(posID)

	world.Relations().Exchange(
		child,             // The entity to modify
		[]ecs.ID{childID}, // Component(s) to add
		nil,               // Component(s) to remove
		childID,           // The relation component of the added components
		parent,            // The target entity
	)
}

func TestAddRelationGeneric(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap1[Position](&world)
	adder := generic.NewMap1[ChildOf](&world, generic.T[ChildOf]())

	parent := world.NewEntity()
	child := builder.New()

	adder.Add(child, parent)
}

func TestAddRelationGenericExchange(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap1[Position](&world)
	adder := generic.NewExchange(&world).
		Adds(generic.T[ChildOf]()).        // Component(s) to add
		WithRelation(generic.T[ChildOf]()) // The relation component of the added components

	parent := world.NewEntity()
	child := builder.New()

	adder.Add(child, parent)
}

func TestSetRelation(t *testing.T) {
	world := ecs.NewWorld()

	childID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()
	child := world.NewEntity(childID)

	world.Relations().Set(child, childID, parent)
}

func TestSetRelationGeneric(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap1[ChildOf](&world)
	mapper := generic.NewMap[ChildOf](&world)

	parent := world.NewEntity()
	child := builder.New()

	mapper.SetRelation(child, parent)
}

func TestGetRelation(t *testing.T) {
	world := ecs.NewWorld()
	childID := ecs.ComponentID[ChildOf](&world)

	child := world.NewEntity(childID)

	_ = world.Relations().Get(child, childID)
}

func TestGetRelationGeneric(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap1[ChildOf](&world)
	mapper := generic.NewMap[ChildOf](&world)

	child := builder.New()

	_ = mapper.GetRelation(child)
}

func TestRelationQuery(t *testing.T) {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	childID := ecs.ComponentID[ChildOf](&world)

	// Two parent entities.
	parent1 := world.NewEntity()
	parent2 := world.NewEntity()

	// A builder with a relation
	builder := ecs.NewBuilder(&world, posID, childID).
		WithRelation(childID)

	// Create 10 entities for each parent.
	for i := 0; i < 10; i++ {
		builder.New(parent1)
		builder.New(parent2)
	}

	// A filter for all entities with Position,
	// and ChildOf with target parent1.
	filter := ecs.NewRelationFilter(ecs.All(posID, childID), parent1)

	query := world.Query(&filter)
	fmt.Println(query.Count()) // Prints 10

	query.Close()
}

func TestRelationQueryGeneric(t *testing.T) {
	world := ecs.NewWorld()

	// Two parent entities.
	parent1 := world.NewEntity()
	parent2 := world.NewEntity()

	// A builder with a relation
	builder := generic.NewMap2[Position, ChildOf](&world, generic.T[ChildOf]())

	// Create 10 entities for each parent.
	for i := 0; i < 10; i++ {
		builder.New(parent1)
		builder.New(parent2)
	}

	// A filter for all entities with Position,
	// and a ChildOf relation.
	filter := generic.NewFilter2[Position, ChildOf]().
		WithRelation(generic.T[ChildOf]())

	// We specify the target when querying.
	// Alternatively, a fixed target can be specified in WithRelation above.
	query := filter.Query(&world, parent1)
	fmt.Println(query.Count()) // Prints 10

	query.Close()
}
