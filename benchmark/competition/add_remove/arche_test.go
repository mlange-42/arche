package addremove

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func BenchmarkIterArche(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)
	ids := []ecs.ID{velID}

	world.Batch().NewEntities(nEntities, posID)

	var filterPos ecs.Filter = ecs.All(posID)
	var filterPosVel ecs.Filter = ecs.All(posID, velID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		entities := make([]ecs.Entity, 0, nEntities)
		query := world.Query(filterPos)
		for query.Next() {
			entities = append(entities, query.Entity())
		}

		for _, e := range entities {
			world.Add(e, ids...)
		}

		entities = entities[:0]
		query = world.Query(filterPosVel)
		for query.Next() {
			entities = append(entities, query.Entity())
		}

		for _, e := range entities {
			world.Remove(e, ids...)
		}
	}
}

func BenchmarkBuildArche(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

		posID := ecs.ComponentID[Position](&world)
		velID := ecs.ComponentID[Velocity](&world)
		ids := []ecs.ID{velID}

		for i := 0; i < nEntities; i++ {
			world.NewEntity(ids...)
		}

		var filterPos ecs.Filter = ecs.All(posID)
		var filterPosVel ecs.Filter = ecs.All(posID, velID)

		_ = filterPos
		_ = filterPosVel
	}
}

func BenchmarkIterArcheGeneric(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	posMapper := generic.NewMap1[Position](&world)
	velMapper := generic.NewMap1[Velocity](&world)

	posMapper.NewEntities(nEntities)

	filterPos := generic.NewFilter1[Position]()
	filterPosVel := generic.NewFilter2[Position, Velocity]()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		entities := make([]ecs.Entity, 0, nEntities)
		query := filterPos.Query(&world)
		for query.Next() {
			entities = append(entities, query.Entity())
		}

		for _, e := range entities {
			velMapper.Add(e)
		}

		entities = entities[:0]
		query2 := filterPosVel.Query(&world)
		for query2.Next() {
			entities = append(entities, query2.Entity())
		}

		for _, e := range entities {
			velMapper.Remove(e)
		}
	}
}

func BenchmarkBuildArcheGeneric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

		posID := ecs.ComponentID[Position](&world)
		velID := ecs.ComponentID[Velocity](&world)

		world.Batch().NewEntities(nEntities, posID)

		var filterPos ecs.Filter = ecs.All(posID)
		var filterPosVel ecs.Filter = ecs.All(posID, velID)

		_ = filterPos
		_ = filterPosVel
	}
}
