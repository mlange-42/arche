package ecs

import (
	"reflect"
	"unsafe"
)

// archetype represents an ECS archetype
type archetype struct {
	Mask BitMask
	Ids  []ID
	// Indirection to avoid a fixed-size array of storages
	// Increases access time by 50-100%
	indices    [MaskTotalBits]uint8
	entities   storage
	components []storage
	toAdd      map[ID]*archetype
	toRemove   map[ID]*archetype
}

var entityType = reflect.TypeOf(Entity{})

// Init initializes an archetype
func (a *archetype) Init(capacityIncrement int, components ...componentType) {
	var mask BitMask
	a.Ids = make([]ID, len(components))
	comps := make([]storage, len(components))

	prev := -1
	for i, c := range components {
		if int(c.ID) <= prev {
			panic("component arguments must be sorted by ID")
		}
		prev = int(c.ID)

		mask.Set(c.ID, true)
		a.Ids[i] = c.ID
		a.indices[c.ID] = uint8(i)
		comps[i] = storage{}
		comps[i].Init(c.Type, capacityIncrement)
	}

	a.Mask = mask
	a.components = comps
	a.entities = storage{}
	a.toAdd = map[ID]*archetype{}
	a.toRemove = map[ID]*archetype{}
	a.entities.Init(entityType, capacityIncrement)
}

// GetEntity returns the entity at the given index
func (a *archetype) GetEntity(index int) Entity {
	return *(*Entity)(a.entities.Get(uint32(index)))
}

// Get returns the component with the given ID at the given index
func (a *archetype) Get(index int, id ID) unsafe.Pointer {
	if !a.Mask.Get(id) {
		return nil
	}
	return a.components[a.indices[id]].Get(uint32(index))
}

// GetUnsafe returns the component with the given ID at the given index,
// without checking if the entity contains that component.
//
// This is used by queries, where the entity is guaranteed to be in the archetype.
func (a *archetype) GetUnsafe(index int, id ID) unsafe.Pointer {
	return a.components[a.indices[id]].Get(uint32(index))
}

// Add adds an entity with components to the archetype
func (a *archetype) Add(entity Entity, components ...Component) uint32 {
	if len(components) != len(a.Ids) {
		panic("Invalid number of components")
	}
	idx := a.entities.Add(&entity)
	for _, c := range components {
		a.components[a.indices[c.ID]].Add(c.Component)
	}
	return idx
}

// AddPointer adds an entity with components to the archetype, using pointers
func (a *archetype) AddPointer(entity Entity, components ...componentPointer) uint32 {
	if len(components) != len(a.Ids) {
		panic("Invalid number of components")
	}
	idx := a.entities.Add(&entity)
	for _, c := range components {
		if c.Pointer == nil {
			a.components[a.indices[c.ID]].Alloc()
		} else {
			a.components[a.indices[c.ID]].AddPointer(c.Pointer)
		}
	}
	return idx
}

// Remove removes an entity from the archetype
func (a *archetype) Remove(index int) bool {
	swapped := a.entities.Remove(uint32(index))
	len := len(a.components)
	for i := 0; i < len; i++ {
		a.components[i].Remove(uint32(index))
	}
	return swapped
}

// Components returns the component IDs for this archetype
func (a *archetype) Components() []ID {
	return a.Ids
}

// HasComponent returns whether the archetype contains the given component ID
func (a *archetype) HasComponent(id ID) bool {
	return a.Mask.Get(id)
}

// Len reports the number of entities in the archetype
func (a *archetype) Len() uint32 {
	return a.entities.Len()
}

// Cap reports the current capacity of the archetype
func (a *archetype) Cap() uint32 {
	return a.entities.Cap()
}

// Set overwrites a component with the data behind the given pointer
func (a *archetype) Set(index uint32, id ID, comp interface{}) unsafe.Pointer {
	return a.components[a.indices[id]].Set(index, comp)
}

// GetTransitionAdd returns the archetype resulting from adding a component
func (a *archetype) GetTransitionAdd(id ID) (*archetype, bool) {
	p, ok := a.toAdd[id]
	return p, ok
}

// GetTransitionRemove returns the archetype resulting from removing a component
func (a *archetype) GetTransitionRemove(id ID) (*archetype, bool) {
	p, ok := a.toRemove[id]
	return p, ok
}

// SetTransitionAdd sets the archetype resulting from adding a component
func (a *archetype) SetTransitionAdd(id ID, to *archetype) {
	a.toAdd[id] = to
}

// SetTransitionRemove sets the archetype resulting from removing a component
func (a *archetype) SetTransitionRemove(id ID, to *archetype) {
	a.toRemove[id] = to
}
