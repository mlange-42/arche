package stats

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// Position component
type Position struct {
	X float64
	Y float64
}

// Heading component
type Heading struct {
	Angle float64
}

func TestWorldStats(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap2[Position, Heading](&world)
	builder.NewBatch(100)

	stats := world.Stats()
	fmt.Println(stats)
}
