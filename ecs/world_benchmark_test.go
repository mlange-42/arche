package ecs

import (
	"testing"

	"github.com/mlange-42/arche/ecs/stats"
)

func BenchmarkEntityAlive_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld(NewConfig().WithCapacityIncrement(1024))
	posID := ComponentID[Position](&world)

	entities := make([]Entity, 0, 1000)
	q := world.newEntitiesQuery(1000, ID{}, false, Entity{}, posID)
	for q.Next() {
		entities = append(entities, q.Entity())
	}

	b.StartTimer()

	var alive bool
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			alive = world.Alive(e)
		}
	}

	_ = alive
}

func BenchmarkGetResource(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	AddResource(&w, &Position{1, 2})
	posID := ResourceID[Position](&w)

	b.StartTimer()

	var res *Position
	for i := 0; i < b.N; i++ {
		res = w.Resources().Get(posID).(*Position)
	}

	_ = res
}

func BenchmarkGetResourceShortcut(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	AddResource(&w, &Position{1, 2})

	b.StartTimer()

	var res *Position
	for i := 0; i < b.N; i++ {
		res = GetResource[Position](&w)
	}

	_ = res
}

func BenchmarkNewEntities_10_000_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := NewWorld(NewConfig().WithCapacityIncrement(1024))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		for i := 0; i < 10000; i++ {
			_ = world.NewEntity(posID, velID)
		}
	}
}

func BenchmarkNewEntitiesBatch_10_000_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := NewWorld(NewConfig().WithCapacityIncrement(1024))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		world.newEntities(10000, ID{}, false, Entity{}, posID, velID)
	}
}

func BenchmarkNewEntities_10_000_Reset(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(1024))

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	for i := 0; i < 10000; i++ {
		_ = world.NewEntity(posID, velID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		world.Reset()
		for i := 0; i < 10000; i++ {
			_ = world.NewEntity(posID, velID)
		}
	}
}

func BenchmarkNewEntitiesBatch_10_000_Reset(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(1024))

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	for i := 0; i < 10000; i++ {
		_ = world.NewEntity(posID, velID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		world.Reset()
		world.newEntities(10000, ID{}, false, Entity{}, posID, velID)
	}
}

func BenchmarkRemoveEntities_10_000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := NewWorld(NewConfig().WithCapacityIncrement(10000))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		entities := make([]Entity, 10000)
		q := world.newEntitiesQuery(10000, ID{}, false, Entity{}, posID, velID)

		cnt := 0
		for q.Next() {
			entities[cnt] = q.Entity()
			cnt++
		}

		b.StartTimer()

		for _, e := range entities {
			world.RemoveEntity(e)
		}
	}
}

func BenchmarkWorldNewQuery(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(10000))
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	NewBuilder(&world, posID, velID).NewBatch(25)

	filter := All(posID, velID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q := world.Query(filter)
		q.Close()
	}
}

func BenchmarkWorldNewQueryNext(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(10000))

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	NewBuilder(&world, posID, velID).NewBatch(25)

	filter := All(posID, velID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q := world.Query(filter)
		q.Next()
		q.Close()
	}
}

func BenchmarkWorldNewQueryCached(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(10000))
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	NewBuilder(&world, posID, velID).NewBatch(25)

	filter := All(posID, velID)
	cf := world.Cache().Register(filter)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q := world.Query(&cf)
		q.Close()
	}
}

func BenchmarkWorldNewQueryNextCached(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(10000))

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	NewBuilder(&world, posID, velID).NewBatch(25)

	filter := All(posID, velID)
	cf := world.Cache().Register(filter)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q := world.Query(&cf)
		q.Next()
		q.Close()
	}
}

func BenchmarkRemoveEntitiesBatch_10_000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := NewWorld(NewConfig().WithCapacityIncrement(10000))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		q := world.newEntitiesQuery(10000, ID{}, false, Entity{}, posID, velID)
		q.Close()
		b.StartTimer()
		world.Batch().RemoveEntities(All(posID, velID))
	}
}

func BenchmarkWorldStats_1Arch(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	w.NewEntity()

	b.StartTimer()

	var st *stats.WorldStats
	for i := 0; i < b.N; i++ {
		st = w.Stats()
	}
	_ = st
}

func BenchmarkWorldStats_10Arch(b *testing.B) {
	b.StopTimer()

	w := NewWorld()

	ids := []ID{
		ComponentID[testStruct0](&w),
		ComponentID[testStruct1](&w),
		ComponentID[testStruct2](&w),
		ComponentID[testStruct3](&w),
		ComponentID[testStruct4](&w),
		ComponentID[testStruct5](&w),
		ComponentID[testStruct6](&w),
		ComponentID[testStruct7](&w),
		ComponentID[testStruct8](&w),
		ComponentID[testStruct9](&w),
	}

	for _, id := range ids {
		w.NewEntity(id)
	}

	b.StartTimer()

	var st *stats.WorldStats
	for i := 0; i < b.N; i++ {
		st = w.Stats()
	}
	_ = st
}
