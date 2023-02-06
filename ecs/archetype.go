package ecs

import (
	"unsafe"
)

// Archetype is the interface for archetypes
type Archetype interface {
	GetEntity(index int) Entity
	Get(index int, id ID) unsafe.Pointer
	Add(entity Entity, components ...Component)
	Remove(index int)
}

// NewArchetype creates a new Archetype for the given components
// Component arguments must be sorted by ID!
func NewArchetype(components ...Component) Archetype {
	return newArchetype(components...)
}

type archetype struct {
	mask       Mask
	indices    []ID
	entities   ReflectStorage
	components []ReflectStorage
}

func newArchetype(components ...Component) *archetype {
	var mask Mask
	indices := make([]ID, len(components))
	comps := make([]ReflectStorage, MaskTotalBits)
	for i, c := range components {
		mask.Set(c.ID, true)
		indices[i] = c.ID
		comps[c.ID] = NewReflectStorage(c.Component, 32)
	}

	return &archetype{
		mask:       mask,
		indices:    indices,
		entities:   NewReflectStorage(Entity{}, 32),
		components: comps,
	}
}

// GetEntity returns the entity at the given index
func (a *archetype) GetEntity(index int) Entity {
	return *(*Entity)(a.entities.Get(uint32(index)))
}

// GetEntity returns the component with the given ID at the given index
func (a *archetype) Get(index int, id ID) unsafe.Pointer {
	return a.components[id].Get(uint32(index))
}

// Add adds an entity with components to the archetype
func (a *archetype) Add(entity Entity, components ...Component) {
	if len(components) != len(a.indices) {
		panic("Invalid number of components")
	}
	a.entities.Add(&entity)
	for _, c := range components {
		a.components[c.ID].Add(c.Component)
	}
}

// Remove removes an entity from the archetype
func (a *archetype) Remove(index int) {
	a.entities.Remove(uint32(index))
	for _, c := range a.indices {
		a.components[c].Remove(uint32(index))
	}
}
