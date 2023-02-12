package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

type testStruct0 struct{ val int8 }
type testStruct1 struct{ val int32 }
type testStruct2 struct {
	val  int32
	val2 int32
}
type testStruct3 struct {
	val  int32
	val2 int32
	val3 int32
}
type testStruct4 struct {
	val  int32
	val2 int32
	val3 int32
	val4 int32
}
type testStruct5 struct {
	val  int32
	val2 int32
	val3 int32
	val4 int32
	val5 int32
}
type testStruct6 struct {
	val  int32
	val2 int32
	val3 int32
	val4 int32
	val5 int32
	val6 int32
}
type testStruct7 struct {
	val  int32
	val2 int32
	val3 int32
	val4 int32
	val5 int32
	val6 int32
	val7 int32
}
type testStruct8 struct {
	val  int32
	val2 int32
	val3 int32
	val4 int32
	val5 int32
	val6 int32
	val7 int32
	val8 int32
}
type testStruct9 struct {
	val  int32
	val2 int32
	val3 int32
	val4 int32
	val5 int32
	val6 int32
	val7 int32
	val8 int32
	val9 int32
}

func registerAll(w *ecs.World) []ecs.ID {
	_ = testStruct0{}
	_ = testStruct1{}
	_ = testStruct2{}
	_ = testStruct3{}
	_ = testStruct4{}
	_ = testStruct5{}
	_ = testStruct6{}
	_ = testStruct7{}
	_ = testStruct8{}
	_ = testStruct9{}

	ids := make([]ecs.ID, 10)
	ids[0] = ecs.ComponentID[testStruct0](w)
	ids[1] = ecs.ComponentID[testStruct1](w)
	ids[2] = ecs.ComponentID[testStruct2](w)
	ids[3] = ecs.ComponentID[testStruct3](w)
	ids[4] = ecs.ComponentID[testStruct4](w)
	ids[5] = ecs.ComponentID[testStruct5](w)
	ids[6] = ecs.ComponentID[testStruct6](w)
	ids[7] = ecs.ComponentID[testStruct7](w)
	ids[8] = ecs.ComponentID[testStruct8](w)
	ids[9] = ecs.ComponentID[testStruct9](w)

	return ids
}

func TestGenericMap(t *testing.T) {
	w := ecs.NewWorld()
	get := NewMap[testStruct0](&w)

	e0 := w.NewEntity()

	Add1[testStruct0](&w, e0)
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
	w := ecs.NewWorld()
	get := NewMap[testStruct0](&w)

	e0 := w.NewEntity()

	Add1[testStruct0](&w, e0)
	_ = get.Has(e0)
	_ = get.Get(e0)
	Remove1[testStruct0](&w, e0)

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
	w := ecs.NewWorld()

	e0 := w.NewEntity()

	Assign1(&w, e0, &testStruct0{})
	Remove1[testStruct0](&w, e0)

	Assign2(&w, e0, &testStruct0{}, &testStruct1{})
	Remove2[testStruct0, testStruct1](&w, e0)

	Assign3(&w, e0, &testStruct0{}, &testStruct1{}, &testStruct2{})
	Remove3[testStruct0, testStruct1, testStruct2](&w, e0)

	Assign4(&w, e0, &testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{})
	Remove4[testStruct0, testStruct1, testStruct2, testStruct3](&w, e0)

	Assign5(&w, e0, &testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{}, &testStruct4{})
	Remove5[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4](&w, e0)
}
