package iterate

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
	TestStruct10ID
	PositionComponentID
	RotationComponentID
)

func runGameEngineEcs(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(1024)
	world.Register(ecs.NewComponentRegistry[c.Position](PositionComponentID))
	world.Register(ecs.NewComponentRegistry[c.Rotation](RotationComponentID))

	for i := 0; i < count; i++ {
		world.NewEntity(PositionComponentID, RotationComponentID)
	}

	mask := ecs.MakeComponentMask(PositionComponentID, RotationComponentID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(mask)
		for query.Next() {
			pos := (*c.Position)(query.Component(PositionComponentID))
			pos.X = 1.0
		}
	}
}

func runGameEngineEcs5C(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(1024)
	world.Register(ecs.NewComponentRegistry[c.TestStruct0](TestStruct0ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct1](TestStruct1ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct2](TestStruct2ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct3](TestStruct3ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct4](TestStruct4ID))

	for i := 0; i < count; i++ {
		_ = world.NewEntity(TestStruct0ID, TestStruct1ID, TestStruct2ID, TestStruct3ID, TestStruct4ID)
	}
	mask := ecs.MakeComponentMask(TestStruct0ID, TestStruct1ID, TestStruct2ID, TestStruct3ID, TestStruct4ID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(mask)
		for query.Next() {
			t1 := (*c.TestStruct0)(query.Component(TestStruct0ID))
			t2 := (*c.TestStruct1)(query.Component(TestStruct1ID))
			t3 := (*c.TestStruct2)(query.Component(TestStruct2ID))
			t4 := (*c.TestStruct3)(query.Component(TestStruct3ID))
			t5 := (*c.TestStruct4)(query.Component(TestStruct4ID))
			t1.Val, t2.Val, t3.Val, t4.Val, t5.Val = 1, 1, 1, 1, 1
		}
	}
}

func runGameEngineEcs1kArch(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(1024)
	world.Register(ecs.NewComponentRegistry[c.TestStruct0](TestStruct0ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct1](TestStruct1ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct2](TestStruct2ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct3](TestStruct3ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct4](TestStruct4ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct5](TestStruct5ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct6](TestStruct6ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct7](TestStruct7ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct8](TestStruct8ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct9](TestStruct9ID))

	perArch := 2 * count / 1000

	for i := 0; i < 1024; i++ {
		mask := i
		add := make([]uint, 0, 10)
		for j := 0; j < 10; j++ {
			id := uint(j)
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			_ = world.NewEntity(add...)
		}
	}

	mask := ecs.MakeComponentMask(TestStruct6ID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(mask)
		for query.Next() {
			t1 := (*c.TestStruct6)(query.Component(TestStruct6ID))
			t1.Val = 1
		}
	}
}

func runGameEngineEcs1Of1kArch(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld(1024)
	world.Register(ecs.NewComponentRegistry[c.TestStruct0](TestStruct0ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct1](TestStruct1ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct2](TestStruct2ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct3](TestStruct3ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct4](TestStruct4ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct5](TestStruct5ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct6](TestStruct6ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct7](TestStruct7ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct8](TestStruct8ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct9](TestStruct9ID))
	world.Register(ecs.NewComponentRegistry[c.TestStruct10](TestStruct10ID))

	perArch := 2 * count / 1000

	for i := 0; i < 1024; i++ {
		mask := i
		add := make([]uint, 0, 10)
		for j := 0; j < 10; j++ {
			id := uint(j)
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			_ = world.NewEntity(add...)
		}
	}
	for i := 0; i < count; i++ {
		world.NewEntity(TestStruct10ID)
	}

	mask := ecs.MakeComponentMask(TestStruct10ID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(mask)
		for query.Next() {
			t1 := (*c.TestStruct10)(query.Component(TestStruct10ID))
			t1.Val = 1
		}
	}
}

func BenchmarkGGEcsIter_1_000(b *testing.B) {
	runGameEngineEcs(b, 1000)
}

func BenchmarkGGEcsIter_10_000(b *testing.B) {
	runGameEngineEcs(b, 10000)
}

func BenchmarkGGEcsIter_100_000(b *testing.B) {
	runGameEngineEcs(b, 100000)
}

func BenchmarkGGEcsIter5C_1_000(b *testing.B) {
	runGameEngineEcs5C(b, 1000)
}

func BenchmarkGGEcsIter5C_10_000(b *testing.B) {
	runGameEngineEcs5C(b, 10000)
}

func BenchmarkGGEcsIter5C_100_000(b *testing.B) {
	runGameEngineEcs5C(b, 100000)
}

func BenchmarkGGEcsIter1kArch_1_000(b *testing.B) {
	runGameEngineEcs1kArch(b, 1000)
}

func BenchmarkGGEcsIter1kArch_10_000(b *testing.B) {
	runGameEngineEcs1kArch(b, 10000)
}

func BenchmarkGGEcsIter1kArch_100_000(b *testing.B) {
	runGameEngineEcs1kArch(b, 100000)
}

func BenchmarkGGEcsIter1Of1kArch_1_000(b *testing.B) {
	runGameEngineEcs1Of1kArch(b, 1000)
}

func BenchmarkGGEcsIter1Of1kArch_10_000(b *testing.B) {
	runGameEngineEcs1Of1kArch(b, 10000)
}

func BenchmarkGGEcsIter1Of1kArch_100_000(b *testing.B) {
	runGameEngineEcs1Of1kArch(b, 100000)
}
