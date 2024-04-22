package ecs

// Cache entry for a [Filter].
type cacheEntry struct {
	Filter     Filter              // The underlying filter.
	Indices    map[*archetype]int  // Map of archetype indices for removal.
	Archetypes pointers[archetype] // Nodes matching the filter.
	ID         uint32              // Filter ID.
}

// Cache provides [Filter] caching to speed up queries.
//
// Access it using [World.Cache].
//
// For registered filters, the relevant archetypes are tracked internally,
// so that there are no mask checks required during iteration.
// This is particularly helpful to avoid query iteration slowdown by a very high number of archetypes.
// If the number of archetypes exceeds approx. 50-100, uncached filters experience a slowdown.
// The relative slowdown increases with lower numbers of entities queried (noticeable below a few thousand entities).
// Cached filters avoid this slowdown.
//
// The overhead of tracking cached filters internally is very low, as updates are required only when new archetypes are created.
type Cache struct {
	indices       map[uint32]int              // Mapping from filter IDs to indices in filters
	filters       []cacheEntry                // The cached filters, indexed by indices
	getArchetypes func(f Filter) []*archetype // Callback for getting archetypes for a new filter from the world
	intPool       intPool[uint32]             // Pool for filter IDs
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
		panic("filter is already registered")
	}
	id := c.intPool.Get()
	c.filters = append(c.filters,
		cacheEntry{
			ID:         id,
			Filter:     f,
			Archetypes: pointers[archetype]{c.getArchetypes(f)},
			Indices:    nil,
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

// Adds a node.
//
// Iterates over all filters and adds the node to the resp. entry where the filter matches.
func (c *Cache) addArchetype(arch *archetype) {
	if !arch.HasRelation() {
		for i := range c.filters {
			e := &c.filters[i]
			if !e.Filter.Matches(&arch.Mask) {
				continue
			}
			e.Archetypes.Add(arch)
		}
		return
	}

	for i := range c.filters {
		e := &c.filters[i]
		if !e.Filter.Matches(&arch.Mask) {
			continue
		}
		if rf, ok := e.Filter.(*RelationFilter); ok {
			if rf.Target == arch.RelationTarget {
				e.Archetypes.Add(arch)
				// Not required: can't add after removing,
				// as the target entity is dead.
				// if e.Indices != nil { e.Indices[arch] = int(e.Archetypes.Len() - 1) }
			}
			continue
		}
		e.Archetypes.Add(arch)
		if e.Indices != nil {
			e.Indices[arch] = int(e.Archetypes.Len() - 1)
		}
	}
}

// Removes an archetype.
//
// Can only be used for archetypes that have a relation target.
// Archetypes without a relation are never removed.
func (c *Cache) removeArchetype(arch *archetype) {
	for i := range c.filters {
		e := &c.filters[i]

		if e.Indices == nil && e.Filter.Matches(&arch.Mask) {
			c.mapArchetypes(e)
		}

		if idx, ok := e.Indices[arch]; ok {
			swap := e.Archetypes.RemoveAt(idx)
			if swap {
				e.Indices[e.Archetypes.Get(int32(idx))] = idx
			}
			delete(e.Indices, arch)
		}
	}
}

func (c *Cache) mapArchetypes(e *cacheEntry) {
	e.Indices = map[*archetype]int{}
	for i, arch := range e.Archetypes.pointers {
		if arch.HasRelation() {
			e.Indices[arch] = i
		}
	}
}
