package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestGenericMap(t *testing.T) {
	w := ecs.NewWorld()
	get := NewMap[testStruct0](&w)
	mut := NewMutate(&w).WithAdd(T[testStruct0]())

	assert.Equal(t, ecs.ComponentID[testStruct0](&w), get.ID())

	e0 := w.NewEntity()

	mut.Add(e0)
	has := get.Has(e0)
	_ = get.Get(e0)
	assert.True(t, has)

	_ = get.Set(e0, &testStruct0{100})
	str := get.Get(e0)

	assert.Equal(t, 100, int(str.val))

	get2 := NewMap[testStruct1](&w)
	assert.Equal(t, ecs.ComponentID[testStruct1](&w), get2.ID())
	assert.Panics(t, func() { get2.Set(e0, &testStruct1{}) })
}

func TestMap1(t *testing.T) {
	w := ecs.NewWorld()

	mut := NewMutate(&w).WithAdd(T8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	]()...)

	mapper := NewMap1[testStruct0](&w)

	e := mut.NewEntity()

	s0 := mapper.Get(e)
	assert.NotNil(t, s0)
}

func TestMap2(t *testing.T) {
	w := ecs.NewWorld()

	mut := NewMutate(&w).WithAdd(T8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	]()...)

	mapper := NewMap2[testStruct0, testStruct1](&w)

	e := mut.NewEntity()

	s0, s1 := mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
}

func TestMap3(t *testing.T) {
	w := ecs.NewWorld()

	mut := NewMutate(&w).WithAdd(T8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	]()...)

	mapper := NewMap3[
		testStruct0, testStruct1, testStruct2,
	](&w)

	e := mut.NewEntity()

	s0, s1, _ := mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
}

func TestMap4(t *testing.T) {
	w := ecs.NewWorld()

	mut := NewMutate(&w).WithAdd(T8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	]()...)

	mapper := NewMap4[
		testStruct0, testStruct1, testStruct2, testStruct3,
	](&w)

	e := mut.NewEntity()

	s0, s1, _, _ := mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
}

func TestMap5(t *testing.T) {
	w := ecs.NewWorld()

	mut := NewMutate(&w).WithAdd(T8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	]()...)

	mapper := NewMap5[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4,
	](&w)

	e := mut.NewEntity()

	s0, s1, _, _, _ := mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
}

func TestMap6(t *testing.T) {
	w := ecs.NewWorld()

	mut := NewMutate(&w).WithAdd(T8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	]()...)

	mapper := NewMap6[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5,
	](&w)

	e := mut.NewEntity()

	s0, s1, _, _, _, _ := mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
}

func TestMap7(t *testing.T) {
	w := ecs.NewWorld()

	mut := NewMutate(&w).WithAdd(T8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	]()...)

	mapper := NewMap7[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6,
	](&w)

	e := mut.NewEntity()

	s0, s1, _, _, _, _, _ := mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
}

func TestMap8(t *testing.T) {
	w := ecs.NewWorld()

	mut := NewMutate(&w).WithAdd(T8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	]()...)

	mapper := NewMap8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	](&w)

	e := mut.NewEntity()

	s0, s1, _, _, _, _, _, _ := mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)
}
