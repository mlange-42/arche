package listener_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/mlange-42/arche/listener"
	"github.com/stretchr/testify/assert"
)

func TestDispatched(t *testing.T) {
	entityEvents := []ecs.EntityEvent{}
	componentEvents := []ecs.EntityEvent{}

	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	l1 := listener.NewCallback(
		func(evt ecs.EntityEvent) { entityEvents = append(entityEvents, evt) },
		event.Entities,
		posID,
	)
	l2 := listener.NewCallback(
		func(evt ecs.EntityEvent) { componentEvents = append(componentEvents, evt) },
		event.Components,
		posID,
	)

	ls := listener.NewDispatched(&l1)
	ls.AddListener(&l2)
	world.SetListener(&ls)

	assert.Equal(t, event.Entities|event.Components, ls.Subscriptions())

	world.NewEntity(posID, velID)
	e := world.NewEntity()
	world.Add(e, posID, velID)

	assert.Equal(t, 1, len(entityEvents))
	assert.Equal(t, 2, len(componentEvents))

	world.NewEntity(velID)
	e = world.NewEntity()
	world.Add(e, velID)

	assert.Equal(t, 1, len(entityEvents))
	assert.Equal(t, 2, len(componentEvents))
}

func TestDispatchedRelations(t *testing.T) {
	relationEvents := []ecs.EntityEvent{}
	targetEvents := []ecs.EntityEvent{}

	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	rel1ID := ecs.ComponentID[Relation1](&world)
	rel2ID := ecs.ComponentID[Relation2](&world)

	l1 := listener.NewCallback(
		func(evt ecs.EntityEvent) { relationEvents = append(relationEvents, evt) },
		event.RelationChanged,
		rel1ID,
	)
	l2 := listener.NewCallback(
		func(evt ecs.EntityEvent) { targetEvents = append(targetEvents, evt) },
		event.TargetChanged,
		rel1ID,
	)

	ls := listener.NewDispatched(&l1)
	ls.AddListener(&l2)
	world.SetListener(&ls)

	parent1 := world.NewEntity(posID)
	parent2 := world.NewEntity(posID)

	builder1 := ecs.NewBuilder(&world, posID, rel1ID).WithRelation(rel1ID)
	builder2 := ecs.NewBuilder(&world, posID, rel2ID).WithRelation(rel2ID)

	assert.Equal(t, 0, len(relationEvents))
	assert.Equal(t, 0, len(targetEvents))

	builder1.NewBatch(10)
	builder1.NewBatch(10, parent1)

	assert.Equal(t, 20, len(relationEvents))
	assert.Equal(t, 10, len(targetEvents))

	builder2.NewBatch(10)
	builder2.NewBatch(10, parent2)

	assert.Equal(t, 20, len(relationEvents))
	assert.Equal(t, 10, len(targetEvents))

	world.Batch().SetRelation(ecs.All(rel1ID), rel1ID, parent2)

	assert.Equal(t, 20, len(relationEvents))
	assert.Equal(t, 30, len(targetEvents))

	world.Batch().RemoveEntities(ecs.All(rel1ID))

	assert.Equal(t, 40, len(relationEvents))
	assert.Equal(t, 50, len(targetEvents))
}

func TestDispatchedAllComps(t *testing.T) {
	events1 := []ecs.EntityEvent{}
	events2 := []ecs.EntityEvent{}

	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Position](&world)
	relID := ecs.ComponentID[Relation1](&world)

	l1 := listener.NewCallback(
		func(evt ecs.EntityEvent) { events1 = append(events1, evt) },
		event.Components,
	)
	l2 := listener.NewCallback(
		func(evt ecs.EntityEvent) { events2 = append(events2, evt) },
		event.Components,
		posID, velID,
	)

	ls := listener.NewDispatched(&l1)
	ls.AddListener(&l2)
	world.SetListener(&ls)

	world.NewEntity(relID)

	assert.Equal(t, 1, len(events1))
	assert.Equal(t, 0, len(events2))

	world.NewEntity(posID)

	assert.Equal(t, 2, len(events1))
	assert.Equal(t, 1, len(events2))
}

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type Relation1 struct {
	ecs.Relation
}

type Relation2 struct {
	ecs.Relation
}

func ExampleDispatched() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	entityListener := listener.NewCallback(
		func(evt ecs.EntityEvent) { fmt.Println("Entity event") },
		event.EntityCreated|event.EntityRemoved,
	)
	componentListener := listener.NewCallback(
		func(evt ecs.EntityEvent) { fmt.Println("Component event") },
		event.ComponentAdded|event.ComponentRemoved,
		posID, // optional restriction of subscribed components
	)

	mainListener := listener.NewDispatched(
		&entityListener,
		&componentListener,
	)

	world.SetListener(&mainListener)

	// Triggers entityListener
	e := world.NewEntity()
	// Triggers componentListener
	world.Add(e, posID)

	// Triggers entityListener and componentListener
	_ = world.NewEntity(posID)

	// Triggers entityListener but not componentListener, as it is restricted to posID
	_ = world.NewEntity(velID)
	// Output: Entity event
	// Component event
	// Entity event
	// Component event
	// Entity event
}
