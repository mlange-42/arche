package ecs

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityPool(t *testing.T) {
	p := newImplicitListEntityPool()

	expectedAll := []Entity{newEntity(0), newEntity(1), newEntity(2), newEntity(3), newEntity(4)}

	for i := 0; i < 5; i++ {
		_ = p.Get()
	}
	assert.Equal(t, expectedAll, p.entities, "Wrong initial entities")

	e0 := p.entities[0]
	p.Recycle(e0)
	assert.False(t, p.Alive(e0), "Dead entity should not be alive")

	e0Old := e0
	e0 = p.Get()
	expectedAll[1].gen++
	assert.True(t, p.Alive(e0), "Recycled entity of new generation should be alive")
	assert.False(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")

	assert.Equal(t, expectedAll, p.entities, "Wrong entities after get/recycle")

	e0Old = p.entities[0]
	for i := 0; i < 5; i++ {
		p.Recycle(p.entities[i])
		expectedAll[i].gen++
	}

	assert.False(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")

	for i := 0; i < 5; i++ {
		_ = p.Get()
	}

	assert.False(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")
}

func TestEntityPoolStochastic(t *testing.T) {
	p := newImplicitListEntityPool()

	alive := map[Entity]bool{}
	for i := 0; i < 100; i++ {
		e := p.Get()
		alive[e] = true
	}
	for e := range alive {
		if rand.Float32() > 0.5 {
			continue
		}
		p.Recycle(e)
		alive[e] = false
	}
	for e, a := range alive {
		assert.Equal(t, a, p.Alive(e), "Wrong alive state of entity %v after removal", e)
	}
	for i := 0; i < 100; i++ {
		e := p.Get()
		alive[e] = true
	}
	for e, a := range alive {
		assert.Equal(t, a, p.Alive(e), "Wrong alive state of entity %v after recycling", e)
	}
	assert.Equal(t, uint32(0), p.available, "No more entities should be available")
}
