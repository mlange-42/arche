package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

//lint:ignore U1000 test type
type testStruct0 struct{ val int8 }

//lint:ignore U1000 test type
type testStruct1 struct{ val int32 }

//lint:ignore U1000 test type
type testStruct2 struct {
	val  int32
	val2 int32
}

//lint:ignore U1000 test type
type testStruct3 struct {
	val  int32
	val2 int32
	val3 int32
}

//lint:ignore U1000 test type
type testStruct4 struct {
	val  int32
	val2 int32
	val3 int32
	val4 int32
}

//lint:ignore U1000 test type
type testStruct5 struct {
	val  int32
	val2 int32
	val3 int32
	val4 int32
	val5 int32
}

//lint:ignore U1000 test type
type testStruct6 struct {
	val  int32
	val2 int32
	val3 int32
	val4 int32
	val5 int32
	val6 int32
}

//lint:ignore U1000 test type
type testStruct7 struct {
	val  int32
	val2 int32
	val3 int32
	val4 int32
	val5 int32
	val6 int32
	val7 int32
}

//lint:ignore U1000 test type
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

//lint:ignore U1000 test type
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

func TestGenericNewEntity(t *testing.T) {
	w := ecs.NewConfig().WithCapacityIncrement(32).Build()

	id0 := ecs.ComponentID[testStruct0](&w)
	id1 := ecs.ComponentID[testStruct1](&w)
	id2 := ecs.ComponentID[testStruct2](&w)

	e0 := w.NewEntity()
	e1, _, _, _ := New3[testStruct0, testStruct1, testStruct2](&w)
	e2, _, _, _ := NewWith3(&w, &testStruct0{1}, &testStruct1{2}, &testStruct2{3, 4})
	e3 := w.NewEntityWith()

	assert.Equal(t, ecs.NewBitMask(), w.Mask(e0))
	assert.Equal(t, ecs.NewBitMask(id0, id1, id2), w.Mask(e1))
	assert.Equal(t, ecs.NewBitMask(id0, id1, id2), w.Mask(e2))
	assert.Equal(t, ecs.NewBitMask(), w.Mask(e3))

	s0 := (*testStruct0)(w.Get(e2, id0))
	s1 := (*testStruct1)(w.Get(e2, id1))
	s2 := (*testStruct2)(w.Get(e2, id2))

	assert.Equal(t, &testStruct0{1}, s0)
	assert.Equal(t, &testStruct1{2}, s1)
	assert.Equal(t, &testStruct2{3, 4}, s2)

	w.RemEntity(e0)
	w.RemEntity(e1)
	w.RemEntity(e2)
	w.RemEntity(e3)

	for i := 0; i < 35; i++ {
		e, _, _, _ := NewWith3(&w,
			&testStruct0{int8(i + 1)},
			&testStruct1{int32(i + 2)},
			&testStruct2{int32(i + 3), int32(i + 4)},
		)

		s0 := (*testStruct0)(w.Get(e, id0))
		s1 := (*testStruct1)(w.Get(e, id1))
		s2 := (*testStruct2)(w.Get(e, id2))

		assert.Equal(t, &testStruct0{int8(i + 1)}, s0)
		assert.Equal(t, &testStruct1{int32(i + 2)}, s1)
		assert.Equal(t, &testStruct2{int32(i + 3), int32(i + 4)}, s2)
	}
}
