package relations

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func benchmarkParentList(b *testing.B, numParents int, numChildren int) {
	b.StopTimer()

	world := ecs.NewWorld()
	parentID := ecs.ComponentID[ParentList](&world)
	childID := ecs.ComponentID[ChildList](&world)

	parentMapper := generic.NewMap1[ParentList](&world)
	childMapper := generic.NewMap1[ChildList](&world)

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
		par := parentMapper.Get(parents[i/numChildren])

		if !par.FirstChild.IsZero() {
			child.Next = par.FirstChild
		}
		par.FirstChild = e
	}

	parentFilter := ecs.All(parentID)
	cf := world.Cache().Register(parentFilter)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&cf)
		for query.Next() {
			par := (*ParentList)(query.Get(parentID))
			child := par.FirstChild
			for !child.IsZero() {
				childList := (*ChildList)(world.Get(child, childID))
				par.Value += childList.Value
				child = childList.Next
			}
		}
	}

	b.StopTimer()
	query := world.Query(&cf)

	expected := numChildren * b.N
	for query.Next() {
		par := (*ParentList)(query.Get(parentID))
		if par.Value != expected {
			panic("wrong number of children")
		}
	}
}

func BenchmarkRelationParentList_1k_10_x_100(b *testing.B) {
	benchmarkParentList(b, 10, 100)
}

func BenchmarkRelationParentList_1k_100_x_10(b *testing.B) {
	benchmarkParentList(b, 100, 10)
}

func BenchmarkRelationParentList_1k_1000_x_1(b *testing.B) {
	benchmarkParentList(b, 1000, 1)
}

func BenchmarkRelationParentList_10k_10_x_1000(b *testing.B) {
	benchmarkParentList(b, 10, 1000)
}

func BenchmarkRelationParentList_10k_100_x_100(b *testing.B) {
	benchmarkParentList(b, 100, 100)
}

func BenchmarkRelationParentList_10k_1000_x_10(b *testing.B) {
	benchmarkParentList(b, 1000, 10)
}

func BenchmarkRelationParentList_100k_10_x_10000(b *testing.B) {
	benchmarkParentList(b, 10, 10000)
}

func BenchmarkRelationParentList_100k_100_x_1000(b *testing.B) {
	benchmarkParentList(b, 100, 1000)
}

func BenchmarkRelationParentList_100k_1000_x_100(b *testing.B) {
	benchmarkParentList(b, 1000, 100)
}

func BenchmarkRelationParentList_100k_10000_x_10(b *testing.B) {
	benchmarkParentList(b, 10000, 10)
}

func BenchmarkRelationParentList_1M_10_x_100000(b *testing.B) {
	benchmarkParentList(b, 10, 100000)
}

func BenchmarkRelationParentList_1M_100_x_10000(b *testing.B) {
	benchmarkParentList(b, 100, 10000)
}

func BenchmarkRelationParentList_1M_1000_x_1000(b *testing.B) {
	benchmarkParentList(b, 1000, 1000)
}

func BenchmarkRelationParentList_1M_10000_x_100(b *testing.B) {
	benchmarkParentList(b, 10000, 100)
}
