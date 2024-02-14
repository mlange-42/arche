package entities

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

func TestZeroEntity(t *testing.T) {
	var zero1 ecs.Entity
	fmt.Println(zero1.IsZero()) // prints true

	zero2 := ecs.Entity{}
	fmt.Println(zero2.IsZero()) // prints true
}

func TestComponentID(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	_, _ = posID, headID
}

func TestEntitiesCreate(t *testing.T) {
	world := ecs.NewWorld()

	entity := world.NewEntity()
	_ = entity
}

func TestEntitiesCreateComponents(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	_ = world.NewEntity(posID)
	_ = world.NewEntity(posID, headID)
}

func TestEntitiesCreateWithComponents(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	_ = world.NewEntityWith(
		ecs.Component{ID: posID, Comp: &Position{X: 1, Y: 2}},
		ecs.Component{ID: headID, Comp: &Heading{Angle: 180}},
	)
}

func TestEntitiesCreateGeneric(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap2[Position, Heading](&world)

	_ = builder.New()
}

func TestEntitiesCreateWithComponentsGeneric(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap2[Position, Heading](&world)

	_ = builder.NewWith(
		&Position{X: 1, Y: 2},
		&Heading{Angle: 180},
	)
}

func TestEntitiesAddRemove(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	entity := world.NewEntity()

	world.Add(entity, posID, headID)
	world.Remove(entity, posID, headID)
}

func TestEntitiesAddRemoveGeneric(t *testing.T) {
	world := ecs.NewWorld()

	mapper := generic.NewMap2[Position, Heading](&world)

	entity := world.NewEntity()

	mapper.Add(entity)
	mapper.Remove(entity)
}

func TestEntitiesExchange(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	entity := world.NewEntity(posID)

	world.Exchange(entity,
		[]ecs.ID{headID}, // Component(s) to add.
		[]ecs.ID{posID},  // Component(s) to remove.
	)
}

func TestEntitiesExchangeGeneric(t *testing.T) {
	world := ecs.NewWorld()

	builder := generic.NewMap1[Position](&world)
	entity := builder.New()

	exchange := generic.NewExchange(&world).
		Adds(generic.T[Heading]()).    // Component(s) to add.
		Removes(generic.T[Position]()) // Component(s) to remove.

	exchange.Exchange(entity)
}

func TestEntitiesAssign(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	headID := ecs.ComponentID[Heading](&world)

	entity := world.NewEntity()

	world.Assign(
		entity,
		ecs.Component{ID: posID, Comp: &Position{X: 1, Y: 2}},
		ecs.Component{ID: headID, Comp: &Heading{Angle: 180}},
	)
}

func TestEntitiesAssignGeneric(t *testing.T) {
	world := ecs.NewWorld()

	mapper := generic.NewMap2[Position, Heading](&world)

	entity := world.NewEntity()

	mapper.Assign(
		entity,
		&Position{X: 1, Y: 2},
		&Heading{Angle: 180},
	)
}

func TestEntitiesRemove(t *testing.T) {
	world := ecs.NewWorld()

	entity := world.NewEntity()
	world.RemoveEntity(entity)
}

func TestEntitiesAlive(t *testing.T) {
	world := ecs.NewWorld()

	entity := world.NewEntity()
	fmt.Println(world.Alive(entity)) // prints true

	world.RemoveEntity(entity)
	fmt.Println(world.Alive(entity)) // prints false
}
