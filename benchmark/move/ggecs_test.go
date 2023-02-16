package move

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
	c "github.com/mlange-42/arche/benchmark/common"
)

func runGameEngineEcsMove(b *testing.B, count int, add, rem []ecs.ID) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		comps := []ecs.ComponentConfig{
			{ID: 0, Component: c.TestStruct0{}},
			{ID: 1, Component: c.TestStruct1{}},
			{ID: 2, Component: c.TestStruct2{}},
			{ID: 3, Component: c.TestStruct3{}},
			{ID: 4, Component: c.TestStruct4{}},
			{ID: 5, Component: c.TestStruct5{}},
			{ID: 6, Component: c.TestStruct6{}},
		}
		world := ecs.NewWorld(comps...)

		for i := 0; i < count; i++ {
			entity := world.NewEntity()
			world.Assign(entity, add...)
		}
		filter := world.NewFilter(add...)

		b.StartTimer()

		for _, e := range filter.Entities() {
			world.Remove(e, rem...)
		}

	}
}

func BenchmarkGGECSMove_1C_1_000(b *testing.B) {
	runGameEngineEcsMove(b, 1000,
		[]ecs.ID{0},
		[]ecs.ID{0},
	)
}

func BenchmarkGGECSMove_5C_1_000(b *testing.B) {
	runGameEngineEcsMove(b, 1000,
		[]ecs.ID{0, 1, 2, 3, 4},
		[]ecs.ID{0},
	)
}
