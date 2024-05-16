package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompResIDs(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	rotID := ComponentID[rotation](&w)

	tPosID := TypeID(&w, reflect.TypeOf(Position{}))
	tRotID := TypeID(&w, reflect.TypeOf(rotation{}))

	res1ID := ResourceID[Position](&w)
	res2ID := ResourceID[Velocity](&w)

	tRes1ID := ResourceTypeID(&w, reflect.TypeOf(Position{}))
	tRes2ID := ResourceTypeID(&w, reflect.TypeOf(Velocity{}))

	assert.Equal(t, posID, tPosID)
	assert.Equal(t, rotID, tRotID)
	assert.Equal(t, res1ID, tRes1ID)
	assert.Equal(t, res2ID, tRes2ID)

	assert.Equal(t, uint8(0), posID.id)
	assert.Equal(t, uint8(1), rotID.id)

	assert.Equal(t, uint8(0), res1ID.id)
	assert.Equal(t, uint8(1), res2ID.id)

	assert.Equal(t, []ID{id(0), id(1)}, ComponentIDs(&w))
	assert.Equal(t, []ResID{{id: 0}, {id: 1}}, ResourceIDs(&w))
}

func TestRegisterComponents(t *testing.T) {
	world := NewWorld()

	ComponentID[Position](&world)

	assert.Equal(t, id(0), ComponentID[Position](&world))
	assert.Equal(t, id(1), ComponentID[rotation](&world))
}

func TestComponentInfo(t *testing.T) {
	w := NewWorld()
	_ = ComponentID[Velocity](&w)
	posID := ComponentID[Position](&w)

	info, ok := ComponentInfo(&w, posID)
	assert.True(t, ok)
	assert.Equal(t, info.Type, reflect.TypeOf(Position{}))

	info, ok = ComponentInfo(&w, ID{id: 3})
	assert.False(t, ok)
	assert.Equal(t, info, CompInfo{})

	resID := ResourceID[Velocity](&w)

	tp, ok := ResourceType(&w, resID)
	assert.True(t, ok)
	assert.Equal(t, tp, reflect.TypeOf(Velocity{}))

	tp, ok = ResourceType(&w, ResID{id: 3})
	assert.False(t, ok)
	assert.Equal(t, tp, nil)
}

func BenchmarkComponentID(b *testing.B) {
	b.StopTimer()
	world := NewWorld()
	id := ComponentID[Position](&world)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		id = ComponentID[Position](&world)
	}
	_ = id
}

func BenchmarkTypeID(b *testing.B) {
	b.StopTimer()
	world := NewWorld()
	id := ComponentID[Position](&world)
	info, _ := ComponentInfo(&world, id)
	tp := info.Type

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		id = TypeID(&world, tp)
	}
	_ = id
}
