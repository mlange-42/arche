package move

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
	c "github.com/mlange-42/arche/benchmark/arche/common"
)

const (
	TestStruct0ID ecs.ComponentID = iota
	TestStruct1ID
	TestStruct2ID
	TestStruct3ID
	TestStruct4ID
	TestStruct5ID
	TestStruct6ID
	TestStruct7ID
	TestStruct8ID
	TestStruct9ID
)

func runGameEngineEcsMove(b *testing.B, count int, add, rem []uint) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := ecs.NewWorld(1024)
		world.Register(ecs.NewComponentRegistry[c.TestStruct0](TestStruct0ID))
		world.Register(ecs.NewComponentRegistry[c.TestStruct1](TestStruct1ID))
		world.Register(ecs.NewComponentRegistry[c.TestStruct2](TestStruct2ID))
		world.Register(ecs.NewComponentRegistry[c.TestStruct3](TestStruct3ID))
		world.Register(ecs.NewComponentRegistry[c.TestStruct4](TestStruct4ID))
		world.Register(ecs.NewComponentRegistry[c.TestStruct5](TestStruct5ID))
		world.Register(ecs.NewComponentRegistry[c.TestStruct6](TestStruct6ID))

		entities := []ecs.EntityID{}
		for i := 0; i < count; i++ {
			entities = append(entities, world.NewEntity(add...))
		}
		b.StartTimer()

		for _, e := range entities {
			for _, comp := range rem {
				world.RemComponent(e, comp)
			}
		}

	}
}

func BenchmarkGGECSMove_1C_1_000(b *testing.B) {
	runGameEngineEcsMove(b, 1000,
		[]uint{0},
		[]uint{0},
	)
}

func BenchmarkGGECSMove_5C_1_000(b *testing.B) {
	runGameEngineEcsMove(b, 1000,
		[]uint{0, 1, 2, 3, 4},
		[]uint{0},
	)
}
