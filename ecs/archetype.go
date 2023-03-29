package ecs

import (
	"math"
	"reflect"
	"unsafe"

	"github.com/mlange-42/arche/ecs/stats"
)

// layoutSize is the size of an archetype column layout in bytes.
var layoutSize = unsafe.Sizeof(layout{})

// archetypeNode is a node in the archetype graph.
type archetypeNode struct {
	mask      Mask             // Mask of the archetype
	archetype *archetype       // The archetype
	toAdd     []*archetypeNode // Mapping from component ID to add to the resulting archetype
	toRemove  []*archetypeNode // Mapping from component ID to remove to the resulting archetype
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

// Interface for an iterator over archetypes
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

// Helper for accessing data from an archetype
type archetypeAccess struct {
	basePointer   unsafe.Pointer // Pointer to the first component column layout.
	entityPointer unsafe.Pointer // Pointer to the entity storage
	Mask          Mask           // Archetype's mask
}

// GetEntity returns the entity at the given index
func (a *archetypeAccess) GetEntity(index uintptr) Entity {
	return *(*Entity)(unsafe.Add(a.entityPointer, entitySize*index))
}

// Get returns the component with the given ID at the given index
func (a *archetypeAccess) Get(index uintptr, id ID) unsafe.Pointer {
	return a.getLayout(id).Get(index)
}

// HasComponent returns whether the archetype contains the given component ID
func (a *archetypeAccess) HasComponent(id ID) bool {
	return a.getLayout(id).pointer != nil
}

// GetLayout returns the column layout for a component.
func (a *archetypeAccess) getLayout(id ID) *layout {
	return (*layout)(unsafe.Add(a.basePointer, layoutSize*uintptr(id)))
}

// layout specification of a component column.
type layout struct {
	pointer  unsafe.Pointer // Pointer to the first element in the component column.
	itemSize uintptr        // Component/step size
}

// Get returns a pointer to the item at the given index.
func (l *layout) Get(index uintptr) unsafe.Pointer {
	if l.pointer == nil {
		return nil
	}
	return unsafe.Add(l.pointer, l.itemSize*index)
}

// archetype represents an ECS archetype
type archetype struct {
	archetypeAccess                   // Access helper, passed to queries.
	graphNode         *archetypeNode  // Node in the archetype graph.
	Ids               []ID            // List of component IDs.
	layouts           []layout        // Column layouts by ID.
	indices           []uint32        // Mapping from IDs to buffer indices.
	buffers           []reflect.Value // Reflection arrays containing component data.
	entityBuffer      reflect.Value   // Reflection array containing entity data.
	len               uint32          // Current number of entities
	cap               uint32          // Current capacity
	capacityIncrement uint32          // Capacity increment
	zeroValue         []byte          // Used as source for setting storage to zero
	zeroPointer       unsafe.Pointer  // Points to zeroValue for fast access
}

// Init initializes an archetype
func (a *archetype) Init(node *archetypeNode, capacityIncrement int, forStorage bool, components ...componentType) {
	var mask Mask
	if len(components) > 0 {
		a.Ids = make([]ID, len(components))
	}

	a.buffers = make([]reflect.Value, len(components))
	a.layouts = make([]layout, MaskTotalBits)
	a.indices = make([]uint32, MaskTotalBits)

	cap := 1
	if forStorage {
		cap = capacityIncrement
	}

	prev := -1
	var maxSize uintptr = 0
	for i, c := range components {
		if int(c.ID) <= prev {
			panic("component arguments must be sorted by ID")
		}
		prev = int(c.ID)
		mask.Set(c.ID, true)

		size, align := c.Type.Size(), uintptr(c.Type.Align())
		size = (size + (align - 1)) / align * align
		if size > maxSize {
			maxSize = size
		}

		a.Ids[i] = c.ID
		a.buffers[i] = reflect.New(reflect.ArrayOf(cap, c.Type)).Elem()
		a.layouts[c.ID] = layout{
			a.buffers[i].Addr().UnsafePointer(),
			size,
		}
		a.indices[c.ID] = uint32(i)
	}
	a.entityBuffer = reflect.New(reflect.ArrayOf(cap, entityType)).Elem()

	a.archetypeAccess = archetypeAccess{
		basePointer:   unsafe.Pointer(&a.layouts[0]),
		entityPointer: a.entityBuffer.Addr().UnsafePointer(),
		Mask:          mask,
	}

	a.graphNode = node

	a.capacityIncrement = uint32(capacityIncrement)
	a.len = 0
	a.cap = uint32(cap)

	if maxSize > 0 {
		a.zeroValue = make([]byte, maxSize)
		a.zeroPointer = unsafe.Pointer(&a.zeroValue[0])
	}
}

// Add adds an entity with optionally zeroed components to the archetype
func (a *archetype) Alloc(entity Entity) uintptr {
	idx := uintptr(a.len)
	a.extend(1)
	a.addEntity(idx, &entity)
	a.len++
	return idx
}

// Add adds storage to the archetype
func (a *archetype) AllocN(count uint32) {
	a.extend(count)
	a.len += count
}

// Add adds an entity with components to the archetype
func (a *archetype) Add(entity Entity, components ...Component) uintptr {
	if len(components) != len(a.Ids) {
		panic("Invalid number of components")
	}
	idx := uintptr(a.len)

	a.extend(1)
	a.addEntity(idx, &entity)
	for _, c := range components {
		lay := a.getLayout(c.ID)
		dst := a.Get(uintptr(idx), c.ID)
		if lay.itemSize == 0 {
			continue
		}
		src := reflect.ValueOf(c.Comp).UnsafePointer()
		a.copy(src, dst, lay.itemSize)
	}
	a.len++
	return idx
}

// Remove removes an entity and its components from the archetype.
//
// Performs a swap-remove and reports whether a swap was necessary
// (i.e. not the last entity that was removed).
func (a *archetype) Remove(index uintptr) bool {
	swapped := a.removeEntity(index)

	old := uintptr(a.len - 1)

	if index != old {
		for _, id := range a.Ids {
			lay := a.getLayout(id)
			if lay.itemSize == 0 {
				continue
			}
			src := unsafe.Add(lay.pointer, old*lay.itemSize)
			dst := unsafe.Add(lay.pointer, index*lay.itemSize)
			a.copy(src, dst, lay.itemSize)
		}
	}
	a.ZeroAll(old)
	a.len--

	return swapped
}

// ZeroAll resets a block of storage in all buffers.
func (a *archetype) ZeroAll(index uintptr) {
	for _, id := range a.Ids {
		a.Zero(index, id)
	}
}

// ZeroAll resets a block of storage in one buffer.
func (a *archetype) Zero(index uintptr, id ID) {
	lay := a.getLayout(id)
	if lay.itemSize == 0 {
		return
	}
	dst := unsafe.Add(lay.pointer, index*lay.itemSize)
	a.copy(a.zeroPointer, dst, lay.itemSize)
}

// SetEntity overwrites an entity
func (a *archetype) SetEntity(index uintptr, entity Entity) {
	a.addEntity(index, &entity)
}

// Set overwrites a component with the data behind the given pointer
func (a *archetype) Set(index uintptr, id ID, comp interface{}) unsafe.Pointer {
	lay := a.getLayout(id)
	dst := a.Get(index, id)
	if lay.itemSize == 0 {
		return dst
	}
	rValue := reflect.ValueOf(comp)

	src := rValue.UnsafePointer()
	a.copy(src, dst, lay.itemSize)
	return dst
}

// SetPointer overwrites a component with the data behind the given pointer
func (a *archetype) SetPointer(index uintptr, id ID, comp unsafe.Pointer) unsafe.Pointer {
	lay := a.getLayout(id)
	dst := a.Get(index, id)
	if lay.itemSize == 0 {
		return dst
	}

	a.copy(comp, dst, lay.itemSize)
	return dst
}

// Reset removes all entities and components.
//
// Does NOT free the reserved memory.
func (a *archetype) Reset() {
	a.len = 0
	for _, buf := range a.buffers {
		buf.SetZero()
	}
}

// Components returns the component IDs for this archetype
func (a *archetype) Components() []ID {
	return a.Ids
}

// Len reports the number of entities in the archetype
func (a *archetype) Len() uint32 {
	return a.len
}

// Cap reports the current capacity of the archetype
func (a *archetype) Cap() uint32 {
	return a.cap
}

// Stats generates statistics for an archetype
func (a *archetype) Stats(reg *componentRegistry[ID]) stats.ArchetypeStats {
	ids := a.Components()
	aCompCount := len(ids)
	aTypes := make([]reflect.Type, aCompCount)
	for j, id := range ids {
		aTypes[j] = reg.ComponentType(id)
	}

	cap := int(a.Cap())
	memPerEntity := 0
	for _, id := range a.Ids {
		lay := a.getLayout(id)
		memPerEntity += int(lay.itemSize)
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

// copy from one pointer to another.
func (a *archetype) copy(src, dst unsafe.Pointer, itemSize uintptr) {
	dstSlice := (*[math.MaxInt32]byte)(dst)[:itemSize:itemSize]
	srcSlice := (*[math.MaxInt32]byte)(src)[:itemSize:itemSize]
	copy(dstSlice, srcSlice)
}

// extend the memory buffers if necessary for adding an entity.
func (a *archetype) extend(by uint32) {
	required := a.len + by
	if a.cap >= required {
		return
	}
	a.cap = capacityU32(required, a.capacityIncrement)

	old := a.entityBuffer
	a.entityBuffer = reflect.New(reflect.ArrayOf(int(a.cap), entityType)).Elem()
	a.entityPointer = a.entityBuffer.Addr().UnsafePointer()
	reflect.Copy(a.entityBuffer, old)

	for _, id := range a.Ids {
		lay := a.getLayout(id)
		if lay.itemSize == 0 {
			continue
		}
		index := a.indices[id]
		old := a.buffers[index]
		a.buffers[index] = reflect.New(reflect.ArrayOf(int(a.cap), old.Type().Elem())).Elem()
		lay.pointer = a.buffers[index].Addr().UnsafePointer()
		reflect.Copy(a.buffers[index], old)
	}
}

// Adds an entity at the given index. Does not extend the entity buffer.
func (a *archetype) addEntity(index uintptr, entity *Entity) {
	dst := unsafe.Add(a.entityPointer, entitySize*index)
	src := unsafe.Pointer(entity)
	a.copy(src, dst, entitySize)
}

// removeEntity removes an entity from tne archetype.
// Components need to be removed separately.
func (a *archetype) removeEntity(index uintptr) bool {
	old := uintptr(a.len - 1)

	if index == old {
		return false
	}

	src := unsafe.Add(a.entityPointer, old*entitySize)
	dst := unsafe.Add(a.entityPointer, index*entitySize)
	a.copy(src, dst, entitySize)

	return true
}
