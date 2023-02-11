package ecs

import (
	"internal/base"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery0(t *testing.T) {
	w := NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e1, 0, &testStruct0{2})
	w.Assign(e2, 0, &testStruct0{2})
	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	query :=
		Query0(&w).
			Not(Mask1[testStruct9](&w))

	for query.Next() {
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func TestQuery1(t *testing.T) {
	w := NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e1, 0, &testStruct0{2})
	w.Assign(e2, 0, &testStruct0{0})
	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	query :=
		Query1[testStruct0](&w).
			Not(Mask1[testStruct9](&w))

	for query.Next() {
		c0 := query.Get1()
		assert.Equal(t, cnt+1, int(c0.val))
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func TestQuery2(t *testing.T) {
	w := NewWorld()

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

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	query :=
		Query2[testStruct0, testStruct1](&w).
			Not(Mask1[testStruct9](&w))

	for query.Next() {
		c1 := query.Get1()
		c2 := query.Get2()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))

		c12, c22 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func TestQuery3(t *testing.T) {
	w := NewWorld()

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

	w.Assign(e0, 2, &testStruct2{3})
	w.Assign(e1, 2, &testStruct2{4})
	w.Assign(e2, 2, &testStruct2{5})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	query :=
		Query3[testStruct0, testStruct1, testStruct2](&w).
			Not(Mask1[testStruct9](&w))

	for query.Next() {
		c1 := query.Get1()
		c2 := query.Get2()
		c3 := query.Get3()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))

		c12, c22, c32 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func TestQuery4(t *testing.T) {
	w := NewWorld()

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

	w.Assign(e0, 2, &testStruct2{3})
	w.Assign(e1, 2, &testStruct2{4})
	w.Assign(e2, 2, &testStruct2{4})

	w.Assign(e0, 3, &testStruct3{4})
	w.Assign(e1, 3, &testStruct3{5})
	w.Assign(e2, 3, &testStruct3{5})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	query :=
		Query4[testStruct0, testStruct1, testStruct2, testStruct3](&w).
			Not(Mask1[testStruct9](&w))
	for query.Next() {
		c1 := query.Get1()
		c2 := query.Get2()
		c3 := query.Get3()
		c4 := query.Get4()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))

		c12, c22, c32, c42 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		assert.Equal(t, c4, c42)
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func TestQuery5(t *testing.T) {
	w := NewWorld()

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

	w.Assign(e0, 2, &testStruct2{3})
	w.Assign(e1, 2, &testStruct2{4})
	w.Assign(e2, 2, &testStruct2{4})

	w.Assign(e0, 3, &testStruct3{4})
	w.Assign(e1, 3, &testStruct3{5})
	w.Assign(e2, 3, &testStruct3{5})

	w.Assign(e0, 4, &testStruct4{5})
	w.Assign(e1, 4, &testStruct4{6})
	w.Assign(e2, 4, &testStruct4{6})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	query :=
		Query5[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4](&w).
			Not(Mask1[testStruct9](&w))
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

		c12, c22, c32, c42, c52 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		assert.Equal(t, c4, c42)
		assert.Equal(t, c5, c52)
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func TestQuery6(t *testing.T) {
	w := NewWorld()

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

	w.Assign(e0, 2, &testStruct2{3})
	w.Assign(e1, 2, &testStruct2{4})
	w.Assign(e2, 2, &testStruct2{4})

	w.Assign(e0, 3, &testStruct3{4})
	w.Assign(e1, 3, &testStruct3{5})
	w.Assign(e2, 3, &testStruct3{5})

	w.Assign(e0, 4, &testStruct4{5})
	w.Assign(e1, 4, &testStruct4{6})
	w.Assign(e2, 4, &testStruct4{6})

	w.Assign(e0, 5, &testStruct5{6})
	w.Assign(e1, 5, &testStruct5{7})
	w.Assign(e2, 5, &testStruct5{7})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	query :=
		Query6[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4, testStruct5](&w).
			Not(Mask1[testStruct9](&w))
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

		c12, c22, c32, c42, c52, c62 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		assert.Equal(t, c4, c42)
		assert.Equal(t, c5, c52)
		assert.Equal(t, c6, c62)
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func TestQuery7(t *testing.T) {
	w := NewWorld()

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

	w.Assign(e0, 2, &testStruct2{3})
	w.Assign(e1, 2, &testStruct2{4})
	w.Assign(e2, 2, &testStruct2{4})

	w.Assign(e0, 3, &testStruct3{4})
	w.Assign(e1, 3, &testStruct3{5})
	w.Assign(e2, 3, &testStruct3{5})

	w.Assign(e0, 4, &testStruct4{5})
	w.Assign(e1, 4, &testStruct4{6})
	w.Assign(e2, 4, &testStruct4{6})

	w.Assign(e0, 5, &testStruct5{6})
	w.Assign(e1, 5, &testStruct5{7})
	w.Assign(e2, 5, &testStruct5{7})

	w.Assign(e0, 6, &testStruct5{7})
	w.Assign(e1, 6, &testStruct5{8})
	w.Assign(e2, 6, &testStruct5{8})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	query :=
		Query7[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4, testStruct5, testStruct6](&w).
			Not(Mask1[testStruct9](&w))
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

		c12, c22, c32, c42, c52, c62, c72 := query.GetAll()
		assert.Equal(t, c1, c12)
		assert.Equal(t, c2, c22)
		assert.Equal(t, c3, c32)
		assert.Equal(t, c4, c42)
		assert.Equal(t, c5, c52)
		assert.Equal(t, c6, c62)
		assert.Equal(t, c7, c72)
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func TestQuery8(t *testing.T) {
	w := NewWorld()

	registerAll(&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, 0, &testStruct0{1})
	w.Assign(e1, 0, &testStruct0{2})

	w.Assign(e0, 1, &testStruct1{2})
	w.Assign(e1, 1, &testStruct1{3})

	w.Assign(e0, 2, &testStruct2{3})
	w.Assign(e1, 2, &testStruct2{4})

	w.Assign(e0, 3, &testStruct3{4})
	w.Assign(e1, 3, &testStruct3{5})

	w.Assign(e0, 4, &testStruct4{5})
	w.Assign(e1, 4, &testStruct4{6})

	w.Assign(e0, 5, &testStruct5{6})
	w.Assign(e1, 5, &testStruct5{7})

	w.Assign(e0, 6, &testStruct5{7})
	w.Assign(e1, 6, &testStruct5{8})

	w.Assign(e0, 7, &testStruct5{8})
	w.Assign(e1, 7, &testStruct5{9})

	w.Assign(e2, 9, &testStruct9{})

	cnt := 0
	query :=
		Query8[
			testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testStruct7](&w).
			Not(Mask1[testStruct9](&w))

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

		c12, c22, c32, c42, c52, c62, c72, c82 := query.GetAll()
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
	assert.Equal(t, 2, cnt)
}

func TestGenericMask(t *testing.T) {
	w := NewWorld()

	assert.Equal(t,
		Mask1[testStruct0](&w).BitMask,
		base.NewBitMask(0),
	)

	assert.Equal(t,
		Mask2[testStruct0, testStruct1](&w).BitMask,
		base.NewBitMask(0, 1),
	)

	assert.Equal(t,
		Mask3[testStruct0, testStruct1, testStruct2](&w).BitMask,
		base.NewBitMask(0, 1, 2),
	)

	assert.Equal(t,
		Mask4[testStruct0, testStruct1, testStruct2, testStruct3](&w).BitMask,
		base.NewBitMask(0, 1, 2, 3),
	)

	assert.Equal(t,
		Mask5[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4](&w).BitMask,
		base.NewBitMask(0, 1, 2, 3, 4),
	)

	assert.Equal(t,
		Mask6[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5](&w).BitMask,
		base.NewBitMask(0, 1, 2, 3, 4, 5),
	)

	assert.Equal(t,
		Mask7[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6](&w).BitMask,
		base.NewBitMask(0, 1, 2, 3, 4, 5, 6),
	)

	assert.Equal(t,
		Mask8[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testStruct7](&w).BitMask,
		base.NewBitMask(0, 1, 2, 3, 4, 5, 6, 7),
	)
}
