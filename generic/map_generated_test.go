package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestMap1Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap1[testStruct0](&w)
	map0 := NewMap[testStruct0](&w)

	e := mut.NewEntity()
	s0 := mut.Get(e)
	assert.NotNil(t, s0)

	mut.Remove(e)
	assert.False(t, map0.Has(e))

	mut.Add(e)
	s0 = mut.Get(e)
	assert.NotNil(t, s0)

	e = w.NewEntity()
	mut.Assign(e, &testStruct0{})
	assert.True(t, map0.Has(e))

	e = mut.NewEntityWith(&testStruct0{})
	assert.True(t, map0.Has(e))

	s0 = mut.Get(e)
	assert.NotNil(t, s0)

	q := mut.NewEntities(2)
	q.Close()
	q = mut.NewEntitiesWith(2, &testStruct0{})
	q.Close()
	mut.RemoveEntities(true)
	mut.RemoveEntities(false)
}

func TestMap2Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap2[testStruct0, testStruct1](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.NewEntity()
	s0, s1 := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1 = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e, &testStruct0{}, &testStruct1{})
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewEntityWith(&testStruct0{}, &testStruct1{})
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1 = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	q := mut.NewEntities(2)
	q.Close()
	q = mut.NewEntitiesWith(2,
		&testStruct0{}, &testStruct1{},
	)
	q.Close()
	mut.RemoveEntities(true)
	mut.RemoveEntities(false)
}

func TestMap3Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap3[
		testStruct0, testStruct1, testStruct2,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.NewEntity()
	s0, s1, _ := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	q := mut.NewEntities(2)
	q.Close()
	q = mut.NewEntitiesWith(2,
		&testStruct0{}, &testStruct1{}, &testStruct2{},
	)
	q.Close()
	mut.RemoveEntities(true)
	mut.RemoveEntities(false)
}

func TestMap4Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap4[
		testStruct0, testStruct1, testStruct2, testStruct3,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.NewEntity()
	s0, s1, _, _ := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	q := mut.NewEntities(2)
	q.Close()
	q = mut.NewEntitiesWith(2,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
	)
	q.Close()
	mut.RemoveEntities(true)
	mut.RemoveEntities(false)
}

func TestMap5Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap5[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.NewEntity()
	s0, s1, _, _, _ := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	q := mut.NewEntities(2)
	q.Close()
	q = mut.NewEntitiesWith(2,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{},
	)
	q.Close()
	mut.RemoveEntities(true)
	mut.RemoveEntities(false)
}

func TestMap6Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap6[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.NewEntity()
	s0, s1, _, _, _, _ := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	q := mut.NewEntities(2)
	q.Close()
	q = mut.NewEntitiesWith(2,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{},
	)
	q.Close()
	mut.RemoveEntities(true)
	mut.RemoveEntities(false)
}

func TestMap7Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap7[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.NewEntity()
	s0, s1, _, _, _, _, _ := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	q := mut.NewEntities(2)
	q.Close()
	q = mut.NewEntitiesWith(2,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{},
	)
	q.Close()
	mut.RemoveEntities(true)
	mut.RemoveEntities(false)
}

func TestMap8Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.NewEntity()
	s0, s1, _, _, _, _, _, _ := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1, _, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	q := mut.NewEntities(2)
	q.Close()
	q = mut.NewEntitiesWith(2,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
	)
	q.Close()
	mut.RemoveEntities(true)
	mut.RemoveEntities(false)
}
