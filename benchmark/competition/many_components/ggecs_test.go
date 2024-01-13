package manycomponents

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
)

const (
	PositionComponentID ecs.ComponentID = iota
	Comp1ID
	Comp2ID
	Comp3ID
	Comp4ID
	Comp5ID
	Comp6ID
	Comp7ID
	Comp8ID
	Comp9ID
)

func BenchmarkIterGGEcs(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(1024)
	world.Register(ecs.NewComponentRegistry[Position](PositionComponentID))
	world.Register(ecs.NewComponentRegistry[Comp1](Comp1ID))
	world.Register(ecs.NewComponentRegistry[Comp2](Comp2ID))
	world.Register(ecs.NewComponentRegistry[Comp3](Comp3ID))
	world.Register(ecs.NewComponentRegistry[Comp4](Comp4ID))
	world.Register(ecs.NewComponentRegistry[Comp5](Comp5ID))
	world.Register(ecs.NewComponentRegistry[Comp6](Comp6ID))
	world.Register(ecs.NewComponentRegistry[Comp7](Comp7ID))
	world.Register(ecs.NewComponentRegistry[Comp8](Comp8ID))
	world.Register(ecs.NewComponentRegistry[Comp9](Comp9ID))

	for i := 0; i < nPos; i++ {
		_ = world.NewEntity(PositionComponentID)
	}
	for i := 0; i < nPosAll; i++ {
		_ = world.NewEntity(PositionComponentID, Comp1ID, Comp2ID, Comp3ID, Comp4ID, Comp5ID, Comp6ID, Comp7ID, Comp8ID, Comp9ID)
	}

	mask := ecs.MakeComponentMask(PositionComponentID, Comp1ID, Comp2ID, Comp3ID, Comp4ID, Comp5ID, Comp6ID, Comp7ID, Comp8ID, Comp9ID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(mask)
		for query.Next() {
			pos := (*Position)(query.Component(PositionComponentID))
			c1 := (*Comp1)(query.Component(Comp1ID))
			c2 := (*Comp1)(query.Component(Comp1ID))
			c3 := (*Comp1)(query.Component(Comp1ID))
			c4 := (*Comp1)(query.Component(Comp1ID))
			c5 := (*Comp1)(query.Component(Comp1ID))
			c6 := (*Comp1)(query.Component(Comp1ID))
			c7 := (*Comp1)(query.Component(Comp1ID))
			c8 := (*Comp1)(query.Component(Comp1ID))
			c9 := (*Comp1)(query.Component(Comp1ID))
			pos.X += c1.X + c2.X + c3.X + c4.X + c5.X + c6.X + c7.X + c8.X + c9.X
			pos.Y += c1.Y + c2.Y + c3.Y + c4.Y + c5.Y + c6.Y + c7.Y + c8.Y + c9.Y
		}
	}
}
