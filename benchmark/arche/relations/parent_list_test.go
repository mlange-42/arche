package relations

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchmarkParentList(b *testing.B, numParents int) {
	b.StopTimer()

	world := ecs.NewWorld(ecs.NewConfig().WithCapacityIncrement(1024))

	parentMapper := generic.NewMap1[ParentList](&world)
	parentMapperFull := generic.NewMap2[ParentData, ParentList](&world)
	childMapper := generic.NewMap2[ChildData, ChildList](&world)

	spawnedPar := parentMapperFull.NewEntitiesQuery(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		parents = append(parents, spawnedPar.Entity())
	}

	spawnedChild := childMapper.NewEntitiesQuery(numParents * numChildren)
	cnt := 0
	for spawnedChild.Next() {
		childEntity := spawnedChild.Entity()
		data, child := spawnedChild.Get()
		data.Value = 1
		par := parentMapper.Get(parents[cnt/numChildren])

		if !par.FirstChild.IsZero() {
			child.Next = par.FirstChild
		}
		par.FirstChild = childEntity
		par.NumChildren++

		cnt++
	}

	parentFilter := generic.NewFilter2[ParentData, ParentList]()
	parentFilter.Register(&world)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := parentFilter.Query(&world)
		for query.Next() {
			data, par := query.Get()
			child := par.FirstChild
			for !child.IsZero() {
				childData, childList := childMapper.Get(child)
				data.Value += childData.Value
				child = childList.Next
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

func BenchmarkRelationParentList_100_x_10(b *testing.B) {
	benchmarkParentList(b, 100)
}

func BenchmarkRelationParentList_1000_x_10(b *testing.B) {
	benchmarkParentList(b, 1000)
}

func BenchmarkRelationParentList_10000_x_10(b *testing.B) {
	benchmarkParentList(b, 10000)
}
