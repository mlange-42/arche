package iterate

import (
	"testing"

	c "github.com/mlange-42/arche/benchmark/arche/common"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/filter"
)

func runQuery1kArch(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	c.RegisterAll(&world)

	perArch := count / 1000

	for i := 0; i < 1024; i++ {
		mask := i
		add := make([]ecs.ID, 0, 11)
		add = append(add, 10)
		for j := 0; j < 10; j++ {
			id := ecs.ID(j)
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			world.NewEntity(add...)
		}
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(ecs.All(10))
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct10)(query.Get(10))
			pos.Val = 1
		}
	}
}

func runQuery1kArchCached(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	c.RegisterAll(&world)

	perArch := count / 1000

	for i := 0; i < 1024; i++ {
		mask := i
		add := make([]ecs.ID, 0, 11)
		add = append(add, 10)
		for j := 0; j < 10; j++ {
			id := ecs.ID(j)
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			world.NewEntity(add...)
		}
	}

	cf := world.Cache().Register(ecs.All(10))
	var filter ecs.Filter = &cf

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter)
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct10)(query.Get(10))
			pos.Val = 1
		}
	}
}

func runFilter1kArch(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	c.RegisterAll(&world)

	perArch := count / 1000

	for i := 0; i < 1024; i++ {
		mask := i
		add := make([]ecs.ID, 0, 11)
		add = append(add, 10)
		for j := 0; j < 10; j++ {
			id := ecs.ID(j)
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			entity := world.NewEntity()
			world.Add(entity, add...)
		}
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter.All(10))
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct10)(query.Get(10))
			pos.Val = 1
		}
	}
}

func runQuery1Of1kArch(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	c.RegisterAll(&world)

	perArch := 2 * count / 1000

	for i := 0; i < 1024; i++ {
		mask := i
		add := make([]ecs.ID, 0, 10)
		for j := 0; j < 10; j++ {
			id := ecs.ID(j)
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			entity := world.NewEntity()
			world.Add(entity, add...)
		}
	}

	ecs.NewBuilder(&world, 10).NewBatch(count)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(ecs.All(10))
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct6)(query.Get(10))
			pos.Val = 1
		}
	}
}

func runQuery1Of1kArchCached(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	c.RegisterAll(&world)

	perArch := 2 * count / 1000

	for i := 0; i < 1024; i++ {
		mask := i
		add := make([]ecs.ID, 0, 10)
		for j := 0; j < 10; j++ {
			id := ecs.ID(j)
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			entity := world.NewEntity()
			world.Add(entity, add...)
		}
	}

	ecs.NewBuilder(&world, 10).NewBatch(count)

	cf := world.Cache().Register(ecs.All(10))
	var filter ecs.Filter = &cf

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter)
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct10)(query.Get(10))
			pos.Val = 1
		}
	}
}

func runQuery1kTargets(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	posID := ecs.ComponentID[c.TestStruct0](&world)
	relID := ecs.ComponentID[c.ChildOf](&world)

	perArch := count / 1000

	builder := ecs.NewBuilder(&world)
	targetQuery := builder.NewQuery(1000)
	targets := make([]ecs.Entity, 0, 1000)
	for targetQuery.Next() {
		targets = append(targets, targetQuery.Entity())
	}

	childBuilder := ecs.NewBuilder(&world, posID, relID).WithRelation(relID)
	for _, target := range targets {
		childBuilder.NewBatch(perArch, target)
	}

	filter := ecs.All(posID, relID)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter)
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct0)(query.Get(posID))
			pos.Val = 1
		}
	}
}

func runQuery1kTargetsCached(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	posID := ecs.ComponentID[c.TestStruct0](&world)
	relID := ecs.ComponentID[c.ChildOf](&world)

	perArch := count / 1000

	builder := ecs.NewBuilder(&world)
	targetQuery := builder.NewQuery(1000)
	targets := make([]ecs.Entity, 0, 1000)
	for targetQuery.Next() {
		targets = append(targets, targetQuery.Entity())
	}

	childBuilder := ecs.NewBuilder(&world, posID, relID).WithRelation(relID)
	for _, target := range targets {
		childBuilder.NewBatch(perArch, target)
	}

	cf := world.Cache().Register(ecs.All(posID, relID))
	var filter ecs.Filter = &cf

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter)
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct0)(query.Get(posID))
			pos.Val = 1
		}
	}
}

func runQuery1Of1kTargets(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	posID := ecs.ComponentID[c.TestStruct0](&world)
	relID := ecs.ComponentID[c.ChildOf](&world)

	builder := ecs.NewBuilder(&world)
	targetQuery := builder.NewQuery(1000)
	targets := make([]ecs.Entity, 0, 1000)
	for targetQuery.Next() {
		targets = append(targets, targetQuery.Entity())
	}

	childBuilder := ecs.NewBuilder(&world, posID, relID).WithRelation(relID)
	for _, target := range targets {
		childBuilder.New(target)
	}
	target := targets[0]
	childBuilder.NewBatch(count, target)

	filter := ecs.RelationFilter(ecs.All(posID, relID), target)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter)
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct0)(query.Get(posID))
			pos.Val = 1
		}
	}
}

func runQuery1Of1kTargetsCached(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()
	posID := ecs.ComponentID[c.TestStruct0](&world)
	relID := ecs.ComponentID[c.ChildOf](&world)

	builder := ecs.NewBuilder(&world)
	targetQuery := builder.NewQuery(1000)
	targets := make([]ecs.Entity, 0, 1000)
	for targetQuery.Next() {
		targets = append(targets, targetQuery.Entity())
	}

	childBuilder := ecs.NewBuilder(&world, posID, relID).WithRelation(relID)
	for _, target := range targets {
		childBuilder.New(target)
	}

	target := targets[0]
	childBuilder.NewBatch(count, target)

	cf := world.Cache().Register(ecs.RelationFilter(ecs.All(posID, relID), target))
	var filter ecs.Filter = &cf

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := world.Query(filter)
		b.StartTimer()
		for query.Next() {
			pos := (*c.TestStruct0)(query.Get(posID))
			pos.Val = 1
		}
	}
}

func BenchmarkIter1kArchID_1_000(b *testing.B) {
	runQuery1kArch(b, 1000)
}

func BenchmarkIter1kArchID_10_000(b *testing.B) {
	runQuery1kArch(b, 10000)
}

func BenchmarkIter1kArchID_100_000(b *testing.B) {
	runQuery1kArch(b, 100000)
}

func BenchmarkIter1kArchIDCached_1_000(b *testing.B) {
	runQuery1kArchCached(b, 1000)
}

func BenchmarkIter1kArchIDCached_10_000(b *testing.B) {
	runQuery1kArchCached(b, 10000)
}

func BenchmarkIter1kArchIDCached_100_000(b *testing.B) {
	runQuery1kArchCached(b, 100000)
}

func BenchmarkFilter1kArchID_1_000(b *testing.B) {
	runFilter1kArch(b, 1000)
}

func BenchmarkFilter1kArchID_10_000(b *testing.B) {
	runFilter1kArch(b, 10000)
}

func BenchmarkFilter1kArchID_100_000(b *testing.B) {
	runFilter1kArch(b, 100000)
}

func BenchmarkIter1Of1kArch_1_000(b *testing.B) {
	runQuery1Of1kArch(b, 1000)
}

func BenchmarkIter1Of1kArch_10_000(b *testing.B) {
	runQuery1Of1kArch(b, 10000)
}

func BenchmarkIter1Of1kArch_100_000(b *testing.B) {
	runQuery1Of1kArch(b, 100000)
}

func BenchmarkIter1Of1kArchCached_1_000(b *testing.B) {
	runQuery1Of1kArchCached(b, 1000)
}

func BenchmarkIter1Of1kArchCached_10_000(b *testing.B) {
	runQuery1Of1kArchCached(b, 10000)
}

func BenchmarkIter1Of1kArchCached_100_000(b *testing.B) {
	runQuery1Of1kArchCached(b, 100000)
}

func BenchmarkIter1kTargets_1_000(b *testing.B) {
	runQuery1kTargets(b, 1000)
}

func BenchmarkIter1kTargets_10_000(b *testing.B) {
	runQuery1kTargets(b, 10000)
}

func BenchmarkIter1kTargets_100_000(b *testing.B) {
	runQuery1kTargets(b, 100000)
}

func BenchmarkIter1kTargetsCached_1_000(b *testing.B) {
	runQuery1kTargetsCached(b, 1000)
}

func BenchmarkIter1kTargetsCached_10_000(b *testing.B) {
	runQuery1kTargetsCached(b, 10000)
}

func BenchmarkIter1kTargetsCached_100_000(b *testing.B) {
	runQuery1kTargetsCached(b, 100000)
}

func BenchmarkIter1Of1kTargets_1_000(b *testing.B) {
	runQuery1Of1kTargets(b, 1000)
}

func BenchmarkIter1Of1kTargets_10_000(b *testing.B) {
	runQuery1Of1kTargets(b, 10000)
}

func BenchmarkIter1Of1kTargets_100_000(b *testing.B) {
	runQuery1Of1kTargets(b, 100000)
}

func BenchmarkIter1Of1kTargetsCached_1_000(b *testing.B) {
	runQuery1Of1kTargetsCached(b, 1000)
}

func BenchmarkIter1Of1kTargetsCached_10_000(b *testing.B) {
	runQuery1kTargetsCached(b, 10000)
}

func BenchmarkIter1Of1kTargetsCached_100_000(b *testing.B) {
	runQuery1kTargetsCached(b, 100000)
}
