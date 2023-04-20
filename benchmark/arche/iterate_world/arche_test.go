package iterate

import (
	"math/rand"
	"testing"

	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func runArcheWorldGet(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	entities := make([]ecs.Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = world.NewEntity(posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			pos := (*c.Position)(world.Get(e, posID))
			pos.X = 1
		}
	}
}

func runArcheWorldGetGeneric(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	get := generic.NewMap1[c.Position](&world)

	entities := make([]ecs.Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = world.NewEntity(posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			pos := get.Get(e)
			pos.X = 1
		}
	}
}

func runArcheWorldGetUnchecked(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	entities := make([]ecs.Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = world.NewEntity(posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			pos := (*c.Position)(world.GetUnchecked(e, posID))
			pos.X = 1
		}
	}
}

func runArcheWorldGetGenericUnchecked(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	get := generic.NewMap1[c.Position](&world)

	entities := make([]ecs.Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = world.NewEntity(posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			pos := get.GetUnchecked(e)
			pos.X = 1
		}
	}
}

func runArcheWorldGetUncheckedShuffled(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	entities := make([]ecs.Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = world.NewEntity(posID, rotID)
	}
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			pos := (*c.Position)(world.GetUnchecked(e, posID))
			pos.X = 1
		}
	}
}

func runArcheWorldGetGenericUncheckedShuffled(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[c.Position](&world)
	rotID := ecs.ComponentID[c.Rotation](&world)

	get := generic.NewMap1[c.Position](&world)

	entities := make([]ecs.Entity, count)
	for i := 0; i < count; i++ {
		entities[i] = world.NewEntity(posID, rotID)
	}
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			pos := get.GetUnchecked(e)
			pos.X = 1
		}
	}
}

func BenchmarkArcheIterWorldID_1_000(b *testing.B) {
	runArcheWorldGet(b, 1000)
}

func BenchmarkArcheIterWorldID_10_000(b *testing.B) {
	runArcheWorldGet(b, 10000)
}

func BenchmarkArcheIterWorldID_100_000(b *testing.B) {
	runArcheWorldGet(b, 100000)
}

func BenchmarkArcheIterWorldGeneric_1_000(b *testing.B) {
	runArcheWorldGetGeneric(b, 1000)
}

func BenchmarkArcheIterWorldGeneric_10_000(b *testing.B) {
	runArcheWorldGetGeneric(b, 10000)
}

func BenchmarkArcheIterWorldGeneric_100_000(b *testing.B) {
	runArcheWorldGetGeneric(b, 100000)
}

func BenchmarkArcheIterWorldIDUnchecked_1_000(b *testing.B) {
	runArcheWorldGetUnchecked(b, 1000)
}

func BenchmarkArcheIterWorldIDUnchecked_10_000(b *testing.B) {
	runArcheWorldGetUnchecked(b, 10000)
}

func BenchmarkArcheIterWorldIDUnchecked_100_000(b *testing.B) {
	runArcheWorldGetUnchecked(b, 100000)
}

func BenchmarkArcheIterWorldGenericUnchecked_1_000(b *testing.B) {
	runArcheWorldGetGenericUnchecked(b, 1000)
}

func BenchmarkArcheIterWorldGenericUnchecked_10_000(b *testing.B) {
	runArcheWorldGetGenericUnchecked(b, 10000)
}

func BenchmarkArcheIterWorldGenericUnchecked_100_000(b *testing.B) {
	runArcheWorldGetGenericUnchecked(b, 100000)
}

func BenchmarkArcheIterWorldIDUncheckedShuffled_1_000(b *testing.B) {
	runArcheWorldGetUncheckedShuffled(b, 1000)
}

func BenchmarkArcheIterWorldIDUncheckedShuffled_10_000(b *testing.B) {
	runArcheWorldGetUncheckedShuffled(b, 10000)
}

func BenchmarkArcheIterWorldIDUncheckedShuffled_100_000(b *testing.B) {
	runArcheWorldGetUncheckedShuffled(b, 100000)
}

func BenchmarkArcheIterWorldGenericUncheckedShuffled_1_000(b *testing.B) {
	runArcheWorldGetGenericUncheckedShuffled(b, 1000)
}

func BenchmarkArcheIterWorldGenericUncheckedShuffled_10_000(b *testing.B) {
	runArcheWorldGetGenericUncheckedShuffled(b, 10000)
}

func BenchmarkArcheIterWorldGenericUncheckedShuffled_100_000(b *testing.B) {
	runArcheWorldGetGenericUncheckedShuffled(b, 100000)
}
