// Demonstrates that ECS can be mixed with non-ECS data structures, as long as they store entities.
// Uses the generic API.
package main

import (
	"fmt"
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// CellCoord component.
type CellCoord struct {
	X int
	Y int
}

// Grid resource / data structure.
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

func main() {
	// Create a new World.
	world := ecs.NewWorld()

	// Create a non-ECS grid data structure,
	// and add is as a resource.
	grid := NewGrid(30, 20)
	ecs.AddResource(&world, &grid)

	// Create enities on the grid.
	createGridEntities(&world, 250)

	// Run a simulation
	run(&world)
}

func createGridEntities(world *ecs.World, count int) {
	// Get the grid resource.
	gridRes := generic.NewResource[Grid](world)
	grid := gridRes.Get()

	// Create a generic Map as a builder.
	builder := generic.NewMap1[CellCoord](world)

	// Put some entities into the grid.
	cnt := 0
	for cnt < count {
		// Draw random coordinates.
		x, y := rand.Intn(grid.Width), rand.Intn(grid.Height)
		// Skip if there is already an entity.
		if !grid.Data[x][y].IsZero() {
			continue
		}
		// Create an entity.
		entity := builder.New()
		// Place the entity in the grid.
		grid.Data[x][y] = entity
		// Initialize entity components.
		coord := builder.Get(entity)
		coord.X, coord.Y = x, y

		cnt++
	}
}

func run(world *ecs.World) {
	// Get the grid resource.
	gridRes := generic.NewResource[Grid](world)
	grid := gridRes.Get()

	// Create a generic Map for component access.
	mapper := generic.NewMap1[CellCoord](world)

	// Print random entities from the grid.
	for i := 0; i < 25; i++ {
		// Draw random coordinates.
		x, y := rand.Intn(grid.Width), rand.Intn(grid.Height)
		// Get the entity.
		entity := grid.Data[x][y]
		// Print zero entity
		if entity.IsZero() {
			fmt.Printf("(%2d,%2d): zero entity\n", x, y)
			continue
		}
		// Print CellCoord component of non-zero entity.
		coord := mapper.Get(entity)
		fmt.Printf("(%2d,%2d): %+v\n", x, y, coord)
	}
}
