package base

import (
	"reflect"
	"unsafe"
)

// Archetype represents an ECS Archetype
type Archetype struct {
	Mask BitMask
	Ids  []ID
	// Indirection to avoid a fixed-size array of storages
	// Increases access time by 50-100%
	indices    [MaskTotalBits]uint8
	entities   Storage
	components []Storage
	toAdd      map[ID]*Archetype
	toRemove   map[ID]*Archetype
}

var entityType = reflect.TypeOf(Entity{})

// Init initializes an archetype
func (a *Archetype) Init(capacityIncrement int, components ...ComponentType) {
	var mask BitMask
	a.Ids = make([]ID, len(components))
	comps := make([]Storage, len(components))

	prev := -1
	for i, c := range components {
		if int(c.ID) <= prev {
			panic("component arguments must be sorted by ID")
		}
		prev = int(c.ID)

		mask.Set(c.ID, true)
		a.Ids[i] = c.ID
		a.indices[c.ID] = uint8(i)
		comps[i] = Storage{}
		comps[i].init(c.Type, capacityIncrement)
	}

	a.Mask = mask
	a.components = comps
	a.entities = Storage{}
	a.toAdd = map[ID]*Archetype{}
	a.toRemove = map[ID]*Archetype{}
	a.entities.init(entityType, capacityIncrement)
}

// GetEntity returns the entity at the given index
func (a *Archetype) GetEntity(index int) Entity {
	return *(*Entity)(a.entities.Get(uint32(index)))
}

// Get returns the component with the given ID at the given index
func (a *Archetype) Get(index int, id ID) unsafe.Pointer {
	if !a.Mask.Get(id) {
		return nil
	}
	return a.components[a.indices[id]].Get(uint32(index))
}

// GetUnsafe returns the component with the given ID at the given index,
// without checking if the entity contains that component.
//
// This is used by queries, where the entity is guaranteed to be in the archetype.
func (a *Archetype) GetUnsafe(index int, id ID) unsafe.Pointer {
	return a.components[a.indices[id]].Get(uint32(index))
}

// Add adds an entity with components to the archetype
func (a *Archetype) Add(entity Entity, components ...Component) uint32 {
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
func (a *Archetype) AddPointer(entity Entity, components ...ComponentPointer) uint32 {
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
func (a *Archetype) Remove(index int) bool {
	swapped := a.entities.Remove(uint32(index))
	len := len(a.components)
	for i := 0; i < len; i++ {
		a.components[i].Remove(uint32(index))
	}
	return swapped
}

// Components returns the component IDs for this archetype
func (a *Archetype) Components() []ID {
	return a.Ids
}

// HasComponent returns whether the archetype contains the given component ID
func (a *Archetype) HasComponent(id ID) bool {
	return a.Mask.Get(id)
}

// Len reports the number of entities in the archetype
func (a *Archetype) Len() uint32 {
	return a.entities.Len()
}

// Set overwrites a component with the data behind the given pointer
func (a *Archetype) Set(index uint32, id ID, comp interface{}) unsafe.Pointer {
	return a.components[a.indices[id]].set(index, comp)
}

// GetTransitionAdd returns the archetype resulting from adding a component
func (a *Archetype) GetTransitionAdd(id ID) (*Archetype, bool) {
	p, ok := a.toAdd[id]
	return p, ok
}

// GetTransitionRemove returns the archetype resulting from removing a component
func (a *Archetype) GetTransitionRemove(id ID) (*Archetype, bool) {
	p, ok := a.toRemove[id]
	return p, ok
}

// SetTransitionAdd sets the archetype resulting from adding a component
func (a *Archetype) SetTransitionAdd(id ID, to *Archetype) {
	a.toAdd[id] = to
}

// SetTransitionRemove sets the archetype resulting from removing a component
func (a *Archetype) SetTransitionRemove(id ID, to *Archetype) {
	a.toRemove[id] = to
}
