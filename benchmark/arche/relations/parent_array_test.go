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
	parentMapperFull := generic.NewMap2[ParentData, ParentArr](&world)
	childMapper := generic.NewMap1[ChildData](&world)

	spawnedPar := parentMapperFull.NewEntitiesQuery(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		parents = append(parents, spawnedPar.Entity())
	}

	spawnedChild := childMapper.NewEntitiesQuery(numParents * numChildren)
	cnt := 0
	for spawnedChild.Next() {
		data := spawnedChild.Get()
		data.Value = 1
		par := parentMapper.Get(parents[cnt/numChildren])
		par.Children[cnt%numChildren] = spawnedChild.Entity()
		cnt++
	}

	parentFilter := generic.NewFilter2[ParentData, ParentArr]()
	parentFilter.Register(&world)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := parentFilter.Query(&world)
		for query.Next() {
			data, par := query.Get()
			for i := 0; i < len(par.Children); i++ {
				childData := childMapper.Get(par.Children[i])
				data.Value += childData.Value
			}
		}
	}

	b.StopTimer()
	parentFilterSimple := generic.NewFilter1[ParentData]()
	parentFilterSimple.Register(&world)
	querySimple := parentFilterSimple.Query(&world)

	expected := numChildren * b.N
	for querySimple.Next() {
		par := querySimple.Get()
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
