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
	_ = ecs.NewWorld(c)
}

func ExampleConfig() {
	config := ecs.NewConfig().WithCapacityIncrement(1024)
	world := ecs.NewWorld(config)

	world.NewEntity()
	// Output:
}
