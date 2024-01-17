package ecs_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/stretchr/testify/assert"
)

func TestEntityEvent(t *testing.T) {
	e := ecs.EntityEvent{EventTypes: event.ComponentAdded}

	assert.False(t, e.Contains(event.EntityCreated))
	assert.False(t, e.Contains(event.EntityRemoved))

	e = ecs.EntityEvent{EventTypes: event.EntityCreated | event.ComponentAdded}

	assert.True(t, e.Contains(event.EntityCreated))
	assert.False(t, e.Contains(event.EntityRemoved))

	e = ecs.EntityEvent{EventTypes: event.EntityRemoved | event.ComponentRemoved}

	assert.False(t, e.Contains(event.EntityCreated))
	assert.True(t, e.Contains(event.EntityRemoved))
}

type eventHandler struct {
	LastEntity ecs.Entity
}

func (h *eventHandler) ListenCopy(e ecs.EntityEvent) {
	h.LastEntity = e.Entity
}

func (h *eventHandler) ListenPointer(e *ecs.EntityEvent) {
	h.LastEntity = e.Entity
}

func BenchmarkEntityEventCreate(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	e := world.NewEntity()
	mask := ecs.All(posID)
	added := []ecs.ID{posID}

	var event ecs.EntityEvent

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		event = ecs.EntityEvent{Entity: e, Added: mask, Removed: mask, AddedIDs: added, RemovedIDs: nil}
	}
	b.StopTimer()
	_ = event
}

func BenchmarkEntityEventHeapPointer(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	e := world.NewEntity()
	mask := ecs.All(posID)
	added := []ecs.ID{posID}

	var event *ecs.EntityEvent

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		event = &ecs.EntityEvent{Entity: e, Added: mask, Removed: mask, AddedIDs: added, RemovedIDs: nil}
	}
	b.StopTimer()
	_ = event
}

func BenchmarkEntityEventCopy(b *testing.B) {
	handler := eventHandler{}

	for i := 0; i < b.N; i++ {
		handler.ListenCopy(ecs.EntityEvent{Entity: ecs.Entity{}, Added: ecs.Mask{}, Removed: ecs.Mask{}, AddedIDs: nil, RemovedIDs: nil})
	}
}

func BenchmarkEntityEventCopyReuse(b *testing.B) {
	handler := eventHandler{}
	event := ecs.EntityEvent{Entity: ecs.Entity{}, Added: ecs.Mask{}, Removed: ecs.Mask{}, AddedIDs: nil, RemovedIDs: nil}

	for i := 0; i < b.N; i++ {
		handler.ListenCopy(event)
	}
}

func BenchmarkEntityEventPointer(b *testing.B) {
	handler := eventHandler{}

	for i := 0; i < b.N; i++ {
		handler.ListenPointer(&ecs.EntityEvent{Entity: ecs.Entity{}, Added: ecs.Mask{}, Removed: ecs.Mask{}, AddedIDs: nil, RemovedIDs: nil})
	}
}

func BenchmarkEntityEventPointerReuse(b *testing.B) {
	handler := eventHandler{}
	event := ecs.EntityEvent{Entity: ecs.Entity{}, Added: ecs.Mask{}, Removed: ecs.Mask{}, AddedIDs: nil, RemovedIDs: nil}

	for i := 0; i < b.N; i++ {
		handler.ListenPointer(&event)
	}
}

func ExampleEntityEvent() {
	world := ecs.NewWorld()

	listener := TestListener{
		Callback: func(world *ecs.World, evt ecs.EntityEvent) { fmt.Println(evt.Entity) },
	}
	world.SetListener(&listener)

	world.NewEntity()
	// Output: {1 0}
}

func ExampleEntityEvent_Contains() {
	evt := ecs.EntityEvent{EventTypes: event.EntityCreated | event.EntityRemoved}

	fmt.Println(evt.Contains(event.EntityRemoved))
	fmt.Println(evt.Contains(event.RelationChanged))
	// Output: true
	// false
}
