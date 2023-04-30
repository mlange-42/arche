package relations

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

type ChildFilter struct {
	Filter generic.Filter1[ChildRelation]
}

func benchmarkRelationCached(b *testing.B, numParents int, numChildren int) {
	b.StopTimer()

	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	parentMapper := generic.NewMap2[ParentList, ChildFilter](&world)
	childMapper := generic.NewMap1[ChildRelation](&world, generic.T[ChildRelation]())

	spawnedPar := parentMapper.NewQuery(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		par := spawnedPar.Entity()
		_, fl := spawnedPar.Get()
		fl.Filter = *generic.NewFilter1[ChildRelation]().WithRelation(generic.T[ChildRelation](), par)
		fl.Filter.Register(&world)

		parents = append(parents, par)
	}

	for _, par := range parents {
		spawnedChild := childMapper.NewQuery(numChildren, par)
		for spawnedChild.Next() {
			child := spawnedChild.Get()
			child.Value = 1
		}
	}

	parentFilter := generic.NewFilter2[ParentList, ChildFilter]()
	parentFilter.Register(&world)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := parentFilter.Query(&world)
		for query.Next() {
			parData, filter := query.Get()

			childQuery := filter.Filter.Query(&world)
			for childQuery.Next() {
				child := childQuery.Get()
				parData.Value += child.Value
			}
		}
	}

	b.StopTimer()

	parQuery := parentFilter.Query(&world)

	expected := numChildren * b.N
	for parQuery.Next() {
		par, _ := parQuery.Get()
		if par.Value != expected {
			panic("wrong number of children")
		}
	}
}

func BenchmarkRelationCached_10_x_10(b *testing.B) {
	benchmarkRelationCached(b, 10, 10)
}

func BenchmarkRelationCached_100_x_10(b *testing.B) {
	benchmarkRelationCached(b, 100, 10)
}

func BenchmarkRelationCached_1000_x_10(b *testing.B) {
	benchmarkRelationCached(b, 1000, 10)
}

func BenchmarkRelationCached_10000_x_10(b *testing.B) {
	benchmarkRelationCached(b, 10000, 10)
}

func BenchmarkRelationCached_10_x_100(b *testing.B) {
	benchmarkRelationCached(b, 10, 100)
}

func BenchmarkRelationCached_100_x_100(b *testing.B) {
	benchmarkRelationCached(b, 100, 100)
}

func BenchmarkRelationCached_1000_x_100(b *testing.B) {
	benchmarkRelationCached(b, 1000, 100)
}

func BenchmarkRelationCached_10000_x_100(b *testing.B) {
	benchmarkRelationCached(b, 10000, 100)
}

func BenchmarkRelationCached_10_x_1000(b *testing.B) {
	benchmarkRelationCached(b, 10, 1000)
}

func BenchmarkRelationCached_100_x_1000(b *testing.B) {
	benchmarkRelationCached(b, 100, 1000)
}

func BenchmarkRelationCached_1000_x_1000(b *testing.B) {
	benchmarkRelationCached(b, 1000, 1000)
}

func BenchmarkRelationCached_10_x_10000(b *testing.B) {
	benchmarkRelationCached(b, 10, 10000)
}

func BenchmarkRelationCached_100_x_10000(b *testing.B) {
	benchmarkRelationCached(b, 100, 10000)
}
