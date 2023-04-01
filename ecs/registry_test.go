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

type relationComp struct {
	Relation
}

type noRelationComp1 struct {
	rel Relation
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

	id1 := registry.ComponentID(relCompTp)
	id2 := registry.ComponentID(noRelCompTp1)
	id3 := registry.ComponentID(noRelCompTp2)
	id4 := registry.ComponentID(noRelCompTp3)

	assert.True(t, registry.IsRelation.Get(id1))
	assert.False(t, registry.IsRelation.Get(id2))
	assert.False(t, registry.IsRelation.Get(id3))
	assert.False(t, registry.IsRelation.Get(id4))
}
