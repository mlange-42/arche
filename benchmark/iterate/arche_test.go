package iterate

import (
	"reflect"
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func runArcheArchetype(b *testing.B, count int) {
	world := ecs.NewWorld()

	comps := []ecs.ComponentType{
		{ID: 0, Type: reflect.TypeOf(position{})},
		{ID: 1, Type: reflect.TypeOf(rotation{})},
	}

	arch := ecs.NewArchetype(comps...)

	for i := 0; i < count; i++ {
		arch.Add(
			world.NewEntity(),
			ecs.Component{ID: 0, Component: &position{1, 2}},
			ecs.Component{ID: 1, Component: &rotation{3}},
		)
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < count; j++ {
			pos := (*position)(arch.Get(i, ecs.ID(0)))
			_ = pos
		}
	}
}

func BenchmarkIterArcheArchetype100(b *testing.B) {
	runArcheArchetype(b, 100)
}

func BenchmarkIterArcheArchetype1000(b *testing.B) {
	runArcheArchetype(b, 1000)
}

func BenchmarkIterArcheArchetype10000(b *testing.B) {
	runArcheArchetype(b, 10000)
}

func runArcheWorld(b *testing.B, count int) {
	world := ecs.NewWorld()

	posID := ecs.RegisterComponent[position](&world)
	rotID := ecs.RegisterComponent[rotation](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(posID, rotID)
		b.StartTimer()
		for query.Next() {
			pos := (*position)(query.Get(posID))
			_ = pos
		}
	}
}

func BenchmarkIterArcheWorld100(b *testing.B) {
	runArcheWorld(b, 100)
}

func BenchmarkIterArcheWorld1000(b *testing.B) {
	runArcheWorld(b, 1000)
}

func BenchmarkIterArcheWorld10000(b *testing.B) {
	runArcheWorld(b, 10000)
}
