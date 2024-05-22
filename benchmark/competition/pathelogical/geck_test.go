package manycomponents

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/mlange-42/arche/benchmark/competition/pathelogical/geckecs"
)

func BenchmarkPathelogicalGECK(b *testing.B) {
	run := func(b *testing.B, entityCount int) {
		b.StopTimer()
		w := geckecs.NewWorld()

		start := time.Now()

		// Create entities
		entities := w.Entities(entityCount)

		for _, e := range entities {
			if flipCoin() {
				e.TagWithComp1()
			}
			if flipCoin() {
				e.TagWithComp2()
			}
			if flipCoin() {
				e.TagWithComp3()
			}
			if flipCoin() {
				e.TagWithComp4()
			}
			if flipCoin() {
				e.TagWithComp5()
			}
			if flipCoin() {
				e.TagWithComp6()
			}
			if flipCoin() {
				e.TagWithComp7()
			}
			if flipCoin() {
				e.TagWithComp8()
			}
			if flipCoin() {
				e.TagWithComp9()
			}
			if flipCoin() {
				e.TagWithComp10()
			}
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

		b.StartTimer()
		iter := w.Comp1Comp2Comp3Set.NewIterator()
		for s := 0; s < sampleCount; s++ {
			start = time.Now()
			iter.Reset()
			for iter.HasNext() {
				iter.Next()
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
			report("geck", sumD, minD, maxD, entitySum)
		}

	}

	for i := minEntityCount; i <= maxEntityCount; i *= 2 {
		run(b, i)
	}
}
