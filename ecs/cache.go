package ecs

type cacheEntry struct {
	ID         ID
	Filter     Filter
	Archetypes pagedPointerArr32[archetype]
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
	bitPool       bitPool
	indices       []int
	filters       []cacheEntry
	getArchetypes func(f Filter) pagedPointerArr32[archetype]
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
//	filter := world.Cache().Register(ecs.All(posID, velID))
//	query := world.Query(&filter)
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

func (c *Cache) get(f *CachedFilter) *cacheEntry {
	idx := c.indices[f.id]
	if idx < MaskTotalBits {
		return &c.filters[idx]
	}
	panic("no filter for id found")
}

func (c *Cache) addArchetype(arch *archetype) {
	ln := len(c.filters)
	for i := 0; i < ln; i++ {
		e := &c.filters[i]
		if e.Filter.Matches(arch.Mask) {
			e.Archetypes.Add(arch)
		}
	}
}
