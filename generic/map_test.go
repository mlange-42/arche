package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestGenericMap(t *testing.T) {
	w := ecs.NewWorld()
	get := NewMap[testStruct0](&w)
	mut := NewExchange(&w).Adds(T[testStruct0]())

	assert.Equal(t, ecs.ComponentID[testStruct0](&w), get.ID())

	e0 := w.NewEntity()

	mut.Add(e0)
	has := get.Has(e0)
	_ = get.Get(e0)
	assert.True(t, has)

	_ = get.Set(e0, &testStruct0{100})
	str := get.Get(e0)

	assert.Equal(t, 100, int(str.val))

	get2 := NewMap[testStruct1](&w)
	assert.Equal(t, ecs.ComponentID[testStruct1](&w), get2.ID())
	assert.Panics(t, func() { get2.Set(e0, &testStruct1{}) })
}
