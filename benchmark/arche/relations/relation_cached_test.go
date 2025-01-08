package relations

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

type ChildFilter struct {
	RelFilter ecs.RelationFilter
	Filter    ecs.CachedFilter
}

func benchmarkRelationCached(b *testing.B, numParents int, numChildren int) {
	b.StopTimer()

	world := ecs.NewWorld()
	parentID := ecs.ComponentID[ParentList](&world)
	filterID := ecs.ComponentID[ChildFilter](&world)
	childID := ecs.ComponentID[ChildRelation](&world)

	parentMapper := generic.NewMap2[ParentList, ChildFilter](&world)
	childMapper := generic.NewMap1[ChildRelation](&world, generic.T[ChildRelation]())
	targetMapper := generic.NewMap[ChildRelation](&world)

	spawnedPar := parentMapper.NewBatchQ(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		par := spawnedPar.Entity()
		_, fl := spawnedPar.Get()
		fl.RelFilter = ecs.NewRelationFilter(ecs.All(childID), par)
		fl.Filter = world.Cache().Register(&fl.RelFilter)

		parents = append(parents, par)
	}

	spawnedChild := childMapper.NewBatchQ(numParents * numChildren)
	children := make([]ecs.Entity, 0, numParents*numChildren)
	for spawnedChild.Next() {
		children = append(children, spawnedChild.Entity())
	}
	rand.Shuffle(len(children), func(i, j int) { children[i], children[j] = children[j], children[i] })

	for i, e := range children {
		child := childMapper.Get(e)
		child.Value = 1
		parent := parents[i/numChildren]
		targetMapper.SetRelation(e, parent)
	}

	parentFilter := ecs.All(parentID, filterID)
	cf := world.Cache().Register(parentFilter)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&cf)
		for query.Next() {
			parData, filter := (*ParentList)(query.Get(parentID)), (*ChildFilter)(query.Get(filterID))

			childQuery := world.Query(&filter.Filter)
			for childQuery.Next() {
				child := (*ChildRelation)(childQuery.Get(childID))
				parData.Value += child.Value
			}
		}
	}

	b.StopTimer()

	parQuery := world.Query(&parentFilter)

	expected := numChildren * b.N
	for parQuery.Next() {
		par := (*ParentList)(parQuery.Get(parentID))
		if par.Value != expected {
			panic("wrong number of children")
		}
	}
}

func BenchmarkRelationCached_1k_10_x_100(b *testing.B) {
	benchmarkRelationCached(b, 10, 100)
}

func BenchmarkRelationCached_1k_100_x_10(b *testing.B) {
	benchmarkRelationCached(b, 100, 10)
}

func BenchmarkRelationCached_1k_1000_x_1(b *testing.B) {
	benchmarkRelationCached(b, 1000, 1)
}

func BenchmarkRelationCached_10k_10_x_1000(b *testing.B) {
	benchmarkRelationCached(b, 10, 1000)
}

func BenchmarkRelationCached_10k_100_x_100(b *testing.B) {
	benchmarkRelationCached(b, 100, 100)
}

func BenchmarkRelationCached_10k_1000_x_10(b *testing.B) {
	benchmarkRelationCached(b, 1000, 10)
}

func BenchmarkRelationCached_100k_10_x_10000(b *testing.B) {
	benchmarkRelationCached(b, 10, 10000)
}

func BenchmarkRelationCached_100k_100_x_1000(b *testing.B) {
	benchmarkRelationCached(b, 100, 1000)
}

func BenchmarkRelationCached_100k_1000_x_100(b *testing.B) {
	benchmarkRelationCached(b, 1000, 100)
}

func BenchmarkRelationCached_100k_10000_x_10(b *testing.B) {
	benchmarkRelationCached(b, 10000, 10)
}

func BenchmarkRelationCached_1M_10_x_100000(b *testing.B) {
	benchmarkRelationCached(b, 10, 100000)
}

func BenchmarkRelationCached_1M_100_x_10000(b *testing.B) {
	benchmarkRelationCached(b, 100, 10000)
}

func BenchmarkRelationCached_1M_1000_x_1000(b *testing.B) {
	benchmarkRelationCached(b, 1000, 1000)
}

func BenchmarkRelationCached_1M_10000_x_100(b *testing.B) {
	benchmarkRelationCached(b, 10000, 100)
}
