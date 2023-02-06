package iterate

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

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

func BenchmarkIterArcheWorld1000(b *testing.B) {
	runArcheWorld(b, 1000)
}

func BenchmarkIterArcheWorld10000(b *testing.B) {
	runArcheWorld(b, 10000)
}
