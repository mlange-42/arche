package ecs_test

import (
	"fmt"
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

func ExampleComponentID() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)

	world.NewEntity(posID)
	// Output:
}

func ExampleTypeID() {
	world := ecs.NewWorld()
	posID := ecs.TypeID(&world, reflect.TypeOf(Position{}))

	world.NewEntity(posID)
	// Output:
}

func ExampleResourceID() {
	world := ecs.NewWorld()
	resID := ecs.ResourceID[Position](&world)

	world.Resources().Add(resID, &Position{100, 100})
	// Output:
}

func ExampleGetResource() {
	world := ecs.NewWorld()

	myRes := Position{100, 100}

	ecs.AddResource(&world, &myRes)
	res := ecs.GetResource[Position](&world)
	fmt.Println(res)
	// Output: &{100 100}
}

func ExampleAddResource() {
	world := ecs.NewWorld()

	myRes := Position{100, 100}
	ecs.AddResource(&world, &myRes)

	res := ecs.GetResource[Position](&world)
	fmt.Println(res)
	// Output: &{100 100}
}

func ExampleWorld() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	_ = world.NewEntity(posID, velID)
	// Output:
}

func ExampleNewWorld() {
	defaultWorld := ecs.NewWorld()

	configWorld := ecs.NewWorld(
		ecs.NewConfig().
			WithCapacityIncrement(1024).
			WithRelationCapacityIncrement(64),
	)

	_, _ = defaultWorld, configWorld
	// Output:
}

func ExampleWorld_NewEntity() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	_ = world.NewEntity(posID, velID)
	// Output:
}

func ExampleWorld_NewEntityWith() {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	_ = world.NewEntityWith(
		ecs.Component{ID: posID, Comp: &Position{X: 0, Y: 0}},
		ecs.Component{ID: velID, Comp: &Velocity{X: 10, Y: 2}},
	)
	// Output:
}

func ExampleWorld_RemoveEntity() {
	world := ecs.NewWorld()
	e := world.NewEntity()
	world.RemoveEntity(e)
	// Output:
}

func ExampleWorld_Alive() {
	world := ecs.NewWorld()

	e := world.NewEntity()
	fmt.Println(world.Alive(e))

	world.RemoveEntity(e)
	fmt.Println(world.Alive(e))
	// Output:
	// true
	// false
}

func ExampleWorld_Get() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)

	e := world.NewEntity(posID)

	pos := (*Position)(world.Get(e, posID))
	pos.X, pos.Y = 10, 5
	// Output:
}

func ExampleWorld_Has() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)

	e := world.NewEntity(posID)

	if world.Has(e, posID) {
		world.Remove(e, posID)
	}
	// Output:
}

func ExampleWorld_Add() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	e := world.NewEntity()

	world.Add(e, posID, velID)
	// Output:
}

func ExampleWorld_Assign() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	e := world.NewEntity()

	world.Assign(e,
		ecs.Component{ID: posID, Comp: &Position{X: 0, Y: 0}},
		ecs.Component{ID: velID, Comp: &Velocity{X: 10, Y: 2}},
	)
	// Output:
}

func ExampleWorld_Set() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)

	e := world.NewEntity(posID)

	world.Set(e, posID, &Position{X: 0, Y: 0})
	// Output:
}

func ExampleWorld_Remove() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	e := world.NewEntity(posID, velID)

	world.Remove(e, posID, velID)
	// Output:
}

func ExampleWorld_Exchange() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	e := world.NewEntity(posID)

	world.Exchange(e, []ecs.ID{velID}, []ecs.ID{posID})
	// Output:
}

func ExampleWorld_Reset() {
	world := ecs.NewWorld()
	_ = world.NewEntity()

	world.Reset()
	// Output:
}

func ExampleWorld_Query() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	filter := ecs.All(posID, velID)
	query := world.Query(filter)
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		vel := (*Velocity)(query.Get(velID))
		pos.X += vel.X
		pos.Y += vel.Y
	}
	// Output:
}

func ExampleWorld_Relations() {
	world := ecs.NewWorld()

	relID := ecs.ComponentID[ChildOf](&world)

	parent := world.NewEntity()
	child := world.NewEntity(relID)

	world.Relations().Set(child, relID, parent)
	fmt.Println(world.Relations().Get(child, relID))
	// Output: {1 0}
}

func ExampleWorld_Resources() {
	world := ecs.NewWorld()

	resID := ecs.ResourceID[Position](&world)

	myRes := Position{}
	world.Resources().Add(resID, &myRes)

	res := (world.Resources().Get(resID)).(*Position)
	res.X, res.Y = 10, 5
	// Output:
}

func ExampleWorld_Cache() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)

	filter := ecs.All(posID)
	cached := world.Cache().Register(filter)
	query := world.Query(&cached)

	for query.Next() {
		// handle entities...
	}
	// Output:
}

func ExampleWorld_Batch() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)

	filter := ecs.All(posID)
	world.Batch().RemoveEntities(filter)
	// Output:
}

func ExampleWorld_SetListener() {
	world := ecs.NewWorld()

	listener := func(evt *ecs.EntityEvent) {
		fmt.Println(evt)
	}
	world.SetListener(listener)

	world.NewEntity()
	// Output: &{{1 0} {[0 0 0 0]} [] [] <nil> <nil> {0 0} 1 false false}
}

func ExampleWorld_Stats() {
	world := ecs.NewWorld()
	stats := world.Stats()
	fmt.Println(stats.Entities.String())
	// Output: Entities -- Used: 0, Recycled: 0, Total: 0, Capacity: 128
}
