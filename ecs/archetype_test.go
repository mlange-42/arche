package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArchetype(t *testing.T) {
	comps := []componentType{
		{ID: 0, Type: reflect.TypeOf(position{})},
		{ID: 1, Type: reflect.TypeOf(rotation{})},
	}

	arch := archetype{}
	arch.Init(nil, 32, false, comps...)

	arch.Add(
		newEntity(0),
		Component{ID: 0, Comp: &position{1, 2}},
		Component{ID: 1, Comp: &rotation{3}},
	)

	arch.Add(
		newEntity(1),
		Component{ID: 0, Comp: &position{4, 5}},
		Component{ID: 1, Comp: &rotation{6}},
	)

	assert.Equal(t, 2, int(arch.entities.Len()))
	assert.Equal(t, 2, int(arch.Len()))

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
	assert.Equal(t, 1, int(arch.Len()))

	pos0 = (*position)(arch.Get(0, ID(0)))
	rot0 = (*rotation)(arch.Get(0, ID(1)))
	assert.Equal(t, 4, pos0.X)
	assert.Equal(t, 5, pos0.Y)
	assert.Equal(t, 6, rot0.Angle)

	assert.Panics(t, func() {
		arch.Add(
			newEntity(1),
			Component{ID: 0, Comp: &position{4, 5}},
		)
	})
}

func TestNewArchetype(t *testing.T) {
	comps := []componentType{
		{ID: 0, Type: reflect.TypeOf(position{})},
		{ID: 1, Type: reflect.TypeOf(rotation{})},
	}
	arch := archetype{}
	arch.Init(nil, 32, true, comps...)
	assert.Equal(t, 32, int(arch.Cap()))

	arch = archetype{}
	arch.Init(nil, 32, false, comps...)
	assert.Equal(t, 1, int(arch.Cap()))

	comps = []componentType{
		{ID: 1, Type: reflect.TypeOf(rotation{})},
		{ID: 0, Type: reflect.TypeOf(position{})},
	}
	assert.Panics(t, func() {
		arch := archetype{}
		arch.Init(nil, 32, true, comps...)
	})
}

func TestArchetypeAddGetSet(t *testing.T) {
	a := archetype{}

	comps := []componentType{
		{ID: 0, Type: reflect.TypeOf(testStruct0{})},
		{ID: 1, Type: reflect.TypeOf(label{})},
	}
	a.Init(nil, 1, true, comps...)

	assert.Equal(t, 1, int(a.Cap()))
	assert.Equal(t, 0, int(a.Len()))

	a.Add(Entity{1, 0}, Component{ID: 0, Comp: &testStruct0{100}}, Component{ID: 1, Comp: &label{}})
	a.Add(Entity{2, 0}, Component{ID: 0, Comp: &testStruct0{200}}, Component{ID: 1, Comp: &label{}})

	ts := (*testStruct0)(a.Get(0, 0))
	assert.Equal(t, 100, int(ts.Val))

	a.Set(1, 0, &testStruct0{200})
	a.Set(1, 1, &label{})

	_ = (*testStruct0)(a.Get(1, 0))
	_ = (*label)(a.Get(1, 1))
}

func BenchmarkArchetypeAccess1_1000(b *testing.B) {
	BenchmarkArchetypeAccess_1000(b)
}

func BenchmarkArchetypeAccess2_1000(b *testing.B) {
	BenchmarkArchetypeAccess_1000(b)
}

func BenchmarkArchetypeAccess3_1000(b *testing.B) {
	BenchmarkArchetypeAccess_1000(b)
}

func BenchmarkArchetypeAccess4_1000(b *testing.B) {
	BenchmarkArchetypeAccess_1000(b)
}

func BenchmarkArchetypeAccess_1000(b *testing.B) {
	b.StopTimer()
	comps := []componentType{
		{ID: 0, Type: reflect.TypeOf(testStruct0{})},
	}

	arch := archetype{}
	arch.Init(nil, 32, true, comps...)

	for i := 0; i < 1000; i++ {
		arch.Alloc(newEntity(eid(i)), true)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		len := uintptr(arch.Len())
		id := ID(0)
		var j uintptr
		for j = 0; j < len; j++ {
			pos := (*testStruct0)(arch.Get(j, id))
			pos.Val = 1
		}
	}
}
