package ecs

// Cache entry for a [Filter].
type cacheEntry struct {
	ID         uint32            // Filter ID.
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
// The relative slowdown increases with lower numbers of entities queried (below a few thousand).
// Cached filters avoid this slowdown.
//
// Further, cached filters should be used for complex queries built with package [github.com/mlange-42/arche/filter].
//
// The overhead of tracking cached filters internally is very low, as updates are required only when new archetypes are created.
type Cache struct {
	indices       map[uint32]int                   // Mapping from filter IDs to indices in filters
	filters       []cacheEntry                     // The cached filters, indexed by indices
	getArchetypes func(f Filter) archetypePointers // Callback for getting archetypes for a new filter from the world
	intPool       intPool[uint32]                  // Pool for filter IDs
}

// newCache creates a new [Cache].
func newCache() Cache {
	return Cache{
		intPool: newIntPool[uint32](128),
		indices: map[uint32]int{},
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
	if _, ok := f.(*CachedFilter); ok {
		panic("filter is already cached")
	}
	id := c.intPool.Get()
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
	idx, ok := c.indices[f.id]
	if !ok {
		panic("no filter for id found to unregister")
	}
	filter := c.filters[idx].Filter
	delete(c.indices, f.id)

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
	if idx, ok := c.indices[f.id]; ok {
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
		if arch.Matches(e.Filter) {
			e.Archetypes.Add(arch)
		}
	}
}

// Removes an archetype.
//
// Iterates over all filters and removes the archetype from the resp. entry where the filter matches.
func (c *Cache) removeArchetype(arch *archetype) {
	ln := len(c.filters)
	for i := 0; i < ln; i++ {
		e := &c.filters[i]
		if arch.Matches(e.Filter) {
			e.Archetypes.Remove(arch)
		}
	}
}
