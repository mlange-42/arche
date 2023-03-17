package ecs

import (
	"reflect"
	"unsafe"

	"github.com/mlange-42/arche/ecs/stats"
)

// archetypeNode is a node in the archetype graph
type archetypeNode struct {
	mask      Mask
	archetype *archetype
	toAdd     []*archetypeNode
	toRemove  []*archetypeNode
}

// Creates a new archetypeNode
func newArchetypeNode(mask Mask) archetypeNode {
	return archetypeNode{
		mask:     mask,
		toAdd:    make([]*archetypeNode, MaskTotalBits),
		toRemove: make([]*archetypeNode, MaskTotalBits),
	}
}

// GetTransitionAdd returns the archetypeNode resulting from adding a component
func (a *archetypeNode) GetTransitionAdd(id ID) (*archetypeNode, bool) {
	p := a.toAdd[id]
	return p, p != nil
}

// GetTransitionRemove returns the archetypeNode resulting from removing a component
func (a *archetypeNode) GetTransitionRemove(id ID) (*archetypeNode, bool) {
	p := a.toRemove[id]
	return p, p != nil
}

// SetTransitionAdd sets the archetypeNode resulting from adding a component
func (a *archetypeNode) SetTransitionAdd(id ID, to *archetypeNode) {
	a.toAdd[id] = to
}

// SetTransitionRemove sets the archetypeNode resulting from removing a component
func (a *archetypeNode) SetTransitionRemove(id ID, to *archetypeNode) {
	a.toRemove[id] = to
}

// archetype represents an ECS archetype
type archetype struct {
	Mask Mask
	Ids  []ID
	// Indirection to avoid a fixed-size array of storages
	// Increases access time by 50-100%
	references []*storage
	entities   genericStorage[Entity]
	components []storage
	graphNode  *archetypeNode
}

// Init initializes an archetype
func (a *archetype) Init(node *archetypeNode, capacityIncrement int, forStorage bool, components ...componentType) {
	var mask Mask
	if len(components) > 0 {
		a.Ids = make([]ID, len(components))
	}
	a.components = make([]storage, len(components))
	a.references = make([]*storage, MaskTotalBits)

	prev := -1
	for i, c := range components {
		if int(c.ID) <= prev {
			panic("component arguments must be sorted by ID")
		}
		prev = int(c.ID)

		mask.Set(c.ID, true)
		a.Ids[i] = c.ID
		a.components[i] = storage{}
		a.components[i].Init(c.Type, capacityIncrement, forStorage)
		a.references[c.ID] = &a.components[i]
	}

	a.graphNode = node
	a.Mask = mask
	a.entities = genericStorage[Entity]{}
	a.entities.Init(capacityIncrement, forStorage)
}

// GetEntity returns the entity at the given index
func (a *archetype) GetEntity(index uint32) Entity {
	return a.entities.Get(index)
}

// Get returns the component with the given ID at the given index
func (a *archetype) Get(index uint32, id ID) unsafe.Pointer {
	ref := a.references[id]
	if ref != nil {
		return ref.Get(index)
	}
	return nil
}

// Add adds an entity with zeroed components to the archetype
func (a *archetype) Alloc(entity Entity, zero bool) uint32 {
	idx := a.entities.Add(entity)
	len := len(a.components)

	for i := 0; i < len; i++ {
		comp := &a.components[i]
		idx := comp.Alloc()
		if zero {
			comp.Zero(idx)
		}
	}
	return idx
}

// Add adds an entity with components to the archetype
func (a *archetype) Add(entity Entity, components ...Component) uint32 {
	if len(components) != len(a.Ids) {
		panic("Invalid number of components")
	}
	idx := a.entities.Add(entity)
	for _, c := range components {
		a.references[c.ID].Add(c.Comp)
	}
	return idx
}

// Remove removes an entity from the archetype
func (a *archetype) Remove(index uint32) bool {
	swapped := a.entities.Remove(index)
	len := len(a.components)
	for i := 0; i < len; i++ {
		a.components[i].Remove(index)
	}
	return swapped
}

// Components returns the component IDs for this archetype
func (a *archetype) Components() []ID {
	return a.Ids
}

// HasComponent returns whether the archetype contains the given component ID
func (a *archetype) HasComponent(id ID) bool {
	return a.references[id] != nil
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
	return a.references[id].Set(index, comp)
}

// SetPointer overwrites a component with the data behind the given pointer
func (a *archetype) SetPointer(index uint32, id ID, comp unsafe.Pointer) unsafe.Pointer {
	return a.references[id].SetPointer(index, comp)
}

// Zero resets th memory at the given position
func (a *archetype) Zero(index uint32, id ID) {
	a.references[id].Zero(index)
}

// Stats generates statistics for an archetype
func (a *archetype) Stats(reg *componentRegistry) stats.ArchetypeStats {
	ids := a.Components()
	aCompCount := len(ids)
	aTypes := make([]reflect.Type, aCompCount)
	for j, id := range ids {
		aTypes[j] = reg.ComponentType(id)
	}

	cap := int(a.Cap())
	memPerEntity := 0
	for i := 0; i < len(a.components); i++ {
		comp := &a.components[i]
		memPerEntity += int(comp.itemSize)
	}
	memory := cap * (int(entitySize) + memPerEntity)

	return stats.ArchetypeStats{
		Size:            int(a.Len()),
		Capacity:        cap,
		Components:      aCompCount,
		ComponentIDs:    ids,
		ComponentTypes:  aTypes,
		Memory:          memory,
		MemoryPerEntity: memPerEntity,
	}
}
