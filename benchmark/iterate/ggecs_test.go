package iterate

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
	c "github.com/mlange-42/arche/benchmark/common"
)

func runGameEngineEcs(b *testing.B, count int) {
	b.StopTimer()
	comps := []ecs.ComponentConfig{
		{ID: 0, Component: c.Position{}},
		{ID: 1, Component: c.Rotation{}},
	}
	world := ecs.NewWorld(comps...)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Assign(entity, 0, 1)
	}
	filter := world.NewFilter(0, 1)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range filter.Entities() {
			pos := (*c.Position)(world.Component(e, 0))
			_ = pos
		}
	}
}

func runGameEngineEcs5C(b *testing.B, count int) {
	b.StopTimer()
	comps := []ecs.ComponentConfig{
		{ID: 0, Component: c.TestStruct0{}},
		{ID: 1, Component: c.TestStruct1{}},
		{ID: 2, Component: c.TestStruct2{}},
		{ID: 3, Component: c.TestStruct3{}},
		{ID: 4, Component: c.TestStruct4{}},
	}
	world := ecs.NewWorld(comps...)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Assign(entity, 0, 1, 2, 3, 4)
	}
	filter := world.NewFilter(0, 1, 2, 3, 4)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range filter.Entities() {
			t1 := (*c.TestStruct0)(world.Component(e, 0))
			t2 := (*c.TestStruct1)(world.Component(e, 1))
			t3 := (*c.TestStruct2)(world.Component(e, 2))
			t4 := (*c.TestStruct3)(world.Component(e, 3))
			t5 := (*c.TestStruct4)(world.Component(e, 4))
			_, _, _, _, _ = t1, t2, t3, t4, t5
		}
	}
}

func runGameEngineEcs1kArch(b *testing.B, count int) {
	b.StopTimer()
	comps := []ecs.ComponentConfig{
		{ID: 0, Component: c.TestStruct0{}},
		{ID: 1, Component: c.TestStruct1{}},
		{ID: 2, Component: c.TestStruct2{}},
		{ID: 3, Component: c.TestStruct3{}},
		{ID: 4, Component: c.TestStruct4{}},
		{ID: 5, Component: c.TestStruct5{}},
		{ID: 6, Component: c.TestStruct6{}},
		{ID: 7, Component: c.TestStruct7{}},
		{ID: 8, Component: c.TestStruct8{}},
		{ID: 9, Component: c.TestStruct9{}},
	}
	world := ecs.NewWorld(comps...)

	perArch := 2 * count / 1000

	for i := 0; i < 1024; i++ {
		mask := i
		add := make([]ecs.ID, 0, 10)
		for j := 0; j < 10; j++ {
			id := ecs.ID(j)
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			entity := world.NewEntity()
			world.Assign(entity, add...)
		}
	}

	filter := world.NewFilter(6)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range filter.Entities() {
			pos := (*c.Position)(world.Component(e, 6))
			_ = pos
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
