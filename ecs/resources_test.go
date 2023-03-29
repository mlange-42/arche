package ecs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResources(t *testing.T) {
	res := newResources()

	posID := res.registry.ComponentID(reflect.TypeOf(Position{}))
	rotID := res.registry.ComponentID(reflect.TypeOf(rotation{}))

	assert.False(t, res.Has(posID))
	assert.Nil(t, res.Get(posID))

	res.Add(posID, &Position{1, 2})

	assert.True(t, res.Has(posID))
	pos, ok := res.Get(posID).(*Position)
	assert.True(t, ok)
	assert.Equal(t, Position{1, 2}, *pos)

	assert.Panics(t, func() { res.Add(posID, &Position{1, 2}) })

	pos, ok = res.Get(posID).(*Position)
	assert.True(t, ok)
	assert.Equal(t, Position{1, 2}, *pos)

	res.Add(rotID, &rotation{5})
	assert.True(t, res.Has(rotID))
	res.Remove(rotID)
	assert.False(t, res.Has(rotID))
	assert.Panics(t, func() { res.Remove(rotID) })
}

func TestResourcesReset(t *testing.T) {
	res := newResources()

	posID := res.registry.ComponentID(reflect.TypeOf(Position{}))
	rotID := res.registry.ComponentID(reflect.TypeOf(rotation{}))

	res.Add(posID, &Position{1, 2})
	res.Add(rotID, &rotation{5})

	pos, ok := res.Get(posID).(*Position)
	assert.True(t, ok)
	assert.Equal(t, Position{1, 2}, *pos)

	rot, ok := res.Get(rotID).(*rotation)
	assert.True(t, ok)
	assert.Equal(t, rotation{5}, *rot)

	res.reset()

	assert.False(t, res.Has(posID))
	assert.False(t, res.Has(rotID))

	res.Add(posID, &Position{10, 20})
	res.Add(rotID, &rotation{50})

	pos, ok = res.Get(posID).(*Position)
	assert.True(t, ok)
	assert.Equal(t, Position{10, 20}, *pos)

	rot, ok = res.Get(rotID).(*rotation)
	assert.True(t, ok)
	assert.Equal(t, rotation{50}, *rot)
}

func ExampleResources() {
	world := NewWorld()

	resID := ResourceID[Position](&world)

	myRes := Position{100, 100}
	world.Resources().Add(resID, &myRes)

	res := (world.Resources().Get(resID)).(*Position)
	fmt.Println(res)

	if world.Resources().Has(resID) {
		world.Resources().Remove(resID)
	}
	// Output: &{100 100}
}
