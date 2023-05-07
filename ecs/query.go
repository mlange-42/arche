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
	archetypes     archetypes       // The query's archetypes (can be all, unfiltered archetypes).
	nodes          []*archNode      // The query's nodes
	world          *World           // The [World].
	access         *archetypeAccess // Access helper for the archetype currently being iterated.
	entityIndex    uintptr          // Iteration index of the current [Entity] current archetype.
	entityIndexMax uintptr          // Maximum entity index in the current archetype.
	archIndex      int32            // Iteration index of the current archetype.
	nodeIndex      int32            // Iteration index of the current archetype.
	count          int32            // Cached entity count.
	lockBit        uint8            // The bit that was used to lock the [World] when the query was created.
	isFiltered     bool             // Whether the list of archetypes is already filtered.
	isBatch        bool
}

// newQuery creates a new Filter
func newQuery(world *World, filter Filter, lockBit uint8, nodes []*archNode, isFiltered bool) Query {
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

// newQuery creates a query on a single archetype
func newArchQuery(world *World, lockBit uint8, archetype *batchArchetype) Query {
	arch := archetype.Archetype
	if archetype.StartIndex > 0 {
		return Query{
			filter:         nil,
			isFiltered:     true,
			isBatch:        true,
			world:          world,
			archetypes:     archetype,
			access:         &arch.archetypeAccess,
			archIndex:      0,
			lockBit:        lockBit,
			count:          int32(archetype.EndIndex - archetype.StartIndex),
			entityIndex:    uintptr(archetype.StartIndex - 1),
			entityIndexMax: uintptr(archetype.EndIndex - 1),
		}
	}
	return Query{
		filter:     nil,
		isFiltered: true,
		isBatch:    true,
		world:      world,
		archetypes: archetype,
		archIndex:  -1,
		lockBit:    lockBit,
		count:      int32(archetype.EndIndex),
	}
}

// Next proceeds to the next [Entity] in the Query.
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
	if q.access.RelationComponent != int8(comp) {
		panic(fmt.Sprintf("entity has no component %v, or it is not a relation component", q.world.registry.Types[comp]))
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
func (q *Query) Count() int {
	if q.count >= 0 {
		return int(q.count)
	}
	q.count = int32(q.countEntities())
	return int(q.count)
}

// Mask returns the archetype [Mask] for the [Entity] at the iterator's current position.
//
// Can be used for fast checks of the entity composition, e.g. using a [Filter].
func (q *Query) Mask() Mask {
	return q.access.Mask
}

// Close closes the Query and unlocks the world.
//
// Automatically called when iteration finishes.
// Needs to be called only if breaking out of the query iteration.
func (q *Query) Close() {
	q.world.closeQuery(q)
}

// nextArchetype proceeds to the next archetype, and returns whether this was successful/possible.
func (q *Query) nextArchetype() bool {
	if !q.isBatch {
		return q.nextNode()
	}
	if q.nextArchetypeBatch() {
		return true
	}
	q.world.closeQuery(q)
	return false
}

func (q *Query) nextArchetypeBatch() bool {
	len := int32(q.archetypes.Len()) - 1
	for q.archIndex < len {
		q.archIndex++
		a := q.archetypes.Get(q.archIndex)
		aLen := a.Len()
		if aLen > 0 {
			q.access = &a.archetypeAccess
			q.entityIndex = 0
			if batch, ok := q.archetypes.(*batchArchetype); ok {
				q.entityIndexMax = uintptr(batch.EndIndex) - 1
			}
			return true
		}
	}
	return false
}

func (q *Query) nextArchetypeSimple() bool {
	len := int32(q.archetypes.Len()) - 1
	for q.archIndex < len {
		q.archIndex++
		a := q.archetypes.Get(q.archIndex)
		aLen := a.Len()
		if aLen > 0 {
			q.access = &a.archetypeAccess
			q.entityIndex = 0
			q.entityIndexMax = uintptr(aLen) - 1
			return true
		}
	}
	return false
}

func (q *Query) nextNode() bool {
	if q.archetypes != nil && q.nextArchetypeSimple() {
		return true
	}

	len := int32(len(q.nodes)) - 1
	for q.nodeIndex < len {
		q.nodeIndex++
		n := q.nodes[q.nodeIndex]

		if !n.IsActive {
			continue
		}
		if !q.isFiltered && !n.Matches(q.filter) {
			continue
		}

		arches := n.Archetypes()

		if !n.HasRelation {
			// There should be at least one archetype.
			// Otherwise, the node would be inactive.
			arch := arches.Get(0)
			if arch.Len() > 0 {
				q.archetypes = nil
				q.archIndex = arch.index
				q.access = &arch.archetypeAccess
				q.entityIndex = 0
				q.entityIndexMax = uintptr(arch.Len()) - 1
				return true
			}
			continue
		}

		if rf, ok := q.filter.(*relationFilter); ok {
			target := rf.Target
			if arch, ok := n.archetypeMap[target]; ok && arch.Len() > 0 {
				q.archetypes = nil
				q.archIndex = arch.index
				q.access = &arch.archetypeAccess
				q.entityIndex = 0
				q.entityIndexMax = uintptr(arch.Len()) - 1
				return true
			}
			continue
		}

		q.archetypes = arches
		q.archIndex = -1
		q.entityIndex = 0
		q.entityIndexMax = 0
		if q.nextArchetypeSimple() {
			return true
		}
	}
	q.archetypes = nil
	q.world.closeQuery(q)
	return false
}

func (q *Query) stepArchetype(step uint32) (int, bool) {
	q.entityIndex += uintptr(step)
	if q.entityIndex <= q.entityIndexMax {
		return 0, true
	}
	return int(q.entityIndex) - int(q.entityIndexMax) - 1, false
}

func (q *Query) countEntities() int {
	// This is not necessary as batch queries get their count upon construction.
	// if q.isBatch {}

	len := int32(len(q.nodes))
	var count uint32 = 0
	var i int32
	for i = 0; i < len; i++ {
		nd := q.nodes[i]
		if !nd.IsActive || (!q.isFiltered && !nd.Matches(q.filter)) {
			continue
		}

		if rf, ok := q.filter.(*relationFilter); ok {
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
			if a.IsActive() && a.Matches(q.filter) {
				count += a.Len()
			}
		}
	}
	return int(count)
}
