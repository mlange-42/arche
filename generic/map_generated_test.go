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
	mut.Assign(e, &testStruct0{val: 42})
	assert.True(t, map0.Has(e))
	assert.Equal(t, testStruct0{val: 42}, *map0.Get(e))

	e = mut.New()
	assert.True(t, map0.Has(e))

	e = mut.NewWith(&testStruct0{val: 23})

	assert.Panics(t, func() {
		mut.NewWith(&testStruct0{}, e)
	})

	s0 = mut.Get(e)
	assert.NotNil(t, s0)
	assert.Equal(t, testStruct0{val: 23}, *s0)

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

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap1[testRelationA](&w)
	mut4 := NewMap1[testStruct0](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
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
	mut.Assign(e, &testStruct0{val: 23}, &testStruct1{val: 42})
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

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap2[
		testRelationA, testStruct1,
	](&w)
	mut4 := NewMap2[
		testStruct0, testStruct1,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
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

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap3[
		testRelationA, testStruct1, testStruct2,
	](&w)
	mut4 := NewMap3[
		testStruct0, testStruct1, testStruct2,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
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

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap4[
		testRelationA, testStruct1, testStruct2, testStruct3,
	](&w)
	mut4 := NewMap4[
		testStruct0, testStruct1, testStruct2, testStruct3,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
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

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap5[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4,
	](&w)
	mut4 := NewMap5[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
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

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap6[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5,
	](&w)
	mut4 := NewMap6[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
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

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap7[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6,
	](&w)
	mut4 := NewMap7[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
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

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap8[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	](&w)
	mut4 := NewMap8[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
}

func TestMap9Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap9[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.New()
	s0, s1, _, _, _, _, _, _, _ := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1, _, _, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{},
	)

	assert.Panics(t, func() {
		mut.NewWith(
			&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
			&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
			&testStruct8{},
			e,
		)
	})

	s0, s1, _, _, _, _, _, _, _ = mut.Get(e)
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
	mut = NewMap9[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap9[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(
		&testRelationA{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{},
		target,
	)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap9[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8,
	](&w)
	mut4 := NewMap9[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
}

func TestMap10Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap10[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.New()
	s0, s1, _, _, _, _, _, _, _, _ := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1, _, _, _, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{}, &testStruct9{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{}, &testStruct9{},
	)

	assert.Panics(t, func() {
		mut.NewWith(
			&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
			&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
			&testStruct8{}, &testStruct9{},
			e,
		)
	})

	s0, s1, _, _, _, _, _, _, _, _ = mut.Get(e)
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
	mut = NewMap10[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap10[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(
		&testRelationA{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{}, &testStruct9{},
		target,
	)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap10[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9,
	](&w)
	mut4 := NewMap10[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
}

func TestMap11Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap11[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9, testStruct10,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.New()
	s0, s1, _, _, _, _, _, _, _, _, _ := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1, _, _, _, _, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{}, &testStruct9{}, &testStruct10{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{}, &testStruct9{}, &testStruct10{},
	)

	assert.Panics(t, func() {
		mut.NewWith(
			&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
			&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
			&testStruct8{}, &testStruct9{}, &testStruct10{},
			e,
		)
	})

	s0, s1, _, _, _, _, _, _, _, _, _ = mut.Get(e)
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
	mut = NewMap11[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9, testStruct10,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap11[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9, testStruct10,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(
		&testRelationA{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{}, &testStruct9{}, &testStruct10{},
		target,
	)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap11[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9, testStruct10,
	](&w)
	mut4 := NewMap11[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9, testStruct10,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
}

func TestMap12Generated(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMap12[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9, testStruct10, testStruct11,
	](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.New()
	s0, s1, _, _, _, _, _, _, _, _, _, _ := mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mut.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1, _, _, _, _, _, _, _, _, _, _ = mut.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	e = w.NewEntity()
	mut.Assign(e,
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{}, &testStruct9{}, &testStruct10{}, &testStruct11{},
	)
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.New()
	assert.True(t, map0.Has(e))
	assert.True(t, map1.Has(e))

	e = mut.NewWith(
		&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{}, &testStruct9{}, &testStruct10{}, &testStruct11{},
	)

	assert.Panics(t, func() {
		mut.NewWith(
			&testStruct0{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
			&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
			&testStruct8{}, &testStruct9{}, &testStruct10{}, &testStruct11{},
			e,
		)
	})

	s0, s1, _, _, _, _, _, _, _, _, _, _ = mut.Get(e)
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
	mut = NewMap12[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9, testStruct10, testStruct11,
	](&w)

	assert.Panics(t, func() { mut.New(target) })
	assert.Panics(t, func() { mut.NewBatch(5, target) })
	assert.Panics(t, func() { mut.NewBatchQ(5, target) })

	mut2 := NewMap12[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9, testStruct10, testStruct11,
	](&w, T[testRelationA]())

	mut2.New(target)
	mut2.NewBatch(5, target)
	q2 := mut2.NewBatchQ(5, target)
	assert.Equal(t, 5, q2.Count())
	q2.Close()

	e = mut2.NewWith(
		&testRelationA{}, &testStruct1{}, &testStruct2{}, &testStruct3{},
		&testStruct4{}, &testStruct5{}, &testStruct6{}, &testStruct7{},
		&testStruct8{}, &testStruct9{}, &testStruct10{}, &testStruct11{},
		target,
	)
	mapper := NewMap[testRelationA](&w)
	assert.Equal(t, target, mapper.GetRelation(e))
	assert.Equal(t, target, mapper.GetRelationUnchecked(e))

	// === Batch operations ====

	w.Batch().RemoveEntities(ecs.All())

	mut3 := NewMap12[
		testRelationA, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9, testStruct10, testStruct11,
	](&w)
	mut4 := NewMap12[
		testStruct0, testStruct1, testStruct2, testStruct3,
		testStruct4, testStruct5, testStruct6, testStruct7,
		testStruct8, testStruct9, testStruct10, testStruct11,
	](&w, T[testRelationB]())
	rel2ID := ecs.ComponentID[testRelationB](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	assert.Panics(t, func() {
		mut3.Add(e1, e2)
	})
	mut2.Add(e1, e2)

	assert.Panics(t, func() {
		mut3.Remove(e1, e2)
	})

	mut.Add(e3)
	w.Add(e3, rel2ID)
	mut4.Remove(e3, e1)

	w.Batch().RemoveEntities(ecs.All())

	e1 = w.NewEntity()
	e2 = w.NewEntity()
	_ = w.NewEntity()
	assert.Panics(t, func() {
		mut3.AddBatch(ecs.All(), e2)
	})
	assert.Panics(t, func() {
		_ = mut3.AddBatchQ(ecs.All(), e2)
	})
	mut2.AddBatch(ecs.All(), e2)
	mut2.RemoveBatch(ecs.All())
	query := mut2.AddBatchQ(ecs.All(), e2)
	assert.Equal(t, 3, query.Count())
	query.Close()

	assert.Panics(t, func() {
		mut3.RemoveBatch(ecs.All(), e2)
	})

	assert.Panics(t, func() {
		_ = mut3.RemoveBatchQ(ecs.All(), e2)
	})

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	mut.AddBatch(ecs.All())
	w.Batch().Add(ecs.All(), rel2ID)
	mut4.RemoveBatch(ecs.All(), e1)

	w.Batch().RemoveEntities(ecs.All())
	e1 = w.NewEntity()
	_ = w.NewEntity()
	_ = w.NewEntity()

	query2 := mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	w.Batch().Add(ecs.All(), rel2ID)

	query3 := mut4.RemoveBatchQ(ecs.All(), e1)
	assert.Equal(t, 3, query3.Count())
	query2.Close()

	query2 = mut.AddBatchQ(ecs.All())
	assert.Equal(t, 3, query2.Count())
	query2.Close()

	query3 = mut4.RemoveBatchQ(ecs.All())
	assert.Equal(t, 3, query3.Count())
	query2.Close()
}
