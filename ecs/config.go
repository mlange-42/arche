package ecs

// Config provides configuration for an ECS [World].
type Config struct {
	CapacityIncrement int
}

// NewConfig creates a new default [World] config.
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

// Build creates a [World] from this Config.
func (c Config) Build() World {
	return FromConfig(c)
}
