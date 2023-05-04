package ecs

// Interface for an iterator over multiple archetypes.
type nodes interface {
	Get(index int32) *archNode
	Len() int32
}

// Interface for an iterator over archetypes.
type archetypes interface {
	Get(index int32) *archetype
	Len() int32
}

// Implementation of an archetype iterator for a single archetype.
// Implements [archetypes].
type singleArchetype struct {
	Archetype *archetype
}

// Get returns the value at the given index.
func (s singleArchetype) Get(index int32) *archetype {
	return s.Archetype
}

// Len returns the current number of items in the paged array.
func (s singleArchetype) Len() int32 {
	return 1
}

// Implementation of an archetype iterator for a single archetype and partial iteration.
// Implements [archetypes].
//
// Used for the [Query] returned by entity batch creation methods.
type batchArchetype struct {
	Archetype    *archetype
	StartIndex   uint32
	EndIndex     uint32
	OldArchetype *archetype
	Added        []ID
	Removed      []ID
}

// Get returns the value at the given index.
func (s *batchArchetype) Get(index int32) *archetype {
	return s.Archetype
}

// Len returns the current number of items in the paged array.
func (s *batchArchetype) Len() int32 {
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

// Add an element.
func (a *archetypePointers) Add(arch *archetype) {
	a.pointers = append(a.pointers, arch)
}

// Remove an element.
func (a *archetypePointers) Remove(arch *archetype) {
	ln := len(a.pointers)
	for i := 0; i < ln; i++ {
		arch2 := a.pointers[i]
		if arch == arch2 {
			a.pointers[i], a.pointers[ln-1] = a.pointers[ln-1], nil
			a.pointers = a.pointers[:ln-1]
			return
		}
	}
}

// Len returns the current number of items in the paged array.
func (a *archetypePointers) Len() int32 {
	return int32(len(a.pointers))
}
