package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[position](&w)
	rotID := ComponentID[rotation](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	e4 := w.NewEntity()

	w.Add(e0, posID)
	w.Add(e1, posID, rotID)
	w.Add(e2, posID, rotID)
	w.Add(e3, rotID)
	w.Add(e4, rotID)

	q := w.Query(posID, rotID)
	cnt := 0
	for q.Next() {
		ent := q.Entity()
		pos := (*position)(q.Get(posID))
		rot := (*rotation)(q.Get(rotID))
		_ = ent
		_ = pos
		_ = rot
		cnt++
	}
	assert.Equal(t, 2, cnt)

	q = w.Query(posID)
	cnt = 0
	for q.Next() {
		ent := q.Entity()
		pos := (*position)(q.Get(posID))
		_ = ent
		_ = pos
		cnt++
	}
	assert.Equal(t, 3, cnt)

	q = w.Query(rotID)
	cnt = 0
	for q.Next() {
		ent := q.Entity()
		rot := (*rotation)(q.Get(rotID))
		_ = ent
		_ = rot
		hasPos := q.Has(posID)
		_ = hasPos
		cnt++
	}
	assert.Equal(t, 4, cnt)

	assert.Panics(t, func() { q.Next() })
}

func TestQuery1(t *testing.T) {
	w := NewWorld()

	id0 := ComponentID[testStruct0](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()

	w.Assign(e0, id0, &testStruct0{1})
	w.Assign(e1, id0, &testStruct0{2})

	cnt := 0
	query := Query1[testStruct0](&w)
	for query.Next() {
		c0 := query.Get1()
		assert.Equal(t, cnt+1, int(c0.val))
		cnt++
	}
	assert.Equal(t, 2, cnt)
}

func TestQuery2(t *testing.T) {
	w := NewWorld()

	id0 := ComponentID[testStruct0](&w)
	id1 := ComponentID[testStruct1](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()

	w.Assign(e0, id0, &testStruct0{1})
	w.Assign(e1, id0, &testStruct0{2})

	w.Assign(e0, id1, &testStruct1{2})
	w.Assign(e1, id1, &testStruct1{3})

	cnt := 0
	query := Query2[testStruct0, testStruct1](&w)
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

	id0 := ComponentID[testStruct0](&w)
	id1 := ComponentID[testStruct1](&w)
	id2 := ComponentID[testStruct2](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()

	w.Assign(e0, id0, &testStruct0{1})
	w.Assign(e1, id0, &testStruct0{2})

	w.Assign(e0, id1, &testStruct1{2})
	w.Assign(e1, id1, &testStruct1{3})

	w.Assign(e0, id2, &testStruct2{3})
	w.Assign(e1, id2, &testStruct2{4})

	cnt := 0
	query := Query3[testStruct0, testStruct1, testStruct2](&w)
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

	id0 := ComponentID[testStruct0](&w)
	id1 := ComponentID[testStruct1](&w)
	id2 := ComponentID[testStruct2](&w)
	id3 := ComponentID[testStruct3](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()

	w.Assign(e0, id0, &testStruct0{1})
	w.Assign(e1, id0, &testStruct0{2})

	w.Assign(e0, id1, &testStruct1{2})
	w.Assign(e1, id1, &testStruct1{3})

	w.Assign(e0, id2, &testStruct2{3})
	w.Assign(e1, id2, &testStruct2{4})

	w.Assign(e0, id3, &testStruct3{4})
	w.Assign(e1, id3, &testStruct3{5})

	cnt := 0
	query := Query4[testStruct0, testStruct1, testStruct2, testStruct3](&w)
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

	id0 := ComponentID[testStruct0](&w)
	id1 := ComponentID[testStruct1](&w)
	id2 := ComponentID[testStruct2](&w)
	id3 := ComponentID[testStruct3](&w)
	id4 := ComponentID[testStruct4](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()

	w.Assign(e0, id0, &testStruct0{1})
	w.Assign(e1, id0, &testStruct0{2})

	w.Assign(e0, id1, &testStruct1{2})
	w.Assign(e1, id1, &testStruct1{3})

	w.Assign(e0, id2, &testStruct2{3})
	w.Assign(e1, id2, &testStruct2{4})

	w.Assign(e0, id3, &testStruct3{4})
	w.Assign(e1, id3, &testStruct3{5})

	w.Assign(e0, id4, &testStruct4{5})
	w.Assign(e1, id4, &testStruct4{6})

	cnt := 0
	query := Query5[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4](&w)
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
