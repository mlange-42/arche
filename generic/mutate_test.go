package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestMutate(t *testing.T) {
	w := ecs.NewWorld()
	registerAll(&w)

	mut := NewMutate(&w).Adds(T2[testStruct0, testStruct1]()...)
	mutR := NewMutate(&w).Removes(T2[testStruct0, testStruct1]()...)
	mutEx := NewMutate(&w).
		Removes(T2[testStruct0, testStruct1]()...).
		Adds(T2[testStruct2, testStruct3]()...)

	mapper := NewMap2[testStruct0, testStruct1](&w)
	mapper2 := NewMap2[testStruct2, testStruct3](&w)
	map0 := NewMap[testStruct0](&w)
	map1 := NewMap[testStruct1](&w)

	e := mut.NewEntity()
	s0, s1 := mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutR.Remove(e)
	assert.False(t, map0.Has(e))
	assert.False(t, map1.Has(e))

	mut.Add(e)
	s0, s1 = mapper.Get(e)
	assert.NotNil(t, s0)
	assert.NotNil(t, s1)

	mutEx.Exchange(e)
	s0, s1 = mapper.Get(e)
	assert.Nil(t, s0)
	assert.Nil(t, s1)
	s2, s3 := mapper2.Get(e)
	assert.NotNil(t, s2)
	assert.NotNil(t, s3)
}
