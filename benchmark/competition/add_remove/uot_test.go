package addremove

import (
	"testing"

	"github.com/unitoftime/ecs"
)

func BenchmarkIterUot(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld()

	queryPos := ecs.Query1[Position](world)
	queryPosVel := ecs.Query2[Position, Velocity](world)
	comp := ecs.C(Velocity{0, 0})

	for i := 0; i < nEntities; i++ {
		id := world.NewId()
		ecs.Write(world, id,
			ecs.C(Position{0, 0}),
		)
	}

	entities := make([]ecs.Id, 0, nEntities)

	// Iterate once for more fairness
	queryPos.MapId(func(id ecs.Id, pos *Position) {
		entities = append(entities, id)
	})

	for _, e := range entities {
		ecs.Write(world, e,
			ecs.C(Velocity{0, 0}),
		)
	}

	entities = entities[:0]

	queryPosVel.MapId(func(id ecs.Id, pos *Position, vel *Velocity) {
		entities = append(entities, id)
	})

	for _, e := range entities {
		ecs.DeleteComponent(world, e, comp)
	}

	entities = entities[:0]

	b.StartTimer()

	for i := 0; i < b.N; i++ {

		queryPos.MapId(func(id ecs.Id, pos *Position) {
			entities = append(entities, id)
		})

		for _, e := range entities {
			ecs.Write(world, e,
				ecs.C(Velocity{0, 0}),
			)
		}

		entities = entities[:0]

		queryPosVel.MapId(func(id ecs.Id, pos *Position, vel *Velocity) {
			entities = append(entities, id)
		})

		for _, e := range entities {
			ecs.DeleteComponent(world, e, comp)
		}

		entities = entities[:0]
	}
}

func BenchmarkBuildUot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := ecs.NewWorld()

		for i := 0; i < nEntities; i++ {
			id := world.NewId()
			ecs.Write(world, id,
				ecs.C(Position{0, 0}),
			)
		}

		queryPos := ecs.Query1[Position](world)
		queryPosVel := ecs.Query2[Position, Velocity](world)

		_ = queryPos
		_ = queryPosVel
	}
}
