package ecs

// Builder for more flexible entity creation.
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

// WithRelation sets the relation component for the builder.
func (b *Builder) WithRelation(comp ID) *Builder {
	b.hasTarget = true
	b.targetID = comp
	return b
}

// New creates an entity.
func (b *Builder) New() Entity {
	if b.comps == nil {
		return b.world.NewEntity(b.ids...)
	}
	return b.world.NewEntityWith(b.comps...)
}

// NewRelation creates an entity with a relation target.
func (b *Builder) NewRelation(target Entity) Entity {
	if !b.hasTarget {
		panic("entity builder has no target")
	}
	if b.comps == nil {
		return b.world.newEntityTarget(b.targetID, target, b.ids...)
	}
	return b.world.newEntityTargetWith(b.targetID, target, b.comps...)
}

// NewBatch creates many entities.
func (b *Builder) NewBatch(count int) {
	if b.comps == nil {
		b.world.newEntities(count, -1, Entity{}, b.ids...)
	} else {
		b.world.newEntitiesWith(count, -1, Entity{}, b.comps...)
	}
}

// NewBatchRelation creates many entities with a relation target.
func (b *Builder) NewBatchRelation(count int, target Entity) {
	if !b.hasTarget {
		panic("entity builder has no target")
	}
	if b.comps == nil {
		b.world.newEntities(count, int8(b.targetID), target, b.ids...)
	} else {
		b.world.newEntitiesWith(count, int8(b.targetID), target, b.comps...)
	}
}

// NewQuery creates many entities and returns a query over them.
func (b *Builder) NewQuery(count int) Query {
	if b.comps == nil {
		return b.world.newEntitiesQuery(count, -1, Entity{}, b.ids...)
	}
	return b.world.newEntitiesWithQuery(count, -1, Entity{}, b.comps...)
}

// NewQueryRelation creates many entities with a relation target and returns a query over them.
func (b *Builder) NewQueryRelation(count int, target Entity) Query {
	if !b.hasTarget {
		panic("entity builder has no target")
	}
	if b.comps == nil {
		return b.world.newEntitiesQuery(count, int8(b.targetID), target, b.ids...)
	}
	return b.world.newEntitiesWithQuery(count, int8(b.targetID), target, b.comps...)
}
