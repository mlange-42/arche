package ecs

import (
	"reflect"
	"unsafe"
)

// Archetype represents an ECS archetype
type Archetype struct {
	mask       Mask
	indices    []ID
	entities   Storage
	components [MaskTotalBits]Storage
}

var entityType = reflect.TypeOf(Entity{})

// NewArchetype creates a new Archetype for the given components
// Component arguments must be sorted by ID!
func NewArchetype(components ...ComponentType) Archetype {
	var mask Mask
	indices := make([]ID, len(components))
	comps := [MaskTotalBits]Storage{}

	prev := -1
	for i, c := range components {
		if int(c.ID) <= prev {
			panic("component arguments must be sorted by ID")
		}
		prev = int(c.ID)

		mask.Set(c.ID, true)
		indices[i] = c.ID
		comps[c.ID] = NewReflectStorage(c.Type, 32)
	}

	return Archetype{
		mask:       mask,
		indices:    indices,
		entities:   NewReflectStorage(entityType, 32),
		components: comps,
	}
}

// GetEntity returns the entity at the given index
func (a *Archetype) GetEntity(index int) Entity {
	return *(*Entity)(a.entities.Get(uint32(index)))
}

// Get returns the component with the given ID at the given index
func (a *Archetype) Get(index int, id ID) unsafe.Pointer {
	return a.components[id].Get(uint32(index))
}

// Add adds an entity with components to the archetype
func (a *Archetype) Add(entity Entity, components ...Component) uint32 {
	if len(components) != len(a.indices) {
		panic("Invalid number of components")
	}
	idx := a.entities.Add(&entity)
	for _, c := range components {
		a.components[c.ID].Add(c.Component)
	}
	return idx
}

// AddPointer adds an entity with components to the archetype, using pointers
func (a *Archetype) AddPointer(entity Entity, components ...ComponentPointer) uint32 {
	if len(components) != len(a.indices) {
		panic("Invalid number of components")
	}
	idx := a.entities.Add(&entity)
	for _, c := range components {
		if c.Pointer == nil {
			a.components[c.ID].Alloc()
		} else {
			a.components[c.ID].AddPointer(c.Pointer)
		}
	}
	return idx
}

// Remove removes an entity from the archetype
func (a *Archetype) Remove(index int) bool {
	swapped := a.entities.Remove(uint32(index))
	for _, c := range a.indices {
		a.components[c].Remove(uint32(index))
	}
	return swapped
}

// Components returns the component IDs for this archetype
func (a *Archetype) Components() []ID {
	return a.indices
}

// HasComponent returns whether the archetype contains the given component ID
func (a *Archetype) HasComponent(id ID) bool {
	return a.mask.Get(id)
}

// Len reports the number of entities in the archetype
func (a *Archetype) Len() uint32 {
	return a.entities.Len()
}
