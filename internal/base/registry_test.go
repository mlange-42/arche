package base

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type position struct {
	X int
	Y int
}

type rotation struct {
	Angle int
}

func TestComponentRegistry(t *testing.T) {
	reg := NewComponentRegistry()

	posType := reflect.TypeOf((*position)(nil)).Elem()
	rotType := reflect.TypeOf((*rotation)(nil)).Elem()

	reg.RegisterComponent(posType, MaskTotalBits)

	assert.Equal(t, ID(0), reg.ComponentID(posType))
	assert.Equal(t, ID(1), reg.ComponentID(rotType))

	assert.Equal(t, posType, reg.ComponentType(ID(0)))
	assert.Equal(t, rotType, reg.ComponentType(ID(1)))
}

func TestComponentRegistryOverflow(t *testing.T) {
	reg := NewComponentRegistry()

	reg.RegisterComponent(reflect.TypeOf((*position)(nil)).Elem(), 1)

	assert.Panics(t, func() {
		reg.RegisterComponent(reflect.TypeOf((*rotation)(nil)).Elem(), 1)
	})
}
