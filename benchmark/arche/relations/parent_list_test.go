package relations

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchmarkParentList(b *testing.B, numParents int, numChildren int) {
	b.StopTimer()

	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	parentMapper := generic.NewMap1[ParentList](&world)
	childMapper := generic.NewMap1[ChildList](&world)

	spawnedPar := parentMapper.NewQuery(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		parents = append(parents, spawnedPar.Entity())
	}

	spawnedChild := childMapper.NewQuery(numParents * numChildren)
	cnt := 0
	for spawnedChild.Next() {
		childEntity := spawnedChild.Entity()
		child := spawnedChild.Get()
		child.Value = 1
		par := parentMapper.Get(parents[cnt/numChildren])

		if !par.FirstChild.IsZero() {
			child.Next = par.FirstChild
		}
		par.FirstChild = childEntity

		cnt++
	}

	parentFilter := generic.NewFilter1[ParentList]()
	parentFilter.Register(&world)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := parentFilter.Query(&world)
		for query.Next() {
			par := query.Get()
			child := par.FirstChild
			for !child.IsZero() {
				childList := childMapper.Get(child)
				par.Value += childList.Value
				child = childList.Next
			}
		}
	}

	b.StopTimer()
	query := parentFilter.Query(&world)

	expected := numChildren * b.N
	for query.Next() {
		par := query.Get()
		if par.Value != expected {
			panic("wrong number of children")
		}
	}
}

func BenchmarkRelationParentList_10_x_10(b *testing.B) {
	benchmarkParentList(b, 10, 10)
}

func BenchmarkRelationParentList_100_x_10(b *testing.B) {
	benchmarkParentList(b, 100, 10)
}

func BenchmarkRelationParentList_1000_x_10(b *testing.B) {
	benchmarkParentList(b, 1000, 10)
}

func BenchmarkRelationParentList_10000_x_100(b *testing.B) {
	benchmarkParentList(b, 10000, 100)
}

func BenchmarkRelationParentList_10_x_100(b *testing.B) {
	benchmarkParentList(b, 10, 100)
}

func BenchmarkRelationParentList_100_x_100(b *testing.B) {
	benchmarkParentList(b, 100, 100)
}

func BenchmarkRelationParentList_1000_x_100(b *testing.B) {
	benchmarkParentList(b, 1000, 100)
}

func BenchmarkRelationParentList_10000_x_1000(b *testing.B) {
	benchmarkParentList(b, 10000, 1000)
}

func BenchmarkRelationParentList_10_x_1000(b *testing.B) {
	benchmarkParentList(b, 10, 1000)
}

func BenchmarkRelationParentList_100_x_1000(b *testing.B) {
	benchmarkParentList(b, 100, 1000)
}

func BenchmarkRelationParentList_1000_x_1000(b *testing.B) {
	benchmarkParentList(b, 1000, 1000)
}

func BenchmarkRelationParentList_10_x_10000(b *testing.B) {
	benchmarkParentList(b, 10, 10000)
}

func BenchmarkRelationParentList_100_x_10000(b *testing.B) {
	benchmarkParentList(b, 100, 10000)
}
