package addremove

import (
	"testing"

	ecs "github.com/mlange-42/arche/benchmark/competition/add_remove/Entitas"
)

func BenchmarkIterEntitas(b *testing.B) {
	b.StopTimer()

	contexts := ecs.CreateContexts()
	game := contexts.Entitas()

	// System registration
	systems := ecs.CreateSystemPool()
	systems.Add(&AddVel{})
	systems.Add(&RemVel{})

	for i := 0; i < nEntities; i++ {
		e := game.CreateEntity()
		e.AddPosition(0, 0)
	}

	// GameLoop
	systems.Init(contexts)

	// Iterate once for more fairness
	systems.Execute()
	systems.Clean()

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
		systems.Add(&AddVel{})
		systems.Add(&RemVel{})

		for i := 0; i < nEntities; i++ {
			e := game.CreateEntity()
			e.AddPosition(0, 0)
		}

		// GameLoop
		systems.Init(contexts)
		systems.Exit(contexts)
	}
}

type AddVel struct {
	group ecs.Group
}

func (s *AddVel) Initer(contexts ecs.Contexts) {
	game := contexts.Entitas()
	s.group = game.Group(ecs.NewMatcher().AllOf(ecs.Position))
}

func (s *AddVel) Executer() {
	for _, e := range s.group.GetEntities() {
		e.AddVelocity(0, 0)
	}
}

type RemVel struct {
	group ecs.Group
}

func (s *RemVel) Initer(contexts ecs.Contexts) {
	game := contexts.Entitas()
	s.group = game.Group(ecs.NewMatcher().AllOf(ecs.Position, ecs.Velocity))
}

func (s *RemVel) Executer() {
	for _, e := range s.group.GetEntities() {
		e.RemoveVelocity()
	}
}
