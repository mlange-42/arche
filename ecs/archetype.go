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
	mask              Mask                  // Mask of the archetype
	Ids               []ID                  // List of component IDs
	archetype         *archetype            // The single archetype for nodes without entity
	archetypes        pagedSlice[archetype] // Storage for archetypes in nodes with entity relation
	archetypeMap      map[Entity]*archetype // Mapping from relation targets to archetypes
	freeIndices       []int32               // Indices of free/inactive archetypes
	TransitionAdd     idMap[*archetypeNode] // Mapping from component ID to add to the resulting archetype
	TransitionRemove  idMap[*archetypeNode] // Mapping from component ID to remove to the resulting archetype
	relation          int8                  // The node's relation component ID. Negative value stands for no relation
	zeroValue         []byte                // Used as source for setting storage to zero
	zeroPointer       unsafe.Pointer        // Points to zeroValue for fast access
	capacityIncrement uint32                // Capacity increment
}

// Creates a new archetypeNode
func newArchetypeNode(mask Mask, relation int8, capacityIncrement int) archetypeNode {
	var arch map[Entity]*archetype
	if relation >= 0 {
		arch = map[Entity]*archetype{}
	}
	return archetypeNode{
		mask:              mask,
		archetypeMap:      arch,
		TransitionAdd:     newIDMap[*archetypeNode](),
		TransitionRemove:  newIDMap[*archetypeNode](),
		relation:          relation,
		capacityIncrement: uint32(capacityIncrement),
	}
}

// Matches the archetype node against a filter.
// Ignores the relation target.
func (a *archetypeNode) Matches(f Filter) bool {
	return f.Matches(a.mask, nil)
}

// Archetypes of the node.
// Returns a single wrapped archetype if there are no relations.
// Returns nil if the node has no archetype(s).
func (a *archetypeNode) Archetypes() archetypes {
	if a.HasRelation() {
		return &a.archetypes
	}
	if a.archetype == nil {
		return nil
	}
	return batchArchetype{Archetype: a.archetype, StartIndex: 0}
}

// GetArchetype returns the archetype for the given relation target.
//
// The target is ignored if the node has no relation component.
func (a *archetypeNode) GetArchetype(target Entity) *archetype {
	if a.relation >= 0 {
		return a.archetypeMap[target]
	}
	return a.archetype
}

// SetArchetype sets the archetype for a node without a relation.
//
// Do not use on nodes without a relation component!
func (a *archetypeNode) SetArchetype(arch *archetype) {
	a.archetype = arch
}

// CreateArchetype creates a new archetype in nodes with relation component.
func (a *archetypeNode) CreateArchetype(target Entity, components ...componentType) *archetype {
	var arch *archetype
	var archIndex int32
	lenFree := len(a.freeIndices)
	if lenFree > 0 {
		archIndex = a.freeIndices[lenFree-1]
		arch = a.archetypes.Get(archIndex)
		a.freeIndices = a.freeIndices[:lenFree-1]
		arch.Activate(target, archIndex)
	} else {
		a.archetypes.Add(archetype{})
		archIndex := a.archetypes.Len() - 1
		arch = a.archetypes.Get(archIndex)
		arch.Init(a, archIndex, true, target, components...)
	}
	a.archetypeMap[target] = arch
	return arch
}

// RemoveArchetype de-activates an archetype.
// The archetype will be re-used by CreateArchetype.
func (a *archetypeNode) RemoveArchetype(arch *archetype) {
	delete(a.archetypeMap, arch.Relation)
	idx := arch.index
	a.freeIndices = append(a.freeIndices, idx)
	a.archetypes.Get(idx).Deactivate()
}

// HasRelation returns whether the node has a relation component.
func (a *archetypeNode) HasRelation() bool {
	return a.relation >= 0
}

// IsActive returns whether the node is active, i.e. has archetypes.
func (a *archetypeNode) IsActive() bool {
	return a.Ids != nil
}

// Stats generates statistics for an archetype node.
func (a *archetypeNode) Stats(reg *componentRegistry[ID]) stats.NodeStats {
	ids := a.Ids
	aCompCount := len(ids)
	aTypes := make([]reflect.Type, aCompCount)
	for j, id := range ids {
		aTypes[j], _ = reg.ComponentType(id)
	}

	arches := a.Archetypes()
	var numArches int32
	cap := 0
	count := 0
	memory := 0
	var archStats []stats.ArchetypeStats
	if arches != nil {
		numArches = arches.Len()
		archStats = make([]stats.ArchetypeStats, numArches)
		var i int32
		for i = 0; i < numArches; i++ {
			archStats[i] = arches.Get(i).Stats(reg)
			stats := &archStats[i]
			cap += stats.Capacity
			count += stats.Size
			memory += stats.Memory
		}
	}

	memPerEntity := 0
	for j := range ids {
		memPerEntity += int(aTypes[j].Size())
	}

	return stats.NodeStats{
		ArchetypeCount:       int(numArches),
		ActiveArchetypeCount: int(numArches) - len(a.freeIndices),
		IsActive:             a.IsActive(),
		HasRelation:          a.HasRelation(),
		Components:           aCompCount,
		ComponentIDs:         ids,
		ComponentTypes:       aTypes,
		Memory:               memory,
		MemoryPerEntity:      memPerEntity,
		Size:                 count,
		Capacity:             cap,
		Archetypes:           archStats,
	}
}

// UpdateStats updates statistics for an archetype node.
func (a *archetypeNode) UpdateStats(stats *stats.NodeStats, reg *componentRegistry[ID]) {
	if !a.IsActive() {
		return
	}

	arches := a.Archetypes()

	if !stats.IsActive {
		temp := a.Stats(reg)
		*stats = temp
		return
	}

	cap := 0
	count := 0
	memory := 0

	cntOld := int32(len(stats.Archetypes))
	cntNew := int32(arches.Len())
	var i int32
	for i = 0; i < cntOld; i++ {
		arch := &stats.Archetypes[i]
		arches.Get(i).UpdateStats(stats, arch, reg)
		cap += arch.Capacity
		count += arch.Size
		memory += arch.Memory
	}
	for i = cntOld; i < cntNew; i++ {
		arch := arches.Get(i).Stats(reg)
		stats.Archetypes = append(stats.Archetypes, arch)
		cap += arch.Capacity
		count += arch.Size
		memory += arch.Memory
	}

	stats.IsActive = true
	stats.ArchetypeCount = int(cntNew)
	stats.ActiveArchetypeCount = int(cntNew) - len(a.freeIndices)
	stats.Capacity = cap
	stats.Size = count
	stats.Memory = memory
}

// Helper for accessing data from an archetype
type archetypeAccess struct {
	Mask              Mask           // Archetype's mask
	basePointer       unsafe.Pointer // Pointer to the first component column layout.
	entityPointer     unsafe.Pointer // Pointer to the entity storage
	Relation          Entity
	RelationComponent int8
}

// Matches checks if the archetype matches the given mask.
func (a *archetype) Matches(f Filter) bool {
	return f.Matches(a.Mask, &a.Relation)
}

// GetEntity returns the entity at the given index
func (a *archetypeAccess) GetEntity(index uintptr) Entity {
	return *(*Entity)(unsafe.Add(a.entityPointer, entitySize*index))
}

// Get returns the component with the given ID at the given index
func (a *archetypeAccess) Get(index uintptr, id ID) unsafe.Pointer {
	return a.getLayout(id).Get(index)
}

// GetEntity returns the entity at the given index
func (a *archetypeAccess) GetRelation() Entity {
	return a.Relation
}

// HasComponent returns whether the archetype contains the given component ID.
func (a *archetypeAccess) HasComponent(id ID) bool {
	return a.getLayout(id).pointer != nil
}

// HasRelation returns whether the archetype has a relation component.
func (a *archetypeAccess) HasRelation() bool {
	return a.RelationComponent >= 0
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
	archetypeAccess                 // Access helper, passed to queries.
	graphNode       *archetypeNode  // Node in the archetype graph.
	layouts         []layout        // Column layouts by ID.
	indices         idMap[uint32]   // Mapping from IDs to buffer indices.
	buffers         []reflect.Value // Reflection arrays containing component data.
	entityBuffer    reflect.Value   // Reflection array containing entity data.
	index           int32           // Index of the archetype in the world.
	len             uint32          // Current number of entities
	cap             uint32          // Current capacity
}

// Init initializes an archetype
func (a *archetype) Init(node *archetypeNode, index int32, forStorage bool, relation Entity, components ...componentType) {
	var mask Mask
	if !node.IsActive() {
		node.Ids = make([]ID, len(components))

		var maxSize uintptr = 0
		for i, c := range components {
			node.Ids[i] = c.ID
			size, align := c.Type.Size(), uintptr(c.Type.Align())
			size = (size + (align - 1)) / align * align
			if size > maxSize {
				maxSize = size
			}
		}

		if maxSize > 0 {
			node.zeroValue = make([]byte, maxSize)
			node.zeroPointer = unsafe.Pointer(&node.zeroValue[0])
		}
	}

	a.buffers = make([]reflect.Value, len(components))
	a.layouts = make([]layout, MaskTotalBits)
	a.indices = newIDMap[uint32]()
	a.index = index

	cap := 1
	if forStorage {
		cap = int(node.capacityIncrement)
	}

	prev := -1
	for i, c := range components {
		if int(c.ID) <= prev {
			panic("component arguments must be sorted by ID")
		}
		prev = int(c.ID)
		mask.Set(c.ID, true)

		size, align := c.Type.Size(), uintptr(c.Type.Align())
		size = (size + (align - 1)) / align * align

		a.buffers[i] = reflect.New(reflect.ArrayOf(cap, c.Type)).Elem()
		a.layouts[c.ID] = layout{
			a.buffers[i].Addr().UnsafePointer(),
			size,
		}
		a.indices.Set(c.ID, uint32(i))
	}
	a.entityBuffer = reflect.New(reflect.ArrayOf(cap, entityType)).Elem()

	a.archetypeAccess = archetypeAccess{
		basePointer:       unsafe.Pointer(&a.layouts[0]),
		entityPointer:     a.entityBuffer.Addr().UnsafePointer(),
		Mask:              mask,
		Relation:          relation,
		RelationComponent: node.relation,
	}

	a.graphNode = node

	a.len = 0
	a.cap = uint32(cap)
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
	if len(components) != len(a.graphNode.Ids) {
		panic("Invalid number of components")
	}
	idx := uintptr(a.len)

	a.extend(1)
	a.addEntity(idx, &entity)
	for _, c := range components {
		lay := a.getLayout(c.ID)
		size := lay.itemSize
		if size == 0 {
			continue
		}
		src := reflect.ValueOf(c.Comp).UnsafePointer()
		dst := a.Get(uintptr(idx), c.ID)
		a.copy(src, dst, size)
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
		for _, id := range a.graphNode.Ids {
			lay := a.getLayout(id)
			size := lay.itemSize
			if size == 0 {
				continue
			}
			src := unsafe.Add(lay.pointer, old*size)
			dst := unsafe.Add(lay.pointer, index*size)
			a.copy(src, dst, size)
		}
	}
	a.ZeroAll(old)
	a.len--

	return swapped
}

// ZeroAll resets a block of storage in all buffers.
func (a *archetype) ZeroAll(index uintptr) {
	for _, id := range a.graphNode.Ids {
		a.Zero(index, id)
	}
}

// ZeroAll resets a block of storage in one buffer.
func (a *archetype) Zero(index uintptr, id ID) {
	lay := a.getLayout(id)
	size := lay.itemSize
	if size == 0 {
		return
	}
	dst := unsafe.Add(lay.pointer, index*size)
	a.copy(a.graphNode.zeroPointer, dst, size)
}

// SetEntity overwrites an entity
func (a *archetype) SetEntity(index uintptr, entity Entity) {
	a.addEntity(index, &entity)
}

// Set overwrites a component with the data behind the given pointer
func (a *archetype) Set(index uintptr, id ID, comp interface{}) unsafe.Pointer {
	lay := a.getLayout(id)
	dst := a.Get(index, id)
	size := lay.itemSize
	if size == 0 {
		return dst
	}
	rValue := reflect.ValueOf(comp)

	src := rValue.UnsafePointer()
	a.copy(src, dst, size)
	return dst
}

// SetPointer overwrites a component with the data behind the given pointer
func (a *archetype) SetPointer(index uintptr, id ID, comp unsafe.Pointer) unsafe.Pointer {
	lay := a.getLayout(id)
	dst := a.Get(index, id)
	size := lay.itemSize
	if size == 0 {
		return dst
	}

	a.copy(comp, dst, size)
	return dst
}

// Reset removes all entities and components.
//
// Does NOT free the reserved memory.
func (a *archetype) Reset() {
	if a.len == 0 {
		return
	}
	a.len = 0
	for _, buf := range a.buffers {
		buf.SetZero()
	}
}

// Deactivate the archetype for later re-use.
func (a *archetype) Deactivate() {
	a.Reset()
	a.index = -1
}

// Activate reactivates a de-activated archetype.
func (a *archetype) Activate(target Entity, index int32) {
	a.index = index
	a.Relation = target
}

// IsActive returns whether the archetype is active.
// Otherwise, it is eligible for re-use.
func (a *archetype) IsActive() bool {
	return a.index >= 0
}

// Components returns the component IDs for this archetype
func (a *archetype) Components() []ID {
	return a.graphNode.Ids
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
		aTypes[j], _ = reg.ComponentType(id)
	}

	cap := int(a.Cap())
	memPerEntity := 0
	for _, id := range a.graphNode.Ids {
		lay := a.getLayout(id)
		memPerEntity += int(lay.itemSize)
	}
	memory := cap * (int(entitySize) + memPerEntity)

	return stats.ArchetypeStats{
		IsActive: a.IsActive(),
		Size:     int(a.Len()),
		Capacity: cap,
		Memory:   memory,
	}
}

// UpdateStats updates statistics for an archetype
func (a *archetype) UpdateStats(node *stats.NodeStats, stats *stats.ArchetypeStats, reg *componentRegistry[ID]) {
	cap := int(a.Cap())
	memory := cap * (int(entitySize) + node.MemoryPerEntity)

	stats.IsActive = a.IsActive()
	stats.Size = int(a.Len())
	stats.Capacity = cap
	stats.Memory = memory
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
	a.cap = capacityU32(required, a.graphNode.capacityIncrement)

	old := a.entityBuffer
	a.entityBuffer = reflect.New(reflect.ArrayOf(int(a.cap), entityType)).Elem()
	a.entityPointer = a.entityBuffer.Addr().UnsafePointer()
	reflect.Copy(a.entityBuffer, old)

	for _, id := range a.graphNode.Ids {
		lay := a.getLayout(id)
		if lay.itemSize == 0 {
			continue
		}
		index, _ := a.indices.Get(id)
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
