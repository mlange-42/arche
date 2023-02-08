package ecs

import (
	"reflect"
	"unsafe"
)

// archetype represents an ECS archetype
type archetype struct {
	mask       Mask
	indices    []ID
	entities   storage
	components [MaskTotalBits]storage
}

var entityType = reflect.TypeOf(Entity{})

func (a *archetype) init(capacityIncrement int, components ...componentType) {
	var mask Mask
	a.indices = make([]ID, len(components))
	comps := [MaskTotalBits]storage{}

	prev := -1
	for i, c := range components {
		if int(c.ID) <= prev {
			panic("component arguments must be sorted by ID")
		}
		prev = int(c.ID)

		mask.Set(c.ID, true)
		a.indices[i] = c.ID
		comps[c.ID] = storage{}
		comps[c.ID].init(c.Type, capacityIncrement)
	}

	a.mask = mask
	a.components = comps
	a.entities = storage{}
	a.entities.init(entityType, capacityIncrement)
}

// GetEntity returns the entity at the given index
func (a *archetype) GetEntity(index int) Entity {
	return *(*Entity)(a.entities.Get(uint32(index)))
}

// Get returns the component with the given ID at the given index
func (a *archetype) Get(index int, id ID) unsafe.Pointer {
	return a.components[id].Get(uint32(index))
}

// Add adds an entity with components to the archetype
func (a *archetype) Add(entity Entity, components ...component) uint32 {
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
func (a *archetype) AddPointer(entity Entity, components ...componentPointer) uint32 {
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
func (a *archetype) Remove(index int) bool {
	swapped := a.entities.Remove(uint32(index))
	for _, c := range a.indices {
		a.components[c].Remove(uint32(index))
	}
	return swapped
}

// Components returns the component IDs for this archetype
func (a *archetype) Components() []ID {
	return a.indices
}

// HasComponent returns whether the archetype contains the given component ID
func (a *archetype) HasComponent(id ID) bool {
	return a.mask.Get(id)
}

// Len reports the number of entities in the archetype
func (a *archetype) Len() uint32 {
	return a.entities.Len()
}
