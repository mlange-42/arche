package ecs

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/mlange-42/arche/ecs/event"
)

// fromConfig creates a new [World] from a [Config].
func fromConfig(conf Config) World {
	if conf.CapacityIncrement < 1 {
		panic("invalid CapacityIncrement in config, must be > 0")
	}
	if conf.RelationCapacityIncrement < 1 {
		conf.RelationCapacityIncrement = conf.CapacityIncrement
	}
	entities := make([]entityIndex, 1, conf.CapacityIncrement)
	entities[0] = entityIndex{arch: nil, index: 0}
	targetEntities := bitSet{}
	targetEntities.ExtendTo(1)

	w := World{
		config:         conf,
		entities:       entities,
		targetEntities: targetEntities,
		entityPool:     newEntityPool(uint32(conf.CapacityIncrement)),
		registry:       newComponentRegistry(),
		archetypes:     pagedSlice[archetype]{},
		archetypeData:  pagedSlice[archetypeData]{},
		nodes:          pagedSlice[archNode]{},
		relationNodes:  []*archNode{},
		locks:          lockMask{},
		listener:       nil,
		resources:      newResources(),
		filterCache:    newCache(),
	}
	node := w.createArchetypeNode(Mask{}, ID{}, false)
	w.createArchetype(node, Entity{}, false)
	return w
}

// Creates a new entity with a relation and a target entity.
func (w *World) newEntityTarget(targetID ID, target Entity, comps ...ID) Entity {
	w.checkLocked()

	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	arch := w.archetypes.Get(0)

	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil, target)
	}
	w.checkRelation(arch, targetID)

	entity := w.createEntity(arch)

	if !target.IsZero() {
		w.targetEntities.Set(target.id, true)
	}

	if w.listener != nil {
		bits := subscription(true, false, len(comps) > 0, false, true, true)
		trigger := w.listener.Subscriptions() & bits
		if trigger != 0 && subscribes(trigger, &arch.Mask, nil, w.listener.Components(), nil, &targetID) {
			w.listener.Notify(w, EntityEvent{Entity: entity, Added: arch.Mask, AddedIDs: comps, NewRelation: &targetID, EventTypes: bits})
		}
	}
	return entity
}

// Creates a new entity with a relation and a target entity.
func (w *World) newEntityTargetWith(targetID ID, target Entity, comps ...Component) Entity {
	w.checkLocked()

	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}

	arch := w.archetypes.Get(0)
	arch = w.findOrCreateArchetype(arch, ids, nil, target)
	w.checkRelation(arch, targetID)

	entity := w.createEntity(arch)

	if !target.IsZero() {
		w.targetEntities.Set(target.id, true)
	}

	for _, c := range comps {
		w.copyTo(entity, c.ID, c.Comp)
	}

	if w.listener != nil {
		bits := subscription(true, false, len(comps) > 0, false, true, true)
		trigger := w.listener.Subscriptions() & bits
		if trigger != 0 && subscribes(trigger, &arch.Mask, nil, w.listener.Components(), nil, &targetID) {
			w.listener.Notify(w, EntityEvent{Entity: entity, Added: arch.Mask, AddedIDs: ids, NewRelation: &targetID, EventTypes: bits})
		}
	}
	return entity
}

// Creates new entities without returning a query over them.
// Used via [World.Batch].
func (w *World) newEntities(count int, targetID ID, hasTarget bool, target Entity, comps ...ID) (*archetype, uint32) {
	arch, startIdx := w.newEntitiesNoNotify(count, targetID, hasTarget, target, comps...)

	if w.listener != nil {
		var newRel *ID
		if arch.HasRelationComponent {
			newRel = &arch.RelationComponent
		}
		bits := subscription(true, false, len(comps) > 0, false, newRel != nil, newRel != nil)
		trigger := w.listener.Subscriptions() & bits
		if trigger != 0 && subscribes(trigger, &arch.Mask, nil, w.listener.Components(), nil, newRel) {
			cnt := uint32(count)
			var i uint32
			for i = 0; i < cnt; i++ {
				idx := startIdx + i
				entity := arch.GetEntity(idx)
				w.listener.Notify(w, EntityEvent{Entity: entity, Added: arch.Mask, AddedIDs: comps, NewRelation: newRel, EventTypes: bits})
			}
		}
	}

	return arch, startIdx
}

// Creates new entities and returns a query over them.
// Used via [World.Batch].
func (w *World) newEntitiesQuery(count int, targetID ID, hasTarget bool, target Entity, comps ...ID) Query {
	arch, startIdx := w.newEntitiesNoNotify(count, targetID, hasTarget, target, comps...)
	lock := w.lock()

	batches := batchArchetypes{
		Added:   arch.Components(),
		Removed: nil,
	}
	batches.Add(arch, nil, startIdx, arch.Len())
	return newBatchQuery(w, lock, &batches)
}

// Creates new entities with component values without returning a query over them.
// Used via [World.Batch].
func (w *World) newEntitiesWith(count int, targetID ID, hasTarget bool, target Entity, comps ...Component) (*archetype, uint32) {
	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}

	arch, startIdx := w.newEntitiesWithNoNotify(count, targetID, hasTarget, target, ids, comps...)

	if w.listener != nil {
		var newRel *ID
		if arch.HasRelationComponent {
			newRel = &arch.RelationComponent
		}
		bits := subscription(true, false, len(comps) > 0, false, newRel != nil, newRel != nil)
		trigger := w.listener.Subscriptions() & bits
		if trigger != 0 && subscribes(trigger, &arch.Mask, nil, w.listener.Components(), nil, newRel) {
			var i uint32
			cnt := uint32(count)
			for i = 0; i < cnt; i++ {
				idx := startIdx + i
				entity := arch.GetEntity(idx)
				w.listener.Notify(w, EntityEvent{Entity: entity, Added: arch.Mask, AddedIDs: ids, NewRelation: newRel, EventTypes: bits})
			}
		}
	}

	return arch, startIdx
}

// Creates new entities with component values and returns a query over them.
// Used via [World.Batch].
func (w *World) newEntitiesWithQuery(count int, targetID ID, hasTarget bool, target Entity, comps ...Component) Query {
	ids := make([]ID, len(comps))
	for i, c := range comps {
		ids[i] = c.ID
	}

	arch, startIdx := w.newEntitiesWithNoNotify(count, targetID, hasTarget, target, ids, comps...)
	lock := w.lock()
	batches := batchArchetypes{
		Added:   arch.Components(),
		Removed: nil,
	}
	batches.Add(arch, nil, startIdx, arch.Len())
	return newBatchQuery(w, lock, &batches)
}

// Internal method to create new entities.
func (w *World) newEntitiesNoNotify(count int, targetID ID, hasTarget bool, target Entity, comps ...ID) (*archetype, uint32) {
	w.checkLocked()

	if count < 1 {
		panic("can only create a positive number of entities")
	}

	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, comps, nil, target)
	}
	if hasTarget {
		w.checkRelation(arch, targetID)
		if !target.IsZero() {
			w.targetEntities.Set(target.id, true)
		}
	}

	startIdx := arch.Len()
	w.createEntities(arch, uint32(count))

	return arch, startIdx
}

// Internal method to create new entities with component values.
func (w *World) newEntitiesWithNoNotify(count int, targetID ID, hasTarget bool, target Entity, ids []ID, comps ...Component) (*archetype, uint32) {
	w.checkLocked()

	if count < 1 {
		panic("can only create a positive number of entities")
	}

	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	if len(comps) == 0 {
		return w.newEntitiesNoNotify(count, targetID, hasTarget, target)
	}

	cnt := uint32(count)

	arch := w.archetypes.Get(0)
	if len(comps) > 0 {
		arch = w.findOrCreateArchetype(arch, ids, nil, target)
	}
	if hasTarget {
		w.checkRelation(arch, targetID)
		if !target.IsZero() {
			w.targetEntities.Set(target.id, true)
		}
	}

	startIdx := arch.Len()
	w.createEntities(arch, uint32(count))

	var i uint32
	for i = 0; i < cnt; i++ {
		idx := startIdx + i
		entity := arch.GetEntity(idx)
		for _, c := range comps {
			w.copyTo(entity, c.ID, c.Comp)
		}
	}

	return arch, startIdx
}

// createEntity creates an Entity and adds it to the given archetype.
func (w *World) createEntity(arch *archetype) Entity {
	entity := w.entityPool.Get()
	idx := arch.Alloc(entity)
	len := len(w.entities)
	if int(entity.id) == len {
		if len == cap(w.entities) {
			old := w.entities
			w.entities = make([]entityIndex, len, len+w.config.CapacityIncrement)
			copy(w.entities, old)

		}
		w.entities = append(w.entities, entityIndex{arch: arch, index: idx})
		w.targetEntities.ExtendTo(len + w.config.CapacityIncrement)
	} else {
		w.entities[entity.id] = entityIndex{arch: arch, index: idx}
		w.targetEntities.Set(entity.id, false)
	}
	return entity
}

// createEntity creates multiple Entities and adds them to the given archetype.
func (w *World) createEntities(arch *archetype, count uint32) {
	startIdx := arch.Len()
	arch.AllocN(count)

	len := len(w.entities)
	required := len + int(count) - w.entityPool.Available()
	capacity := capacity(required, w.config.CapacityIncrement)
	if required > cap(w.entities) {
		old := w.entities
		w.entities = make([]entityIndex, required, capacity)
		copy(w.entities, old)
	} else if required > len {
		w.entities = w.entities[:required]
	}
	w.targetEntities.ExtendTo(capacity)

	var i uint32
	for i = 0; i < count; i++ {
		idx := startIdx + i
		entity := w.entityPool.Get()
		arch.SetEntity(idx, entity)
		w.entities[entity.id] = entityIndex{arch: arch, index: idx}
		w.targetEntities.Set(entity.id, false)
	}
}

// RemoveEntities removes and recycles all entities matching a filter.
//
// Returns the number of removed entities.
//
// Panics when called on a locked world.
// Do not use during [Query] iteration!
func (w *World) removeEntities(filter Filter) int {
	w.checkLocked()

	lock := w.lock()

	var bits event.Subscription
	var listen bool

	var count uint32

	arches := w.getArchetypes(filter)
	numArches := int32(len(arches))
	var i int32
	for i = 0; i < numArches; i++ {
		arch := arches[i]
		ln := arch.Len()
		if ln == 0 {
			continue
		}

		count += ln

		var oldRel *ID
		var oldIds []ID
		if w.listener != nil {
			if arch.HasRelationComponent {
				oldRel = &arch.RelationComponent
			}
			if len(arch.node.Ids) > 0 {
				oldIds = arch.node.Ids
			}
			bits = subscription(false, true, false, len(oldIds) > 0, oldRel != nil, oldRel != nil)
			trigger := w.listener.Subscriptions() & bits
			listen = trigger != 0 && subscribes(trigger, nil, &arch.Mask, w.listener.Components(), oldRel, nil)
		}

		var j uint32
		for j = 0; j < ln; j++ {
			entity := arch.GetEntity(j)
			if listen {
				w.listener.Notify(w, EntityEvent{Entity: entity, Removed: arch.Mask, RemovedIDs: oldIds, OldRelation: oldRel, OldTarget: arch.RelationTarget, EventTypes: bits})
			}
			index := &w.entities[entity.id]
			index.arch = nil

			if w.targetEntities.Get(entity.id) {
				w.cleanupArchetypes(entity)
				w.targetEntities.Set(entity.id, false)
			}

			w.entityPool.Recycle(entity)
		}
		arch.Reset()
		w.cleanupArchetype(arch)
	}
	w.unlock(lock)

	return int(count)
}

// assign with relation target.
func (w *World) assign(entity Entity, relation ID, hasRelation bool, target Entity, comps ...Component) {
	len := len(comps)
	if len == 0 {
		panic("no components given to assign")
	}
	ids := make([]ID, len)
	for i, c := range comps {
		ids[i] = c.ID
	}
	arch, oldMask, oldTarget, oldRel := w.exchangeNoNotify(entity, ids, nil, relation, hasRelation, target)
	for _, c := range comps {
		w.copyTo(entity, c.ID, c.Comp)
	}
	if w.listener != nil {
		w.notifyExchange(arch, oldMask, entity, ids, nil, oldTarget, oldRel)
	}
}

// exchange with relation target.
func (w *World) exchange(entity Entity, add []ID, rem []ID, relation ID, hasRelation bool, target Entity) {
	if w.listener != nil {
		arch, oldMask, oldTarget, oldRel := w.exchangeNoNotify(entity, add, rem, relation, hasRelation, target)
		w.notifyExchange(arch, oldMask, entity, add, rem, oldTarget, oldRel)
		return
	}
	w.exchangeNoNotify(entity, add, rem, relation, hasRelation, target)
}

// perform exchange operation without notifying listeners.
func (w *World) exchangeNoNotify(entity Entity, add []ID, rem []ID, relation ID, hasRelation bool, target Entity) (*archetype, *Mask, Entity, *ID) {
	w.checkLocked()

	if !w.entityPool.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}

	if len(add) == 0 && len(rem) == 0 {
		if hasRelation {
			panic("exchange operation has no effect, but a relation is specified. Use World.Relation instead")
		}
		return nil, nil, Entity{}, nil
	}
	index := &w.entities[entity.id]
	oldArch := index.arch

	oldMask := oldArch.Mask
	mask := w.getExchangeMask(oldMask, add, rem)

	if hasRelation {
		if !mask.Get(relation) {
			tp, _ := w.registry.ComponentType(relation.id)
			panic(fmt.Sprintf("can't add relation: resulting entity has no component %s", tp.Name()))
		}
		if !w.registry.IsRelation.Get(relation) {
			tp, _ := w.registry.ComponentType(relation.id)
			panic(fmt.Sprintf("can't add relation: %s is not a relation component", tp.Name()))
		}
	} else {
		target = oldArch.RelationTarget
		if !oldArch.RelationTarget.IsZero() && oldArch.Mask.ContainsAny(&w.registry.IsRelation) {
			for _, id := range rem {
				// Removing a relation
				if w.registry.IsRelation.Get(id) {
					target = Entity{}
					break
				}
			}
		}
	}

	oldIDs := oldArch.Components()

	arch := w.findOrCreateArchetype(oldArch, add, rem, target)
	newIndex := arch.Alloc(entity)

	for _, id := range oldIDs {
		if mask.Get(id) {
			comp := oldArch.Get(index.index, id)
			arch.SetPointer(newIndex, id, comp)
		}
	}

	swapped := oldArch.Remove(index.index)

	if swapped {
		swapEntity := oldArch.GetEntity(index.index)
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{arch: arch, index: newIndex}

	var oldRel *ID
	if oldArch.HasRelationComponent {
		oldRel = &oldArch.RelationComponent
	}
	oldTarget := oldArch.RelationTarget

	if !target.IsZero() {
		w.targetEntities.Set(target.id, true)
	}

	w.cleanupArchetype(oldArch)

	return arch, &oldMask, oldTarget, oldRel
}

// notify listeners for an exchange.
func (w *World) notifyExchange(arch *archetype, oldMask *Mask, entity Entity, add []ID, rem []ID, oldTarget Entity, oldRel *ID) {
	var newRel *ID
	if arch.HasRelationComponent {
		newRel = &arch.RelationComponent
	}
	relChanged := false
	if oldRel != nil || newRel != nil {
		relChanged = (oldRel == nil) != (newRel == nil) || *oldRel != *newRel
	}
	targChanged := oldTarget != arch.RelationTarget

	bits := subscription(false, false, len(add) > 0, len(rem) > 0, relChanged, relChanged || targChanged)
	trigger := w.listener.Subscriptions() & bits
	if trigger != 0 {
		changed := oldMask.Xor(&arch.Mask)
		added := arch.Mask.And(&changed)
		removed := oldMask.And(&changed)
		if subscribes(trigger, &added, &removed, w.listener.Components(), oldRel, newRel) {
			w.listener.Notify(w,
				EntityEvent{Entity: entity, Added: added, Removed: removed,
					AddedIDs: add, RemovedIDs: rem, OldRelation: oldRel, NewRelation: newRel,
					OldTarget: oldTarget, EventTypes: bits},
			)
		}
	}
}

// Modify a mask by adding and removing IDs.
func (w *World) getExchangeMask(mask Mask, add []ID, rem []ID) Mask {
	for _, comp := range rem {
		if !mask.Get(comp) {
			panic(fmt.Sprintf("entity does not have a component of type %v, can't remove", w.registry.Types[comp.id]))
		}
		mask.Set(comp, false)
	}
	for _, comp := range add {
		if mask.Get(comp) {
			panic(fmt.Sprintf("entity already has component of type %v, can't add", w.registry.Types[comp.id]))
		}
		mask.Set(comp, true)
	}
	return mask
}

// ExchangeBatch exchanges components for many entities, matching a filter.
//
// If the callback argument is given, it is called with a [Query] over the affected entities,
// one Query for each affected archetype.
//
// Panics:
//   - when called with components that can't be added or removed because they are already present/not present, respectively.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See also [World.Exchange].
func (w *World) exchangeBatch(filter Filter, add []ID, rem []ID, relation ID, hasRelation bool, target Entity) int {
	batches := batchArchetypes{
		Added:   add,
		Removed: rem,
	}

	count := w.exchangeBatchNoNotify(filter, add, rem, relation, hasRelation, target, &batches)

	if w.listener != nil {
		w.notifyQuery(&batches)
	}
	return count
}

func (w *World) exchangeBatchQuery(filter Filter, add []ID, rem []ID, relation ID, hasRelation bool, target Entity) Query {
	batches := batchArchetypes{
		Added:   add,
		Removed: rem,
	}

	w.exchangeBatchNoNotify(filter, add, rem, relation, hasRelation, target, &batches)

	lock := w.lock()
	return newBatchQuery(w, lock, &batches)
}

func (w *World) exchangeBatchNoNotify(filter Filter, add []ID, rem []ID, relation ID, hasRelation bool, target Entity, batches *batchArchetypes) int {
	w.checkLocked()

	if len(add) == 0 && len(rem) == 0 {
		if hasRelation {
			panic("exchange operation has no effect, but a relation is specified. Use Batch.SetRelation instead")
		}
		return 0
	}

	arches := w.getArchetypes(filter)
	lengths := make([]uint32, len(arches))
	var totalEntities uint32 = 0
	for i, arch := range arches {
		lengths[i] = arch.Len()
		totalEntities += arch.Len()
	}

	for i, arch := range arches {
		archLen := lengths[i]

		if archLen == 0 {
			continue
		}

		newArch, start := w.exchangeArch(arch, archLen, add, rem, relation, hasRelation, target)
		batches.Add(newArch, arch, start, newArch.Len())
	}

	return int(totalEntities)
}

func (w *World) exchangeArch(oldArch *archetype, oldArchLen uint32, add []ID, rem []ID, relation ID, hasRelation bool, target Entity) (*archetype, uint32) {
	mask := w.getExchangeMask(oldArch.Mask, add, rem)
	oldIDs := oldArch.Components()

	if hasRelation {
		if !mask.Get(relation) {
			tp, _ := w.registry.ComponentType(relation.id)
			panic(fmt.Sprintf("can't add relation: resulting entity has no component %s", tp.Name()))
		}
		if !w.registry.IsRelation.Get(relation) {
			tp, _ := w.registry.ComponentType(relation.id)
			panic(fmt.Sprintf("can't add relation: %s is not a relation component", tp.Name()))
		}
	} else {
		target = oldArch.RelationTarget
		if !target.IsZero() && oldArch.Mask.ContainsAny(&w.registry.IsRelation) {
			for _, id := range rem {
				// Removing a relation
				if w.registry.IsRelation.Get(id) {
					target = Entity{}
					break
				}
			}
		}
	}

	arch := w.findOrCreateArchetype(oldArch, add, rem, target)

	startIdx := arch.Len()
	count := oldArchLen
	arch.AllocN(uint32(count))

	var i uint32
	for i = 0; i < count; i++ {
		idx := startIdx + i
		entity := oldArch.GetEntity(i)
		index := &w.entities[entity.id]
		arch.SetEntity(idx, entity)
		index.arch = arch
		index.index = idx

		for _, id := range oldIDs {
			if mask.Get(id) {
				comp := oldArch.Get(i, id)
				arch.SetPointer(idx, id, comp)
			}
		}
	}

	if !target.IsZero() {
		w.targetEntities.Set(target.id, true)
	}

	// Theoretically, it could be oldArchLen < oldArch.Len(),
	// which means we can't reset the archetype.
	// However, this should not be possible as processing an entity twice
	// would mean an illegal component addition/removal.
	oldArch.Reset()
	w.cleanupArchetype(oldArch)

	return arch, startIdx
}

// getRelation returns the target entity for an entity relation.
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//
// See [Relation] for details and examples.
func (w *World) getRelation(entity Entity, comp ID) Entity {
	if !w.entityPool.Alive(entity) {
		panic("can't get relation of a dead entity")
	}

	index := &w.entities[entity.id]
	w.checkRelation(index.arch, comp)

	return index.arch.RelationTarget
}

// getRelationUnchecked returns the target entity for an entity relation.
//
// getRelationUnchecked is an optimized version of [World.getRelation].
// Does not check if the entity is alive or that the component ID is applicable.
func (w *World) getRelationUnchecked(entity Entity, comp ID) Entity {
	_ = comp
	index := &w.entities[entity.id]
	return index.arch.RelationTarget
}

// setRelation sets the target entity for an entity relation.
//
// Panics:
//   - when called for a removed (and potentially recycled) entity.
//   - when called for a removed (and potentially recycled) target.
//   - when called for a missing component.
//   - when called for a component that is not a relation.
//   - when called on a locked world. Do not use during [Query] iteration!
//
// See [Relation] for details and examples.
func (w *World) setRelation(entity Entity, comp ID, target Entity) {
	w.checkLocked()

	if !w.entityPool.Alive(entity) {
		panic("can't set relation for a dead entity")
	}
	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	index := &w.entities[entity.id]
	w.checkRelation(index.arch, comp)

	oldArch := index.arch

	if oldArch.RelationTarget == target {
		return
	}

	arch := oldArch.node.GetArchetype(target)
	if arch == nil {
		arch = w.createArchetype(oldArch.node, target, true)
	}

	newIndex := arch.Alloc(entity)
	for _, id := range oldArch.node.Ids {
		comp := oldArch.Get(index.index, id)
		arch.SetPointer(newIndex, id, comp)
	}

	swapped := oldArch.Remove(index.index)

	if swapped {
		swapEntity := oldArch.GetEntity(index.index)
		w.entities[swapEntity.id].index = index.index
	}
	w.entities[entity.id] = entityIndex{arch: arch, index: newIndex}

	if !target.IsZero() {
		w.targetEntities.Set(target.id, true)
	}

	oldTarget := oldArch.RelationTarget
	w.cleanupArchetype(oldArch)

	if w.listener != nil {
		trigger := w.listener.Subscriptions() & event.TargetChanged
		if trigger != 0 && subscribes(trigger, nil, nil, w.listener.Components(), &comp, &comp) {
			w.listener.Notify(w, EntityEvent{Entity: entity, OldRelation: &comp, NewRelation: &comp, OldTarget: oldTarget, EventTypes: event.TargetChanged})
		}
	}
}

// set relation target in batches.
func (w *World) setRelationBatch(filter Filter, comp ID, target Entity) int {
	batches := batchArchetypes{}
	count := w.setRelationBatchNoNotify(filter, comp, target, &batches)
	if w.listener != nil && w.listener.Subscriptions().Contains(event.TargetChanged) {
		w.notifyQuery(&batches)
	}
	return count
}

func (w *World) setRelationBatchQuery(filter Filter, comp ID, target Entity) Query {
	batches := batchArchetypes{}
	w.setRelationBatchNoNotify(filter, comp, target, &batches)
	lock := w.lock()
	return newBatchQuery(w, lock, &batches)
}

func (w *World) setRelationBatchNoNotify(filter Filter, comp ID, target Entity, batches *batchArchetypes) int {
	w.checkLocked()

	if !target.IsZero() && !w.entityPool.Alive(target) {
		panic("can't make a dead entity a relation target")
	}

	arches := w.getArchetypes(filter)
	lengths := make([]uint32, len(arches))
	var totalEntities uint32 = 0
	for i, arch := range arches {
		lengths[i] = arch.Len()
		totalEntities += arch.Len()
	}

	for i, arch := range arches {
		archLen := lengths[i]

		if archLen == 0 {
			continue
		}

		if arch.RelationTarget == target {
			continue
		}

		newArch, start, end := w.setRelationArch(arch, archLen, comp, target)
		batches.Add(newArch, arch, start, end)
	}
	return int(totalEntities)
}

func (w *World) setRelationArch(oldArch *archetype, oldArchLen uint32, comp ID, target Entity) (*archetype, uint32, uint32) {
	w.checkRelation(oldArch, comp)

	// Before, entities with unchanged target were included in the query,
	// and events were emitted for them. Seems better to skip them completely,
	// which is done in World.setRelationBatchNoNotify.
	//if oldArch.RelationTarget == target {
	//	return oldArch, 0, oldArchLen
	//}

	oldIDs := oldArch.Components()

	arch := oldArch.node.GetArchetype(target)
	if arch == nil {
		arch = w.createArchetype(oldArch.node, target, true)
	}

	startIdx := arch.Len()
	count := oldArchLen
	arch.AllocN(count)

	var i uint32
	for i = 0; i < count; i++ {
		idx := startIdx + i
		entity := oldArch.GetEntity(i)
		index := &w.entities[entity.id]
		arch.SetEntity(idx, entity)
		index.arch = arch
		index.index = idx

		for _, id := range oldIDs {
			comp := oldArch.Get(i, id)
			arch.SetPointer(idx, id, comp)
		}
	}

	if !target.IsZero() {
		w.targetEntities.Set(target.id, true)
	}

	// Theoretically, it could be oldArchLen < oldArch.Len(),
	// which means we can't reset the archetype.
	// However, this should not be possible as processing an entity twice
	// would mean an illegal component addition/removal.
	oldArch.Reset()
	w.cleanupArchetype(oldArch)

	return arch, uint32(startIdx), arch.Len()
}

func (w *World) checkRelation(arch *archetype, comp ID) {
	if arch.node.Relation.id != comp.id {
		w.relationError(arch, comp)
	}
}

func (w *World) relationError(arch *archetype, comp ID) {
	if !arch.HasComponent(comp) {
		panic(fmt.Sprintf("entity does not have relation component %v", w.registry.Types[comp.id]))
	}
	panic(fmt.Sprintf("not a relation component: %v", w.registry.Types[comp.id]))
}

// lock the world and get the lock bit for later unlocking.
func (w *World) lock() uint8 {
	return w.locks.Lock()
}

// unlock unlocks the given lock bit.
func (w *World) unlock(l uint8) {
	w.locks.Unlock(l)
}

// checkLocked checks if the world is locked, and panics if so.
func (w *World) checkLocked() {
	if w.IsLocked() {
		panic("attempt to modify a locked world")
	}
}

// Copies a component to an entity
func (w *World) copyTo(entity Entity, id ID, comp interface{}) unsafe.Pointer {
	if !w.Has(entity, id) {
		panic("can't copy component into entity that has no such component type")
	}
	index := &w.entities[entity.id]
	arch := index.arch

	return arch.Set(index.index, id, comp)
}

// Tries to find an archetype by traversing the archetype graph,
// searching by mask and extending the graph if necessary.
// A new archetype is created for the final graph node if not already present.
func (w *World) findOrCreateArchetype(start *archetype, add []ID, rem []ID, target Entity) *archetype {
	curr := start.node
	mask := start.Mask
	relation := start.RelationComponent
	hasRelation := start.HasRelationComponent
	for _, id := range rem {
		// Not required, as removing happens only via exchange,
		// which calls getExchangeMask, which does the same check.
		//if !mask.Get(id) {
		//	panic(fmt.Sprintf("entity does not have a component of type %v, or it was removed twice", w.registry.Types[id.id]))
		//}
		mask.Set(id, false)
		if w.registry.IsRelation.Get(id) {
			relation = ID{}
			hasRelation = false
		}
		if next, ok := curr.TransitionRemove.Get(id.id); ok {
			curr = next
		} else {
			next, _ := w.findOrCreateArchetypeSlow(mask, relation, hasRelation)
			next.TransitionAdd.Set(id.id, curr)
			curr.TransitionRemove.Set(id.id, next)
			curr = next
		}
	}
	for _, id := range add {
		if mask.Get(id) {
			panic(fmt.Sprintf("entity already has component of type %v, or it was added twice", w.registry.Types[id.id]))
		}
		if start.Mask.Get(id) {
			panic(fmt.Sprintf("component of type %v added and removed in the same exchange operation", w.registry.Types[id.id]))
		}
		mask.Set(id, true)
		if w.registry.IsRelation.Get(id) {
			if hasRelation {
				panic("entity already has a relation component")
			}
			relation = id
			hasRelation = true
		}
		if next, ok := curr.TransitionAdd.Get(id.id); ok {
			curr = next
		} else {
			next, _ := w.findOrCreateArchetypeSlow(mask, relation, hasRelation)
			next.TransitionRemove.Set(id.id, curr)
			curr.TransitionAdd.Set(id.id, next)
			curr = next
		}
	}
	arch := curr.GetArchetype(target)
	if arch == nil {
		arch = w.createArchetype(curr, target, true)
	}
	return arch
}

// Tries to find an archetype for a mask, when it can't be reached through the archetype graph.
// Creates an archetype graph node.
func (w *World) findOrCreateArchetypeSlow(mask Mask, relation ID, hasRelation bool) (*archNode, bool) {
	if arch, ok := w.findArchetypeSlow(mask); ok {
		return arch, false
	}
	return w.createArchetypeNode(mask, relation, hasRelation), true
}

// Searches for an archetype by a mask.
func (w *World) findArchetypeSlow(mask Mask) (*archNode, bool) {
	length := w.nodes.Len()
	var i int32
	for i = 0; i < length; i++ {
		nd := w.nodes.Get(i)
		if nd.Mask == mask {
			return nd, true
		}
	}
	return nil, false
}

// Creates a node in the archetype graph.
func (w *World) createArchetypeNode(mask Mask, relation ID, hasRelation bool) *archNode {
	capInc := w.config.CapacityIncrement
	if hasRelation {
		capInc = w.config.RelationCapacityIncrement
	}

	types := mask.toTypes(&w.registry)

	w.nodeData.Add(nodeData{})
	w.nodes.Add(newArchNode(mask, w.nodeData.Get(w.nodeData.Len()-1), relation, hasRelation, capInc, types))
	nd := w.nodes.Get(w.nodes.Len() - 1)
	w.relationNodes = append(w.relationNodes, nd)
	w.nodePointers = append(w.nodePointers, nd)

	return nd
}

// Creates an archetype for the given archetype graph node.
// Initializes the archetype with a capacity according to CapacityIncrement if forStorage is true,
// and with a capacity of 1 otherwise.
func (w *World) createArchetype(node *archNode, target Entity, forStorage bool) *archetype {
	var arch *archetype
	layouts := capacityNonZero(w.registry.Count(), int(layoutChunkSize))

	if node.HasRelation {
		arch = node.CreateArchetype(uint8(layouts), target)
	} else {
		w.archetypes.Add(archetype{})
		w.archetypeData.Add(archetypeData{})
		archIndex := w.archetypes.Len() - 1
		arch = w.archetypes.Get(archIndex)
		arch.Init(node, w.archetypeData.Get(archIndex), archIndex, forStorage, uint8(layouts), Entity{})
		node.SetArchetype(arch)
	}
	w.filterCache.addArchetype(arch)
	return arch
}

// Returns all archetypes that match the given filter.
func (w *World) getArchetypes(filter Filter) []*archetype {
	if cached, ok := filter.(*CachedFilter); ok {
		return w.filterCache.get(cached).Archetypes.pointers
	}

	arches := []*archetype{}
	nodes := w.nodePointers

	for _, nd := range nodes {
		if !nd.IsActive || !nd.Matches(filter) {
			continue
		}

		if rf, ok := filter.(*RelationFilter); ok {
			target := rf.Target
			if arch, ok := nd.archetypeMap[target]; ok {
				arches = append(arches, arch)
			}
			continue
		}

		nodeArches := nd.Archetypes()
		ln2 := int32(nodeArches.Len())
		var j int32
		for j = 0; j < ln2; j++ {
			a := nodeArches.Get(j)
			if a.IsActive() {
				arches = append(arches, a)
			}
		}
	}

	return arches
}

// Removes the archetype if it is empty, and has a relation to a dead target.
func (w *World) cleanupArchetype(arch *archetype) {
	if arch.Len() > 0 || !arch.node.HasRelation {
		return
	}
	target := arch.RelationTarget
	if target.IsZero() || w.Alive(target) {
		return
	}

	w.removeArchetype(arch)
}

// Removes empty archetypes that have a target relation to the given entity.
func (w *World) cleanupArchetypes(target Entity) {
	for _, node := range w.relationNodes {
		if arch, ok := node.archetypeMap[target]; ok && arch.Len() == 0 {
			w.removeArchetype(arch)
		}
	}
}

// Removes/de-activates a relation archetype.
func (w *World) removeArchetype(arch *archetype) {
	arch.node.RemoveArchetype(arch)
	w.Cache().removeArchetype(arch)
}

// Extend the number of access layouts in archetypes.
func (w *World) extendArchetypeLayouts(count uint8) {
	len := w.nodes.Len()
	var i int32
	for i = 0; i < len; i++ {
		w.nodes.Get(i).ExtendArchetypeLayouts(count)
	}
}

// componentID returns the ID for a component type, and registers it if not already registered.
func (w *World) componentID(tp reflect.Type) ID {
	id, newID := w.registry.ComponentID(tp)
	if newID {
		if w.IsLocked() {
			w.registry.unregisterLastComponent()
			panic("attempt to register a new component in a locked world")
		}
		if id > 0 && id%layoutChunkSize == 0 {
			w.extendArchetypeLayouts(id + layoutChunkSize)
		}
	}
	return ID{id: id}
}

// resourceID returns the ID for a resource type, and registers it if not already registered.
func (w *World) resourceID(tp reflect.Type) ResID {
	id, _ := w.resources.registry.ComponentID(tp)
	return ResID{id: id}
}

// closeQuery closes a query and unlocks the world.
func (w *World) closeQuery(query *Query) {
	query.nodeIndex = -2
	query.archIndex = -2
	w.unlock(query.lockBit)

	if w.listener != nil {
		if arch, ok := query.nodeArchetypes.(*batchArchetypes); ok {
			w.notifyQuery(arch)
		}
	}
}

// notifies the listener for all entities on a batch query.
func (w *World) notifyQuery(batchArch *batchArchetypes) {
	count := batchArch.Len()
	var i int32
	for i = 0; i < count; i++ {
		arch := batchArch.Get(i)

		var newRel *ID
		if arch.HasRelationComponent {
			newRel = &arch.RelationComponent
		}

		event := EntityEvent{
			Entity: Entity{}, Added: arch.Mask, Removed: Mask{}, AddedIDs: batchArch.Added, RemovedIDs: batchArch.Removed,
			OldRelation: nil, NewRelation: newRel,
			OldTarget: Entity{}, EventTypes: 0,
		}

		oldArch := batchArch.OldArchetype[i]
		relChanged := newRel != nil
		targChanged := !arch.RelationTarget.IsZero()

		if oldArch != nil {
			var oldRel *ID
			if oldArch.HasRelationComponent {
				oldRel = &oldArch.RelationComponent
			}
			relChanged = false
			if oldRel != nil || newRel != nil {
				relChanged = (oldRel == nil) != (newRel == nil) || *oldRel != *newRel
			}
			targChanged = oldArch.RelationTarget != arch.RelationTarget
			changed := event.Added.Xor(&oldArch.node.Mask)
			event.Added = changed.And(&event.Added)
			event.Removed = changed.And(&oldArch.node.Mask)
			event.OldTarget = oldArch.RelationTarget
			event.OldRelation = oldRel
		}

		bits := subscription(oldArch == nil, false, len(batchArch.Added) > 0, len(batchArch.Removed) > 0, relChanged, relChanged || targChanged)
		event.EventTypes = bits

		trigger := w.listener.Subscriptions() & bits
		if trigger != 0 && subscribes(trigger, &event.Added, &event.Removed, w.listener.Components(), event.OldRelation, event.NewRelation) {
			start, end := batchArch.StartIndex[i], batchArch.EndIndex[i]
			var e uint32
			for e = start; e < end; e++ {
				entity := arch.GetEntity(e)
				event.Entity = entity
				w.listener.Notify(w, event)
			}
		}
	}
}
