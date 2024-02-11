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
	assert.PanicsWithValue(t, "can't copy component into entity that has no such component type", func() { get2.Set(e0, &testStruct1{}) })

	w.RemoveEntity(e0)
	_ = w.NewEntity()

	assert.PanicsWithValue(t, "can't check for component of a dead entity", func() { get.Has(e0) })
	assert.PanicsWithValue(t, "can't get component of a dead entity", func() { get.Get(e0) })
}

func TestGenericMapRelations(t *testing.T) {
	w := ecs.NewWorld()
	get := NewMap[testRelationA](&w)
	genTarg := NewMap1[Position](&w)
	gen := NewMap2[testRelationA, Position](&w)

	targ := genTarg.New()
	e0 := gen.New()

	targ2 := get.GetRelation(e0)
	assert.Equal(t, ecs.Entity{}, targ2)

	get.SetRelation(e0, targ)

	targ2 = get.GetRelation(e0)
	assert.Equal(t, targ, targ2)

	assert.Equal(t, targ, get.GetRelationUnchecked(e0))
}

func TestGenericMapBatchRelations(t *testing.T) {
	w := ecs.NewWorld()
	get := NewMap[testRelationA](&w)
	genTarg := NewMap1[Position](&w)
	gen := NewMap2[testRelationA, Position](&w)

	targ1 := genTarg.New()
	targ2 := genTarg.New()
	gen.NewBatch(10)

	filter := NewFilter2[testRelationA, Position]().Filter(&w)
	get.SetRelationBatch(filter, targ1)

	filter = NewFilter2[testRelationA, Position]().Filter(&w, targ1)
	query := get.SetRelationBatchQ(filter, targ2)
	assert.Equal(t, 10, query.Count())

	for query.Next() {
		assert.Equal(t, targ2, query.Relation())
	}
}

func ExampleMap() {
	// Create a world.
	world := ecs.NewWorld()

	// Spawn some entities using the generic API.
	spawner := NewMap2[Position, Velocity](&world)
	entity := spawner.New()

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
