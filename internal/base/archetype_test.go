package base

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type position struct {
	X int
	Y int
}

type rotation struct {
	Angle int
}

func TestArchetype(t *testing.T) {
	comps := []ComponentType{
		{ID: 0, Type: reflect.TypeOf(position{})},
		{ID: 1, Type: reflect.TypeOf(rotation{})},
	}

	arch := Archetype{}
	arch.Init(32, comps...)

	arch.Add(
		NewEntity(0),
		Component{ID: 0, Component: &position{1, 2}},
		Component{ID: 1, Component: &rotation{3}},
	)

	arch.Add(
		NewEntity(1),
		Component{ID: 0, Component: &position{4, 5}},
		Component{ID: 1, Component: &rotation{6}},
	)

	assert.Equal(t, 2, int(arch.entities.Len()))
	assert.Equal(t, 2, int(arch.components[0].Len()))
	assert.Equal(t, 2, int(arch.components[1].Len()))

	e0 := arch.GetEntity(0)
	e1 := arch.GetEntity(1)
	assert.Equal(t, Entity{0, 0}, e0)
	assert.Equal(t, Entity{1, 0}, e1)

	pos0 := (*position)(arch.Get(0, ID(0)))
	rot0 := (*rotation)(arch.Get(0, ID(1)))
	pos1 := (*position)(arch.Get(1, ID(0)))
	rot1 := (*rotation)(arch.Get(1, ID(1)))

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

	pos0 = (*position)(arch.Get(0, ID(0)))
	rot0 = (*rotation)(arch.Get(0, ID(1)))
	assert.Equal(t, 4, pos0.X)
	assert.Equal(t, 5, pos0.Y)
	assert.Equal(t, 6, rot0.Angle)

	assert.Panics(t, func() {
		arch.Add(
			NewEntity(1),
			Component{ID: 0, Component: &position{4, 5}},
		)
	})

	assert.Panics(t, func() {
		arch.AddPointer(NewEntity(1))
	})
}

func TestNewArchetype(t *testing.T) {
	comps := []ComponentType{
		{ID: 0, Type: reflect.TypeOf(position{})},
		{ID: 1, Type: reflect.TypeOf(rotation{})},
	}
	arch := Archetype{}
	arch.Init(32, comps...)

	comps = []ComponentType{
		{ID: 1, Type: reflect.TypeOf(rotation{})},
		{ID: 0, Type: reflect.TypeOf(position{})},
	}
	assert.Panics(t, func() {
		arch := Archetype{}
		arch.Init(32, comps...)
	})
}

func BenchmarkArchetypeAccess_1000(b *testing.B) {
	b.StopTimer()
	comps := []ComponentType{
		{ID: 0, Type: reflect.TypeOf(testStruct{})},
	}

	arch := Archetype{}
	arch.Init(32, comps...)

	for i := 0; i < 1000; i++ {
		arch.Add(
			NewEntity(Eid(i)),
			Component{ID: 0, Component: &testStruct{}},
		)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		len := int(arch.Len())
		id := ID(0)
		for j := 0; j < len; j++ {
			pos := (*testStruct)(arch.Get(i, id))
			_ = pos
		}
	}
}

func BenchmarkArchetypeAccessUnsafe_1000(b *testing.B) {
	b.StopTimer()
	comps := []ComponentType{
		{ID: 0, Type: reflect.TypeOf(testStruct{})},
	}

	arch := Archetype{}
	arch.Init(32, comps...)

	for i := 0; i < 1000; i++ {
		arch.Add(
			NewEntity(Eid(i)),
			Component{ID: 0, Component: &testStruct{}},
		)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		len := int(arch.Len())
		id := ID(0)
		for j := 0; j < len; j++ {
			pos := (*testStruct)(arch.GetUnsafe(i, id))
			_ = pos
		}
	}
}
