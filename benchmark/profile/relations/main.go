package main

// Profiling:
// go build ./benchmark/profile/relations
// relations
// go tool pprof -http=":8000" -nodefraction=0.001 relations cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 relations mem.pprof

import (
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/pkg/profile"
)

// ParentList component
type ParentList struct {
	FirstChild ecs.Entity
	Value      int
}

// ChildRelation component
type ChildRelation struct {
	ecs.Relation
	Value int
}

func main() {

	count := 10
	iters := 1000
	parents := 10000
	children := 1

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count, iters, parents, children)
	stop.Stop()

	stop = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("."))
	run(count, iters, parents, children)
	stop.Stop()
}

func run(rounds, iters, numParents, numChildren int) {
	for i := 0; i < rounds; i++ {
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

		var childF ecs.Filter = ecs.All(childID)
		chf := ecs.RelationFilter{
			Filter: childF,
			Target: ecs.Entity{},
		}

		for i := 0; i < iters; i++ {
			query := world.Query(&cf)
			for query.Next() {
				parData := (*ParentList)(query.Get(parentID))
				par := query.Entity()

				chf.Target = par

				childQuery := world.Query(&chf)
				for childQuery.Next() {
					child := (*ChildRelation)(childQuery.Get(childID))
					parData.Value += child.Value
				}
			}
		}

		parQuery := world.Query(&parentFilter)

		expected := numChildren * iters
		for parQuery.Next() {
			par := (*ParentList)(parQuery.Get(parentID))
			if par.Value != expected {
				panic("wrong number of children")
			}
		}
	}
}
