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

	w.Assign(e0, ecs.Component{ID: 0, Comp: &testStruct0{1}})
	w.Assign(e0, ecs.Component{ID: 1, Comp: &testStruct1{1}})

	w.Assign(e1, ecs.Component{ID: 0, Comp: &testStruct0{2}})
	w.Assign(e1, ecs.Component{ID: 1, Comp: &testStruct1{2}})
	w.Assign(e1, ecs.Component{ID: 2, Comp: &testStruct2{1, 1}})
	w.Assign(e1, ecs.Component{ID: 8, Comp: &testStruct8{}})

	w.Assign(e2, ecs.Component{ID: 0, Comp: &testStruct0{3}})
	w.Assign(e2, ecs.Component{ID: 1, Comp: &testStruct1{3}})
	w.Assign(e2, ecs.Component{ID: 2, Comp: &testStruct2{1, 1}})
	w.Assign(e2, ecs.Component{ID: 9, Comp: &testStruct9{}})

	query2 := NewFilter2[testStruct0, testStruct1]().Query(&w)
	cnt := 0
	for query2.Next() {
		_, _ = query2.Get()
		cnt++
	}
	assert.Equal(t, 3, cnt)

	query2 = NewFilter2[testStruct0, testStruct1]().Without(T[testStruct9]()).Query(&w)
	cnt = 0
	for query2.Next() {
		_, _ = query2.Get()
		cnt++
	}
	assert.Equal(t, 2, cnt)

	query2 = NewFilter2[testStruct0, testStruct1]().Without(T[testStruct8](), T[testStruct9]()).Query(&w)
	cnt = 0
	for query2.Next() {
		_, _ = query2.Get()
		cnt++
	}
	assert.Equal(t, 1, cnt)

	query2 = NewFilter2[testStruct0, testStruct1]().With(T[testStruct2]()).Without(T[testStruct9]()).Query(&w)
	cnt = 0
	for query2.Next() {
		_, _ = query2.Get()
		cnt++
	}
	assert.Equal(t, 1, cnt)

	query3 := NewFilter3[testStruct0, testStruct1, testStruct9]().Query(&w)
	cnt = 0
	for query3.Next() {
		_, _, _ = query3.Get()
		cnt++
	}
	assert.Equal(t, 1, cnt)

	query3 = NewFilter3[testStruct0, testStruct1, testStruct9]().Optional(T[testStruct9]()).Query(&w)
	cnt = 0
	for query3.Next() {
		_, _, _ = query3.Get()
		cnt++
	}
	assert.Equal(t, 3, cnt)
}

func TestQuery0(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)
	relID := ecs.ComponentID[testRelationA](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, ecs.Component{ID: 0, Comp: &testStruct0{1}})
	w.Assign(e0, ecs.Component{ID: 8, Comp: &testStruct8{}})

	w.Assign(e1, ecs.Component{ID: 0, Comp: &testStruct0{2}})

	w.Assign(e2, ecs.Component{ID: 0, Comp: &testStruct0{2}})
	w.Assign(e2, ecs.Component{ID: 9, Comp: &testStruct9{}})

	cnt := 0
	filter :=
		NewFilter0().
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	f := ecs.All(8).Without(9)
	assert.Equal(t, f, filter.Filter(&w))
	filter.Register(&w)

	query := filter.Query(&w)
	for query.Next() {
		_ = query.Entity()
		assert.Panics(t, func() { query.Relation() })
		cnt++
	}
	assert.Equal(t, 1, cnt)

	filter.Unregister(&w)
	assert.Panics(t, func() { filter.Unregister(&w) })

	targ := w.NewEntity(0)

	w.Add(e0, relID)
	w.Add(e1, relID)
	w.Add(e2, relID)

	w.SetRelation(e0, relID, targ)

	filter2 :=
		NewFilter0().
			With(T[testRelationA]()).
			WithRelation(T[testRelationA](), targ)

	q := filter2.Query(&w)
	assert.Equal(t, 1, q.Count())
	for q.Next() {
		trg := q.Relation()
		assert.Equal(t, targ, trg)
	}

	filter2.Register(&w)
	assert.Panics(t, func() { filter2.With(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.Without(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.WithRelation(T[testRelationA](), targ) })
	filter2.Unregister(&w)
}

func TestQuery1(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)
	relID := ecs.ComponentID[testRelationA](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, ecs.Component{ID: 0, Comp: &testStruct0{1}})

	w.Assign(e1, ecs.Component{ID: 0, Comp: &testStruct0{2}})

	w.Assign(e0, ecs.Component{ID: 8, Comp: &testStruct8{9, 0, 0, 0, 0, 0, 0, 0}})

	w.Assign(e2, ecs.Component{ID: 0, Comp: &testStruct0{0}})
	w.Assign(e2, ecs.Component{ID: 9, Comp: &testStruct9{}})

	cnt := 0
	filter :=
		NewFilter1[testStruct0]().
			Optional(T[testStruct9]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	f := ecs.All(0, 8).Without(9)
	assert.Equal(t, f, filter.Filter(&w))
	filter.Register(&w)

	query := filter.Query(&w)
	for query.Next() {
		c0 := query.Get()
		assert.Equal(t, cnt+1, int(c0.val))
		assert.Panics(t, func() { query.Relation() })
		cnt++
	}
	assert.Equal(t, 1, cnt)

	filter.Unregister(&w)
	assert.Panics(t, func() { filter.Unregister(&w) })

	targ := w.NewEntity(0)

	w.Add(e0, relID)
	w.Add(e1, relID)
	w.Add(e2, relID)

	w.SetRelation(e0, relID, targ)

	filter2 :=
		NewFilter1[testRelationA]().
			WithRelation(T[testRelationA](), targ)

	q := filter2.Query(&w)
	assert.Equal(t, 1, q.Count())
	for q.Next() {
		trg := q.Relation()
		assert.Equal(t, targ, trg)
	}

	filter2.Register(&w)
	assert.Panics(t, func() { filter2.Optional(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.With(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.Without(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.WithRelation(T[testRelationA](), targ) })
	filter2.Unregister(&w)
}

func TestQuery2(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)
	relID := ecs.ComponentID[testRelationA](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, ecs.Component{ID: 0, Comp: &testStruct0{1}})
	w.Assign(e1, ecs.Component{ID: 0, Comp: &testStruct0{2}})
	w.Assign(e2, ecs.Component{ID: 0, Comp: &testStruct0{3}})

	w.Assign(e0, ecs.Component{ID: 1, Comp: &testStruct1{2}})
	w.Assign(e1, ecs.Component{ID: 1, Comp: &testStruct1{3}})
	w.Assign(e2, ecs.Component{ID: 1, Comp: &testStruct1{4}})

	w.Assign(e0, ecs.Component{ID: 8, Comp: &testStruct8{9, 0, 0, 0, 0, 0, 0, 0}})

	w.Assign(e2, ecs.Component{ID: 9, Comp: &testStruct9{}})

	filter :=
		NewFilter2[testStruct0, testStruct1]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	f := ecs.All(0, 8).Without(9)
	assert.Equal(t, f, filter.Filter(&w))
	filter.Register(&w)

	for i := 0; i < 10; i++ {
		cnt := 0
		q := filter.Query(&w)
		for q.Next() {
			c1, c2 := q.Get()
			assert.Equal(t, cnt+1, int(c1.val))
			assert.Equal(t, cnt+2, int(c2.val))
			assert.Panics(t, func() { q.Relation() })
			cnt++
		}
		assert.Equal(t, 1, cnt)
	}

	filter.Unregister(&w)
	assert.Panics(t, func() { filter.Unregister(&w) })

	targ := w.NewEntity(0)

	w.Add(e0, relID)
	w.Add(e1, relID)
	w.Add(e2, relID)

	w.SetRelation(e0, relID, targ)

	filter2 :=
		NewFilter2[testStruct0, testRelationA]().
			WithRelation(T[testRelationA](), targ)

	q := filter2.Query(&w)
	assert.Equal(t, 1, q.Count())
	for q.Next() {
		trg := q.Relation()
		assert.Equal(t, targ, trg)
	}

	filter2.Register(&w)
	assert.Panics(t, func() { filter2.Optional(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.With(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.Without(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.WithRelation(T[testRelationA](), targ) })
	filter2.Unregister(&w)
}

func TestQuery3(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)
	relID := ecs.ComponentID[testRelationA](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, ecs.Component{ID: 0, Comp: &testStruct0{1}})
	w.Assign(e1, ecs.Component{ID: 0, Comp: &testStruct0{2}})
	w.Assign(e2, ecs.Component{ID: 0, Comp: &testStruct0{3}})

	w.Assign(e0, ecs.Component{ID: 1, Comp: &testStruct1{2}})
	w.Assign(e1, ecs.Component{ID: 1, Comp: &testStruct1{3}})
	w.Assign(e2, ecs.Component{ID: 1, Comp: &testStruct1{4}})

	w.Assign(e0, ecs.Component{ID: 2, Comp: &testStruct2{3, 0}})
	w.Assign(e1, ecs.Component{ID: 2, Comp: &testStruct2{4, 0}})
	w.Assign(e2, ecs.Component{ID: 2, Comp: &testStruct2{5, 0}})

	w.Assign(e0, ecs.Component{ID: 8, Comp: &testStruct8{9, 0, 0, 0, 0, 0, 0, 0}})

	w.Assign(e2, ecs.Component{ID: 9, Comp: &testStruct9{}})

	cnt := 0
	filter :=
		NewFilter3[testStruct0, testStruct1, testStruct2]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	f := ecs.All(0, 2, 8).Without(9)
	assert.Equal(t, f, filter.Filter(&w))
	filter.Register(&w)

	query := filter.Query(&w)
	for query.Next() {
		c1, c2, c3 := query.Get()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Panics(t, func() { query.Relation() })
		cnt++
	}
	assert.Equal(t, 1, cnt)

	filter.Unregister(&w)
	assert.Panics(t, func() { filter.Unregister(&w) })

	targ := w.NewEntity(0)

	w.Add(e0, relID)
	w.Add(e1, relID)
	w.Add(e2, relID)

	w.SetRelation(e0, relID, targ)

	filter2 :=
		NewFilter3[testStruct0, testStruct1, testRelationA]().
			WithRelation(T[testRelationA](), targ)

	q := filter2.Query(&w)
	assert.Equal(t, 1, q.Count())
	for q.Next() {
		trg := q.Relation()
		assert.Equal(t, targ, trg)
	}

	filter2.Register(&w)
	assert.Panics(t, func() { filter2.Optional(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.With(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.Without(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.WithRelation(T[testRelationA](), targ) })
	filter2.Unregister(&w)
}

func TestQuery4(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)
	relID := ecs.ComponentID[testRelationA](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, ecs.Component{ID: 0, Comp: &testStruct0{1}})
	w.Assign(e1, ecs.Component{ID: 0, Comp: &testStruct0{2}})
	w.Assign(e2, ecs.Component{ID: 0, Comp: &testStruct0{2}})

	w.Assign(e0, ecs.Component{ID: 1, Comp: &testStruct1{2}})
	w.Assign(e1, ecs.Component{ID: 1, Comp: &testStruct1{3}})
	w.Assign(e2, ecs.Component{ID: 1, Comp: &testStruct1{3}})

	w.Assign(e0, ecs.Component{ID: 2, Comp: &testStruct2{3, 0}})
	w.Assign(e1, ecs.Component{ID: 2, Comp: &testStruct2{4, 0}})
	w.Assign(e2, ecs.Component{ID: 2, Comp: &testStruct2{4, 0}})

	w.Assign(e0, ecs.Component{ID: 3, Comp: &testStruct3{4, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 3, Comp: &testStruct3{5, 0, 0}})
	w.Assign(e2, ecs.Component{ID: 3, Comp: &testStruct3{5, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 8, Comp: &testStruct8{9, 0, 0, 0, 0, 0, 0, 0}})

	w.Assign(e2, ecs.Component{ID: 9, Comp: &testStruct9{}})

	cnt := 0
	filter :=
		NewFilter4[testStruct0, testStruct1, testStruct2, testStruct3]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	f := ecs.All(0, 2, 3, 8).Without(9)
	assert.Equal(t, f, filter.Filter(&w))
	filter.Register(&w)

	query := filter.Query(&w)
	for query.Next() {
		c1, c2, c3, c4 := query.Get()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))
		assert.Panics(t, func() { query.Relation() })
		cnt++
	}
	assert.Equal(t, 1, cnt)

	filter.Unregister(&w)
	assert.Panics(t, func() { filter.Unregister(&w) })

	targ := w.NewEntity(0)

	w.Add(e0, relID)
	w.Add(e1, relID)
	w.Add(e2, relID)

	w.SetRelation(e0, relID, targ)

	filter2 :=
		NewFilter4[
			testStruct0, testStruct1, testStruct2, testRelationA,
		]().
			WithRelation(T[testRelationA](), targ)

	q := filter2.Query(&w)
	assert.Equal(t, 1, q.Count())
	for q.Next() {
		trg := q.Relation()
		assert.Equal(t, targ, trg)
	}

	filter2.Register(&w)
	assert.Panics(t, func() { filter2.Optional(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.With(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.Without(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.WithRelation(T[testRelationA](), targ) })
	filter2.Unregister(&w)
}

func TestQuery5(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)
	relID := ecs.ComponentID[testRelationA](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, ecs.Component{ID: 0, Comp: &testStruct0{1}})
	w.Assign(e1, ecs.Component{ID: 0, Comp: &testStruct0{2}})
	w.Assign(e2, ecs.Component{ID: 0, Comp: &testStruct0{2}})

	w.Assign(e0, ecs.Component{ID: 1, Comp: &testStruct1{2}})
	w.Assign(e1, ecs.Component{ID: 1, Comp: &testStruct1{3}})
	w.Assign(e2, ecs.Component{ID: 1, Comp: &testStruct1{3}})

	w.Assign(e0, ecs.Component{ID: 2, Comp: &testStruct2{3, 0}})
	w.Assign(e1, ecs.Component{ID: 2, Comp: &testStruct2{4, 0}})
	w.Assign(e2, ecs.Component{ID: 2, Comp: &testStruct2{4, 0}})

	w.Assign(e0, ecs.Component{ID: 3, Comp: &testStruct3{4, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 3, Comp: &testStruct3{5, 0, 0}})
	w.Assign(e2, ecs.Component{ID: 3, Comp: &testStruct3{5, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 4, Comp: &testStruct4{5, 0, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 4, Comp: &testStruct4{6, 0, 0, 0}})
	w.Assign(e2, ecs.Component{ID: 4, Comp: &testStruct4{6, 0, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 8, Comp: &testStruct8{9, 0, 0, 0, 0, 0, 0, 0}})

	w.Assign(e2, ecs.Component{ID: 9, Comp: &testStruct9{}})

	cnt := 0
	filter :=
		NewFilter5[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	f := ecs.All(0, 2, 3, 4, 8).Without(9)
	assert.Equal(t, f, filter.Filter(&w))
	filter.Register(&w)

	query := filter.Query(&w)
	for query.Next() {
		c1, c2, c3, c4, c5 := query.Get()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))
		assert.Equal(t, cnt+5, int(c5.val))
		assert.Panics(t, func() { query.Relation() })
		cnt++
	}
	assert.Equal(t, 1, cnt)

	filter.Unregister(&w)
	assert.Panics(t, func() { filter.Unregister(&w) })

	targ := w.NewEntity(0)

	w.Add(e0, relID)
	w.Add(e1, relID)
	w.Add(e2, relID)

	w.SetRelation(e0, relID, targ)

	filter2 :=
		NewFilter5[
			testStruct0, testStruct1, testStruct2, testStruct3,
			testRelationA,
		]().
			WithRelation(T[testRelationA](), targ)

	q := filter2.Query(&w)
	assert.Equal(t, 1, q.Count())
	for q.Next() {
		trg := q.Relation()
		assert.Equal(t, targ, trg)
	}

	filter2.Register(&w)
	assert.Panics(t, func() { filter2.Optional(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.With(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.Without(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.WithRelation(T[testRelationA](), targ) })
	filter2.Unregister(&w)
}

func TestQuery6(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)
	relID := ecs.ComponentID[testRelationA](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, ecs.Component{ID: 0, Comp: &testStruct0{1}})
	w.Assign(e1, ecs.Component{ID: 0, Comp: &testStruct0{2}})
	w.Assign(e2, ecs.Component{ID: 0, Comp: &testStruct0{2}})

	w.Assign(e0, ecs.Component{ID: 1, Comp: &testStruct1{2}})
	w.Assign(e1, ecs.Component{ID: 1, Comp: &testStruct1{3}})
	w.Assign(e2, ecs.Component{ID: 1, Comp: &testStruct1{3}})

	w.Assign(e0, ecs.Component{ID: 2, Comp: &testStruct2{3, 0}})
	w.Assign(e1, ecs.Component{ID: 2, Comp: &testStruct2{4, 0}})
	w.Assign(e2, ecs.Component{ID: 2, Comp: &testStruct2{4, 0}})

	w.Assign(e0, ecs.Component{ID: 3, Comp: &testStruct3{4, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 3, Comp: &testStruct3{5, 0, 0}})
	w.Assign(e2, ecs.Component{ID: 3, Comp: &testStruct3{5, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 4, Comp: &testStruct4{5, 0, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 4, Comp: &testStruct4{6, 0, 0, 0}})
	w.Assign(e2, ecs.Component{ID: 4, Comp: &testStruct4{6, 0, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 5, Comp: &testStruct5{6, 0, 0, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 5, Comp: &testStruct5{7, 0, 0, 0, 0}})
	w.Assign(e2, ecs.Component{ID: 5, Comp: &testStruct5{7, 0, 0, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 8, Comp: &testStruct8{9, 0, 0, 0, 0, 0, 0, 0}})

	w.Assign(e2, ecs.Component{ID: 9, Comp: &testStruct9{}})

	cnt := 0
	filter :=
		NewFilter6[testStruct0, testStruct1, testStruct2, testStruct3, testStruct4, testStruct5]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	f := ecs.All(0, 2, 3, 4, 5, 8).Without(9)
	assert.Equal(t, f, filter.Filter(&w))
	filter.Register(&w)

	query := filter.Query(&w)
	for query.Next() {
		c1, c2, c3, c4, c5, c6 := query.Get()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))
		assert.Equal(t, cnt+5, int(c5.val))
		assert.Equal(t, cnt+6, int(c6.val))
		assert.Panics(t, func() { query.Relation() })
		cnt++
	}
	assert.Equal(t, 1, cnt)

	filter.Unregister(&w)
	assert.Panics(t, func() { filter.Unregister(&w) })

	targ := w.NewEntity(0)

	w.Add(e0, relID)
	w.Add(e1, relID)
	w.Add(e2, relID)

	w.SetRelation(e0, relID, targ)

	filter2 :=
		NewFilter6[
			testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testRelationA,
		]().
			WithRelation(T[testRelationA](), targ)

	q := filter2.Query(&w)
	assert.Equal(t, 1, q.Count())
	for q.Next() {
		trg := q.Relation()
		assert.Equal(t, targ, trg)
	}

	filter2.Register(&w)
	assert.Panics(t, func() { filter2.Optional(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.With(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.Without(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.WithRelation(T[testRelationA](), targ) })
	filter2.Unregister(&w)
}

func TestQuery7(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)
	relID := ecs.ComponentID[testRelationA](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, ecs.Component{ID: 0, Comp: &testStruct0{1}})
	w.Assign(e1, ecs.Component{ID: 0, Comp: &testStruct0{2}})
	w.Assign(e2, ecs.Component{ID: 0, Comp: &testStruct0{2}})

	w.Assign(e0, ecs.Component{ID: 1, Comp: &testStruct1{2}})
	w.Assign(e1, ecs.Component{ID: 1, Comp: &testStruct1{3}})
	w.Assign(e2, ecs.Component{ID: 1, Comp: &testStruct1{3}})

	w.Assign(e0, ecs.Component{ID: 2, Comp: &testStruct2{3, 0}})
	w.Assign(e1, ecs.Component{ID: 2, Comp: &testStruct2{4, 0}})
	w.Assign(e2, ecs.Component{ID: 2, Comp: &testStruct2{4, 0}})

	w.Assign(e0, ecs.Component{ID: 3, Comp: &testStruct3{4, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 3, Comp: &testStruct3{5, 0, 0}})
	w.Assign(e2, ecs.Component{ID: 3, Comp: &testStruct3{5, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 4, Comp: &testStruct4{5, 0, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 4, Comp: &testStruct4{6, 0, 0, 0}})
	w.Assign(e2, ecs.Component{ID: 4, Comp: &testStruct4{6, 0, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 5, Comp: &testStruct5{6, 0, 0, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 5, Comp: &testStruct5{7, 0, 0, 0, 0}})
	w.Assign(e2, ecs.Component{ID: 5, Comp: &testStruct5{7, 0, 0, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 6, Comp: &testStruct6{7, 0, 0, 0, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 6, Comp: &testStruct6{8, 0, 0, 0, 0, 0}})
	w.Assign(e2, ecs.Component{ID: 6, Comp: &testStruct6{8, 0, 0, 0, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 8, Comp: &testStruct8{9, 0, 0, 0, 0, 0, 0, 0}})

	w.Assign(e2, ecs.Component{ID: 9, Comp: &testStruct9{}})

	cnt := 0
	filter :=
		NewFilter7[
			testStruct0, testStruct1, testStruct2, testStruct3, testStruct4,
			testStruct5, testStruct6,
		]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	f := ecs.All(0, 2, 3, 4, 5, 6, 8).Without(9)
	assert.Equal(t, f, filter.Filter(&w))
	filter.Register(&w)

	query := filter.Query(&w)
	for query.Next() {
		c1, c2, c3, c4, c5, c6, c7 := query.Get()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))
		assert.Equal(t, cnt+5, int(c5.val))
		assert.Equal(t, cnt+6, int(c6.val))
		assert.Equal(t, cnt+7, int(c7.val))
		assert.Panics(t, func() { query.Relation() })
		cnt++
	}
	assert.Equal(t, 1, cnt)

	filter.Unregister(&w)
	assert.Panics(t, func() { filter.Unregister(&w) })

	targ := w.NewEntity(0)

	w.Add(e0, relID)
	w.Add(e1, relID)
	w.Add(e2, relID)

	w.SetRelation(e0, relID, targ)

	filter2 :=
		NewFilter7[
			testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testRelationA,
		]().
			WithRelation(T[testRelationA](), targ)

	q := filter2.Query(&w)
	assert.Equal(t, 1, q.Count())
	for q.Next() {
		trg := q.Relation()
		assert.Equal(t, targ, trg)
	}

	filter2.Register(&w)
	assert.Panics(t, func() { filter2.Optional(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.With(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.Without(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.WithRelation(T[testRelationA](), targ) })
	filter2.Unregister(&w)
}

func TestQuery8(t *testing.T) {
	w := ecs.NewWorld()

	registerAll(&w)
	relID := ecs.ComponentID[testRelationA](&w)

	e0 := w.NewEntity()
	e1 := w.NewEntity()
	e2 := w.NewEntity()

	w.Assign(e0, ecs.Component{ID: 0, Comp: &testStruct0{1}})
	w.Assign(e1, ecs.Component{ID: 0, Comp: &testStruct0{2}})

	w.Assign(e0, ecs.Component{ID: 1, Comp: &testStruct1{2}})
	w.Assign(e1, ecs.Component{ID: 1, Comp: &testStruct1{3}})

	w.Assign(e0, ecs.Component{ID: 2, Comp: &testStruct2{3, 0}})
	w.Assign(e1, ecs.Component{ID: 2, Comp: &testStruct2{4, 0}})

	w.Assign(e0, ecs.Component{ID: 3, Comp: &testStruct3{4, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 3, Comp: &testStruct3{5, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 4, Comp: &testStruct4{5, 0, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 4, Comp: &testStruct4{6, 0, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 5, Comp: &testStruct5{6, 0, 0, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 5, Comp: &testStruct5{7, 0, 0, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 6, Comp: &testStruct6{7, 0, 0, 0, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 6, Comp: &testStruct6{8, 0, 0, 0, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 7, Comp: &testStruct7{8, 0, 0, 0, 0, 0, 0}})
	w.Assign(e1, ecs.Component{ID: 7, Comp: &testStruct7{9, 0, 0, 0, 0, 0, 0}})

	w.Assign(e0, ecs.Component{ID: 8, Comp: &testStruct8{9, 0, 0, 0, 0, 0, 0, 0}})

	w.Assign(e2, ecs.Component{ID: 9, Comp: &testStruct9{}})

	cnt := 0
	filter :=
		NewFilter8[
			testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testStruct7,
		]().
			Optional(T[testStruct1]()).
			With(T[testStruct8]()).
			Without(T[testStruct9]())

	f := ecs.All(0, 2, 3, 4, 5, 6, 7, 8).Without(9)
	assert.Equal(t, f, filter.Filter(&w))
	filter.Register(&w)

	query := filter.Query(&w)
	for query.Next() {
		c1, c2, c3, c4, c5, c6, c7, c8 := query.Get()
		assert.Equal(t, cnt+1, int(c1.val))
		assert.Equal(t, cnt+2, int(c2.val))
		assert.Equal(t, cnt+3, int(c3.val))
		assert.Equal(t, cnt+4, int(c4.val))
		assert.Equal(t, cnt+5, int(c5.val))
		assert.Equal(t, cnt+6, int(c6.val))
		assert.Equal(t, cnt+7, int(c7.val))
		assert.Equal(t, cnt+8, int(c8.val))
		assert.Panics(t, func() { query.Relation() })
		cnt++
	}
	assert.Equal(t, 1, cnt)

	filter.Unregister(&w)
	assert.Panics(t, func() { filter.Unregister(&w) })

	targ := w.NewEntity(0)

	w.Add(e0, relID)
	w.Add(e1, relID)
	w.Add(e2, relID)

	w.SetRelation(e0, relID, targ)

	filter2 :=
		NewFilter8[
			testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testRelationA,
		]().
			WithRelation(T[testRelationA](), targ)

	q := filter2.Query(&w)
	assert.Equal(t, 1, q.Count())
	for q.Next() {
		trg := q.Relation()
		assert.Equal(t, targ, trg)
	}

	filter2.Register(&w)
	assert.Panics(t, func() { filter2.Optional(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.With(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.Without(T[testStruct0]()) })
	assert.Panics(t, func() { filter2.WithRelation(T[testRelationA](), targ) })
	filter2.Unregister(&w)
}

func TestQueryGeneric(t *testing.T) {
	count := 1000
	world := ecs.NewWorld()

	posID := ecs.ComponentID[testStruct2](&world)
	rotID := ecs.ComponentID[testStruct3](&world)
	relID := ecs.ComponentID[testRelationA](&world)

	entities := make([]ecs.Entity, count)
	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID, relID)
		entities[i] = entity
	}
	query := NewFilter2[testStruct2, testStruct3]()

	q := query.Query(&world)
	cnt := 0
	for q.Next() {
		s1, _ := q.Get()
		*s1 = testStruct2{int32(cnt), int32(cnt)}

		s1, _ = q.Get()
		assert.Equal(t, &testStruct2{int32(cnt), int32(cnt)}, s1)
		cnt++
	}
	assert.Equal(t, count, cnt)

	target := world.NewEntity(rotID)

	world.SetRelation(entities[0], relID, target)

	filter := NewFilter2[testStruct2, testRelationA]().WithRelation(T[testRelationA](), target)
	q2 := filter.Query(&world)
	assert.Equal(t, 1, q2.Count())
	for q2.Next() {
		trg := q2.Relation()
		assert.Equal(t, target, trg)
	}

	filter = NewFilter2[testStruct2, testRelationA]().WithRelation(T[testRelationB](), target)
	assert.Panics(t, func() { filter.Query(&world) })

	filter = NewFilter2[testStruct2, testRelationA]().WithRelation(T[testStruct2](), target)
	assert.Panics(t, func() { filter.Query(&world) })
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
