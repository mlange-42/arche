package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelation(t *testing.T) {
	rel := Relation{}

	assert.True(t, rel.Target().IsZero())

	rel.setTarget(Entity{100, 1000})
	assert.Equal(t, Entity{100, 1000}, rel.Target())
}
