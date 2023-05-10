package ecs

// Filter is the interface for logic filters.
// Filters are required to query entities using [World.Query].
//
// See [Mask] and [MaskFilter] for basic filters.
// For type-safe generics queries, see package [github.com/mlange-42/arche/generic].
// For advanced filtering, see package [github.com/mlange-42/arche/filter].
type Filter interface {
	// Matches the filter against a bitmask, i.e. a component composition.
	Matches(bits Mask) bool
}

// MaskFilter is a [Filter] for including and excluding certain components.
// See [All] and [Mask.Without].
type MaskFilter struct {
	Include Mask // Components to include.
	Exclude Mask // Components to exclude.
}

// Matches matches a filter against a mask.
func (f *MaskFilter) Matches(bits Mask) bool {
	return bits.Contains(f.Include) && (f.Exclude.IsZero() || !bits.ContainsAny(f.Exclude))
}

// RelationFilter is a [Filter] for a [Relation] target, in addition to components.
type RelationFilter struct {
	Filter Filter // Components filter.
	Target Entity // Relation target entity.
}

// NewRelationFilter creates a new [Relation] filter.
// It is a [Filter] for a [Relation] target, in addition to components.
//
// Logic filters ignore relation targets. Thus, a relation filter should be the outermost filter.
//
// See [Relation] for details and examples.
func NewRelationFilter(filter Filter, target Entity) RelationFilter {
	return RelationFilter{
		Filter: filter,
		Target: target,
	}
}

// Matches matches a filter against a mask.
func (f *RelationFilter) Matches(bits Mask) bool {
	return f.Filter.Matches(bits)
}

// CachedFilter is a filter that is cached by the world.
//
// Create it using [Cache.Register].
type CachedFilter struct {
	filter Filter
	id     uint32
}

// Matches matches a filter against a mask.
func (f *CachedFilter) Matches(bits Mask) bool {
	return f.filter.Matches(bits)
}
