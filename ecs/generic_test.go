package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericMap(t *testing.T) {
	w := NewWorld()
	get := NewMap[testStruct0](&w)

	e0 := w.NewEntity()

	Add[testStruct0](&w, e0)
	has := get.Has(e0)
	_ = get.Get(e0)
	assert.True(t, has)

	_ = get.Set(e0, &testStruct0{100})
	str := get.Get(e0)

	assert.Equal(t, 100, int(str.val))

	get2 := NewMap[testStruct1](&w)
	assert.Panics(t, func() { get2.Set(e0, &testStruct1{}) })
}

func TestGenericAddRemove(t *testing.T) {
	w := NewWorld()
	get := NewMap[testStruct0](&w)

	e0 := w.NewEntity()

	Add[testStruct0](&w, e0)
	_ = get.Has(e0)
	_ = get.Get(e0)
	Remove[testStruct0](&w, e0)

	Add2[testStruct0, testStruct1](&w, e0)
	Remove2[testStruct0, testStruct1](&w, e0)

	Add3[testStruct0, testStruct1, testStruct2](&w, e0)
	Remove3[testStruct0, testStruct1, testStruct2](&w, e0)

	Add4[testStruct0, testStruct1, testStruct2, testStruct3](&w, e0)
	Remove4[testStruct0, testStruct1, testStruct2, testStruct3](&w, e0)

	Add5[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4](&w, e0)
	Remove5[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4](&w, e0)
}

func TestGenericAssignRemove(t *testing.T) {
	w := NewWorld()

	e0 := w.NewEntity()

	Assign(&w, e0, &testStruct0{})
	Remove[testStruct0](&w, e0)

	Assign2(&w, e0, &testStruct0{}, &testStruct1{})
	Remove2[testStruct0, testStruct1](&w, e0)

	Assign3(&w, e0, &testStruct0{}, &testStruct1{}, &testStruct2{})
	Remove3[testStruct0, testStruct1, testStruct2](&w, e0)

	Assign4(&w, e0, &testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{})
	Remove4[testStruct0, testStruct1, testStruct2, testStruct3](&w, e0)

	Assign5(&w, e0, &testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{}, &testStruct4{})
	Remove5[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4](&w, e0)
}
