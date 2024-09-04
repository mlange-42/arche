package ecs

import (
	"testing"

	"github.com/mlange-42/arche/ecs/event"
	"github.com/stretchr/testify/assert"
)

func TestWorldListener(t *testing.T) {
	w := NewWorld()

	events := []EntityEvent{}
	listener := newTestListener(func(world *World, e EntityEvent) {
		events = append(events, e)
	})
	w.SetListener(&listener)

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	rotID := ComponentID[rotation](&w)
	relID := ComponentID[relationComp](&w)

	e0 := w.NewEntity()
	assert.Equal(t, 1, len(events))
	assert.Equal(t, EntityEvent{
		Entity:     e0,
		EventTypes: event.EntityCreated,
	}, events[len(events)-1])

	w.RemoveEntity(e0)
	assert.Equal(t, 2, len(events))
	assert.Equal(t, EntityEvent{
		Entity:     e0,
		EventTypes: event.EntityRemoved,
	}, events[len(events)-1])

	e0 = w.NewEntity(posID, velID)
	assert.Equal(t, 3, len(events))
	assert.Equal(t, EntityEvent{
		Entity:     e0,
		Added:      All(posID, velID),
		AddedIDs:   []ID{posID, velID},
		EventTypes: event.EntityCreated | event.ComponentAdded,
	}, events[len(events)-1])

	w.RemoveEntity(e0)
	assert.Equal(t, 4, len(events))
	assert.Equal(t, EntityEvent{
		Entity:     e0,
		Removed:    All(posID, velID),
		RemovedIDs: []ID{posID, velID},
		EventTypes: event.EntityRemoved | event.ComponentRemoved,
	}, events[len(events)-1])

	e0 = w.NewEntityWith(Component{posID, &Position{}}, Component{velID, &Velocity{}}, Component{relID, &relationComp{}})
	assert.Equal(t, 5, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      e0,
		Added:       All(posID, velID, relID),
		AddedIDs:    []ID{posID, velID, relID},
		NewRelation: &relID,
		EventTypes:  event.EntityCreated | event.ComponentAdded | event.RelationChanged | event.TargetChanged,
	}, events[len(events)-1])

	w.Add(e0, rotID)
	assert.Equal(t, 6, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      e0,
		Added:       All(rotID),
		AddedIDs:    []ID{rotID},
		OldRelation: &relID,
		NewRelation: &relID,
		EventTypes:  event.ComponentAdded,
	}, events[len(events)-1])

	w.Remove(e0, posID)
	assert.Equal(t, 7, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      e0,
		Removed:     All(posID),
		RemovedIDs:  []ID{posID},
		OldRelation: &relID,
		NewRelation: &relID,
		EventTypes:  event.ComponentRemoved,
	}, events[len(events)-1])

	e1 := w.NewEntity(posID)
	w.Relations().Set(e0, relID, e1)
	assert.Equal(t, 9, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      e0,
		OldRelation: &relID,
		NewRelation: &relID,
		EventTypes:  event.TargetChanged,
	}, events[len(events)-1])

	w.Remove(e0, relID)
	assert.Equal(t, 10, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      e0,
		Removed:     All(relID),
		RemovedIDs:  []ID{relID},
		OldRelation: &relID,
		NewRelation: nil,
		OldTarget:   e1,
		EventTypes:  event.ComponentRemoved | event.RelationChanged | event.TargetChanged,
	}, events[len(events)-1])

	w.Assign(e0, Component{ID: posID, Comp: &Position{X: 1, Y: 2}})
	assert.Equal(t, 11, len(events))
	assert.Equal(t, EntityEvent{
		Entity:     e0,
		Added:      All(posID),
		AddedIDs:   []ID{posID},
		EventTypes: event.ComponentAdded,
	}, events[len(events)-1])
}

func TestWorldListenerBuilder(t *testing.T) {
	w := NewWorld()

	events := []EntityEvent{}
	listener := newTestListener(func(world *World, e EntityEvent) {
		events = append(events, e)
	})
	w.SetListener(&listener)

	posID := ComponentID[Position](&w)
	relID := ComponentID[relationComp](&w)

	parent := w.NewEntity(posID)

	builder := NewBuilder(&w, posID, relID).WithRelation(relID)
	builder.NewBatch(10)

	assert.Equal(t, 11, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      Entity{11, 0},
		Added:       All(posID, relID),
		AddedIDs:    []ID{posID, relID},
		NewRelation: &relID,
		EventTypes:  event.EntityCreated | event.ComponentAdded | event.RelationChanged | event.TargetChanged,
	}, events[len(events)-1])

	query := builder.NewBatchQ(10)
	query.Close()

	assert.Equal(t, 21, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      Entity{21, 0},
		Added:       All(posID, relID),
		AddedIDs:    []ID{posID, relID},
		NewRelation: &relID,
		EventTypes:  event.EntityCreated | event.ComponentAdded | event.RelationChanged | event.TargetChanged,
	}, events[len(events)-1])

	builder.NewBatch(10, parent)

	assert.Equal(t, 31, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      Entity{31, 0},
		Added:       All(posID, relID),
		AddedIDs:    []ID{posID, relID},
		NewRelation: &relID,
		EventTypes:  event.EntityCreated | event.ComponentAdded | event.RelationChanged | event.TargetChanged,
	}, events[len(events)-1])

	query = builder.NewBatchQ(10, parent)
	query.Close()

	assert.Equal(t, 41, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      Entity{41, 0},
		Added:       All(posID, relID),
		AddedIDs:    []ID{posID, relID},
		NewRelation: &relID,
		EventTypes:  event.EntityCreated | event.ComponentAdded | event.RelationChanged | event.TargetChanged,
	}, events[len(events)-1])

	builder = NewBuilderWith(&w,
		Component{ID: posID, Comp: &Position{}},
		Component{ID: relID, Comp: &relationComp{}},
	).WithRelation(relID)

	builder.NewBatch(10)

	assert.Equal(t, 51, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      Entity{51, 0},
		Added:       All(posID, relID),
		AddedIDs:    []ID{posID, relID},
		NewRelation: &relID,
		EventTypes:  event.EntityCreated | event.ComponentAdded | event.RelationChanged | event.TargetChanged,
	}, events[len(events)-1])

	query = builder.NewBatchQ(10)
	query.Close()

	assert.Equal(t, 61, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      Entity{61, 0},
		Added:       All(posID, relID),
		AddedIDs:    []ID{posID, relID},
		NewRelation: &relID,
		EventTypes:  event.EntityCreated | event.ComponentAdded | event.RelationChanged | event.TargetChanged,
	}, events[len(events)-1])

	builder.NewBatch(10, parent)

	assert.Equal(t, 71, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      Entity{71, 0},
		Added:       All(posID, relID),
		AddedIDs:    []ID{posID, relID},
		NewRelation: &relID,
		EventTypes:  event.EntityCreated | event.ComponentAdded | event.RelationChanged | event.TargetChanged,
	}, events[len(events)-1])

	query = builder.NewBatchQ(10, parent)
	query.Close()

	assert.Equal(t, 81, len(events))
	assert.Equal(t, EntityEvent{
		Entity:      Entity{81, 0},
		Added:       All(posID, relID),
		AddedIDs:    []ID{posID, relID},
		NewRelation: &relID,
		EventTypes:  event.EntityCreated | event.ComponentAdded | event.RelationChanged | event.TargetChanged,
	}, events[len(events)-1])
}
