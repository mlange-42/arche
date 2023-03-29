package ecs

// Config provides configuration for an ECS [World].
//
// # Example
//
//	config := NewConfig().WithCapacityIncrement(1024)
//	world := NewWorld(config)
type Config struct {
	// Capacity increment for archetypes and the entity index.
	// The default value is 128.
	CapacityIncrement int
}

// NewConfig creates a new default [World] configuration.
func NewConfig() Config {
	return Config{
		CapacityIncrement: 128,
	}
}

// WithCapacityIncrement return a new Config with CapacityIncrement set.
// Use with method chaining.
func (c Config) WithCapacityIncrement(inc int) Config {
	c.CapacityIncrement = inc
	return c
}
