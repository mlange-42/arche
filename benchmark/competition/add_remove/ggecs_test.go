package addremove

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
)

const (
	PositionComponentID ecs.ComponentID = iota
	VelocityComponentID
)

func BenchmarkIterGGEcs(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(1024)
	world.Register(ecs.NewComponentRegistry[Position](PositionComponentID))
	world.Register(ecs.NewComponentRegistry[Velocity](VelocityComponentID))

	for i := 0; i < nEntities; i++ {
		_ = world.NewEntity(PositionComponentID)
	}

	posMask := ecs.MakeComponentMask(PositionComponentID)
	posVelMask := ecs.MakeComponentMask(PositionComponentID, VelocityComponentID)

	entities := make([]ecs.EntityID, 0, nEntities)

	// Iterate once for more fairness
	query := world.Query(posMask)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	for _, e := range entities {
		world.AddComponent(e, VelocityComponentID)
	}

	entities = entities[:0]
	query = world.Query(posVelMask)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	for _, e := range entities {
		world.RemComponent(e, VelocityComponentID)
	}
	entities = entities[:0]

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(posMask)
		for query.Next() {
			entities = append(entities, query.Entity())
		}

		for _, e := range entities {
			world.AddComponent(e, VelocityComponentID)
		}

		entities = entities[:0]
		query = world.Query(posVelMask)
		for query.Next() {
			entities = append(entities, query.Entity())
		}

		for _, e := range entities {
			world.RemComponent(e, VelocityComponentID)
		}
		entities = entities[:0]
	}
}

func BenchmarkBuildGGEcs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := ecs.NewWorld(1024)
		world.Register(ecs.NewComponentRegistry[Position](PositionComponentID))
		world.Register(ecs.NewComponentRegistry[Velocity](VelocityComponentID))

		for i := 0; i < nEntities; i++ {
			_ = world.NewEntity(PositionComponentID)
		}
	}
}
