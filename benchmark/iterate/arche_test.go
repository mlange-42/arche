package iterate

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func runArcheQuery(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[position](&world)
	rotID := ecs.ComponentID[rotation](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}
	b.StartTimer()

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

func runArcheQueryGeneric(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[position](&world)
	rotID := ecs.ComponentID[rotation](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := ecs.Query2[position, rotation](&world)
		b.StartTimer()
		for query.Next() {
			pos := query.Get1()
			_ = pos
		}
	}
}

func runArcheQuery1kArch(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	registerAll(&world)

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
			world.Add(entity, add...)
		}
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(6)
		b.StartTimer()
		for query.Next() {
			pos := (*position)(query.Get(6))
			_ = pos
		}
	}
}

func runArcheQueryEntity(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[position](&world)
	rotID := ecs.ComponentID[rotation](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(posID, rotID)
		b.StartTimer()
		for query.Next() {
			e := query.Entity()
			_ = e
		}
	}
}

func runArcheWorld(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[position](&world)
	rotID := ecs.ComponentID[rotation](&world)

	for i := 0; i < count; i++ {
		entity := world.NewEntity()
		world.Add(entity, posID, rotID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(posID, rotID)
		b.StartTimer()
		for query.Next() {
			entity := query.Entity()
			pos := (*position)(world.Get(entity, posID))
			_ = pos
		}
	}
}

func BenchmarkArcheIterQuery_1_000(b *testing.B) {
	runArcheQuery(b, 1000)
}

func BenchmarkArcheIterQuery_10_000(b *testing.B) {
	runArcheQuery(b, 10000)
}

func BenchmarkArcheIterQuery_100_000(b *testing.B) {
	runArcheQuery(b, 100000)
}

func BenchmarkArcheIterQueryGeneric_1_000(b *testing.B) {
	runArcheQueryGeneric(b, 1000)
}

func BenchmarkArcheIterQueryGeneric_10_000(b *testing.B) {
	runArcheQueryGeneric(b, 10000)
}

func BenchmarkArcheIterQueryGeneric_100_000(b *testing.B) {
	runArcheQueryGeneric(b, 100000)
}

func BenchmarkArcheIterQueryEntity_1_000(b *testing.B) {
	runArcheQueryEntity(b, 1000)
}

func BenchmarkArcheIterQueryEntity_10_000(b *testing.B) {
	runArcheQueryEntity(b, 10000)
}

func BenchmarkArcheIterQueryEntity_100_000(b *testing.B) {
	runArcheQueryEntity(b, 100000)
}

func BenchmarkArcheIterWorld_1_000(b *testing.B) {
	runArcheWorld(b, 1000)
}

func BenchmarkArcheIterWorld_10_000(b *testing.B) {
	runArcheWorld(b, 10000)
}

func BenchmarkArcheIterWorld_100_000(b *testing.B) {
	runArcheWorld(b, 100000)
}

func BenchmarkArcheIter1kArch_1_000(b *testing.B) {
	runArcheQuery1kArch(b, 1000)
}

func BenchmarkArcheIter1kArch_10_000(b *testing.B) {
	runArcheQuery1kArch(b, 10000)
}

func BenchmarkArcheIter1kArch_100_000(b *testing.B) {
	runArcheQuery1kArch(b, 100000)
}
