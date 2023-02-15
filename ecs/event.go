package ecs

// ChangeEvent contains information about component changes.
//
// To receive change events, register a function func(e ChangeEvent) with [World.RegisterListener].
type ChangeEvent struct {
	// The entity that was changed.
	Entity Entity
	// The old and new component masks.
	OldMask, NewMask BitMask
	// Components added, removed, and after the change.
	Added, Removed, Current []ID
	// Whether the entity itself was added (> 0), removed (< 0), or only changed (= 0).
	AddedRemoved int
}

// EntityAdded reports whether the entity was newly added.
func (e *ChangeEvent) EntityAdded() bool {
	return e.AddedRemoved > 0
}

// EntityRemoved reports whether the entity was removed.
func (e *ChangeEvent) EntityRemoved() bool {
	return e.AddedRemoved < 0
}
