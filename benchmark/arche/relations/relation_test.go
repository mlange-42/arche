package relations

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchmarkRelation(b *testing.B, numParents int, numChildren int) {
	b.StopTimer()

	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	parentMapper := generic.NewMap1[ParentList](&world)
	childMapper := generic.NewMap1[ChildRelation](&world, generic.T[ChildRelation]())

	spawnedPar := parentMapper.NewQuery(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		parents = append(parents, spawnedPar.Entity())
	}

	for _, par := range parents {
		spawnedChild := childMapper.NewQuery(numChildren, par)
		for spawnedChild.Next() {
			child := spawnedChild.Get()
			child.Value = 1
		}
	}

	comp := generic.T[ChildRelation]()
	parentFilter := generic.NewFilter1[ParentList]()
	parentFilter.Register(&world)

	childFilter := generic.NewFilter1[ChildRelation]().WithRelation(comp, ecs.Entity{})
	childF := childFilter.Filter(&world)
	childID := ecs.ComponentID[ChildRelation](&world)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := parentFilter.Query(&world)
		for query.Next() {
			parData := query.Get()
			par := query.Entity()

			cf := ecs.RelationFilter{Filter: &childF, Target: par}

			childQuery := world.Query(&cf)
			for childQuery.Next() {
				child := (*ChildRelation)(childQuery.Get(childID))
				parData.Value += child.Value
			}
		}
	}

	b.StopTimer()

	parQuery := parentFilter.Query(&world)

	expected := numChildren * b.N
	for parQuery.Next() {
		par := parQuery.Get()
		if par.Value != expected {
			panic("wrong number of children")
		}
	}
}

func BenchmarkRelation_10_x_10(b *testing.B) {
	benchmarkRelation(b, 10, 10)
}

func BenchmarkRelation_100_x_10(b *testing.B) {
	benchmarkRelation(b, 100, 10)
}

func BenchmarkRelation_1000_x_10(b *testing.B) {
	benchmarkRelation(b, 1000, 10)
}

func BenchmarkRelation_10000_x_10(b *testing.B) {
	benchmarkRelation(b, 10000, 10)
}

func BenchmarkRelation_10_x_100(b *testing.B) {
	benchmarkRelation(b, 10, 100)
}

func BenchmarkRelation_100_x_100(b *testing.B) {
	benchmarkRelation(b, 100, 100)
}

func BenchmarkRelation_1000_x_100(b *testing.B) {
	benchmarkRelation(b, 1000, 100)
}

func BenchmarkRelation_10000_x_100(b *testing.B) {
	benchmarkRelation(b, 10000, 100)
}

func BenchmarkRelation_10_x_1000(b *testing.B) {
	benchmarkRelation(b, 10, 1000)
}

func BenchmarkRelation_100_x_1000(b *testing.B) {
	benchmarkRelation(b, 100, 1000)
}

func BenchmarkRelation_1000_x_1000(b *testing.B) {
	benchmarkRelation(b, 1000, 1000)
}

func BenchmarkRelation_10_x_10000(b *testing.B) {
	benchmarkRelation(b, 10, 10000)
}

func BenchmarkRelation_100_x_10000(b *testing.B) {
	benchmarkRelation(b, 100, 10000)
}
