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
	// Node statistics
	Nodes []NodeStats
	// Number of active nodes
	ActiveNodeCount int
	// Memory used by entities and components
	Memory int
	// Number of cached filters
	CachedFilters int
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

// NodeStats provide statistics for an archetype graph node.
type NodeStats struct {
	// Total number of archetypes, incl. inactive.
	ArchetypeCount int
	// Number of active archetypes.
	ActiveArchetypeCount int
	// Whether the node is active and contains archetypes.
	IsActive bool
	// Whether the node contains relation archetypes.
	HasRelation bool
	// Number of components
	Components int
	// Component IDs
	ComponentIDs []uint8
	// Component types for ComponentIDs
	ComponentTypes []reflect.Type
	// Memory for components per entity
	MemoryPerEntity int
	// Total reserved memory for entities and components, in bytes
	Memory int
	// Number of entities in the archetype
	Size int
	// Capacity of the archetype
	Capacity int
	// Archetype statistics
	Archetypes []ArchetypeStats
}

// ArchetypeStats provide statistics for an archetype.
type ArchetypeStats struct {
	// Whether the archetype is currently active
	IsActive bool
	// Number of entities in the archetype
	Size int
	// Capacity of the archetype
	Capacity int
	// Total reserved memory for entities and components, in bytes
	Memory int
}

func (s *WorldStats) String() string {
	b := strings.Builder{}

	fmt.Fprintf(
		&b, "World -- Components: %d, Nodes: %d, Filters: %d, Memory: %.1f kB, Locked: %t\n",
		s.ComponentCount, len(s.Nodes), s.CachedFilters, float64(s.Memory)/1024.0, s.Locked,
	)

	typeNames := make([]string, len(s.ComponentTypes))
	for i, tp := range s.ComponentTypes {
		typeNames[i] = tp.Name()
	}
	fmt.Fprintf(&b, "  Components: %s\n", strings.Join(typeNames, ", "))
	fmt.Fprint(&b, s.Entities.String())

	for i := range s.Nodes {
		fmt.Fprint(&b, s.Nodes[i].String())
	}

	return b.String()
}

func (s *EntityStats) String() string {
	return fmt.Sprintf("Entities -- Used: %d, Recycled: %d, Total: %d, Capacity: %d\n", s.Used, s.Recycled, s.Total, s.Capacity)
}

func (s *NodeStats) String() string {
	if !s.IsActive {
		return ""
	}

	typeNames := make([]string, len(s.ComponentTypes))
	for i, tp := range s.ComponentTypes {
		typeNames[i] = tp.Name()
	}

	return fmt.Sprintf(
		"Node -- Components: %2d, Entities: %6d, Capacity: %6d, Memory: %7.1f kB, Per entity: %4d B\n  Components: %s\n",
		s.Components, s.Size, s.Capacity, float64(s.Memory)/1024.0, s.MemoryPerEntity, strings.Join(typeNames, ", "),
	)
}
