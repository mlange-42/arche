package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Position struct {
	X int
	Y int
}

type Rotation struct {
	Angle int
}

func TestArchetype(t *testing.T) {
	comps := []Component{
		{ID: 0, Component: Position{}},
		{ID: 1, Component: Rotation{}},
	}

	arch := newArchetype(comps...)

	arch.Add(
		newEntity(0),
		Component{ID: 0, Component: &Position{1, 2}},
		Component{ID: 1, Component: &Rotation{3}},
	)

	arch.Add(
		newEntity(1),
		Component{ID: 0, Component: &Position{4, 5}},
		Component{ID: 1, Component: &Rotation{6}},
	)

	assert.Equal(t, 2, int(arch.entities.Len()))
	assert.Equal(t, 2, int(arch.components[0].Len()))
	assert.Equal(t, 2, int(arch.components[1].Len()))

	e0 := arch.GetEntity(0)
	e1 := arch.GetEntity(1)
	assert.Equal(t, Entity{0, 0}, e0)
	assert.Equal(t, Entity{1, 0}, e1)

	pos0 := (*Position)(arch.Get(0, ID(0)))
	rot0 := (*Rotation)(arch.Get(0, ID(1)))
	pos1 := (*Position)(arch.Get(1, ID(0)))
	rot1 := (*Rotation)(arch.Get(1, ID(1)))

	assert.Equal(t, 1, pos0.X)
	assert.Equal(t, 2, pos0.Y)
	assert.Equal(t, 3, rot0.Angle)
	assert.Equal(t, 4, pos1.X)
	assert.Equal(t, 5, pos1.Y)
	assert.Equal(t, 6, rot1.Angle)

	arch.Remove(0)
	assert.Equal(t, 1, int(arch.entities.Len()))
	assert.Equal(t, 1, int(arch.components[0].Len()))
	assert.Equal(t, 1, int(arch.components[1].Len()))

	pos0 = (*Position)(arch.Get(0, ID(0)))
	rot0 = (*Rotation)(arch.Get(0, ID(1)))
	assert.Equal(t, 4, pos0.X)
	assert.Equal(t, 5, pos0.Y)
	assert.Equal(t, 6, rot0.Angle)

	assert.Panics(t, func() {
		arch.Add(
			newEntity(1),
			Component{ID: 0, Component: &Position{4, 5}},
		)
	})
}

func TestNewArchetype(t *testing.T) {
	comps := []Component{
		{ID: 0, Component: Position{}},
		{ID: 1, Component: Rotation{}},
	}

	_ = NewArchetype(comps...)
}

func BenchmarkArchetypeAccess(b *testing.B) {
	comps := []Component{
		{ID: 0, Component: Position{}},
		{ID: 1, Component: Rotation{}},
	}

	arch := newArchetype(comps...)

	for i := 0; i < 1000; i++ {
		arch.Add(
			newEntity(i),
			Component{ID: 0, Component: &Position{1, 2}},
			Component{ID: 1, Component: &Rotation{3}},
		)
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			pos := (*Position)(arch.Get(i, ID(0)))
			_ = pos
		}
	}
}
