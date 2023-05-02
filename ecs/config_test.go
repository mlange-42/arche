package ecs_test

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := ecs.NewConfig()
	c = c.WithCapacityIncrement(16)
	assert.Equal(t, 16, c.CapacityIncrement)
	assert.Equal(t, 0, c.RelationCapacityIncrement)

	c = c.WithRelationCapacityIncrement(8)
	assert.Equal(t, 8, c.RelationCapacityIncrement)

	_ = ecs.NewWorld(c)
}

func ExampleConfig() {
	config :=
		ecs.NewConfig().
			WithCapacityIncrement(1024).       // Optionally set capacity increment
			WithRelationCapacityIncrement(128) // Optionally set capacity increment for relations

	world := ecs.NewWorld(config)

	world.NewEntity()
	// Output:
}
