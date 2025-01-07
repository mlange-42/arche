package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArchetype(t *testing.T) {
	comps := []componentType{
		{ID: id(0), Type: reflect.TypeOf(Position{})},
		{ID: id(1), Type: reflect.TypeOf(rotation{})},
	}

	node := newArchNode(All(id(0), id(1)), &nodeData{}, ID{}, false, 32, comps)
	arch := archetype{}
	data := archetypeData{}
	arch.Init(&node, &data, 0, false, 16, Entity{})

	idx := arch.Alloc(newEntity(0))
	arch.Set(idx, id(0), &Position{1, 2})
	arch.Set(idx, id(1), &rotation{3})

	idx = arch.Alloc(newEntity(1))
	arch.Set(idx, id(0), &Position{4, 5})
	arch.Set(idx, id(1), &rotation{6})

	assert.Equal(t, 2, int(arch.Len()))

	e0 := arch.GetEntity(0)
	e1 := arch.GetEntity(1)
	assert.Equal(t, Entity{0, 0}, e0)
	assert.Equal(t, Entity{1, 0}, e1)

	pos0 := (*Position)(arch.Get(0, id(0)))
	rot0 := (*rotation)(arch.Get(0, id(1)))
	pos1 := (*Position)(arch.Get(1, id(0)))
	rot1 := (*rotation)(arch.Get(1, id(1)))

	assert.Equal(t, 1, pos0.X)
	assert.Equal(t, 2, pos0.Y)
	assert.Equal(t, 3, rot0.Angle)
	assert.Equal(t, 4, pos1.X)
	assert.Equal(t, 5, pos1.Y)
	assert.Equal(t, 6, rot1.Angle)

	arch.Remove(0)
	assert.Equal(t, 1, int(arch.Len()))

	pos0 = (*Position)(arch.Get(0, id(0)))
	rot0 = (*rotation)(arch.Get(0, id(1)))
	assert.Equal(t, 4, pos0.X)
	assert.Equal(t, 5, pos0.Y)
	assert.Equal(t, 6, rot0.Angle)
}

func TestNewArchetype(t *testing.T) {
	comps := []componentType{
		{ID: id(0), Type: reflect.TypeOf(Position{})},
		{ID: id(1), Type: reflect.TypeOf(rotation{})},
	}

	node := newArchNode(All(id(0), id(1)), &nodeData{}, ID{}, false, 32, comps)
	arch := archetype{}
	data := archetypeData{}
	arch.Init(&node, &data, 0, true, 16, Entity{})
	assert.Equal(t, 32, int(arch.Cap()))

	arch = archetype{}
	data = archetypeData{}
	arch.Init(&node, &data, 0, false, 16, Entity{})
	assert.Equal(t, 1, int(arch.Cap()))

	comps = []componentType{
		{ID: id(1), Type: reflect.TypeOf(rotation{})},
		{ID: id(0), Type: reflect.TypeOf(Position{})},
	}
	assert.PanicsWithValue(t, "component arguments must be sorted by ID", func() {
		node := newArchNode(All(id(0), id(1)), &nodeData{}, ID{}, false, 32, comps)
		arch := archetype{}
		data := archetypeData{}
		arch.Init(&node, &data, 0, true, 16, Entity{})
	})
}

func TestArchetypeExtend(t *testing.T) {
	comps := []componentType{
		{ID: id(0), Type: reflect.TypeOf(Position{})},
		{ID: id(1), Type: reflect.TypeOf(rotation{})},
	}

	node := newArchNode(All(id(0), id(1)), &nodeData{}, ID{}, false, 8, comps)
	arch := archetype{}
	data := archetypeData{}
	arch.Init(&node, &data, 0, true, 16, Entity{})

	assert.Equal(t, 8, int(arch.Cap()))
	assert.Equal(t, 0, int(arch.Len()))

	arch.extend(5)
	assert.Equal(t, 8, int(arch.Cap()))

	arch.extend(8)
	assert.Equal(t, 8, int(arch.Cap()))

	arch.extend(17)
	assert.Equal(t, 32, int(arch.Cap()))
}

func TestArchetypeExtendLayouts(t *testing.T) {
	comps := []componentType{
		{ID: id(0), Type: reflect.TypeOf(Position{})},
		{ID: id(1), Type: reflect.TypeOf(rotation{})},
		{ID: id(2), Type: reflect.TypeOf(relationComp{})},
	}
	entity := newEntity(1)

	node := newArchNode(All(id(0), id(1)), &nodeData{}, ID{}, false, 8, comps[:2])
	arch := archetype{}
	data := archetypeData{}
	arch.Init(&node, &data, 0, true, 16, Entity{})
	node.SetArchetype(&arch)

	assert.Equal(t, len(arch.layouts), 16)
	node.ExtendArchetypeLayouts(16)
	assert.Equal(t, len(arch.layouts), 16)
	node.ExtendArchetypeLayouts(32)
	assert.Equal(t, len(arch.layouts), 32)

	node = newArchNode(All(id(0), id(1)), &nodeData{}, id(2), true, 8, comps)
	node.CreateArchetype(16, entity)
	arch2, ok := node.GetArchetype(entity)

	assert.True(t, ok)
	assert.Equal(t, len(arch2.layouts), 16)
	node.ExtendArchetypeLayouts(16)
	assert.Equal(t, len(arch2.layouts), 16)
	node.ExtendArchetypeLayouts(32)
	assert.Equal(t, len(arch2.layouts), 32)
}

func TestArchetypeAlloc(t *testing.T) {
	comps := []componentType{
		{ID: id(0), Type: reflect.TypeOf(Position{})},
		{ID: id(1), Type: reflect.TypeOf(rotation{})},
	}
	node := newArchNode(All(id(0), id(1)), &nodeData{}, ID{}, false, 8, comps)
	arch := archetype{}
	data := archetypeData{}
	arch.Init(&node, &data, 0, true, 16, Entity{})

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
	comps := []componentType{
		{ID: id(0), Type: reflect.TypeOf(testStruct0{})},
		{ID: id(1), Type: reflect.TypeOf(label{})},
	}

	node := newArchNode(All(id(0), id(1)), &nodeData{}, ID{}, false, 1, comps)
	a := archetype{}
	data := archetypeData{}
	a.Init(&node, &data, 0, true, 16, Entity{})

	assert.Equal(t, 1, int(a.Cap()))
	assert.Equal(t, 0, int(a.Len()))

	idx := a.Alloc(Entity{1, 0})
	a.Set(idx, id(0), &testStruct0{100})
	a.Set(idx, id(1), &label{})

	idx = a.Alloc(Entity{2, 0})
	a.Set(idx, id(0), &testStruct0{200})
	a.Set(idx, id(1), &label{})

	ts := (*testStruct0)(a.Get(0, id(0)))
	assert.Equal(t, 100, int(ts.Val))

	a.Set(1, id(0), &testStruct0{200})
	a.Set(1, id(1), &label{})

	_ = (*testStruct0)(a.Get(1, id(0)))
	_ = (*label)(a.Get(1, id(1)))

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
		{ID: id(0), Type: reflect.TypeOf(Position{})},
		{ID: id(1), Type: reflect.TypeOf(rotation{})},
	}

	node := newArchNode(All(id(0), id(1)), &nodeData{}, ID{}, false, 32, comps)
	arch := archetype{}
	data := archetypeData{}
	arch.Init(&node, &data, 0, false, 16, Entity{})

	idx := arch.Alloc(newEntity(0))
	arch.Set(idx, id(0), &Position{1, 2})
	arch.Set(idx, id(1), &rotation{3})

	idx = arch.Alloc(newEntity(1))
	arch.Set(idx, id(0), &Position{4, 5})
	arch.Set(idx, id(1), &rotation{6})

	assert.Equal(t, Position{1, 2}, *(*Position)(arch.Get(0, id(0))))
	assert.Equal(t, Position{4, 5}, *(*Position)(arch.Get(1, id(0))))
	assert.Equal(t, 2, int(arch.Len()))

	arch.Reset()
	assert.Equal(t, 0, int(arch.Len()))

	idx = arch.Alloc(newEntity(0))
	arch.Set(idx, id(0), &Position{10, 20})
	arch.Set(idx, id(1), &rotation{3})

	idx = arch.Alloc(newEntity(1))
	arch.Set(idx, id(0), &Position{40, 50})
	arch.Set(idx, id(1), &rotation{6})

	assert.Equal(t, Position{10, 20}, *(*Position)(arch.Get(0, id(0))))
	assert.Equal(t, Position{40, 50}, *(*Position)(arch.Get(1, id(0))))
	assert.Equal(t, 2, int(arch.Len()))
}

func TestArchetypeZero(t *testing.T) {
	comps := []componentType{
		{ID: id(0), Type: reflect.TypeOf(Position{})},
		{ID: id(1), Type: reflect.TypeOf(PointerComp{})},
		{ID: id(2), Type: reflect.TypeOf(label{})},
	}

	node := newArchNode(All(id(0), id(1), id(2)), &nodeData{}, ID{}, false, 32, comps)
	arch := archetype{}
	data := archetypeData{}
	arch.Init(&node, &data, 0, false, 16, Entity{})

	arch.Alloc(newEntity(0))
	arch.Alloc(newEntity(1))

	assert.Equal(t, Position{0, 0}, *(*Position)(arch.Get(0, id(0))))
	assert.Equal(t, Position{0, 0}, *(*Position)(arch.Get(1, id(0))))

	assert.Equal(t, PointerComp{nil, 0}, *(*PointerComp)(arch.Get(0, id(1))))
	assert.Equal(t, PointerComp{nil, 0}, *(*PointerComp)(arch.Get(1, id(1))))

	pos := (*Position)(arch.Get(0, id(0)))
	pos.X = 100
	pos = (*Position)(arch.Get(1, id(0)))
	pos.X = 100

	ptr := (*PointerComp)(arch.Get(0, id(1)))
	ptr.Value = 99
	ptr = (*PointerComp)(arch.Get(1, id(1)))
	ptr.Value = 99

	assert.Equal(t, Position{100, 0}, *(*Position)(arch.Get(0, id(0))))
	assert.Equal(t, Position{100, 0}, *(*Position)(arch.Get(1, id(0))))

	assert.Equal(t, PointerComp{nil, 99}, *(*PointerComp)(arch.Get(0, id(1))))
	assert.Equal(t, PointerComp{nil, 99}, *(*PointerComp)(arch.Get(1, id(1))))

	arch.Remove(0)
	arch.Remove(0)
	arch.Alloc(newEntity(0))
	arch.Alloc(newEntity(1))

	assert.Equal(t, Position{0, 0}, *(*Position)(arch.Get(0, id(0))))
	assert.Equal(t, Position{0, 0}, *(*Position)(arch.Get(1, id(0))))

	assert.Equal(t, PointerComp{nil, 0}, *(*PointerComp)(arch.Get(0, id(1))))
	assert.Equal(t, PointerComp{nil, 0}, *(*PointerComp)(arch.Get(1, id(1))))
}

func BenchmarkIterArchetype_1000(b *testing.B) {
	b.StopTimer()
	comps := []componentType{
		{ID: id(0), Type: reflect.TypeOf(testStruct0{})},
	}

	node := newArchNode(All(id(0)), &nodeData{}, ID{}, false, 32, comps)
	arch := archetype{}
	data := archetypeData{}
	arch.Init(&node, &data, 0, true, 16, Entity{})

	for i := 0; i < 1000; i++ {
		arch.Alloc(newEntity(eid(i)))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		len := arch.Len()
		id := id(0)
		var j uint32
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
