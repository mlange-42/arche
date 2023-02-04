package ecs

import (
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

	e1 := p.entities[1]
	p.Recycle(e1)
	assert.False(t, p.Alive(e1), "Dead entity should not be alive")

	e1Old := e1
	e1 = p.Get()
	expectedAll[1].gen++
	assert.True(t, p.Alive(e1), "Recycled entity of new generation should be alive")
	assert.False(t, p.Alive(e1Old), "Recycled entity of old generation should not be alive")

	assert.Equal(t, expectedAll, p.entities, "Wrong entities after get/recycle")

	e0Old := p.entities[0]
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
