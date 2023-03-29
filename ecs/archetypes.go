package ecs

// Interface for an iterator over archetypes.
type archetypes interface {
	Get(index int) *archetype
	Len() int
}

// Implementation of an archetype iterator for a single archetype.
// Implements [archetypes].
//
// Used for the [Query] returned by entity batch creation methods.
type batchArchetype struct {
	Archetype  *archetype
	StartIndex uint32
}

// Get returns the value at the given index.
func (s batchArchetype) Get(index int) *archetype {
	return s.Archetype
}

// Len returns the current number of items in the paged array.
func (s batchArchetype) Len() int {
	return 1
}

// Implementation of an archetype iterator for pointers.
// Implements [archetypes].
//
// Used for tracking filter archetypes in [Cache].
type archetypePointers struct {
	pointers []*archetype
}

// Get returns the value at the given index.
func (a *archetypePointers) Get(index int) *archetype {
	return a.pointers[index]
}

// Add adds an element.
func (a *archetypePointers) Add(arch *archetype) {
	a.pointers = append(a.pointers, arch)
}

// Len returns the current number of items in the paged array.
func (a *archetypePointers) Len() int {
	return len(a.pointers)
}
