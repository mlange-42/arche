package ecs_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestEntityEvent(t *testing.T) {
	e := ecs.EntityEvent{AddedRemoved: 0}

	assert.False(t, e.EntityAdded())
	assert.False(t, e.EntityRemoved())

	e = ecs.EntityEvent{AddedRemoved: 1}

	assert.True(t, e.EntityAdded())
	assert.False(t, e.EntityRemoved())

	e = ecs.EntityEvent{AddedRemoved: -1}

	assert.False(t, e.EntityAdded())
	assert.True(t, e.EntityRemoved())
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

func BenchmarkEntityEventCopy(b *testing.B) {
	handler := eventHandler{}

	for i := 0; i < b.N; i++ {
		handler.ListenCopy(ecs.EntityEvent{Entity: ecs.Entity{}, OldMask: ecs.Mask{}, NewMask: ecs.Mask{}, Added: nil, Removed: nil, Current: nil, AddedRemoved: 0})
	}
}

func BenchmarkEntityEventCopyReuse(b *testing.B) {
	handler := eventHandler{}
	event := ecs.EntityEvent{Entity: ecs.Entity{}, OldMask: ecs.Mask{}, NewMask: ecs.Mask{}, Added: nil, Removed: nil, Current: nil, AddedRemoved: 0}

	for i := 0; i < b.N; i++ {
		handler.ListenCopy(event)
	}
}

func BenchmarkEntityEventPointer(b *testing.B) {
	handler := eventHandler{}

	for i := 0; i < b.N; i++ {
		handler.ListenPointer(&ecs.EntityEvent{Entity: ecs.Entity{}, OldMask: ecs.Mask{}, NewMask: ecs.Mask{}, Added: nil, Removed: nil, Current: nil, AddedRemoved: 0})
	}
}

func BenchmarkEntityEventPointerReuse(b *testing.B) {
	handler := eventHandler{}
	event := ecs.EntityEvent{Entity: ecs.Entity{}, OldMask: ecs.Mask{}, NewMask: ecs.Mask{}, Added: nil, Removed: nil, Current: nil, AddedRemoved: 0}

	for i := 0; i < b.N; i++ {
		handler.ListenPointer(&event)
	}
}

func ExampleEntityEvent() {
	world := ecs.NewWorld()

	listener := func(evt *ecs.EntityEvent) {
		fmt.Println(evt)
	}
	world.SetListener(listener)

	world.NewEntity()
	// Output: &{{1 0} {0 0} {0 0} [] [] [] 1 {0 0} {0 0} false}
}
