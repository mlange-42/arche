package ecs

// Filter is the interface for logic filters.
// Filters are required to query entities using [World.Query].
//
// See [Mask], [MaskFilter] anf [RelationFilter] for basic filters.
// For type-safe generics queries, see package [github.com/mlange-42/arche/generic].
// For advanced filtering, see package [github.com/mlange-42/arche/filter].
type Filter interface {
	// Matches the filter against a mask, i.e. a component composition.
	Matches(bits *Mask) bool
}

// MaskFilter is a [Filter] for including and excluding certain components.
//
// See [All], [Mask.Without] and [Mask.Exclusive].
type MaskFilter struct {
	Include Mask // Components to include.
	Exclude Mask // Components to exclude.
}

// Matches the filter against a mask.
func (f *MaskFilter) Matches(bits *Mask) bool {
	return bits.Contains(&f.Include) && (!bits.ContainsAny(&f.Exclude) || f.Exclude.IsZero())
}

// RelationFilter is a [Filter] for a [Relation] target, in addition to components.
//
// See [Relation] for details and examples.
type RelationFilter struct {
	Filter Filter // Components filter.
	Target Entity // Relation target entity.
}

// NewRelationFilter creates a new [RelationFilter].
// It is a [Filter] for a [Relation] target, in addition to components.
func NewRelationFilter(filter Filter, target Entity) RelationFilter {
	return RelationFilter{
		Filter: filter,
		Target: target,
	}
}

// Matches the filter against a mask.
func (f *RelationFilter) Matches(bits *Mask) bool {
	return f.Filter.Matches(bits)
}

// CachedFilter is a filter that is cached by the world.
//
// Create a cached filter from any other filter using [Cache.Register].
// For details on caching, see [Cache].
type CachedFilter struct {
	filter Filter
	id     uint32
}

// Matches the filter against a mask.
func (f *CachedFilter) Matches(bits *Mask) bool {
	return f.filter.Matches(bits)
}
