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

	e := mut.New()
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

	e = mut.New()
	assert.True(t, map0.Has(e))

	e = mut.NewWith(&testStruct0{})

	assert.Panics(t, func() {
		mut.NewWith(&testStruct0{}, e)
	})

	s0 = mut.Get(e)
	assert.NotNil(t, s0)

	mut.NewBatch(2)

	q := mut.NewBatchQ(2)
	q.Close()

	mut.GetUnchecked(e)

	cnt := mut.RemoveEntities(true)
	assert.Equal(t, 8, cnt)
	cnt = mut.RemoveEntities(false)
	assert.Equal(t, 0, cnt)

	assert.Panics(t, func() { mut.Get(e) })

	target := w.NewEntity()
	mut = NewMap1[testStruct0](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap1[testRelationA](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(&testRelationA{}, target)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))
}

func TestMap2Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap2[testStruct0, testStruct1](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.New()
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

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(&testStruct0{}, &testStruct1{})

	assert.Panics(t, func() {
		mut.NewWith(&testStruct0{}, &testStruct1{}, e)
	})

	s0, s1 = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.NewBatch(2)

	q := mut.NewBatchQ(2)
	q.Close()

	mut.GetUnchecked(e)

	mut.RemoveEntities(true)
	mut.RemoveEntities(false)

	assert.Panics(t, func() { mut.Get(e) })

	target := w.NewEntity()
	mut = NewMap2[
		testStruct0, testStruct1,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap2[
		testRelationA, testStruct1,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(&testRelationA{}, &testStruct1{}, target)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))
}

func TestMap3Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap3[
		testStruct0, testStruct1, testStruct2,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.New()
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

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(&testStruct0{}, &testStruct1{}, &testStruct2{})

	assert.Panics(t, func() {
		mut.NewWith(&testStruct0{}, &testStruct1{}, &testStruct2{}, e)
	})

	s0, s1, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.NewBatch(2)

	q := mut.NewBatchQ(2)
	q.Close()

	mut.GetUnchecked(e)

	mut.RemoveEntities(true)
	mut.RemoveEntities(false)

	assert.Panics(t, func() { mut.Get(e) })

	target := w.NewEntity()
	mut = NewMap3[
		testStruct0, testStruct1, testStruct2,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap3[
		testRelationA, testStruct1, testStruct2,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(&testRelationA{}, &testStruct1{}, &testStruct2{}, target)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))
}

func TestMap4Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap4[
		testStruct0, testStruct1, testStruct2, testStruct3,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.New()
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

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{})

	assert.Panics(t, func() {
		mut.NewWith(
			&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
			e,
		)
	})

	s0, s1, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.NewBatch(2)

	q := mut.NewBatchQ(2)
	q.Close()

	mut.GetUnchecked(e)

	mut.RemoveEntities(true)
	mut.RemoveEntities(false)

	assert.Panics(t, func() { mut.Get(e) })

	target := w.NewEntity()
	mut = NewMap4[
		testStruct0, testStruct1, testStruct2, testStruct3,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap4[
		testRelationA, testStruct1, testStruct2, testStruct3,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(&testRelationA{}, &testStruct1{}, &testStruct2{}, &testStruct3{}, target)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))
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

	e := mut.New()
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

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{},
	)

	assert.Panics(t, func() {
		mut.NewWith(
			&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
			&testStruct4{},
			e,
		)
	})

	s0, s1, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.NewBatch(2)

	q := mut.NewBatchQ(2)
	q.Close()

	mut.GetUnchecked(e)

	mut.RemoveEntities(true)
	mut.RemoveEntities(false)

	assert.Panics(t, func() { mut.Get(e) })

	target := w.NewEntity()
	mut = NewMap5[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap5[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(
		&testRelationA{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{},
		target,
	)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))
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

	e := mut.New()
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

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{},
	)

	assert.Panics(t, func() {
		mut.NewWith(
			&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
			&testStruct4{}, &testStruct5{},
			e,
		)
	})

	s0, s1, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.NewBatch(2)

	q := mut.NewBatchQ(2)
	q.Close()

	mut.GetUnchecked(e)

	mut.RemoveEntities(true)
	mut.RemoveEntities(false)

	assert.Panics(t, func() { mut.Get(e) })

	target := w.NewEntity()
	mut = NewMap6[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap6[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(
		&testRelationA{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{},
		target,
	)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))
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

	e := mut.New()
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

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{},
	)

	assert.Panics(t, func() {
		mut.NewWith(
			&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
			&testStruct4{}, &testStruct5{}, &testStruct6{},
			e,
		)
	})

	s0, s1, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.NewBatch(2)

	q := mut.NewBatchQ(2)
	q.Close()

	mut.GetUnchecked(e)

	mut.RemoveEntities(true)
	mut.RemoveEntities(false)

	assert.Panics(t, func() { mut.Get(e) })

	target := w.NewEntity()
	mut = NewMap7[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap7[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(
		&testRelationA{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{},
		target,
	)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))
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

	e := mut.New()
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

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
	)

	assert.Panics(t, func() {
		mut.NewWith(
			&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
			&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
			e,
		)
	})

	s0, s1, _, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.NewBatch(2)

	q := mut.NewBatchQ(2)
	q.Close()

	mut.GetUnchecked(e)

	mut.RemoveEntities(true)
	mut.RemoveEntities(false)

	assert.Panics(t, func() { mut.Get(e) })

	target := w.NewEntity()
	mut = NewMap8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap8[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(
		&testRelationA{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		target,
	)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))
}
