package relations

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchmarkParentSlice(b *testing.B, numParents int) {
	b.StopTimer()

	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	parentMapper := generic.NewMap1[ParentSlice](&world)
	childMapper := generic.NewMap1[Child](&world)

	spawnedPar := parentMapper.NewQuery(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		parents = append(parents, spawnedPar.Entity())
	}

	spawnedChild := childMapper.NewQuery(numParents * numChildren)
	cnt := 0
	for spawnedChild.Next() {
		data := spawnedChild.Get()
		data.Value = 1
		par := parentMapper.Get(parents[cnt/numChildren])
		par.Children = append(par.Children, spawnedChild.Entity())
		cnt++
	}

	parentFilter := generic.NewFilter1[ParentSlice]()
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

	expected := numChildren * b.N
	for query.Next() {
		par := query.Get()
		if par.Value != expected {
			panic("wrong number of children")
		}
	}
}

func BenchmarkRelationParentSlice_100_x_10(b *testing.B) {
	benchmarkParentSlice(b, 100)
}

func BenchmarkRelationParentSlice_1000_x_10(b *testing.B) {
	benchmarkParentSlice(b, 1000)
}

func BenchmarkRelationParentSlice_10000_x_10(b *testing.B) {
	benchmarkParentSlice(b, 10000)
}

func BenchmarkRelationParentSlice_100000_x_10(b *testing.B) {
	benchmarkParentSlice(b, 100000)
}
