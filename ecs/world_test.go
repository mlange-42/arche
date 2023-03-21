package ecs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldConfig(t *testing.T) {
	_ = NewWorld(NewConfig())

	assert.Panics(t, func() { _ = NewWorld(Config{}) })
	assert.Panics(t, func() { _ = NewWorld(Config{}, Config{}) })
}

func TestWorldEntites(t *testing.T) {
	w := NewWorld()

	assert.Equal(t, newEntityGen(1, 0), w.NewEntity())
	assert.Equal(t, newEntityGen(2, 0), w.NewEntity())
	assert.Equal(t, newEntityGen(3, 0), w.NewEntity())

	assert.Equal(t, 0, int(w.entities[0].index))
	assert.Equal(t, 0, int(w.entities[1].index))
	assert.Equal(t, 1, int(w.entities[2].index))
	assert.Equal(t, 2, int(w.entities[3].index))
	w.RemoveEntity(newEntityGen(2, 0))
	assert.False(t, w.Alive(newEntityGen(2, 0)))

	assert.Equal(t, 0, int(w.entities[1].index))
	assert.Equal(t, 1, int(w.entities[3].index))

	assert.Equal(t, newEntityGen(2, 1), w.NewEntity())
	assert.False(t, w.Alive(newEntityGen(2, 0)))
	assert.True(t, w.Alive(newEntityGen(2, 1)))

	assert.Equal(t, 2, int(w.entities[2].index))

	w.RemoveEntity(newEntityGen(3, 0))
	w.RemoveEntity(newEntityGen(2, 1))
	w.RemoveEntity(newEntityGen(1, 0))

	assert.Panics(t, func() { w.RemoveEntity(newEntityGen(3, 0)) })
	assert.Panics(t, func() { w.RemoveEntity(newEntityGen(2, 1)) })
	assert.Panics(t, func() { w.RemoveEntity(newEntityGen(1, 0)) })
}

func TestWorldNewEntites(t *testing.T) {
	w := NewWorld(NewConfig().WithCapacityIncrement(32))

	posID := ComponentID[position](&w)
	velID := ComponentID[velocity](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity(posID, velID, rotID)
	e2 := w.NewEntityWith(
		Component{posID, &position{1, 2}},
		Component{velID, &velocity{3, 4}},
		Component{rotID, &rotation{5}},
	)
	e3 := w.NewEntityWith()

	assert.Equal(t, All(), w.Mask(e0))
	assert.Equal(t, All(posID, velID, rotID), w.Mask(e1))
	assert.Equal(t, All(posID, velID, rotID), w.Mask(e2))
	assert.Equal(t, All(), w.Mask(e3))

	pos := (*position)(w.Get(e2, posID))
	vel := (*velocity)(w.Get(e2, velID))
	rot := (*rotation)(w.Get(e2, rotID))

	assert.Equal(t, &position{1, 2}, pos)
	assert.Equal(t, &velocity{3, 4}, vel)
	assert.Equal(t, &rotation{5}, rot)

	w.RemoveEntity(e0)
	w.RemoveEntity(e1)
	w.RemoveEntity(e2)
	w.RemoveEntity(e3)

	for i := 0; i < 35; i++ {
		e := w.NewEntityWith(
			Component{posID, &position{i + 1, i + 2}},
			Component{velID, &velocity{i + 3, i + 4}},
			Component{rotID, &rotation{i + 5}},
		)

		pos := (*position)(w.Get(e, posID))
		vel := (*velocity)(w.Get(e, velID))
		rot := (*rotation)(w.Get(e, rotID))

		assert.Equal(t, &position{i + 1, i + 2}, pos)
		assert.Equal(t, &velocity{i + 3, i + 4}, vel)
		assert.Equal(t, &rotation{i + 5}, rot)
	}
}

func TestWorldComponents(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)

	tPosID := TypeID(&w, reflect.TypeOf(position{}))
	tRotID := TypeID(&w, reflect.TypeOf(rotation{}))

	assert.Equal(t, posID, tPosID)
	assert.Equal(t, rotID, tRotID)

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

	assert.Equal(t, All(posID), w.Mask(e0))
	assert.Equal(t, All(posID, rotID), w.Mask(e1))

	w.Remove(e2, posID)

	maskNone := All()
	maskPos := All(posID)
	maskRot := All(rotID)
	maskPosRot := All(posID, rotID)

	archNone, ok := w.findArchetypeSlow(maskNone)
	assert.True(t, ok)
	archPos, ok := w.findArchetypeSlow(maskPos)
	assert.True(t, ok)
	archRot, ok := w.findArchetypeSlow(maskRot)
	assert.True(t, ok)
	archPosRot, ok := w.findArchetypeSlow(maskPosRot)
	assert.True(t, ok)

	assert.Equal(t, 0, int(archNone.archetype.Len()))
	assert.Equal(t, 1, int(archPos.archetype.Len()))
	assert.Equal(t, 1, int(archRot.archetype.Len()))
	assert.Equal(t, 1, int(archPosRot.archetype.Len()))

	w.Remove(e1, posID)

	assert.Equal(t, 0, int(archNone.archetype.Len()))
	assert.Equal(t, 1, int(archPos.archetype.Len()))
	assert.Equal(t, 2, int(archRot.archetype.Len()))
	assert.Equal(t, 0, int(archPosRot.archetype.Len()))

	w.Add(e0, rotID)
	assert.Equal(t, 0, int(archPos.archetype.Len()))
	assert.Equal(t, 1, int(archPosRot.archetype.Len()))

	w.Remove(e2, rotID)
	// No-op add/remove
	w.Add(e0)
	w.Remove(e0)

	w.RemoveEntity(e0)
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

	assert.Panics(t, func() { w.Assign(e0) })

	w.Assign(e0, Component{posID, &position{2, 3}})
	pos := (*position)(w.Get(e0, posID))
	assert.Equal(t, 2, pos.X)
	pos.X = 5

	pos = (*position)(w.Get(e0, posID))
	assert.Equal(t, 5, pos.X)

	assert.Panics(t, func() { w.Assign(e0, Component{posID, &position{2, 3}}) })
	assert.Panics(t, func() { _ = (*position)(w.copyTo(e1, posID, &position{2, 3})) })

	e2 := w.NewEntity()
	w.Assign(e2,
		Component{posID, &position{4, 5}},
		Component{velID, &velocity{1, 2}},
		Component{rotID, &rotation{3}},
	)
	assert.True(t, w.Has(e2, velID))
	assert.True(t, w.Has(e2, rotID))
	assert.True(t, w.Has(e2, posID))

	pos = (*position)(w.Get(e2, posID))
	rot := (*rotation)(w.Get(e2, rotID))
	vel := (*velocity)(w.Get(e2, velID))
	assert.Equal(t, &position{4, 5}, pos)
	assert.Equal(t, &rotation{3}, rot)
	assert.Equal(t, &velocity{1, 2}, vel)

	_ = (*position)(w.Set(e2, posID, &position{7, 8}))
	pos = (*position)(w.Get(e2, posID))
	assert.Equal(t, 7, pos.X)

	*pos = position{8, 9}
	pos = (*position)(w.Get(e2, posID))
	assert.Equal(t, 8, pos.X)
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

	w.RemoveEntity(e0)

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
		query := world.Query(All(posID, rotID))
		for query.Next() {
			pos := (*position)(query.Get(posID))
			_ = pos
		}
		assert.Panics(t, func() { query.Next() })
	}

	for i := 0; i < MaskTotalBits-1; i++ {
		query := world.Query(All(posID, rotID))
		for query.Next() {
			pos := (*position)(query.Get(posID))
			_ = pos
			break
		}
	}
	query := world.Query(All(posID, rotID))

	assert.Panics(t, func() { world.Query(All(posID, rotID)) })

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

	query1 := world.Query(All(posID))
	query2 := world.Query(All(posID))
	assert.True(t, world.IsLocked())
	query1.Close()
	assert.True(t, world.IsLocked())
	query2.Close()
	assert.False(t, world.IsLocked())

	query1 = world.Query(All(posID))

	assert.Panics(t, func() { world.NewEntity() })
	assert.Panics(t, func() { world.RemoveEntity(entity) })
	assert.Panics(t, func() { world.Add(entity, rotID) })
	assert.Panics(t, func() { world.Remove(entity, posID) })
}

func TestWorldStats(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Add(e0, posID)
	w.Add(e1, posID, rotID)
	w.Add(e2, posID, rotID)

	stats := w.Stats()
	fmt.Println(stats.String())
}

func TestWorldResources(t *testing.T) {
	w := NewWorld()

	posID := ResourceID[position](&w)
	rotID := ResourceID[rotation](&w)

	assert.False(t, w.HasResource(posID))
	assert.Nil(t, w.GetResource(posID))

	AddResource(&w, &position{1, 2})

	fmt.Println(w.resources.registry.Components)

	assert.True(t, w.HasResource(posID))
	pos, ok := w.GetResource(posID).(*position)

	assert.True(t, ok)
	assert.Equal(t, position{1, 2}, *pos)

	assert.Panics(t, func() { w.AddResource(posID, &position{1, 2}) })

	pos = GetResource[position](&w)
	assert.Equal(t, position{1, 2}, *pos)

	w.AddResource(rotID, &rotation{5})
	assert.True(t, w.HasResource(rotID))
	w.RemoveResource(rotID)
	assert.False(t, w.HasResource(rotID))
	assert.Panics(t, func() { w.RemoveResource(rotID) })
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
	arch0 := world.findOrCreateArchetype(archEmpty, []ID{posID, velID}, []ID{})
	archEmpty2 := world.findOrCreateArchetype(arch0, []ID{}, []ID{velID, posID})
	assert.Equal(t, archEmpty, archEmpty2)
	assert.Equal(t, 2, world.archetypes.Len())
	assert.Equal(t, 3, world.graph.Len())

	archEmpty3 := world.findOrCreateArchetype(arch0, []ID{}, []ID{posID, velID})
	assert.Equal(t, archEmpty, archEmpty3)
	assert.Equal(t, 2, world.archetypes.Len())
	assert.Equal(t, 4, world.graph.Len())

	arch01 := world.findOrCreateArchetype(arch0, []ID{velID}, []ID{})
	arch012 := world.findOrCreateArchetype(arch01, []ID{rotID}, []ID{})

	assert.Equal(t, []ID{0, 1, 2}, arch012.Ids)

	archEmpty4 := world.findOrCreateArchetype(arch012, []ID{}, []ID{posID, rotID, velID})
	assert.Equal(t, archEmpty, archEmpty4)
}

func TestWorldListener(t *testing.T) {
	events := []EntityEvent{}
	listen := func(e EntityEvent) {
		events = append(events, e)
	}

	w := NewWorld()

	w.SetListener(listen)

	posID := ComponentID[position](&w)
	velID := ComponentID[velocity](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	assert.Equal(t, 1, len(events))
	assert.Equal(t, EntityEvent{
		Entity: e0, AddedRemoved: 1,
	}, events[len(events)-1])

	w.RemoveEntity(e0)
	assert.Equal(t, 2, len(events))
	assert.Equal(t, EntityEvent{
		Entity: e0, AddedRemoved: -1,
	}, events[len(events)-1])

	e0 = w.NewEntity(posID, velID)
	assert.Equal(t, 3, len(events))
	assert.Equal(t, EntityEvent{
		Entity:       e0,
		NewMask:      All(posID, velID),
		Added:        []ID{posID, velID},
		Current:      []ID{posID, velID},
		AddedRemoved: 1,
	}, events[len(events)-1])

	w.RemoveEntity(e0)
	assert.Equal(t, 4, len(events))
	assert.Equal(t, EntityEvent{
		Entity:       e0,
		OldMask:      All(posID, velID),
		NewMask:      Mask{},
		Removed:      []ID{posID, velID},
		Current:      nil,
		AddedRemoved: -1,
	}, events[len(events)-1])

	e0 = w.NewEntityWith(Component{posID, &position{}}, Component{velID, &velocity{}})
	assert.Equal(t, 5, len(events))
	assert.Equal(t, EntityEvent{
		Entity:       e0,
		NewMask:      All(posID, velID),
		Added:        []ID{posID, velID},
		Current:      []ID{posID, velID},
		AddedRemoved: 1,
	}, events[len(events)-1])

	w.Add(e0, rotID)
	assert.Equal(t, 6, len(events))
	assert.Equal(t, EntityEvent{
		Entity:       e0,
		OldMask:      All(posID, velID),
		NewMask:      All(posID, velID, rotID),
		Added:        []ID{rotID},
		Current:      []ID{posID, velID, rotID},
		AddedRemoved: 0,
	}, events[len(events)-1])

	w.Remove(e0, posID)
	assert.Equal(t, 7, len(events))
	assert.Equal(t, EntityEvent{
		Entity:       e0,
		OldMask:      All(posID, velID, rotID),
		NewMask:      All(velID, rotID),
		Removed:      []ID{posID},
		Current:      []ID{velID, rotID},
		AddedRemoved: 0,
	}, events[len(events)-1])

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
		mask := Mask{uint64(i), 0}
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
	query := w.Query(All(0, 7))
	for query.Next() {
		cnt++
	}

	assert.Equal(t, 256, cnt)
}

func TestTypeSizes(t *testing.T) {
	printTypeSize[Entity]()
	printTypeSize[entityIndex]()
	printTypeSize[Mask]()
	printTypeSize[World]()
	printTypeSizeName[pagedArr32[archetype]]("pagedArr32")
	printTypeSize[archetype]()
	printTypeSize[archetypeAccess]()
	printTypeSize[archetypeNode]()
	printTypeSize[entityPool]()
	printTypeSizeName[componentRegistry[ID]]("componentRegistry")
	printTypeSize[bitPool]()
	printTypeSize[Query]()
	printTypeSize[archetypeIter]()
	printTypeSize[resources]()
	printTypeSizeName[reflect.Value]("reflect.Value")
}

func printTypeSize[T any]() {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%18s: %5d B\n", tp.Name(), tp.Size())
}

func printTypeSizeName[T any](name string) {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%18s: %5d B\n", name, tp.Size())
}

func BenchmarkGetResource(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	AddResource(&w, &position{1, 2})
	posID := ResourceID[position](&w)

	b.StartTimer()

	var res *position
	for i := 0; i < b.N; i++ {
		res = w.GetResource(posID).(*position)
	}

	_ = res
}

func BenchmarkGetResourceShortcut(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	AddResource(&w, &position{1, 2})

	b.StartTimer()

	var res *position
	for i := 0; i < b.N; i++ {
		res = GetResource[position](&w)
	}

	_ = res
}
