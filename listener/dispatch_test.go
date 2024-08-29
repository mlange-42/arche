package listener_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/mlange-42/arche/listener"
	"github.com/stretchr/testify/assert"
)

type EventHandler struct {
	events []ecs.EntityEvent
}

func (h *EventHandler) Notify(w *ecs.World, e ecs.EntityEvent) {
	h.events = append(h.events, e)
}

func TestDispatch(t *testing.T) {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	h1 := EventHandler{}
	l1 := listener.NewCallback(
		h1.Notify,
		event.Entities,
		posID,
	)
	h2 := EventHandler{}
	l2 := listener.NewCallback(
		h2.Notify,
		event.Components,
		posID,
	)
	h3 := EventHandler{}
	l3 := listener.NewCallback(
		h3.Notify,
		event.Entities,
	)

	ls := listener.NewDispatch(&l1)
	ls.AddListener(&l2)
	world.SetListener(&ls)

	assert.Equal(t, event.Entities|event.Components, ls.Subscriptions())
	assert.NotNil(t, ls.Components())

	world.NewEntity(posID, velID)
	e := world.NewEntity()
	world.Add(e, posID, velID)

	assert.Equal(t, 1, len(h1.events))
	assert.Equal(t, 2, len(h2.events))

	world.NewEntity(velID)
	e = world.NewEntity()
	world.Add(e, velID)

	assert.Equal(t, 1, len(h1.events))
	assert.Equal(t, 2, len(h2.events))

	ls.AddListener(&l3)

	assert.Equal(t, event.Entities|event.Components, ls.Subscriptions())
	assert.Nil(t, ls.Components())
}

func TestDispatchRelations(t *testing.T) {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	rel1ID := ecs.ComponentID[Relation1](&world)
	rel2ID := ecs.ComponentID[Relation2](&world)

	h1 := EventHandler{}
	l1 := listener.NewCallback(
		h1.Notify,
		event.RelationChanged,
		rel1ID,
	)
	h2 := EventHandler{}
	l2 := listener.NewCallback(
		h2.Notify,
		event.TargetChanged,
		rel1ID,
	)

	ls := listener.NewDispatch(&l1)
	ls.AddListener(&l2)
	world.SetListener(&ls)

	parent1 := world.NewEntity(posID)
	parent2 := world.NewEntity(posID)

	builder1 := ecs.NewBuilder(&world, posID, rel1ID).WithRelation(rel1ID)
	builder2 := ecs.NewBuilder(&world, posID, rel2ID).WithRelation(rel2ID)

	assert.Equal(t, 0, len(h1.events))
	assert.Equal(t, 0, len(h2.events))

	builder1.NewBatch(10)
	builder1.NewBatch(10, parent1)

	assert.Equal(t, 20, len(h1.events))
	assert.Equal(t, 20, len(h2.events))

	builder2.NewBatch(10)
	builder2.NewBatch(10, parent2)

	assert.Equal(t, 20, len(h1.events))
	assert.Equal(t, 20, len(h2.events))

	world.Batch().SetRelation(ecs.All(rel1ID), rel1ID, parent2)

	assert.Equal(t, 20, len(h1.events))
	assert.Equal(t, 40, len(h2.events))

	world.Batch().RemoveEntities(ecs.All(rel1ID))

	assert.Equal(t, 40, len(h1.events))
	assert.Equal(t, 60, len(h2.events))
}

func TestDispatchAllComps(t *testing.T) {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Position](&world)
	relID := ecs.ComponentID[Relation1](&world)

	h1 := EventHandler{}
	l1 := listener.NewCallback(
		h1.Notify,
		event.Components,
	)
	h2 := EventHandler{}
	l2 := listener.NewCallback(
		h2.Notify,
		event.Components,
		posID, velID,
	)

	ls := listener.NewDispatch(&l1)
	ls.AddListener(&l2)
	world.SetListener(&ls)

	world.NewEntity(relID)

	assert.Equal(t, 1, len(h1.events))
	assert.Equal(t, 0, len(h2.events))

	world.NewEntity(posID)

	assert.Equal(t, 2, len(h1.events))
	assert.Equal(t, 1, len(h2.events))
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

func ExampleDispatch() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	entityListener := listener.NewCallback(
		func(w *ecs.World, ee ecs.EntityEvent) { fmt.Println("Entity event") },
		event.EntityCreated|event.EntityRemoved,
	)
	componentListener := listener.NewCallback(
		func(w *ecs.World, ee ecs.EntityEvent) { fmt.Println("Component event on Position") },
		event.ComponentAdded|event.ComponentRemoved,
		posID, // optional restriction of subscribed components
	)

	mainListener := listener.NewDispatch(
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
	// Component event on Position
	// Entity event
	// Component event on Position
	// Entity event
}
