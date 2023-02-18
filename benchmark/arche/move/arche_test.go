package move

import (
	"testing"

	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
	g "github.com/mlange-42/arche/generic"
)

func runArcheMove(b *testing.B, count int, add, rem []g.Comp) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld(
			ecs.NewConfig().WithCapacityIncrement(1024),
		)
		c.RegisterAll(&world)

		addIDs := make([]ecs.ID, len(add))
		remIDs := make([]ecs.ID, len(rem))
		for i, t := range add {
			addIDs[i] = ecs.TypeID(&world, t)
		}
		for i, t := range rem {
			remIDs[i] = ecs.TypeID(&world, t)
		}

		entities := make([]ecs.Entity, count)
		for i := 0; i < count; i++ {
			e := world.NewEntity(addIDs...)
			entities[i] = e
		}
		b.StartTimer()

		for _, e := range entities {
			world.Remove(e, remIDs...)
		}
	}
}

func BenchmarkArcheMove_1C_1_000(b *testing.B) {
	runArcheMove(b, 1000,
		g.T1[c.TestStruct0](),
		g.T1[c.TestStruct0](),
	)
}

func BenchmarkArcheMove_5C_1_000(b *testing.B) {
	runArcheMove(b, 1000,
		g.T5[c.TestStruct0, c.TestStruct1, c.TestStruct2, c.TestStruct3, c.TestStruct4](),
		g.T1[c.TestStruct0](),
	)
}
