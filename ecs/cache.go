package ecs

// Cache entry for a [Filter].
type cacheEntry struct {
	ID         ID                // Filter ID.
	Filter     Filter            // The underlying filter.
	Archetypes archetypePointers // Archetypes matching the filter.
}

// Cache provides [Filter] caching to speed up queries.
//
// Access it using [World.Cache].
//
// For registered filters, the relevant archetypes are tracked internally,
// so that there are no mask checks required during iteration.
// This is particularly helpful to avoid query iteration slowdown by a very high number of archetypes.
// If the number of archetypes exceeds approx. 50-100, uncached filters experience a slowdown.
// The relative slowdown increases with lower numbers of entities queried (below 10.000).
// Cached filters avoid this slowdown.
//
// The overhead of tracking cached filters is very low, as updates are required only when new archetypes are created.
type Cache struct {
	indices       []int                            // Mapping from filter IDs to indices in filters
	filters       []cacheEntry                     // The cached filters, indexed by indices
	getArchetypes func(f Filter) archetypePointers // Callback for getting archetypes for a new filter from the world
	bitPool       bitPool                          // Pool for filter IDs
}

// newCache creates a new [Cache].
func newCache() Cache {
	indices := make([]int, MaskTotalBits)
	for i := 0; i < MaskTotalBits; i++ {
		indices[i] = MaskTotalBits
	}
	return Cache{
		bitPool: bitPool{},
		indices: indices,
		filters: []cacheEntry{},
	}
}

// Register a [Filter].
//
// Use the returned [CachedFilter] to construct queries:
//
//	filter := All(posID, velID)
//	cached := world.Cache().Register(&filter)
//	query := world.Query(&cached)
func (c *Cache) Register(f Filter) CachedFilter {
	id := c.bitPool.Get()
	c.filters = append(c.filters,
		cacheEntry{
			ID:         id,
			Filter:     f,
			Archetypes: c.getArchetypes(f),
		})
	c.indices[id] = len(c.filters) - 1
	return CachedFilter{f, id}
}

// Unregister a filter.
//
// Returns the original filter.
func (c *Cache) Unregister(f *CachedFilter) Filter {
	idx := c.indices[f.id]
	if idx >= MaskTotalBits {
		panic("no filter for id found to unregister")
	}
	filter := c.filters[idx].Filter
	c.indices[f.id] = MaskTotalBits

	last := len(c.filters) - 1
	if idx != last {
		c.filters[idx], c.filters[last] = c.filters[last], c.filters[idx]
		c.indices[c.filters[idx].ID] = idx
	}
	c.filters[last] = cacheEntry{}
	c.filters = c.filters[:last]

	return filter
}

// Returns the [cacheEntry] for the given filter.
//
// Panics if there is no entry for the filter's ID.
func (c *Cache) get(f *CachedFilter) *cacheEntry {
	idx := c.indices[f.id]
	if idx < MaskTotalBits {
		return &c.filters[idx]
	}
	panic("no filter for id found")
}

// Adds an archetype.
//
// Iterates over all filters and adds the archetype to the resp. entry where the filter matches.
func (c *Cache) addArchetype(arch *archetype) {
	ln := len(c.filters)
	for i := 0; i < ln; i++ {
		e := &c.filters[i]
		if e.Filter.Matches(arch.Mask) {
			e.Archetypes.Add(arch)
		}
	}
}
