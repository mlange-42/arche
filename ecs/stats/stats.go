// Package stats provides the structs returned by ecs.World.Stats().
package stats

import (
	"fmt"
	"reflect"
	"strings"
)

// WorldStats provide statistics for a [World].
type WorldStats struct {
	// Entity statistics
	Entities EntityStats
	// Total number of components
	ComponentCount int
	// Component types, indexed by component ID
	ComponentTypes []reflect.Type
	// Locked state of the world
	Locked bool
	// Archetype statistics
	Archetypes []ArchetypeStats
	// Memory used by entities and components
	Memory int
}

// EntityStats provide statistics about [World] entities.
type EntityStats struct {
	// Currently used/alive entities
	Used int
	// Current capacity of the entity pool
	Total int
	// Recycled/available entities
	Recycled int
	// Current capacity of the entities list
	Capacity int
}

// ArchetypeStats provide statistics for an archetype.
type ArchetypeStats struct {
	// Number of entities in the archetype
	Size int
	// Capacity of the archetype
	Capacity int
	// Number of components
	Components int
	// Component IDs
	ComponentIDs []uint8
	// Component types for ComponentIDs
	ComponentTypes []reflect.Type
	// Total reserved memory for entities and components, in bytes
	Memory int
	// Memory for components per entity
	MemoryPerEntity int
}

func (s *WorldStats) String() string {
	b := strings.Builder{}

	fmt.Fprintf(
		&b, "World -- Components: %d, Archetypes: %d, Memory: %.1f kB, Locked: %t\n",
		s.ComponentCount, len(s.Archetypes), float64(s.Memory)/1024.0, s.Locked,
	)

	typeNames := make([]string, len(s.ComponentTypes))
	for i, tp := range s.ComponentTypes {
		typeNames[i] = tp.Name()
	}
	fmt.Fprintf(&b, "  Components: %s\n", strings.Join(typeNames, ", "))
	fmt.Fprint(&b, s.Entities.String())

	for _, arch := range s.Archetypes {
		fmt.Fprint(&b, arch.String())
	}

	return b.String()
}

func (s *EntityStats) String() string {
	return fmt.Sprintf("Entities -- Used: %d, Recycled: %d, Total: %d, Capacity: %d\n", s.Used, s.Recycled, s.Total, s.Capacity)
}

func (s *ArchetypeStats) String() string {
	typeNames := make([]string, len(s.ComponentTypes))
	for i, tp := range s.ComponentTypes {
		typeNames[i] = tp.Name()
	}
	return fmt.Sprintf(
		"Archetype -- Components: %2d, Entities: %6d, Capacity: %6d, Memory: %7.1f kB, Per entity: %4d B\n  Components: %s\n",
		s.Components, s.Size, s.Capacity, float64(s.Memory)/1024.0, s.MemoryPerEntity, strings.Join(typeNames, ", "),
	)
}
