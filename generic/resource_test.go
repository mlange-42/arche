package generic

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestGenericResource(t *testing.T) {
	w := ecs.NewWorld()
	get := NewResource[testStruct0](&w)

	assert.Equal(t, ecs.ResourceID[testStruct0](&w), get.ID())

	assert.False(t, get.Has())
	get.Add(&testStruct0{100})

	assert.True(t, get.Has())
	res := get.Get()

	assert.Equal(t, testStruct0{100}, *res)

	get.Remove()
	assert.False(t, get.Has())
}

func ExampleResource() {
	world := ecs.NewWorld()
	myRes := Position{}

	resAccess := NewResource[Position](&world)
	resAccess.Add(&myRes)
	res := resAccess.Get()
	res.X, res.Y = 10, 5
	// Output:
}
