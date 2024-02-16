package resources

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// Grid resource
type Grid struct {
	Data   [][]ecs.Entity
	Width  int
	Height int
}

// NewGrid creates a ner Grid of the given size.
func NewGrid(w, h int) Grid {
	grid := make([][]ecs.Entity, w)
	for i := 0; i < w; i++ {
		grid[i] = make([]ecs.Entity, h)
	}
	return Grid{
		Data:   grid,
		Width:  w,
		Height: h,
	}
}

func TestResourceID(t *testing.T) {
	world := ecs.NewWorld()

	gridID := ecs.ResourceID[Grid](&world)
	_ = gridID
}

func TestResourceAdd(t *testing.T) {
	world := ecs.NewWorld()

	// Create the resource, and add a pointer to it to the world.
	grid := NewGrid(30, 20)
	ecs.AddResource(&world, &grid)
}

func TestResourceAdd2(t *testing.T) {
	world := ecs.NewWorld()
	gridID := ecs.ResourceID[Grid](&world)

	// Create the resource, and add a pointer to it to the world.
	grid := NewGrid(30, 20)
	world.Resources().Add(gridID, &grid)
}

func TestResourceGet(t *testing.T) {
	world := ecs.NewWorld()

	// Create the resource, and add a pointer to it to the world.
	grid := NewGrid(30, 20)
	ecs.AddResource(&world, &grid)

	// Then, somewhere else in the code...
	gridID := ecs.ResourceID[Grid](&world)
	grid2 := world.Resources().Get(gridID).(*Grid)

	_ = grid2.Data[1][2]
}

func TestResourceGetGeneric(t *testing.T) {
	world := ecs.NewWorld()

	// Create the resource, and add a pointer to it to the world.
	grid := NewGrid(30, 20)
	ecs.AddResource(&world, &grid)

	// Then, somewhere else in the code...
	gridRes := generic.NewResource[Grid](&world)
	grid2 := gridRes.Get()

	_ = grid2.Data[1][2]
}
