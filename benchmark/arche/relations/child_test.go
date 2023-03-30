package relations

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchmarkChild(b *testing.B, numParents int) {
	b.StopTimer()

	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	parentMapper := generic.NewMap1[ParentData](&world)
	childMapper := generic.NewMap2[ChildData, Child](&world)

	spawnedPar := parentMapper.NewEntitiesQuery(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		parents = append(parents, spawnedPar.Entity())
	}

	spawnedChild := childMapper.NewEntitiesQuery(numParents * numChildren)
	cnt := 0
	for spawnedChild.Next() {
		data, child := spawnedChild.Get()
		data.Value = 1
		child.Parent = parents[cnt/numChildren]
		cnt++
	}

	parentFilter := generic.NewFilter1[ParentData]()
	parentFilter.Register(&world)
	childFilter := generic.NewFilter2[ChildData, Child]()
	childFilter.Register(&world)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := childFilter.Query(&world)
		for query.Next() {
			data, child := query.Get()
			parData := parentMapper.Get(child.Parent)

			parData.Value += data.Value
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

func BenchmarkRelationChild_100_x_10(b *testing.B) {
	benchmarkChild(b, 100)
}

func BenchmarkRelationChild_1000_x_10(b *testing.B) {
	benchmarkChild(b, 1000)
}

func BenchmarkRelationChild_10000_x_10(b *testing.B) {
	benchmarkChild(b, 10000)
}
