package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterCache(t *testing.T) {
	cache := newCache()
	cache.getArchetypes = getArchetypes

	f1 := cache.Register(All(0, 1))
	f2 := cache.Register(All(0, 1, 2))
	assert.Equal(t, 0, f1.id)
	assert.Equal(t, 1, f2.id)

	e1 := cache.get(&f1)
	e2 := cache.get(&f2)

	assert.Equal(t, f1.Filter, e1.Filter)
	assert.Equal(t, f2.Filter, e2.Filter)

	cache.Unregister(&f1)
	cache.Unregister(&f2)

	assert.Panics(t, func() { cache.Unregister(&f1) })
	assert.Panics(t, func() { cache.get(&f1) })
}

func getArchetypes(f Filter) pagedPointerArr32[archetype] {
	return pagedPointerArr32[archetype]{}
}
