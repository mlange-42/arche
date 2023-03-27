package ecs

type cacheEntry struct {
	ID         ID
	Filter     Filter
	Archetypes pagedPointerArr32[archetype]
}

// Cache provides filter caching to speed up queries.
type Cache struct {
	bitPool       bitPool
	indices       []int
	filters       []cacheEntry
	getArchetypes func(f Filter) pagedPointerArr32[archetype]
}

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

// Register registers a new filter.
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

// Unregister un-registers a filter.
func (c *Cache) Unregister(f *CachedFilter) {
	idx := c.indices[f.id]
	if idx >= MaskTotalBits {
		panic("no filter for id found to unregister")
	}
	c.indices[f.id] = MaskTotalBits

	last := len(c.filters) - 1
	if idx != last {
		c.filters[idx], c.filters[last] = c.filters[last], c.filters[idx]
		c.indices[c.filters[idx].ID] = idx
	}
	c.filters[last] = cacheEntry{}
	c.filters = c.filters[:last]
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
