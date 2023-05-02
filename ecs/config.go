package ecs

// Config provides configuration for an ECS [World].
type Config struct {
	// Capacity increment for archetypes and the entity index.
	// The default value is 128.
	CapacityIncrement int
	// Capacity increment for archetypes with a relation component.
	// The default value is CapacityIncrement.
	RelationCapacityIncrement int
}

// NewConfig creates a new default [World] configuration.
func NewConfig() Config {
	return Config{
		CapacityIncrement:         128,
		RelationCapacityIncrement: 0,
	}
}

// WithCapacityIncrement return a new Config with CapacityIncrement set.
// Use with method chaining.
func (c Config) WithCapacityIncrement(inc int) Config {
	c.CapacityIncrement = inc
	return c
}

// WithRelationCapacityIncrement return a new Config with RelationCapacityIncrement set.
// Use with method chaining.
func (c Config) WithRelationCapacityIncrement(inc int) Config {
	c.RelationCapacityIncrement = inc
	return c
}
