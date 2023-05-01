package ecs

// Interface for an iterator over multiple archetypes.
type nodes interface {
	Get(index int32) *archetypeNode
	Len() int32
}

// Interface for an iterator over archetypes.
type archetypes interface {
	Get(index int32) *archetype
	Len() int32
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
func (s batchArchetype) Get(index int32) *archetype {
	return s.Archetype
}

// Len returns the current number of items in the paged array.
func (s batchArchetype) Len() int32 {
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
func (a *archetypePointers) Get(index int32) *archetype {
	return a.pointers[index]
}

// Add adds an element.
func (a *archetypePointers) Add(arch *archetype) {
	a.pointers = append(a.pointers, arch)
}

// Add adds an element.
func (a *archetypePointers) Remove(arch *archetype) {
	for i := 0; i < len(a.pointers); i++ {
		arch2 := a.pointers[i]
		if arch == arch2 {
			a.pointers = append(a.pointers[:i], a.pointers[i+1:]...)
			return
		}
	}
}

// Len returns the current number of items in the paged array.
func (a *archetypePointers) Len() int32 {
	return int32(len(a.pointers))
}
