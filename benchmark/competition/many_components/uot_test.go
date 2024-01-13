package manycomponents

import (
	"testing"

	"github.com/unitoftime/ecs"
)

func BenchmarkIterUot(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld()

	for i := 0; i < nPos; i++ {
		id := world.NewId()
		ecs.Write(world, id,
			ecs.C(Position{0, 0}),
		)
	}
	for i := 0; i < nPosAll; i++ {
		id := world.NewId()
		ecs.Write(world, id,
			ecs.C(Position{0, 0}),
			ecs.C(Comp1{0, 0}),
			ecs.C(Comp2{0, 0}),
			ecs.C(Comp3{0, 0}),
			ecs.C(Comp4{0, 0}),
			ecs.C(Comp5{0, 0}),
			ecs.C(Comp6{0, 0}),
			ecs.C(Comp7{0, 0}),
			ecs.C(Comp8{0, 0}),
			ecs.C(Comp9{0, 0}),
		)
	}
	query := ecs.Query10[Position, Comp1, Comp2, Comp3, Comp4, Comp5, Comp6, Comp7, Comp8, Comp9](world)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query.MapId(func(id ecs.Id, pos *Position,
			c1 *Comp1,
			c2 *Comp2,
			c3 *Comp3,
			c4 *Comp4,
			c5 *Comp5,
			c6 *Comp6,
			c7 *Comp7,
			c8 *Comp8,
			c9 *Comp9,
		) {
			pos.X += c1.X + c2.X + c3.X + c4.X + c5.X + c6.X + c7.X + c8.X + c9.X
			pos.Y += c1.Y + c2.Y + c3.Y + c4.Y + c5.Y + c6.Y + c7.Y + c8.Y + c9.Y
		})
	}
}
