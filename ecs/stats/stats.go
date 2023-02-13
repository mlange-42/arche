package stats

import (
	"fmt"
	"reflect"
	"strings"
)

// WorldStats provide statistics for a [World].
type WorldStats struct {
	Entities       EntityStats
	ComponentCount int
	ComponentTypes []reflect.Type
	Locked         bool
	Archetypes     []ArchetypeStats
}

// EntityStats provide statistics about [World] entities.
type EntityStats struct {
	Used     int
	Capacity int
	Recycled int
}

// ArchetypeStats provide statistics for an archetype.
type ArchetypeStats struct {
	Size           int
	Capacity       int
	Components     int
	ComponentIDs   []uint8
	ComponentTypes []reflect.Type
}

func (s *WorldStats) String() string {
	b := strings.Builder{}

	fmt.Fprintf(&b, "World -- Components: %d, Archetypes: %d, Locked: %t\n", s.ComponentCount, len(s.Archetypes), s.Locked)

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
	return fmt.Sprintf("Entities -- Used: %d, Recycled: %d, Capacity: %d\n", s.Used, s.Recycled, s.Capacity)
}

func (s *ArchetypeStats) String() string {
	typeNames := make([]string, len(s.ComponentTypes))
	for i, tp := range s.ComponentTypes {
		typeNames[i] = tp.Name()
	}
	return fmt.Sprintf(
		"Archetype -- Components: %d, Entities: %d, Capacity: %d\n  Components: %s\n",
		s.Components, s.Size, s.Capacity, strings.Join(typeNames, ", "),
	)
}
