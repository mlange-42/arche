package relations

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchmarkRelation(b *testing.B, numParents int, numChildren int) {
	b.StopTimer()

	world := ecs.NewWorld()
	parentID := ecs.ComponentID[ParentList](&world)
	childID := ecs.ComponentID[ChildRelation](&world)

	parentMapper := generic.NewMap1[ParentList](&world)
	childMapper := generic.NewMap1[ChildRelation](&world, generic.T[ChildRelation]())
	targetMapper := generic.NewMap[ChildRelation](&world)

	spawnedPar := parentMapper.NewBatchQ(numParents)
	parents := make([]ecs.Entity, 0, numParents)
	for spawnedPar.Next() {
		parents = append(parents, spawnedPar.Entity())
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

	parentFilter := ecs.All(parentID)
	cf := world.Cache().Register(parentFilter)

	childF := ecs.All(childID)
	relFilter := ecs.NewRelationFilter(childF, ecs.Entity{})
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&cf)
		for query.Next() {
			parData := (*ParentList)(query.Get(parentID))
			par := query.Entity()

			relFilter.Target = par

			childQuery := world.Query(&relFilter)
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

func BenchmarkRelationDefault_1k_10_x_100(b *testing.B) {
	benchmarkRelation(b, 10, 100)
}

func BenchmarkRelationDefault_1k_100_x_10(b *testing.B) {
	benchmarkRelation(b, 100, 10)
}

func BenchmarkRelationDefault_1k_1000_x_1(b *testing.B) {
	benchmarkRelation(b, 1000, 1)
}

func BenchmarkRelationDefault_10k_10_x_1000(b *testing.B) {
	benchmarkRelation(b, 10, 1000)
}

func BenchmarkRelationDefault_10k_100_x_100(b *testing.B) {
	benchmarkRelation(b, 100, 100)
}

func BenchmarkRelationDefault_10k_1000_x_10(b *testing.B) {
	benchmarkRelation(b, 1000, 10)
}

func BenchmarkRelationDefault_100k_10_x_10000(b *testing.B) {
	benchmarkRelation(b, 10, 10000)
}

func BenchmarkRelationDefault_100k_100_x_1000(b *testing.B) {
	benchmarkRelation(b, 100, 1000)
}

func BenchmarkRelationDefault_100k_1000_x_100(b *testing.B) {
	benchmarkRelation(b, 1000, 100)
}

func BenchmarkRelationDefault_100k_10000_x_10(b *testing.B) {
	benchmarkRelation(b, 10000, 10)
}

func BenchmarkRelationDefault_1M_10_x_100000(b *testing.B) {
	benchmarkRelation(b, 10, 100000)
}

func BenchmarkRelationDefault_1M_100_x_10000(b *testing.B) {
	benchmarkRelation(b, 100, 10000)
}

func BenchmarkRelationDefault_1M_1000_x_1000(b *testing.B) {
	benchmarkRelation(b, 1000, 1000)
}

func BenchmarkRelationDefault_1M_10000_x_100(b *testing.B) {
	benchmarkRelation(b, 10000, 100)
}
