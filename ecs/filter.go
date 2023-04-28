package ecs

// Filter is the interface for logic filters.
// Filters are required to query entities using [World.Query].
//
// See [Mask] and [MaskFilter] for basic filters.
// For type-safe generics queries, see package [github.com/mlange-42/arche/generic].
// For advanced filtering, see package [github.com/mlange-42/arche/filter].
type Filter interface {
	// Matches the filter against a bitmask, i.e. a component composition.
	Matches(bits Mask, relation *Entity) bool
}

// MaskFilter is a [Filter] for including and excluding certain components.
// See [All] and [Mask.Without].
type MaskFilter struct {
	Include Mask // Components to include.
	Exclude Mask // Components to exclude.
}

// Matches matches a filter against a mask.
func (f *MaskFilter) Matches(bits Mask, relation *Entity) bool {
	return bits.Contains(f.Include) && (f.Exclude.IsZero() || !bits.ContainsAny(f.Exclude))
}

// CachedFilter is a filter that is cached by the world.
//
// Create it using [Cache.Register].
type CachedFilter struct {
	filter Filter
	id     ID
}

// Matches matches a filter against a mask.
func (f *CachedFilter) Matches(bits Mask, relation *Entity) bool {
	return f.filter.Matches(bits, relation)
}
