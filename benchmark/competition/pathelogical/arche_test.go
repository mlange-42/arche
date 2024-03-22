package manycomponents

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/mlange-42/arche/benchmark/competition/pathelogical/geckecs"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func BenchmarkPathelogicalArche(b *testing.B) {
	run := func(b *testing.B, entityCount int) {
		b.StopTimer()
		w := ecs.NewWorld()

		start := time.Now()

		// Create component ids
		comp1ID := ecs.ComponentID[geckecs.Comp1](&w)
		comp2ID := ecs.ComponentID[geckecs.Comp2](&w)
		comp3ID := ecs.ComponentID[geckecs.Comp3](&w)
		comp4ID := ecs.ComponentID[geckecs.Comp4](&w)
		comp5ID := ecs.ComponentID[geckecs.Comp5](&w)
		comp6ID := ecs.ComponentID[geckecs.Comp6](&w)
		comp7ID := ecs.ComponentID[geckecs.Comp7](&w)
		comp8ID := ecs.ComponentID[geckecs.Comp8](&w)
		comp9ID := ecs.ComponentID[geckecs.Comp9](&w)
		comp10ID := ecs.ComponentID[geckecs.Comp10](&w)
		tags := []ecs.ID{
			comp1ID, comp2ID, comp3ID, comp4ID, comp5ID,
			comp6ID, comp7ID, comp8ID, comp9ID, comp10ID,
		}

		// Create entities
		tagsToInclude := make([]ecs.ID, 0, len(tags))
		for i := 0; i < entityCount; i++ {
			tagsToInclude = tagsToInclude[:0]
			for _, tag := range tags {
				if flipCoin() {
					tagsToInclude = append(tagsToInclude, tag)
				}
			}
			w.NewEntity(tagsToInclude...)
		}

		fmt.Printf(
			"upsert %s entities with %d components with flip rate of %0.0f%% in %s\n",
			humanize.Comma(int64(entityCount)),
			componentCount,
			flipRate*100,
			time.Since(start),
		)

		fmt.Printf("querying for %d components taking %d samples\n", queryCount, sampleCount)

		var (
			entitySum        int
			sumD, minD, maxD time.Duration
		)

		// f, err := os.Create("profile.out")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// pprof.StartCPUProfile(f)
		// defer pprof.StopCPUProfile()

		filter := generic.NewFilter3[geckecs.Comp1, geckecs.Comp2, geckecs.Comp3]()
		b.StartTimer()
		for s := 0; s < sampleCount; s++ {
			start = time.Now()
			query := filter.Query(&w)
			for query.Next() {
				query.Get()
				entitySum++
			}
			d := time.Since(start)
			sumD += d
			minD = min(minD, d)
			maxD = max(maxD, d)
		}

		if entitySum == 0 {
			fmt.Println("no entities found")
		} else {
			report("arche", sumD, minD, maxD, entitySum)
		}
	}

	for i := minEntityCount; i <= maxEntityCount; i *= 2 {
		run(b, i)
	}
}
