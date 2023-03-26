package addremove

import (
	"testing"

	"github.com/wfranczyk/ento"
)

func BenchmarkIterEnto(b *testing.B) {
	b.StopTimer()
	world := ento.NewWorldBuilder().
		WithDenseComponents(Position{}).
		WithSparseComponents(Velocity{}).
		Build(1024)

	add := AddVelSystem{}
	rem := RemVelSystem{}
	world.AddSystems(&add, &rem)

	for i := 0; i < nEntities; i++ {
		world.AddEntity(Position{})
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkBuildEnto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := ento.NewWorldBuilder().
			WithDenseComponents(Position{}).
			WithSparseComponents(Velocity{}).
			Build(1024)

		add := AddVelSystem{}
		rem := RemVelSystem{}
		world.AddSystems(&add, &rem)

		for i := 0; i < nEntities; i++ {
			world.AddEntity(Position{})
		}
	}
}

type AddVelSystem struct {
	Pos *Position `ento:"required"`
}

func (s *AddVelSystem) Update(entity *ento.Entity) {
	entity.Set(Position{})
}

type RemVelSystem struct {
	Pos *Position `ento:"required"`
	Vel *Velocity `ento:"required"`
}

func (s *RemVelSystem) Update(entity *ento.Entity) {
	entity.Rem(Position{})
}
