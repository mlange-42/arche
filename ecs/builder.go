package ecs

// Builder for more flexible and batched entity creation.
type Builder struct {
	world     *World
	ids       []ID
	comps     []Component
	hasTarget bool
	targetID  ID
}

// NewBuilder creates a builder from component IDs.
func NewBuilder(w *World, comps ...ID) *Builder {
	return &Builder{
		world: w,
		ids:   comps,
		comps: nil,
	}
}

// NewBuilderWith creates a builder from component pointers.
func NewBuilderWith(w *World, comps ...Component) *Builder {
	return &Builder{
		world: w,
		ids:   nil,
		comps: comps,
	}
}

// WithRelation sets the [Relation] component for the builder.
//
// Use in conjunction with the optional target argument of [Builder.New], [Builder.NewBatch] and [Builder.NewQuery].
//
// See [Relation] for details and examples.
func (b *Builder) WithRelation(comp ID) *Builder {
	b.hasTarget = true
	b.targetID = comp
	return b
}

// New creates an entity.
//
// The optional argument can be used to set the target [Entity] for the Builder's [Relation].
// See [Builder.WithRelation].
func (b *Builder) New(target ...Entity) Entity {
	if len(target) > 0 {
		if !b.hasTarget {
			panic("can't set target entity: builder has no relation")
		}
		if b.comps == nil {
			return b.world.newEntityTarget(b.targetID, target[0], b.ids...)
		}
		return b.world.newEntityTargetWith(b.targetID, target[0], b.comps...)
	}
	if b.comps == nil {
		return b.world.NewEntity(b.ids...)
	}
	return b.world.NewEntityWith(b.comps...)
}

// NewBatch creates many entities.
//
// The optional argument can be used to set the target [Entity] for the Builder's [Relation].
// See [Builder.WithRelation].
func (b *Builder) NewBatch(count int, target ...Entity) {
	if len(target) > 0 {
		if !b.hasTarget {
			panic("can't set target entity: builder has no relation")
		}
		if b.comps == nil {
			b.world.newEntities(count, int8(b.targetID), target[0], b.ids...)
			return
		}
		b.world.newEntitiesWith(count, int8(b.targetID), target[0], b.comps...)
		return
	}
	if b.comps == nil {
		b.world.newEntities(count, -1, Entity{}, b.ids...)
	} else {
		b.world.newEntitiesWith(count, -1, Entity{}, b.comps...)
	}
}

// NewQuery creates many entities and returns a query over them.
//
// The optional argument can be used to set the target [Entity] for the Builder's [Relation].
// See [Builder.WithRelation].
func (b *Builder) NewQuery(count int, target ...Entity) Query {
	if len(target) > 0 {
		if !b.hasTarget {
			panic("can't set target entity: builder has no relation")
		}
		if b.comps == nil {
			return b.world.newEntitiesQuery(count, int8(b.targetID), target[0], b.ids...)
		}
		return b.world.newEntitiesWithQuery(count, int8(b.targetID), target[0], b.comps...)
	}
	if b.comps == nil {
		return b.world.newEntitiesQuery(count, -1, Entity{}, b.ids...)
	}
	return b.world.newEntitiesWithQuery(count, -1, Entity{}, b.comps...)
}

// Add the builder's components to an entity.
//
// The optional argument can be used to set the target [Entity] for the Builder's [Relation].
// See [Builder.WithRelation].
func (b *Builder) Add(entity Entity, target ...Entity) {
	if len(target) > 0 {
		if !b.hasTarget {
			panic("can't set target entity: builder has no relation")
		}
		if b.comps == nil {
			b.world.exchange(entity, b.ids, nil, int8(b.targetID), target[0])
			return
		}
		b.world.assign(entity, int8(b.targetID), target[0], b.comps...)
		return
	}
	if b.comps == nil {
		b.world.Exchange(entity, b.ids, nil)
		return
	}
	b.world.Assign(entity, b.comps...)
}
