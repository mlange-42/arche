package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentRegistry(t *testing.T) {
	reg := NewComponentRegistry()

	posType := reflect.TypeOf((*position)(nil)).Elem()
	rotType := reflect.TypeOf((*rotation)(nil)).Elem()

	reg.RegisterComponent(posType)

	assert.Equal(t, ID(0), reg.ComponentID(posType))
	assert.Equal(t, ID(1), reg.ComponentID(rotType))

	assert.Equal(t, posType, reg.ComponentType(ID(0)))
	assert.Equal(t, rotType, reg.ComponentType(ID(1)))
}
