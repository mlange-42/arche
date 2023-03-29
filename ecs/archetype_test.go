package ecs

import (
	"reflect"
	"testing"
	"unsafe"

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

func TestArchetypeExtend(t *testing.T) {
	comps := []componentType{
		{ID: 0, Type: reflect.TypeOf(position{})},
		{ID: 1, Type: reflect.TypeOf(rotation{})},
	}
	arch := archetype{}
	arch.Init(nil, 8, true, comps...)

	assert.Equal(t, 8, int(arch.Cap()))
	assert.Equal(t, 0, int(arch.Len()))

	arch.extend(5)
	assert.Equal(t, 8, int(arch.Cap()))

	arch.extend(8)
	assert.Equal(t, 8, int(arch.Cap()))

	arch.extend(17)
	assert.Equal(t, 24, int(arch.Cap()))
}

func TestArchetypeAlloc(t *testing.T) {
	comps := []componentType{
		{ID: 0, Type: reflect.TypeOf(position{})},
		{ID: 1, Type: reflect.TypeOf(rotation{})},
	}
	arch := archetype{}
	arch.Init(nil, 8, true, comps...)

	assert.Equal(t, 8, int(arch.Cap()))
	assert.Equal(t, 0, int(arch.Len()))

	arch.AllocN(1)
	assert.Equal(t, 1, int(arch.Len()))

	arch.AllocN(7)
	assert.Equal(t, 8, int(arch.Len()))
	assert.Equal(t, 8, int(arch.Cap()))

	arch.AllocN(1)
	assert.Equal(t, 9, int(arch.Len()))
	assert.Equal(t, 16, int(arch.Cap()))
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

	a.SetEntity(1, Entity{100, 200})
	e := a.GetEntity(1)
	assert.Equal(t, Entity{100, 200}, e)

	a.Remove(0)
	assert.Equal(t, 1, int(a.Len()))
	a.Remove(0)
	assert.Equal(t, 0, int(a.Len()))
}

func TestArchetypeReset(t *testing.T) {
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

	assert.Equal(t, position{1, 2}, *(*position)(arch.Get(0, 0)))
	assert.Equal(t, position{4, 5}, *(*position)(arch.Get(1, 0)))
	assert.Equal(t, 2, int(arch.Len()))

	arch.Reset()
	assert.Equal(t, 0, int(arch.Len()))

	arch.Add(
		newEntity(0),
		Component{ID: 0, Comp: &position{10, 20}},
		Component{ID: 1, Comp: &rotation{3}},
	)

	arch.Add(
		newEntity(1),
		Component{ID: 0, Comp: &position{40, 50}},
		Component{ID: 1, Comp: &rotation{6}},
	)

	assert.Equal(t, position{10, 20}, *(*position)(arch.Get(0, 0)))
	assert.Equal(t, position{40, 50}, *(*position)(arch.Get(1, 0)))
	assert.Equal(t, 2, int(arch.Len()))
}

func TestArchetypeZero(t *testing.T) {
	comps := []componentType{
		{ID: 0, Type: reflect.TypeOf(position{})},
		{ID: 1, Type: reflect.TypeOf(rotation{})},
	}

	arch := archetype{}
	arch.Init(nil, 32, false, comps...)

	arch.Alloc(newEntity(0))
	arch.Alloc(newEntity(1))

	assert.Equal(t, position{0, 0}, *(*position)(arch.Get(0, 0)))
	assert.Equal(t, position{0, 0}, *(*position)(arch.Get(1, 0)))

	pos := (*position)(arch.Get(0, 0))
	pos.X = 100
	pos = (*position)(arch.Get(1, 0))
	pos.X = 100

	assert.Equal(t, position{100, 0}, *(*position)(arch.Get(0, 0)))
	assert.Equal(t, position{100, 0}, *(*position)(arch.Get(1, 0)))

	arch.Remove(0)
	arch.Remove(0)
	arch.Alloc(newEntity(0))
	arch.Alloc(newEntity(1))
	assert.Equal(t, position{0, 0}, *(*position)(arch.Get(0, 0)))
	assert.Equal(t, position{0, 0}, *(*position)(arch.Get(1, 0)))
}

func TestArchetypePointers(t *testing.T) {
	pt := archetypePointers{}

	a1 := archetype{}
	a2 := archetype{}
	a3 := archetype{}

	pt.Add(&a1)
	pt.Add(&a2)
	pt.Add(&a3)

	assert.Equal(t, 3, pt.Len())

	for i := 0; i < 45; i++ {
		pt.Add(&a3)
	}

	assert.Equal(t, unsafe.Pointer(&a1), unsafe.Pointer(pt.Get(0)))
	assert.Equal(t, unsafe.Pointer(&a2), unsafe.Pointer(pt.Get(1)))
	assert.Equal(t, unsafe.Pointer(&a3), unsafe.Pointer(pt.Get(2)))
}

func BenchmarkIterArchetype_1000(b *testing.B) {
	b.StopTimer()
	comps := []componentType{
		{ID: 0, Type: reflect.TypeOf(testStruct0{})},
	}

	arch := archetype{}
	arch.Init(nil, 32, true, comps...)

	for i := 0; i < 1000; i++ {
		arch.Alloc(newEntity(eid(i)))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		len := uintptr(arch.Len())
		id := ID(0)
		var j uintptr
		for j = 0; j < len; j++ {
			pos := (*testStruct0)(arch.Get(j, id))
			pos.Val++
		}
	}
}

func BenchmarkIterSlice_1000(b *testing.B) {
	b.StopTimer()
	s := []testStruct0{}
	for i := 0; i < 1000; i++ {
		s = append(s, testStruct0{})
	}
	assert.Equal(b, 1000, len(s))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			a := s[j]
			a.Val++
		}
	}
}

func BenchmarkIterSliceInterface_1000(b *testing.B) {
	b.StopTimer()
	s := []interface{}{}
	for i := 0; i < 1000; i++ {
		s = append(s, testStruct0{})
	}
	assert.Equal(b, 1000, len(s))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			a := s[j].(testStruct0)
			a.Val++
		}
	}
}
