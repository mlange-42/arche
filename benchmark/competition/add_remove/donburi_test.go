package addremove

import (
	"testing"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

func BenchmarkIterDonburi(b *testing.B) {
	b.StopTimer()
	world := donburi.NewWorld()

	var position = donburi.NewComponentType[Position]()
	var velocity = donburi.NewComponentType[Velocity]()

	for i := 0; i < nEntities; i++ {
		world.Create(position)
	}

	queryPos := donburi.NewQuery(filter.Contains(position))
	queryPosVel := donburi.NewQuery(filter.Contains(position, velocity))

	// Iterate once for more fairness
	queryPos.Each(world, func(entry *donburi.Entry) {
		entry.AddComponent(velocity)
	})
	queryPosVel.Each(world, func(entry *donburi.Entry) {
		entry.RemoveComponent(velocity)
	})

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		queryPos.Each(world, func(entry *donburi.Entry) {
			entry.AddComponent(velocity)
		})
		queryPosVel.Each(world, func(entry *donburi.Entry) {
			entry.RemoveComponent(velocity)
		})
	}
}

func BenchmarkBuildDonburi(b *testing.B) {
	var position = donburi.NewComponentType[Position]()

	for i := 0; i < b.N; i++ {
		world := donburi.NewWorld()
		for i := 0; i < nEntities; i++ {
			world.Create(position)
		}
	}
}
