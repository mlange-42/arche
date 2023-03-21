package posvel

import (
	"testing"

	ecs "github.com/mlange-42/arche/benchmark/competition/pos_vel/Entitas"
)

func BenchmarkIterEntitas(b *testing.B) {
	b.StopTimer()

	contexts := ecs.CreateContexts()
	game := contexts.Entitas()

	// System registration
	systems := ecs.CreateSystemPool()
	systems.Add(&Translate{})

	for i := 0; i < nPos; i++ {
		e := game.CreateEntity()
		e.AddPosition(0, 0)
	}
	for i := 0; i < nPosVel; i++ {
		e := game.CreateEntity()
		e.AddPosition(0, 0)
		e.AddVelocity(0, 0)
	}

	// GameLoop
	systems.Init(contexts)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		systems.Execute()
		systems.Clean()
	}
	systems.Exit(contexts)
}

func BenchmarkBuildEntitas(b *testing.B) {
	for i := 0; i < b.N; i++ {
		contexts := ecs.CreateContexts()
		game := contexts.Entitas()

		// System registration
		systems := ecs.CreateSystemPool()
		systems.Add(&Translate{})

		for i := 0; i < nPos; i++ {
			e := game.CreateEntity()
			e.AddPosition(0, 0)
		}
		for i := 0; i < nPosVel; i++ {
			e := game.CreateEntity()
			e.AddPosition(0, 0)
			e.AddVelocity(0, 0)
		}

		// GameLoop
		systems.Init(contexts)
		systems.Exit(contexts)
	}
}

type Translate struct {
	group ecs.Group
}

func (s *Translate) Initer(contexts ecs.Contexts) {
	game := contexts.Entitas()
	s.group = game.Group(ecs.NewMatcher().AllOf(ecs.Position, ecs.Velocity))
}

func (s *Translate) Executer() {
	for _, e := range s.group.GetEntities() {
		pos := e.GetPosition()
		vel := e.GetVelocity()
		e.ReplacePosition(pos.X+vel.X, pos.X+vel.Y)
	}
}
