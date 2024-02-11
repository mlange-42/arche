package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterCache(t *testing.T) {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)
	rotID := ComponentID[rotation](&world)

	cache := world.Cache()

	world.NewEntity()
	world.NewEntity(posID, velID)
	world.NewEntity(posID, velID, rotID)

	all1 := All(posID, velID)
	all2 := All(posID, velID, rotID)

	f1 := cache.Register(all1)
	f2 := cache.Register(all2)
	assert.Equal(t, 0, int(f1.id))
	assert.Equal(t, 1, int(f2.id))

	assert.Equal(t, 2, len(world.getArchetypes(&f1)))
	assert.Equal(t, 1, len(world.getArchetypes(&f2)))

	assert.PanicsWithValue(t, "filter is already registered", func() { cache.Register(&f2) })

	e1 := cache.get(&f1)
	e2 := cache.get(&f2)

	assert.Equal(t, f1.filter, e1.Filter)
	assert.Equal(t, f2.filter, e2.Filter)

	ff1 := cache.Unregister(&f1)
	ff2 := cache.Unregister(&f2)

	assert.Equal(t, all1, ff1)
	assert.Equal(t, all2, ff2)

	assert.PanicsWithValue(t, "no filter for id found to unregister", func() { cache.Unregister(&f1) })
	assert.PanicsWithValue(t, "no filter for id found", func() { cache.get(&f1) })
}

func TestFilterCacheRelation(t *testing.T) {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	rel1ID := ComponentID[testRelationA](&world)
	rel2ID := ComponentID[testRelationB](&world)

	target1 := world.NewEntity()
	target2 := world.NewEntity()
	target3 := world.NewEntity()
	target4 := world.NewEntity()

	cache := world.Cache()

	f1 := All(rel1ID)
	ff1 := cache.Register(f1)

	f2 := NewRelationFilter(f1, target1)
	ff2 := cache.Register(&f2)

	f3 := NewRelationFilter(f1, target2)
	ff3 := cache.Register(&f3)

	c1 := world.Cache().get(&ff1)
	c2 := world.Cache().get(&ff2)
	c3 := world.Cache().get(&ff3)

	NewBuilder(&world, posID).NewBatch(10)

	assert.Equal(t, int32(0), c1.Archetypes.Len())
	assert.Equal(t, int32(0), c2.Archetypes.Len())
	assert.Equal(t, int32(0), c3.Archetypes.Len())

	e1 := NewBuilder(&world, rel1ID).WithRelation(rel1ID).New(target1)
	assert.Equal(t, int32(1), c1.Archetypes.Len())
	assert.Equal(t, int32(1), c2.Archetypes.Len())

	_ = NewBuilder(&world, rel1ID).WithRelation(rel1ID).New(target3)
	assert.Equal(t, int32(2), c1.Archetypes.Len())
	assert.Equal(t, int32(1), c2.Archetypes.Len())

	_ = NewBuilder(&world, rel2ID).WithRelation(rel2ID).New(target2)

	world.RemoveEntity(e1)
	world.RemoveEntity(target1)
	assert.Equal(t, int32(1), c1.Archetypes.Len())
	assert.Equal(t, int32(0), c2.Archetypes.Len())

	_ = NewBuilder(&world, rel1ID).WithRelation(rel1ID).New(target2)
	_ = NewBuilder(&world, rel1ID, posID).WithRelation(rel1ID).New(target2)
	_ = NewBuilder(&world, rel1ID, posID).WithRelation(rel1ID).New(target3)
	_ = NewBuilder(&world, rel1ID, posID).WithRelation(rel1ID).New(target4)
	assert.Equal(t, int32(5), c1.Archetypes.Len())
	assert.Equal(t, int32(2), c3.Archetypes.Len())

	world.Batch().RemoveEntities(All())
	assert.Equal(t, int32(0), c1.Archetypes.Len())
	assert.Equal(t, int32(0), c2.Archetypes.Len())
}

func ExampleCache() {
	world := NewWorld()
	posID := ComponentID[Position](&world)

	filter := All(posID)
	cached := world.Cache().Register(filter)
	query := world.Query(&cached)

	for query.Next() {
		// ...
	}
	// Output:
}
