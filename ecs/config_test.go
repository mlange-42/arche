package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := NewConfig()
	c = c.WithCapacityIncrement(16)
	assert.Equal(t, 16, c.CapacityIncrement)
	_ = NewWorld(c)
}

func ExampleConfig() {
	config := NewConfig().WithCapacityIncrement(1024)
	world := NewWorld(config)

	world.NewEntity()
	// Output:
}
