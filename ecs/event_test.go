package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityEvent(t *testing.T) {
	e := EntityEvent{AddedRemoved: 0}

	assert.False(t, e.EntityAdded())
	assert.False(t, e.EntityRemoved())

	e = EntityEvent{AddedRemoved: 1}

	assert.True(t, e.EntityAdded())
	assert.False(t, e.EntityRemoved())

	e = EntityEvent{AddedRemoved: -1}

	assert.False(t, e.EntityAdded())
	assert.True(t, e.EntityRemoved())
}

type eventHandler struct {
	LastEntity Entity
}

func (h *eventHandler) ListenCopy(e EntityEvent) {
	h.LastEntity = e.Entity
}

func (h *eventHandler) ListenPointer(e *EntityEvent) {
	h.LastEntity = e.Entity
}

func BenchmarkEntityEventCopy(b *testing.B) {
	handler := eventHandler{}

	for i := 0; i < b.N; i++ {
		handler.ListenCopy(EntityEvent{Entity: Entity{}, OldMask: Mask{}, NewMask: Mask{}, Added: nil, Removed: nil, Current: nil, AddedRemoved: 0})
	}
}

func BenchmarkEntityEventCopyReuse(b *testing.B) {
	handler := eventHandler{}
	event := EntityEvent{Entity: Entity{}, OldMask: Mask{}, NewMask: Mask{}, Added: nil, Removed: nil, Current: nil, AddedRemoved: 0}

	for i := 0; i < b.N; i++ {
		handler.ListenCopy(event)
	}
}

func BenchmarkEntityEventPointer(b *testing.B) {
	handler := eventHandler{}

	for i := 0; i < b.N; i++ {
		handler.ListenPointer(&EntityEvent{Entity: Entity{}, OldMask: Mask{}, NewMask: Mask{}, Added: nil, Removed: nil, Current: nil, AddedRemoved: 0})
	}
}

func BenchmarkEntityEventPointerReuse(b *testing.B) {
	handler := eventHandler{}
	event := EntityEvent{Entity: Entity{}, OldMask: Mask{}, NewMask: Mask{}, Added: nil, Removed: nil, Current: nil, AddedRemoved: 0}

	for i := 0; i < b.N; i++ {
		handler.ListenPointer(&event)
	}
}
