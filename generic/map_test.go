package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestGenericMap(t *testing.T) {
	w := ecs.NewWorld()
	get := NewMap[testStruct0](&w)
	mut := NewExchange(&w).Adds(T[testStruct0]())

	assert.Equal(t, ecs.ComponentID[testStruct0](&w), get.ID())

	e0 := w.NewEntity()

	mut.Add(e0)
	has := get.Has(e0)
	_ = get.Get(e0)
	assert.True(t, has)

	has = get.HasUnchecked(e0)
	_ = get.GetUnchecked(e0)
	assert.True(t, has)

	_ = get.Set(e0, &testStruct0{100})
	str := get.Get(e0)

	assert.Equal(t, 100, int(str.val))

	get2 := NewMap[testStruct1](&w)
	assert.Equal(t, ecs.ComponentID[testStruct1](&w), get2.ID())
	assert.Panics(t, func() { get2.Set(e0, &testStruct1{}) })

	w.RemoveEntity(e0)
	_ = w.NewEntity()

	assert.Panics(t, func() { get.Has(e0) })
	assert.Panics(t, func() { get.Get(e0) })
}

func TestGenericMapRelations(t *testing.T) {
	w := ecs.NewWorld()
	get := NewMap[testRelationA](&w)
	genTarg := NewMap1[Position](&w)
	gen := NewMap2[testRelationA, Position](&w)

	targ := genTarg.NewEntity()
	e0 := gen.NewEntity()

	targ2 := get.GetRelation(e0)
	assert.Equal(t, ecs.Entity{}, targ2)

	get.SetRelation(e0, targ)

	targ2 = get.GetRelation(e0)
	assert.Equal(t, targ, targ2)
}

func ExampleMap() {
	// Create a world.
	world := ecs.NewWorld()

	// Spawn some entities using the generic API.
	spawner := NewMap2[Position, Velocity](&world)
	entity := spawner.NewEntity()

	// Create a new Map.
	mapper := NewMap[Position](&world)

	// Get the map's component.
	pos := mapper.Get(entity)
	pos.X, pos.Y = 10, 5

	// Get the map's component, optimized for cases when the entity is guaranteed to be still alive.
	pos = mapper.GetUnchecked(entity)
	pos.X, pos.Y = 10, 5
	// Output:
}
