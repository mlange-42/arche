package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestMutate0(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMutate0(&w)

	_ = mut.NewEntity()
}

func TestMutate1(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMutate1[testStruct0](&w)
	mutR := NewMutate0(&w).WithRemove(T[testStruct0]())
	mutR1 := NewMutate1[testStruct1](&w).WithRemove(T[testStruct0]())

	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0 := mut.NewEntity()
	assert.NotNil(t, s0)

	mutR.Remove(e)
	assert.False(t, map0.Has(e))

	s0 = mut.Add(e)
	assert.NotNil(t, s0)

	e = w.NewEntity()
	_ = mut.Assign(e, &testStruct0{})
	assert.True(t, map0.Has(e))

	e, _ = mut.NewEntityWith(&testStruct0{})
	assert.True(t, map0.Has(e))

	mutR1.Exchange(e)
	assert.False(t, map0.Has(e))
	assert.True(t, map1.Has(e))
}

func TestMutate2(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMutate2[testStruct0, testStruct1](&w).WithRemove()
	mutR := NewMutate0(&w).WithRemove(T2[testStruct0, testStruct1]()...)

	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0, s1 := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutR.Remove(e)
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
}

func TestMutate3(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMutate3[
		testStruct0, testStruct1, testStruct2,
	](&w).WithRemove()

	mutR := NewMutate0(&w).WithRemove(T3[
		testStruct0, testStruct1, testStruct2,
	]()...)

	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0, s1, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutR.Remove(e)
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
}

func TestMutate4(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMutate4[
		testStruct0, testStruct1, testStruct2, testStruct3,
	](&w).WithRemove()

	mutR := NewMutate0(&w).WithRemove(T4[
		testStruct0, testStruct1, testStruct2, testStruct3,
	]()...)

	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0, s1, _, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutR.Remove(e)
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
}

func TestMutate5(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMutate5[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4,
	](&w).WithRemove()

	mutR := NewMutate0(&w).WithRemove(T5[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4,
	]()...)

	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0, s1, _, _, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutR.Remove(e)
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
}

func TestMutate6(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMutate6[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5,
	](&w).WithRemove()

	mutR := NewMutate0(&w).WithRemove(T6[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5,
	]()...)

	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0, s1, _, _, _, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutR.Remove(e)
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
}

func TestMutate7(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMutate7[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6,
	](&w).WithRemove()

	mutR := NewMutate0(&w).WithRemove(T7[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6,
	]()...)

	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0, s1, _, _, _, _, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutR.Remove(e)
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
}

func TestMutate8(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMutate8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	](&w).WithRemove()

	mutR := NewMutate0(&w).WithRemove(T8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	]()...)

	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e, s0, s1, _, _, _, _, _, _ := mut.NewEntity()
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutR.Remove(e)
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
}
