package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestQueryOptionalNot(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e0, 1, &testStruct1{1})

	w.Assign(e1, 0, &testStruct0{2})
	w.Assign(e1, 1, &testStruct1{2})
	w.Assign(e1, 2, &testStruct2{1, 1})
	w.Assign(e1, 8, &testStruct8{})

	w.Assign(e2, 0, &testStruct0{3})
	w.Assign(e2, 1, &testStruct1{3})
	w.Assign(e2, 2, &testStruct2{1, 1})
	w.Assign(e2, 9, &testStruct9{})

	query2 := NewFilter2[testStruct0, testStruct1]().Query(&w)
	cnt := 0
	for query2.Next() {
		cnt++
	}
	assert.Equal(t, 3, cnt)

	query2 = NewFilter2[testStruct0, testStruct1]().Without(T[testStruct9]()).Query(&w)
	cnt = 0
	for query2.Next() {
		cnt++
	}
	assert.Equal(t, 2, cnt)

	query2 = NewFilter2[testStruct0, testStruct1]().Without(T[testStruct8](), T[testStruct9]()).Query(&w)
	cnt = 0
	for query2.Next() {
		cnt++
	}
	assert.Equal(t, 1, cnt)

	query2 = NewFilter2[testStruct0, testStruct1]().With(T[testStruct2]()).Without(T[testStruct9]()).Query(&w)
	cnt = 0
	for query2.Next() {
		cnt++
	}
	assert.Equal(t, 1, cnt)

	query3 := NewFilter3[testStruct0, testStruct1, testStruct9]().Query(&w)
	cnt = 0
	for query3.Next() {
		cnt++
	}
	assert.Equal(t, 1, cnt)

	query3 = NewFilter3[testStruct0, testStruct1, testStruct9]().Optional(T[testStruct9]()).Query(&w)
	cnt = 0
	for query3.Next() {
		cnt++
	}
	assert.Equal(t, 3, cnt)
}

func TestQuery0(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e0, 8, &testStruct8{})

	w.Assign(e1, 0, &testStruct0{2})

	w.Assign(e2, 0, &testStruct0{2})
	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	filter :=
		NewFilter0().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	assert.Equal(t, ecs.All(8).Without(9), filter.Filter(&w))
	query := filter.Query(&w)
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 1, cnt)

}

func TestQuery1(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})

	w.Assign(e1, 0, &testStruct0{2})

	w.Assign(e0, 8, &testStruct8{9, 0, 0, 0, 0, 0, 0, 0})

	w.Assign(e2, 0, &testStruct0{0})
	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	filter :=
		NewFilter1[testStruct0]().
			Optional(T[testStruct9]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	assert.Equal(t, ecs.All(0, 8).Without(9), filter.Filter(&w))
	query := filter.Query(&w)
	for query.Next() {
		c0 := query.Get1()
		_, c02 := query.GetAll()
		assert.Equal(t, c0, c02)
		assert.Equal(t, cnt+1, int(c0.val))
		cnt++
	}
	assert.Equal(t, 1, cnt)
}

func TestQuery2(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e1, 0, &testStruct0{2})
	w.Assign(e2, 0, &testStruct0{3})

	w.Assign(e0, 1, &testStruct1{2})
	w.Assign(e1, 1, &testStruct1{3})
	w.Assign(e2, 1, &testStruct1{4})

	w.Assign(e0, 8, &testStruct8{9, 0, 0, 0, 0, 0, 0, 0})

	w.Assign(e2, 9, &testStruct9{})

	filter :=
		NewFilter2[testStruct0, testStruct1]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	assert.Equal(t, ecs.All(0, 8).Without(9), filter.Filter(&w))

	for i := 0; i < 10; i++ {
		cnt := 0
		q := filter.Query(&w)
		for q.Next() {
			c1 := q.Get1()
			c2 := q.Get2()
			assert.Equal(t, cnt+1, int(c1.val))
			assert.Equal(t, cnt+2, int(c2.val))

			_, c12, c22 := q.GetAll()
			assert.Equal(t, c1, c12)
			assert.Equal(t, c2, c22)
			cnt++
		}
		assert.Equal(t, 1, cnt)
	}
}

func TestQuery3(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e1, 0, &testStruct0{2})
	w.Assign(e2, 0, &testStruct0{3})

	w.Assign(e0, 1, &testStruct1{2})
	w.Assign(e1, 1, &testStruct1{3})
	w.Assign(e2, 1, &testStruct1{4})

	w.Assign(e0, 2, &testStruct2{3, 0})
	w.Assign(e1, 2, &testStruct2{4, 0})
	w.Assign(e2, 2, &testStruct2{5, 0})

	w.Assign(e0, 8, &testStruct8{9, 0, 0, 0, 0, 0, 0, 0})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	filter :=
		NewFilter3[testStruct0, testStruct1, testStruct2]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	assert.Equal(t, ecs.All(0, 2, 8).Without(9), filter.Filter(&w))
	query := filter.Query(&w)
	for query.Next() {
		c1 := query.Get1()
		c2 := query.Get2()
		c3 := query.Get3()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))

		_, c12, c22, c32 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		cnt++
	}
	assert.Equal(t, 1, cnt)
}

func TestQuery4(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e1, 0, &testStruct0{2})
	w.Assign(e2, 0, &testStruct0{2})

	w.Assign(e0, 1, &testStruct1{2})
	w.Assign(e1, 1, &testStruct1{3})
	w.Assign(e2, 1, &testStruct1{3})

	w.Assign(e0, 2, &testStruct2{3, 0})
	w.Assign(e1, 2, &testStruct2{4, 0})
	w.Assign(e2, 2, &testStruct2{4, 0})

	w.Assign(e0, 3, &testStruct3{4, 0, 0})
	w.Assign(e1, 3, &testStruct3{5, 0, 0})
	w.Assign(e2, 3, &testStruct3{5, 0, 0})

	w.Assign(e0, 8, &testStruct8{9, 0, 0, 0, 0, 0, 0, 0})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	filter :=
		NewFilter4[testStruct0, testStruct1, testStruct2, testStruct3]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	assert.Equal(t, ecs.All(0, 2, 3, 8).Without(9), filter.Filter(&w))
	query := filter.Query(&w)
	for query.Next() {
		c1 := query.Get1()
		c2 := query.Get2()
		c3 := query.Get3()
		c4 := query.Get4()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))

		_, c12, c22, c32, c42 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		assert.Equal(t, c4, c42)
		cnt++
	}
	assert.Equal(t, 1, cnt)
}

func TestQuery5(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e1, 0, &testStruct0{2})
	w.Assign(e2, 0, &testStruct0{2})

	w.Assign(e0, 1, &testStruct1{2})
	w.Assign(e1, 1, &testStruct1{3})
	w.Assign(e2, 1, &testStruct1{3})

	w.Assign(e0, 2, &testStruct2{3, 0})
	w.Assign(e1, 2, &testStruct2{4, 0})
	w.Assign(e2, 2, &testStruct2{4, 0})

	w.Assign(e0, 3, &testStruct3{4, 0, 0})
	w.Assign(e1, 3, &testStruct3{5, 0, 0})
	w.Assign(e2, 3, &testStruct3{5, 0, 0})

	w.Assign(e0, 4, &testStruct4{5, 0, 0, 0})
	w.Assign(e1, 4, &testStruct4{6, 0, 0, 0})
	w.Assign(e2, 4, &testStruct4{6, 0, 0, 0})

	w.Assign(e0, 8, &testStruct8{9, 0, 0, 0, 0, 0, 0, 0})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	filter :=
		NewFilter5[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	assert.Equal(t, ecs.All(0, 2, 3, 4, 8).Without(9), filter.Filter(&w))
	query := filter.Query(&w)
	for query.Next() {
		c1 := query.Get1()
		c2 := query.Get2()
		c3 := query.Get3()
		c4 := query.Get4()
		c5 := query.Get5()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))
		assert.Equal(t, cnt+5, int(c5.val))

		_, c12, c22, c32, c42, c52 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		assert.Equal(t, c4, c42)
		assert.Equal(t, c5, c52)
		cnt++
	}
	assert.Equal(t, 1, cnt)
}

func TestQuery6(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e1, 0, &testStruct0{2})
	w.Assign(e2, 0, &testStruct0{2})

	w.Assign(e0, 1, &testStruct1{2})
	w.Assign(e1, 1, &testStruct1{3})
	w.Assign(e2, 1, &testStruct1{3})

	w.Assign(e0, 2, &testStruct2{3, 0})
	w.Assign(e1, 2, &testStruct2{4, 0})
	w.Assign(e2, 2, &testStruct2{4, 0})

	w.Assign(e0, 3, &testStruct3{4, 0, 0})
	w.Assign(e1, 3, &testStruct3{5, 0, 0})
	w.Assign(e2, 3, &testStruct3{5, 0, 0})

	w.Assign(e0, 4, &testStruct4{5, 0, 0, 0})
	w.Assign(e1, 4, &testStruct4{6, 0, 0, 0})
	w.Assign(e2, 4, &testStruct4{6, 0, 0, 0})

	w.Assign(e0, 5, &testStruct5{6, 0, 0, 0, 0})
	w.Assign(e1, 5, &testStruct5{7, 0, 0, 0, 0})
	w.Assign(e2, 5, &testStruct5{7, 0, 0, 0, 0})

	w.Assign(e0, 8, &testStruct8{9, 0, 0, 0, 0, 0, 0, 0})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	filter :=
		NewFilter6[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4, testStruct5]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	assert.Equal(t, ecs.All(0, 2, 3, 4, 5, 8).Without(9), filter.Filter(&w))
	query := filter.Query(&w)
	for query.Next() {
		c1 := query.Get1()
		c2 := query.Get2()
		c3 := query.Get3()
		c4 := query.Get4()
		c5 := query.Get5()
		c6 := query.Get6()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))
		assert.Equal(t, cnt+5, int(c5.val))
		assert.Equal(t, cnt+6, int(c6.val))

		_, c12, c22, c32, c42, c52, c62 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		assert.Equal(t, c4, c42)
		assert.Equal(t, c5, c52)
		assert.Equal(t, c6, c62)
		cnt++
	}
	assert.Equal(t, 1, cnt)
}

func TestQuery7(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e1, 0, &testStruct0{2})
	w.Assign(e2, 0, &testStruct0{2})

	w.Assign(e0, 1, &testStruct1{2})
	w.Assign(e1, 1, &testStruct1{3})
	w.Assign(e2, 1, &testStruct1{3})

	w.Assign(e0, 2, &testStruct2{3, 0})
	w.Assign(e1, 2, &testStruct2{4, 0})
	w.Assign(e2, 2, &testStruct2{4, 0})

	w.Assign(e0, 3, &testStruct3{4, 0, 0})
	w.Assign(e1, 3, &testStruct3{5, 0, 0})
	w.Assign(e2, 3, &testStruct3{5, 0, 0})

	w.Assign(e0, 4, &testStruct4{5, 0, 0, 0})
	w.Assign(e1, 4, &testStruct4{6, 0, 0, 0})
	w.Assign(e2, 4, &testStruct4{6, 0, 0, 0})

	w.Assign(e0, 5, &testStruct5{6, 0, 0, 0, 0})
	w.Assign(e1, 5, &testStruct5{7, 0, 0, 0, 0})
	w.Assign(e2, 5, &testStruct5{7, 0, 0, 0, 0})

	w.Assign(e0, 6, &testStruct6{7, 0, 0, 0, 0, 0})
	w.Assign(e1, 6, &testStruct6{8, 0, 0, 0, 0, 0})
	w.Assign(e2, 6, &testStruct6{8, 0, 0, 0, 0, 0})

	w.Assign(e0, 8, &testStruct8{9, 0, 0, 0, 0, 0, 0, 0})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	filter :=
		NewFilter7[
			testStruct0, testStruct1, testStruct2, testStruct3, testStruct4,
			testStruct5, testStruct6,
		]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	assert.Equal(t, ecs.All(0, 2, 3, 4, 5, 6, 8).Without(9), filter.Filter(&w))
	query := filter.Query(&w)
	for query.Next() {
		c1 := query.Get1()
		c2 := query.Get2()
		c3 := query.Get3()
		c4 := query.Get4()
		c5 := query.Get5()
		c6 := query.Get6()
		c7 := query.Get7()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))
		assert.Equal(t, cnt+5, int(c5.val))
		assert.Equal(t, cnt+6, int(c6.val))
		assert.Equal(t, cnt+7, int(c7.val))

		_, c12, c22, c32, c42, c52, c62, c72 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		assert.Equal(t, c4, c42)
		assert.Equal(t, c5, c52)
		assert.Equal(t, c6, c62)
		assert.Equal(t, c7, c72)
		cnt++
	}
	assert.Equal(t, 1, cnt)
}

func TestQuery8(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e1, 0, &testStruct0{2})

	w.Assign(e0, 1, &testStruct1{2})
	w.Assign(e1, 1, &testStruct1{3})

	w.Assign(e0, 2, &testStruct2{3, 0})
	w.Assign(e1, 2, &testStruct2{4, 0})

	w.Assign(e0, 3, &testStruct3{4, 0, 0})
	w.Assign(e1, 3, &testStruct3{5, 0, 0})

	w.Assign(e0, 4, &testStruct4{5, 0, 0, 0})
	w.Assign(e1, 4, &testStruct4{6, 0, 0, 0})

	w.Assign(e0, 5, &testStruct5{6, 0, 0, 0, 0})
	w.Assign(e1, 5, &testStruct5{7, 0, 0, 0, 0})

	w.Assign(e0, 6, &testStruct6{7, 0, 0, 0, 0, 0})
	w.Assign(e1, 6, &testStruct6{8, 0, 0, 0, 0, 0})

	w.Assign(e0, 7, &testStruct7{8, 0, 0, 0, 0, 0, 0})
	w.Assign(e1, 7, &testStruct7{9, 0, 0, 0, 0, 0, 0})

	w.Assign(e0, 8, &testStruct8{9, 0, 0, 0, 0, 0, 0, 0})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	filter :=
		NewFilter8[
			testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testStruct7,
		]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	assert.Equal(t, ecs.All(0, 2, 3, 4, 5, 6, 7, 8).Without(9), filter.Filter(&w))
	query := filter.Query(&w)
	for query.Next() {
		c1 := query.Get1()
		c2 := query.Get2()
		c3 := query.Get3()
		c4 := query.Get4()
		c5 := query.Get5()
		c6 := query.Get6()
		c7 := query.Get7()
		c8 := query.Get8()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))
		assert.Equal(t, cnt+5, int(c5.val))
		assert.Equal(t, cnt+6, int(c6.val))
		assert.Equal(t, cnt+7, int(c7.val))
		assert.Equal(t, cnt+8, int(c8.val))

		_, c12, c22, c32, c42, c52, c62, c72, c82 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		assert.Equal(t, c4, c42)
		assert.Equal(t, c5, c52)
		assert.Equal(t, c6, c62)
		assert.Equal(t, c7, c72)
		assert.Equal(t, c8, c82)
		cnt++
	}
	assert.Equal(t, 1, cnt)
}

func TestQueryGeneric(t *testing.T) {
	count := 1000
	world := ecs.NewWorld()

	posID := ecs.ComponentID[testStruct2](&world)
	rotID := ecs.ComponentID[testStruct3](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}
	query := NewFilter2[testStruct2, testStruct3]()

	q := query.Query(&world)
	cnt := 0
	for q.Next() {
		s1 := q.Get1()
		s2 := q.Get2()
		_ = s1
		_ = s2
		cnt++
	}
	assert.Equal(t, count, cnt)
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
