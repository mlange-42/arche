package ecs

type cacheEntry struct {
	Filter     Filter
	Archetypes pagedPointerArr32[archetype]
}

// Cache provides filter caching to speed up queries.
type Cache struct {
	// TODO write a better data structure!
	filters       map[int]*cacheEntry
	counter       int
	getArchetypes func(f Filter) pagedPointerArr32[archetype]
}

func newCache() Cache {
	return Cache{
		filters: map[int]*cacheEntry{},
		counter: 0,
	}
}

// Register registers a new filter.
func (c *Cache) Register(f Filter) CachedFilter {
	id := c.counter
	c.filters[id] = &cacheEntry{
		Filter:     f,
		Archetypes: c.getArchetypes(f),
	}
	c.counter++
	return CachedFilter{f, id}
}

// Unregister unregisters a.
func (c *Cache) Unregister(f *CachedFilter) {
	if _, ok := c.filters[f.id]; !ok {
		panic("no filter for id found to unregister")
	}
	delete(c.filters, f.id)
}

func (c *Cache) get(f *CachedFilter) *cacheEntry {
	if e, ok := c.filters[f.id]; ok {
		return e
	}
	panic("no filter for id found")
}

func (c *Cache) addArchetype(arch *archetype) {
	for _, e := range c.filters {
		if e.Filter.Matches(arch.Mask) {
			e.Archetypes.Add(arch)
		}
	}
}
