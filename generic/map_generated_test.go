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

	e, s0 := mut.NewEntity()
	assert.NotNil(t, s0)

	mut.Remove(e)
	assert.False(t, map0.Has(e))

	s0 = mut.Add(e)
	assert.NotNil(t, s0)

	e = w.NewEntity()
	_ = mut.Assign(e, &testStruct0{})
	assert.True(t, map0.Has(e))

	e, _ = mut.NewEntityWith(&testStruct0{})
	assert.True(t, map0.Has(e))

	s0 = mut.Get(e)
	assert.NotNil(t, s0)
}

func TestMap2Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap2[testStruct0, testStruct1](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0, s1 := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	s0, s1 = mut.Add(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	_, _ = mut.Assign(e, &testStruct0{}, &testStruct1{})
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e, _, _ = mut.NewEntityWith(&testStruct0{}, &testStruct1{})
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1 = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
}

func TestMap3Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap3[
		testStruct0, testStruct1, testStruct2,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0, s1, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	s0, s1, _ = mut.Add(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	_, _, _ = mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e, _, _, _ = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
}

func TestMap4Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap4[
		testStruct0, testStruct1, testStruct2, testStruct3,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0, s1, _, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	s0, s1, _, _ = mut.Add(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	_, _, _, _ = mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e, _, _, _, _ = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
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

	e, s0, s1, _, _, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	s0, s1, _, _, _ = mut.Add(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	_, _, _, _, _ = mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e, _, _, _, _, _ = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
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

	e, s0, s1, _, _, _, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	s0, s1, _, _, _, _ = mut.Add(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	_, _, _, _, _, _ = mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e, _, _, _, _, _, _ = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
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

	e, s0, s1, _, _, _, _, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	s0, s1, _, _, _, _, _ = mut.Add(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	_, _, _, _, _, _, _ = mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e, _, _, _, _, _, _, _ = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
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

	e, s0, s1, _, _, _, _, _, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	s0, s1, _, _, _, _, _, _ = mut.Add(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	_, _, _, _, _, _, _, _ = mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e, _, _, _, _, _, _, _, _ = mut.NewEntityWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	s0, s1, _, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
}
