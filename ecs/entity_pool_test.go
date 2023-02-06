package ecs

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityPoolConstructor(t *testing.T) {
	_ = newEntityPool()
}

func TestEntityPool(t *testing.T) {
	p := newEntityPool()

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
	expectedAll[0].gen++
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
	for i := 0; i < 10; i++ {
		p := newEntityPool()

		alive := map[Entity]bool{}
		for i := 0; i < 10; i++ {
			e := p.Get()
			alive[e] = true
		}
		for e, isAlive := range alive {
			assert.Equal(t, isAlive, p.Alive(e), "Wrong alive state of entity %v after initialization", e)
			if rand.Float32() > 0.75 {
				continue
			}
			p.Recycle(e)
			alive[e] = false
		}
		for e, isAlive := range alive {
			assert.Equal(t, isAlive, p.Alive(e), "Wrong alive state of entity %v after 1st removal. Entity is %v", e, p.entities[e.id])
		}
		for i := 0; i < 10; i++ {
			e := p.Get()
			alive[e] = true
		}
		for e, isAlive := range alive {
			assert.Equal(t, isAlive, p.Alive(e), "Wrong alive state of entity %v after 1st recycling. Entity is %v", e, p.entities[e.id])
		}
		assert.Equal(t, uint32(0), p.available, "No more entities should be available")

		for e, isAlive := range alive {
			if !isAlive || rand.Float32() > 0.75 {
				continue
			}
			p.Recycle(e)
			alive[e] = false
		}
		for e, a := range alive {
			assert.Equal(t, a, p.Alive(e), "Wrong alive state of entity %v after 2nd removal. Entity is %v", e, p.entities[e.id])
		}
	}
}
