package geckecs

import "github.com/btvoidx/mint"

type EntityCreatedEvent struct {
	Entity Entity
}

type EntityDestroyedEvent struct {
	Entity Entity
}

type UnsubscribeFunc func()

func (w *World) OnEntityCreated(fn func(EntityCreatedEvent)) UnsubscribeFunc {
	stopCh := mint.On(w.eventBus, fn)
	return func() { stopCh() }
}

func (w *World) OnEntityDestroyed(fn func(EntityDestroyedEvent)) UnsubscribeFunc {
	stopCh := mint.On(w.eventBus, fn)
	return func() { stopCh() }
}

func fireEvent[T any](w *World, event T) {
	mint.Emit(w.eventBus, event)
}
