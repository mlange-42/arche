package ecs

import (
	"fmt"
	"unsafe"
)

// Query is an iterator to iterate entities, filtered by a [Filter].
//
// Create queries through the [World] using [World.Query].
//
// See also the generic alternatives [github.com/mlange-42/arche/generic.Query1],
// [github.com/mlange-42/arche/generic.Query2], etc.
// For advanced filtering, see package [github.com/mlange-42/arche/filter].
type Query struct {
	filter         Filter           // The filter used by the query.
	archetypes     []*archetype     // The query's filtered archetypes.
	nodeArchetypes archetypes       // The query's archetypes of the current node.
	nodes          []*archNode      // The query's nodes.
	world          *World           // The [World].
	access         *archetypeAccess // Access helper for the archetype currently being iterated.
	archetype      *archetype       // The archetype currently being iterated.
	entityIndex    uint32           // Iteration index of the current [Entity] current archetype.
	entityIndexMax uint32           // Maximum entity index in the current archetype.
	archIndex      int32            // Iteration index of the current archetype.
	nodeIndex      int32            // Iteration index of the current archetype.
	count          int32            // Cached entity count.
	lockBit        uint8            // The bit that was used to lock the [World] when the query was created.
	isFiltered     bool             // Whether the list of archetype nodes is already filtered.
	isBatch        bool             // Marks the query as a query over a batch iteration.
}

// newQuery creates a new Filter
func newQuery(world *World, filter Filter, lockBit uint8, nodes []*archNode) Query {
	return Query{
		filter:     filter,
		world:      world,
		nodes:      nodes,
		archIndex:  -1,
		nodeIndex:  -1,
		lockBit:    lockBit,
		count:      -1,
		isFiltered: false,
		isBatch:    false,
	}
}

// newQuery creates a new Filter
func newCachedQuery(world *World, filter Filter, lockBit uint8, archetypes []*archetype) Query {
	return Query{
		filter:     filter,
		world:      world,
		archetypes: archetypes,
		archIndex:  -1,
		nodeIndex:  -1,
		lockBit:    lockBit,
		count:      -1,
		isFiltered: true,
		isBatch:    false,
	}
}

// newQuery creates a query on a single archetype
func newBatchQuery(world *World, lockBit uint8, archetype *batchArchetypes) Query {
	return Query{
		filter:         nil,
		isFiltered:     false,
		isBatch:        true,
		world:          world,
		nodeArchetypes: archetype,
		archIndex:      -1,
		lockBit:        lockBit,
		count:          -1,
	}
}

// Next proceeds to the next [Entity] in the Query.
//
// Returns false if no next entity could be found.
func (q *Query) Next() bool {
	if q.entityIndex < q.entityIndexMax {
		q.entityIndex++
		return true
	}
	// outline to allow inlining of the fast path
	return q.nextArchetype()
}

// Has returns whether the current entity has the given component.
func (q *Query) Has(comp ID) bool {
	return q.access.HasComponent(comp)
}

// Get returns the pointer to the given component at the iterator's position.
func (q *Query) Get(comp ID) unsafe.Pointer {
	return q.access.Get(q.entityIndex, comp)
}

// Entity returns the entity at the iterator's position.
func (q *Query) Entity() Entity {
	return q.access.GetEntity(q.entityIndex)
}

// Relation returns the target entity for an entity relation.
//
// Panics if the entity does not have the given component, or if the component is not a [Relation].
func (q *Query) Relation(comp ID) Entity {
	if q.access.RelationComponent.id != comp.id {
		panic(fmt.Sprintf("entity has no component %v, or it is not a relation component", q.world.registry.Types[comp.id]))
	}
	return q.access.RelationTarget
}

// RelationUnchecked returns the target entity for an entity relation.
//
// Returns the zero entity if the entity does not have the given component,
// or if the component is not a [Relation].
//
// GetRelationUnchecked is an optimized version of [Query.Relation].
// Does not check that the component ID is applicable.
func (q *Query) relationUnchecked(comp ID) Entity {
	return q.access.RelationTarget
}

// Step advances the query iterator by the given number of entities.
//
// Query.Step(1) is equivalent to [Query.Next]().
//
// This method, used together with [Query.Count], can be useful for the selection of random entities.
func (q *Query) Step(step int) bool {
	if step <= 0 {
		panic("step size must be positive")
	}
	var ok bool
	for {
		step, ok = q.stepArchetype(uint32(step))
		if ok {
			return true
		}
		if !q.nextArchetype() {
			return false
		}
		if step == 0 {
			return true
		}
	}
}

// Count counts the entities matching this query.
//
// Involves a small overhead of iterating through archetypes when called the first time.
// However, this is still much faster than manual counting via iteration.
//
// Does not close the query.
func (q *Query) Count() int {
	if q.count >= 0 {
		return int(q.count)
	}
	q.count = int32(q.countEntities())
	return int(q.count)
}

// Mask returns the archetype [Mask] for the [Entity] at the iterator's current position.
func (q *Query) Mask() Mask {
	return q.access.Mask
}

// Ids returns the component IDs for the archetype of the [Entity] at the iterator's current position.
//
// Returns a copy of the archetype's component IDs slice, for safety.
// This means that the result can be manipulated safely,
// but also that calling the method may incur some significant cost.
func (q *Query) Ids() []ID {
	return append([]ID{}, q.archetype.node.Ids...)
}

// Close closes the Query and unlocks the world.
//
// Automatically called when iteration finishes.
// Needs to be called only if breaking out of the query iteration or not iterating at all.
func (q *Query) Close() {
	q.world.closeQuery(q)
}

// nextArchetype proceeds to the next archetype, and returns whether this was successful/possible.
func (q *Query) nextArchetype() bool {
	if q.isFiltered {
		return q.nextArchetypeFiltered()
	}
	if q.isBatch {
		return q.nextBatch()
	}
	return q.nextNodeOrArchetype()
}

func (q *Query) nextBatch() bool {
	if q.nextArchetypeBatch() {
		return true
	}
	q.world.closeQuery(q)
	return false
}

func (q *Query) nextArchetypeBatch() bool {
	len := int32(q.nodeArchetypes.Len()) - 1
	for q.archIndex < len {
		q.archIndex++
		a := q.nodeArchetypes.Get(q.archIndex)
		aLen := a.Len()
		if aLen > 0 {
			q.access = &a.archetypeAccess
			q.archetype = a
			batch := q.nodeArchetypes.(*batchArchetypes)
			q.entityIndex = batch.StartIndex[q.archIndex]
			q.entityIndexMax = batch.EndIndex[q.archIndex] - 1
			return true
		}
	}
	return false
}

func (q *Query) nextArchetypeSimple() bool {
	len := int32(q.nodeArchetypes.Len()) - 1
	for q.archIndex < len {
		q.archIndex++
		a := q.nodeArchetypes.Get(q.archIndex)
		aLen := a.Len()
		if aLen == 0 {
			continue
		}
		q.access = &a.archetypeAccess
		q.archetype = a
		q.entityIndex = 0
		q.entityIndexMax = aLen - 1
		return true
	}
	return false
}

func (q *Query) nextArchetypeFiltered() bool {
	len := int32(len(q.archetypes) - 1)
	for q.archIndex < len {
		q.archIndex++
		a := q.archetypes[q.archIndex]
		aLen := a.Len()
		if aLen == 0 {
			continue
		}
		q.access = &a.archetypeAccess
		q.archetype = a
		q.entityIndex = 0
		q.entityIndexMax = aLen - 1
		return true
	}
	q.world.closeQuery(q)
	return false
}

func (q *Query) nextNodeOrArchetype() bool {
	if q.nodeArchetypes != nil && q.nextArchetypeSimple() {
		return true
	}
	return q.nextNode()
}

func (q *Query) nextNode() bool {
	len := int32(len(q.nodes)) - 1
	for q.nodeIndex < len {
		q.nodeIndex++
		n := q.nodes[q.nodeIndex]

		if !n.IsActive {
			continue
		}
		if !n.Matches(q.filter) {
			continue
		}

		arches := n.Archetypes()

		if !n.HasRelation {
			// There should be at least one archetype.
			// Otherwise, the node would be inactive.
			arch := arches.Get(0)
			archLen := arch.Len()
			if archLen > 0 {
				q.setArchetype(nil, &arch.archetypeAccess, arch, arch.index, archLen-1)
				return true
			}
			continue
		}

		if rf, ok := q.filter.(*RelationFilter); ok {
			target := rf.Target
			if arch, ok := n.archetypeMap[target]; ok && arch.Len() > 0 {
				q.setArchetype(nil, &arch.archetypeAccess, arch, arch.index, arch.Len()-1)
				return true
			}
			continue
		}

		q.setArchetype(arches, nil, nil, -1, 0)
		if q.nextArchetypeSimple() {
			return true
		}
	}
	q.nodeArchetypes = nil
	q.world.closeQuery(q)
	return false
}

func (q *Query) setArchetype(arches archetypes, access *archetypeAccess, arch *archetype, archIndex int32, maxIndex uint32) {
	q.nodeArchetypes = arches
	q.archIndex = archIndex
	q.access = access
	q.archetype = arch
	q.entityIndex = 0
	q.entityIndexMax = maxIndex
}

func (q *Query) stepArchetype(step uint32) (int, bool) {
	q.entityIndex += step
	if q.entityIndex <= q.entityIndexMax {
		return 0, true
	}
	return int(q.entityIndex) - int(q.entityIndexMax) - 1, false
}

func (q *Query) countEntities() int {
	var count uint32 = 0

	if q.isBatch {
		batch := q.nodeArchetypes.(*batchArchetypes)
		nArch := batch.Len()
		var j int32
		for j = 0; j < nArch; j++ {
			count += batch.EndIndex[j] - batch.StartIndex[j]
		}
		return int(count)
	}

	if q.isFiltered {
		ln := int32(len(q.archetypes))
		var i int32
		for i = 0; i < ln; i++ {
			count += q.archetypes[i].Len()
		}
		return int(count)
	}

	for _, nd := range q.nodes {
		if !nd.IsActive || !nd.Matches(q.filter) {
			continue
		}

		if !nd.HasRelation {
			// There should be at least one archetype.
			// Otherwise, the node would be inactive.
			arches := nd.Archetypes()
			arch := arches.Get(0)
			count += arch.Len()
			continue
		}

		if rf, ok := q.filter.(*RelationFilter); ok {
			target := rf.Target
			if arch, ok := nd.archetypeMap[target]; ok {
				count += arch.Len()
			}
			continue
		}

		arches := nd.Archetypes()
		nArch := arches.Len()
		var j int32
		for j = 0; j < nArch; j++ {
			a := arches.Get(j)
			count += a.Len()
		}
	}
	return int(count)
}
