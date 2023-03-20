package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestGenericResourceMap(t *testing.T) {
	w := ecs.NewWorld()
	get := NewResMap[testStruct0](&w)

	assert.Equal(t, ecs.ResourceID[testStruct0](&w), get.ID())

	assert.False(t, get.Has())
	w.AddResource(&testStruct0{100})

	assert.True(t, get.Has())
	res := get.Get()

	assert.Equal(t, testStruct0{100}, *res)
}
