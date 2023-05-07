package relations

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchmarkChild(b *testing.B, numParents int, numChildren int) {
	b.StopTimer()

	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	parentMapper := generic.NewMap1[ParentList](&world)
	childMapper := generic.NewMap1[Child](&world)

	spawnedPar := parentMapper.NewQuery(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		parents = append(parents, spawnedPar.Entity())
	}

	spawnedChild := childMapper.NewQuery(numParents * numChildren)
	cnt := 0
	for spawnedChild.Next() {
		child := spawnedChild.Get()
		child.Value = 1
		child.Parent = parents[cnt/numChildren]
		cnt++
	}

	parentFilter := generic.NewFilter1[ParentList]()
	parentFilter.Register(&world)
	childFilter := generic.NewFilter1[Child]()
	childFilter.Register(&world)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := childFilter.Query(&world)
		for query.Next() {
			child := query.Get()
			parData := parentMapper.Get(child.Parent)

			parData.Value += child.Value
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

func BenchmarkRelationChild_1k_10_x_100(b *testing.B) {
	benchmarkChild(b, 10, 100)
}

func BenchmarkRelationChild_1k_100_x_10(b *testing.B) {
	benchmarkChild(b, 100, 10)
}

func BenchmarkRelationChild_10k_10_x_1000(b *testing.B) {
	benchmarkChild(b, 10, 1000)
}

func BenchmarkRelationChild_10k_100_x_100(b *testing.B) {
	benchmarkChild(b, 100, 100)
}

func BenchmarkRelationChild_10k_1000_x_10(b *testing.B) {
	benchmarkChild(b, 1000, 10)
}

func BenchmarkRelationChild_100k_10_x_10000(b *testing.B) {
	benchmarkChild(b, 10, 10000)
}

func BenchmarkRelationChild_100k_100_x_1000(b *testing.B) {
	benchmarkChild(b, 100, 1000)
}

func BenchmarkRelationChild_100k_1000_x_100(b *testing.B) {
	benchmarkChild(b, 1000, 100)
}

func BenchmarkRelationChild_100k_10000_x_10(b *testing.B) {
	benchmarkChild(b, 10000, 10)
}

func BenchmarkRelationChild_1M_100_x_10000(b *testing.B) {
	benchmarkChild(b, 100, 10000)
}

func BenchmarkRelationChild_1M_1000_x_1000(b *testing.B) {
	benchmarkChild(b, 1000, 1000)
}

func BenchmarkRelationChild_1M_10000_x_100(b *testing.B) {
	benchmarkChild(b, 10000, 100)
}
