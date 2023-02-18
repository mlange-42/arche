package posvel

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
)

func BenchmarkGGEcs(b *testing.B) {
	b.StopTimer()
	comps := []ecs.ComponentConfig{
		{ID: 0, Component: Position{}},
		{ID: 1, Component: Velocity{}},
	}
	world := ecs.NewWorld(comps...)

	for i := 0; i < 9000; i++ {
		entity := world.NewEntity()
		world.Assign(entity, 0)
	}
	for i := 0; i < 1000; i++ {
		entity := world.NewEntity()
		world.Assign(entity, 0, 1)
	}

	filter := world.NewFilter(0, 1)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range filter.Entities() {
			pos := (*Position)(world.Component(e, 0))
			vel := (*Velocity)(world.Component(e, 1))
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
