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
	assert.Equal(t, []uint8{uint8(0)}, reg.IDs)

	reg.registerComponent(rotType, MaskTotalBits)
	reg.unregisterLastComponent()
	assert.Equal(t, []uint8{uint8(0)}, reg.IDs)

	id0, _ := reg.ComponentID(posType)
	id1, _ := reg.ComponentID(rotType)
	assert.Equal(t, uint8(0), id0)
	assert.Equal(t, uint8(1), id1)

	assert.Equal(t, []uint8{uint8(0), uint8(1)}, reg.IDs)

	t1, _ := reg.ComponentType(uint8(0))
	t2, _ := reg.ComponentType(uint8(1))

	assert.Equal(t, posType, t1)
	assert.Equal(t, rotType, t2)
}

func TestComponentRegistryOverflow(t *testing.T) {
	reg := newComponentRegistry()

	reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), 1)

	assert.PanicsWithValue(t, "exceeded the maximum of 1 component types or resource types", func() {
		reg.registerComponent(reflect.TypeOf((*rotation)(nil)).Elem(), 1)
	})
}

type relationComp struct {
	Relation
}

type noRelationComp1 struct {
	Rel Relation
}

type noRelationComp2 struct {
	Position
}

type noRelationComp3 struct{}

func TestRegistryRelations(t *testing.T) {
	registry := newComponentRegistry()

	relCompTp := reflect.TypeOf((*relationComp)(nil)).Elem()
	noRelCompTp1 := reflect.TypeOf((*noRelationComp1)(nil)).Elem()
	noRelCompTp2 := reflect.TypeOf((*noRelationComp2)(nil)).Elem()
	noRelCompTp3 := reflect.TypeOf((*noRelationComp3)(nil)).Elem()

	assert.True(t, registry.isRelation(relCompTp))
	assert.False(t, registry.isRelation(noRelCompTp1))
	assert.False(t, registry.isRelation(noRelCompTp2))
	assert.False(t, registry.isRelation(noRelCompTp3))

	id1, _ := registry.ComponentID(relCompTp)
	id2, _ := registry.ComponentID(noRelCompTp1)
	id3, _ := registry.ComponentID(noRelCompTp2)
	id4, _ := registry.ComponentID(noRelCompTp3)

	assert.True(t, registry.IsRelation.Get(id(id1)))
	assert.False(t, registry.IsRelation.Get(id(id2)))
	assert.False(t, registry.IsRelation.Get(id(id3)))
	assert.False(t, registry.IsRelation.Get(id(id4)))
}
