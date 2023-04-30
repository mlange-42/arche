package relations

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchmarkParentArray(b *testing.B, numParents int) {
	b.StopTimer()

	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	parentMapper := generic.NewMap1[ParentArr](&world)
	childMapper := generic.NewMap1[Child](&world)

	spawnedPar := parentMapper.NewQuery(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		parents = append(parents, spawnedPar.Entity())
	}

	spawnedChild := childMapper.NewQuery(numParents * numArrChildren)
	cnt := 0
	for spawnedChild.Next() {
		data := spawnedChild.Get()
		data.Value = 1
		par := parentMapper.Get(parents[cnt/numArrChildren])
		par.Children[cnt%numArrChildren] = spawnedChild.Entity()
		cnt++
	}

	parentFilter := generic.NewFilter1[ParentArr]()
	parentFilter.Register(&world)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := parentFilter.Query(&world)
		for query.Next() {
			par := query.Get()
			for i := 0; i < len(par.Children); i++ {
				childData := childMapper.Get(par.Children[i])
				par.Value += childData.Value
			}
		}
	}

	b.StopTimer()
	query := parentFilter.Query(&world)

	expected := numArrChildren * b.N
	for query.Next() {
		par := query.Get()
		if par.Value != expected {
			panic("wrong number of children")
		}
	}
}

func BenchmarkRelationParentArray_100_x_10(b *testing.B) {
	benchmarkParentArray(b, 100)
}

func BenchmarkRelationParentArray_1000_x_10(b *testing.B) {
	benchmarkParentArray(b, 1000)
}

func BenchmarkRelationParentArray_10000_x_10(b *testing.B) {
	benchmarkParentArray(b, 10000)
}

func BenchmarkRelationParentArray_100000_x_10(b *testing.B) {
	benchmarkParentArray(b, 100000)
}
