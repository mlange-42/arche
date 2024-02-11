// Demonstrates that ECS can be mixed with non-ECS data structures, as long as they store entities.
// Uses the ID-based API.
package main

import (
	"fmt"
	"math/rand"

	"github.com/mlange-42/arche/ecs"
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
	gridId := ecs.ResourceID[Grid](world)
	grid := world.Resources().Get(gridId).(*Grid)

	// Get the CellCoord component ID.
	coordId := ecs.ComponentID[CellCoord](world)

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
		entity := world.NewEntity(coordId)
		// Place the entity in the grid.
		grid.Data[x][y] = entity
		// Initialize entity components.
		coord := (*CellCoord)(world.Get(entity, coordId))
		coord.X, coord.Y = x, y

		cnt++
	}
}

func run(world *ecs.World) {
	// Get the grid resource.
	gridId := ecs.ResourceID[Grid](world)
	grid := world.Resources().Get(gridId).(*Grid)

	// Get the CellCoord component ID.
	coordId := ecs.ComponentID[CellCoord](world)

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
		coord := (*CellCoord)(world.Get(entity, coordId))
		fmt.Printf("(%2d,%2d): %+v\n", x, y, coord)
	}
}
