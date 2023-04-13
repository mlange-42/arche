package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentRegistry(t *testing.T) {
	reg := newComponentRegistry()

	posType := reflect.TypeOf((*Position)(nil)).Elem()
	rotType := reflect.TypeOf((*rotation)(nil)).Elem()

	reg.registerComponent(posType, MaskTotalBits)

	assert.Equal(t, ID(0), reg.ComponentID(posType))
	assert.Equal(t, ID(1), reg.ComponentID(rotType))

	t1, _ := reg.ComponentType(ID(0))
	t2, _ := reg.ComponentType(ID(1))

	assert.Equal(t, posType, t1)
	assert.Equal(t, rotType, t2)
}

func TestComponentRegistryOverflow(t *testing.T) {
	reg := newComponentRegistry()

	reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), 1)

	assert.Panics(t, func() {
		reg.registerComponent(reflect.TypeOf((*rotation)(nil)).Elem(), 1)
	})
}
