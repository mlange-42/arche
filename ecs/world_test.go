package ecs

import (
	"fmt"
	"internal/base"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldEntites(t *testing.T) {
	w := NewWorld()

	assert.Equal(t, newEntityGen(1, 0), w.NewEntity())
	assert.Equal(t, newEntityGen(2, 0), w.NewEntity())
	assert.Equal(t, newEntityGen(3, 0), w.NewEntity())

	assert.Equal(t, 0, int(w.entities[0].index))
	assert.Equal(t, 0, int(w.entities[1].index))
	assert.Equal(t, 1, int(w.entities[2].index))
	assert.Equal(t, 2, int(w.entities[3].index))
	w.RemEntity(newEntityGen(2, 0))
	assert.False(t, w.Alive(newEntityGen(2, 0)))

	assert.Equal(t, 0, int(w.entities[1].index))
	assert.Equal(t, 1, int(w.entities[3].index))

	assert.Equal(t, newEntityGen(2, 1), w.NewEntity())
	assert.False(t, w.Alive(newEntityGen(2, 0)))
	assert.True(t, w.Alive(newEntityGen(2, 1)))

	assert.Equal(t, 2, int(w.entities[2].index))

	w.RemEntity(newEntityGen(3, 0))
	w.RemEntity(newEntityGen(2, 1))
	w.RemEntity(newEntityGen(1, 0))

	assert.Panics(t, func() { w.RemEntity(newEntityGen(3, 0)) })
	assert.Panics(t, func() { w.RemEntity(newEntityGen(2, 1)) })
	assert.Panics(t, func() { w.RemEntity(newEntityGen(1, 0)) })
}

func TestWorldComponents(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	assert.Equal(t, 1, w.archetypes.Len())

	w.Add(e0, posID)
	assert.Equal(t, 2, w.archetypes.Len())
	w.Add(e1, posID, rotID)
	assert.Equal(t, 3, w.archetypes.Len())
	w.Add(e2, posID, rotID)
	assert.Equal(t, 3, w.archetypes.Len())

	w.Remove(e2, posID)

	maskNone := base.NewBitMask()
	maskPos := base.NewBitMask(posID)
	maskRot := base.NewBitMask(rotID)
	maskPosRot := base.NewBitMask(posID, rotID)

	archNone, ok := w.findArchetype(maskNone)
	assert.True(t, ok)
	archPos, ok := w.findArchetype(maskPos)
	assert.True(t, ok)
	archRot, ok := w.findArchetype(maskRot)
	assert.True(t, ok)
	archPosRot, ok := w.findArchetype(maskPosRot)
	assert.True(t, ok)

	assert.Equal(t, 0, int(archNone.Len()))
	assert.Equal(t, 1, int(archPos.Len()))
	assert.Equal(t, 1, int(archRot.Len()))
	assert.Equal(t, 1, int(archPosRot.Len()))

	w.Remove(e1, posID)

	assert.Equal(t, 0, int(archNone.Len()))
	assert.Equal(t, 1, int(archPos.Len()))
	assert.Equal(t, 2, int(archRot.Len()))
	assert.Equal(t, 0, int(archPosRot.Len()))

	w.Add(e0, rotID)
	assert.Equal(t, 0, int(archPos.Len()))
	assert.Equal(t, 1, int(archPosRot.Len()))

	w.Remove(e2, rotID)
	// No-op add/remove
	w.Add(e0)
	w.Remove(e0)

	w.RemEntity(e0)
	assert.Panics(t, func() { w.Has(newEntityGen(1, 0), posID) })
	assert.Panics(t, func() { w.Get(newEntityGen(1, 0), posID) })
}

func TestWorldLabels(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	labID := ComponentID[label](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Add(e0, posID, labID)
	w.Add(e1, labID)
	w.Add(e1, posID)

	lab0 := (*label)(w.Get(e0, labID))
	assert.NotNil(t, lab0)

	lab1 := (*label)(w.Get(e1, labID))
	assert.NotNil(t, lab1)

	assert.True(t, w.Has(e0, labID))
	assert.True(t, w.Has(e1, labID))
	assert.False(t, w.Has(e2, labID))

	assert.Equal(t, lab0, lab1)
}

func TestWorldExchange(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	velID := ComponentID[velocity](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Exchange(e0, []ID{posID}, []ID{})
	w.Exchange(e1, []ID{posID, rotID}, []ID{})
	w.Exchange(e2, []ID{rotID}, []ID{})

	assert.True(t, w.Has(e0, posID))
	assert.False(t, w.Has(e0, rotID))

	assert.True(t, w.Has(e1, posID))
	assert.True(t, w.Has(e1, rotID))

	assert.False(t, w.Has(e2, posID))
	assert.True(t, w.Has(e2, rotID))

	w.Exchange(e2, []ID{posID}, []ID{})
	assert.True(t, w.Has(e2, posID))
	assert.True(t, w.Has(e2, rotID))

	w.Exchange(e0, []ID{rotID}, []ID{posID})
	assert.False(t, w.Has(e0, posID))
	assert.True(t, w.Has(e0, rotID))

	w.Exchange(e1, []ID{velID}, []ID{posID})
	assert.False(t, w.Has(e1, posID))
	assert.True(t, w.Has(e1, rotID))
	assert.True(t, w.Has(e1, velID))

	assert.Panics(t, func() { w.Exchange(e1, []ID{velID}, []ID{}) })
	assert.Panics(t, func() { w.Exchange(e1, []ID{}, []ID{posID}) })
}

func TestWorldAssignSet(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	velID := ComponentID[velocity](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()

	pos := (*position)(w.Assign(e0, posID, &position{2, 3}))
	assert.Equal(t, 2, pos.X)
	pos.X = 5

	pos = (*position)(w.Get(e0, posID))
	assert.Equal(t, 5, pos.X)

	assert.Panics(t, func() { _ = (*position)(w.Assign(e0, posID, &position{2, 3})) })
	assert.Panics(t, func() { _ = (*position)(w.copyTo(e1, posID, &position{2, 3})) })

	e2 := w.NewEntity()
	w.AssignN(e2,
		Component{velID, &velocity{1, 2}},
		Component{rotID, &rotation{3}},
		Component{posID, &position{4, 5}},
	)
	assert.True(t, w.Has(e2, velID))
	assert.True(t, w.Has(e2, rotID))
	assert.True(t, w.Has(e2, posID))

	pos = (*position)(w.Get(e2, posID))
	assert.Equal(t, 4, pos.X)

	_ = (*position)(w.Set(e2, posID, &position{7, 8}))
	pos = (*position)(w.Get(e2, posID))
	assert.Equal(t, 7, pos.X)
}
func TestWorldGetComponents(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Add(e0, posID, rotID)
	w.Add(e1, posID, rotID)
	w.Add(e2, rotID)

	assert.False(t, w.Has(e2, posID))
	assert.True(t, w.Has(e2, rotID))

	pos1 := (*position)(w.Get(e1, posID))
	assert.Equal(t, &position{}, pos1)

	pos1.X = 100
	pos1.Y = 101

	pos0 := (*position)(w.Get(e0, posID))
	pos1 = (*position)(w.Get(e1, posID))
	assert.Equal(t, &position{}, pos0)
	assert.Equal(t, &position{100, 101}, pos1)

	w.RemEntity(e0)

	pos1 = (*position)(w.Get(e1, posID))
	assert.Equal(t, &position{100, 101}, pos1)

	pos2 := (*position)(w.Get(e2, posID))
	assert.True(t, pos2 == nil)
}

func TestWorldIter(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[position](&world)
	rotID := ComponentID[rotation](&world)

	for i := 0; i < 1000; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}

	for i := 0; i < 100; i++ {
		query := world.Query(posID, rotID)
		for query.Next() {
			pos := (*position)(query.Get(posID))
			_ = pos
		}
		assert.Panics(t, func() { query.Next() })
	}

	for i := 0; i < MaskTotalBits-1; i++ {
		query := world.Query(posID, rotID)
		for query.Next() {
			pos := (*position)(query.Get(posID))
			_ = pos
			break
		}
	}
	query := world.Query(posID, rotID)

	assert.Panics(t, func() { world.Query(posID, rotID) })

	query.Close()
	assert.Panics(t, func() { query.Close() })
}

func TestWorldLock(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[position](&world)
	rotID := ComponentID[rotation](&world)

	var entity Entity
	for i := 0; i < 100; i++ {
		entity = world.NewEntity()
		world.Add(entity, posID)
	}

	query1 := world.Query(posID)
	query2 := world.Query(posID)
	assert.True(t, world.IsLocked())
	query1.Close()
	assert.True(t, world.IsLocked())
	query2.Close()
	assert.False(t, world.IsLocked())

	query1 = world.Query(posID)

	assert.Panics(t, func() { world.NewEntity() })
	assert.Panics(t, func() { world.RemEntity(entity) })
	assert.Panics(t, func() { world.Add(entity, rotID) })
	assert.Panics(t, func() { world.Remove(entity, posID) })
}

func TestRegisterComponents(t *testing.T) {
	world := NewWorld()

	ComponentID[position](&world)

	assert.Equal(t, ID(0), ComponentID[position](&world))
	assert.Equal(t, ID(1), ComponentID[rotation](&world))
}

func TestArchetypeGraph(t *testing.T) {
	world := NewWorld()

	posID := ComponentID[position](&world)
	velID := ComponentID[velocity](&world)
	rotID := ComponentID[rotation](&world)

	archEmpty := world.archetypes.Get(0)
	arch0 := world.findOrCreateArchetype(archEmpty, []ID{posID}, []ID{})
	archEmpty2 := world.findOrCreateArchetype(arch0, []ID{}, []ID{posID})

	assert.Equal(t, archEmpty, archEmpty2)

	arch01 := world.findOrCreateArchetype(arch0, []ID{velID}, []ID{})
	arch012 := world.findOrCreateArchetype(arch01, []ID{rotID}, []ID{})

	assert.Equal(t, []ID{0, 1, 2}, arch012.Ids)

	archEmpty3 := world.findOrCreateArchetype(arch012, []ID{}, []ID{posID, rotID, velID})
	assert.Equal(t, archEmpty, archEmpty3)
}

func Test1000Archetypes(t *testing.T) {
	_ = testStruct0{1}
	_ = testStruct1{1}
	_ = testStruct2{1}
	_ = testStruct3{1}
	_ = testStruct4{1}
	_ = testStruct5{1}
	_ = testStruct6{1}
	_ = testStruct7{1}
	_ = testStruct8{1}
	_ = testStruct9{1}

	w := NewWorld()

	ids := [10]ID{}
	ids[0] = ComponentID[testStruct0](&w)
	ids[1] = ComponentID[testStruct1](&w)
	ids[2] = ComponentID[testStruct2](&w)
	ids[3] = ComponentID[testStruct3](&w)
	ids[4] = ComponentID[testStruct4](&w)
	ids[5] = ComponentID[testStruct5](&w)
	ids[6] = ComponentID[testStruct6](&w)
	ids[7] = ComponentID[testStruct7](&w)
	ids[8] = ComponentID[testStruct8](&w)
	ids[9] = ComponentID[testStruct9](&w)

	for i := 0; i < 1024; i++ {
		mask := bitMask(i)
		add := make([]ID, 0, 10)
		for j := 0; j < 10; j++ {
			id := ID(j)
			if mask.Get(id) {
				add = append(add, id)
			}
		}
		entity := w.NewEntity()
		w.Add(entity, add...)
	}
	assert.Equal(t, 1024, w.archetypes.Len())

	cnt := 0
	query := w.Query(0, 7)
	for query.Next() {
		cnt++
	}

	assert.Equal(t, 256, cnt)
}

func TestTypeSizes(t *testing.T) {
	printTypeSize[World]()
	printTypeSizeName[base.PagedArr32[archetype]]("PagedArr32")
	printTypeSize[archetype]()
	printTypeSize[base.Storage]()
	printTypeSize[Query]()
	printTypeSize[archetypeIter]()
	printTypeSizeName[Q1[testStruct0]]("Q1")
	printTypeSizeName[Q8[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4, testStruct5, testStruct6, testStruct7]]("Q8")
}

func printTypeSize[T any]() {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%16s: %5db\n", tp.Name(), tp.Size())
}

func printTypeSizeName[T any](name string) {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%16s: %5db\n", name, tp.Size())
}
