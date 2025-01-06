package manycomponents

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func BenchmarkIterArche(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	posID := ecs.ComponentID[Position](&world)
	c1ID := ecs.ComponentID[Comp1](&world)
	c2ID := ecs.ComponentID[Comp2](&world)
	c3ID := ecs.ComponentID[Comp3](&world)
	c4ID := ecs.ComponentID[Comp4](&world)
	c5ID := ecs.ComponentID[Comp5](&world)
	c6ID := ecs.ComponentID[Comp6](&world)
	c7ID := ecs.ComponentID[Comp7](&world)
	c8ID := ecs.ComponentID[Comp8](&world)
	c9ID := ecs.ComponentID[Comp9](&world)

	ecs.NewBuilder(&world, posID).NewBatch(nPos)
	ecs.NewBuilder(&world, posID, c1ID, c2ID, c3ID, c4ID, c5ID, c6ID, c7ID, c8ID, c9ID).NewBatch(nPosAll)

	filter := ecs.All(posID, c1ID, c2ID, c3ID, c4ID, c5ID, c6ID, c7ID, c8ID, c9ID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			pos := (*Position)(query.Get(posID))
			c1 := (*Comp1)(query.Get(c1ID))
			c2 := (*Comp1)(query.Get(c2ID))
			c3 := (*Comp1)(query.Get(c3ID))
			c4 := (*Comp1)(query.Get(c4ID))
			c5 := (*Comp1)(query.Get(c5ID))
			c6 := (*Comp1)(query.Get(c6ID))
			c7 := (*Comp1)(query.Get(c7ID))
			c8 := (*Comp1)(query.Get(c8ID))
			c9 := (*Comp1)(query.Get(c9ID))
			pos.X += c1.X + c2.X + c3.X + c4.X + c5.X + c6.X + c7.X + c8.X + c9.X
			pos.Y += c1.Y + c2.Y + c3.Y + c4.Y + c5.Y + c6.Y + c7.Y + c8.Y + c9.Y
		}
	}
}

func BenchmarkIterArcheGeneric(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	posMapper := generic.NewMap1[Position](&world)
	posAllMapper := generic.NewMap10[
		Position, Comp1, Comp2, Comp3, Comp4, Comp5, Comp6, Comp7, Comp8, Comp9,
	](&world)

	posMapper.NewBatch(nPos)
	posAllMapper.NewBatch(nPosAll)

	filter := generic.NewFilter10[Position, Comp1, Comp2, Comp3, Comp4, Comp5, Comp6, Comp7, Comp8, Comp9]()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := filter.Query(&world)
		for query.Next() {
			pos, c1, c2, c3, c4, c5, c6, c7, c8, c9 := query.Get()
			pos.X += c1.X + c2.X + c3.X + c4.X + c5.X + c6.X + c7.X + c8.X + c9.X
			pos.Y += c1.Y + c2.Y + c3.Y + c4.Y + c5.Y + c6.Y + c7.Y + c8.Y + c9.Y
		}
	}
}
