package ecs

import "unsafe"

// Archetype is the interface for archetypes
type Archetype interface {
	GetEntity(index int) Entity
	Get(index int, id ID) unsafe.Pointer
}

// NewArchetype creates a new Archetype for the given components
// Component arguments must be sorted by ID!
func NewArchetype(components ...Component) Archetype {
	var mask Mask
	indices := make([]int, MaskTotalBits)
	comps := make([]Storage, len(components))
	for i, c := range components {
		mask.Set(c.ID, true)
		indices[c.ID] = i
		comps[i] = NewReflectStorage(c.reference, 32)
	}

	return &archetype{
		mask:       mask,
		indices:    indices,
		entities:   NewReflectStorage(Entity{}, 32),
		components: comps,
	}
}

type archetype struct {
	mask       Mask
	indices    []int
	entities   Storage
	components []Storage
}

// GetEntity returns the entity at the given index
func (a *archetype) GetEntity(index int) Entity {
	return *(*Entity)(a.entities.Get(uint32(index)))
}

// GetEntity returns the component with the given ID at the given index
func (a *archetype) Get(index int, id ID) unsafe.Pointer {
	idx := a.indices[id]
	return a.components[idx].Get(uint32(index))
}
