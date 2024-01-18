package ecs

// Builder for more flexible and batched entity creation.
type Builder struct {
	world       *World
	ids         []ID
	comps       []Component
	hasRelation bool
	relationID  ID
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
// Use in conjunction with the optional target argument of [Builder.New], [Builder.NewBatch] and [Builder.NewBatchQ].
//
// See [Relation] for details and examples.
func (b *Builder) WithRelation(comp ID) *Builder {
	b.hasRelation = true
	b.relationID = comp
	return b
}

// New creates an entity.
//
// The optional argument can be used to set the target [Entity] for the Builder's [Relation].
// See [Builder.WithRelation].
func (b *Builder) New(target ...Entity) Entity {
	if len(target) > 0 {
		if !b.hasRelation {
			panic("can't set target entity: builder has no relation")
		}
		if b.comps == nil {
			return b.world.newEntityTarget(b.relationID, target[0], b.ids...)
		}
		return b.world.newEntityTargetWith(b.relationID, target[0], b.comps...)
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
		if !b.hasRelation {
			panic("can't set target entity: builder has no relation")
		}
		if b.comps == nil {
			b.world.newEntities(count, b.relationID, true, target[0], b.ids...)
			return
		}
		b.world.newEntitiesWith(count, b.relationID, true, target[0], b.comps...)
		return
	}
	if b.comps == nil {
		b.world.newEntities(count, ID{}, false, Entity{}, b.ids...)
	} else {
		b.world.newEntitiesWith(count, ID{}, false, Entity{}, b.comps...)
	}
}

// NewBatchQ creates many entities and returns a query over them.
//
// The optional argument can be used to set the target [Entity] for the Builder's [Relation].
// See [Builder.WithRelation].
func (b *Builder) NewBatchQ(count int, target ...Entity) Query {
	if len(target) > 0 {
		if !b.hasRelation {
			panic("can't set target entity: builder has no relation")
		}
		if b.comps == nil {
			return b.world.newEntitiesQuery(count, b.relationID, true, target[0], b.ids...)
		}
		return b.world.newEntitiesWithQuery(count, b.relationID, true, target[0], b.comps...)
	}
	if b.comps == nil {
		return b.world.newEntitiesQuery(count, ID{}, false, Entity{}, b.ids...)
	}
	return b.world.newEntitiesWithQuery(count, ID{}, false, Entity{}, b.comps...)
}

// Add the builder's components to an entity.
//
// The optional argument can be used to set the target [Entity] for the Builder's [Relation].
// See [Builder.WithRelation].
func (b *Builder) Add(entity Entity, target ...Entity) {
	if len(target) > 0 {
		if !b.hasRelation {
			panic("can't set target entity: builder has no relation")
		}
		if b.comps == nil {
			b.world.exchange(entity, b.ids, nil, b.relationID, b.hasRelation, target[0])
			return
		}
		b.world.assign(entity, b.relationID, b.hasRelation, target[0], b.comps...)
		return
	}
	if b.comps == nil {
		b.world.Exchange(entity, b.ids, nil)
		return
	}
	b.world.Assign(entity, b.comps...)
}
