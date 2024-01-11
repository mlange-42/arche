package ecs

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityPoolConstructor(t *testing.T) {
	_ = newEntityPool(128)
}

func TestEntityPool(t *testing.T) {
	p := newEntityPool(128)

	expectedAll := []Entity{newEntity(0), newEntity(1), newEntity(2), newEntity(3), newEntity(4), newEntity(5)}
	expectedAll[0].gen = math.MaxUint32

	for i := 0; i < 5; i++ {
		_ = p.Get()
	}
	assert.Equal(t, expectedAll, p.entities, "Wrong initial entities")

	assert.Panics(t, func() { p.Recycle(p.entities[0]) })

	e0 := p.entities[1]
	p.Recycle(e0)
	assert.False(t, p.Alive(e0), "Dead entity should not be alive")

	e0Old := e0
	e0 = p.Get()
	expectedAll[1].gen++
	assert.True(t, p.Alive(e0), "Recycled entity of new generation should be alive")
	assert.False(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")

	assert.Equal(t, expectedAll, p.entities, "Wrong entities after get/recycle")

	e0Old = p.entities[1]
	for i := 0; i < 5; i++ {
		p.Recycle(p.entities[i+1])
		expectedAll[i+1].gen++
	}

	assert.False(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")

	for i := 0; i < 5; i++ {
		_ = p.Get()
	}

	assert.False(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")
	assert.False(t, p.Alive(Entity{}), "Zero entity should not be alive")
}

func TestEntityPoolStochastic(t *testing.T) {
	p := newEntityPool(128)

	for i := 0; i < 10; i++ {
		p.Reset()
		assert.Equal(t, 0, p.Len())
		assert.Equal(t, 0, p.Available())

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

func TestBitPool(t *testing.T) {
	p := bitPool{}

	for i := 0; i < MaskTotalBits; i++ {
		assert.Equal(t, i, int(p.Get()))
	}

	assert.Panics(t, func() { p.Get() })

	for i := 0; i < 10; i++ {
		p.Recycle(uint8(i))
	}
	for i := 9; i >= 0; i-- {
		assert.Equal(t, i, int(p.Get()))
	}

	assert.Panics(t, func() { p.Get() })

	p.Reset()

	for i := 0; i < MaskTotalBits; i++ {
		assert.Equal(t, i, int(p.Get()))
	}

	assert.Panics(t, func() { p.Get() })

	for i := 0; i < 10; i++ {
		p.Recycle(uint8(i))
	}
	for i := 9; i >= 0; i-- {
		assert.Equal(t, i, int(p.Get()))
	}
}

func TestIntPool(t *testing.T) {
	p := newIntPool[int](16)

	for n := 0; n < 3; n++ {
		for i := 0; i < 32; i++ {
			assert.Equal(t, i, p.Get())
		}

		assert.Equal(t, 32, len(p.pool))

		p.Recycle(3)
		p.Recycle(4)
		assert.Equal(t, 4, p.Get())
		assert.Equal(t, 3, p.Get())

		p.Reset()
	}
}
