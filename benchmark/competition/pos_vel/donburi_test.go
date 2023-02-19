package posvel

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

	for i := 0; i < nPos; i++ {
		world.Create(position)
	}
	for i := 0; i < nPosVel; i++ {
		world.Create(position, velocity)
	}

	query := donburi.NewQuery(filter.Contains(position, velocity))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query.Each(world, func(entry *donburi.Entry) {
			pos := position.Get(entry)
			vel := velocity.Get(entry)

			pos.X += vel.X
			pos.Y += vel.Y
		})
	}
}

func BenchmarkBuildDonburi(b *testing.B) {
	var position = donburi.NewComponentType[Position]()
	var velocity = donburi.NewComponentType[Velocity]()

	for i := 0; i < b.N; i++ {
		world := donburi.NewWorld()

		for i := 0; i < nPos; i++ {
			world.Create(position)
		}
		for i := 0; i < nPosVel; i++ {
			world.Create(position, velocity)
		}
	}
}
