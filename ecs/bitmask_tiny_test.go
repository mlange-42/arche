//go:build tiny

package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitMaskTiny(t *testing.T) {
	mask := Mask{}
	mask.Set(id(100), true)

	assert.True(t, mask.IsZero())
}

func TestMaskTinyToTypes(t *testing.T) {
	w := NewWorld()

	id1 := ComponentID[Position](&w)
	id2 := ComponentID[Velocity](&w)

	mask := All()
	comps := mask.toTypes(&w.registry)
	assert.Equal(t, []componentType{}, comps)

	mask = All(id1, id2)
	comps = mask.toTypes(&w.registry)
	assert.Equal(t, []componentType{
		{ID: id1, Type: reflect.TypeOf((*Position)(nil)).Elem()},
		{ID: id2, Type: reflect.TypeOf((*Velocity)(nil)).Elem()},
	}, comps)
}
