// Package manycomponents benchmarks access to 10 components
//
// Setup:
//   - 1000 entities with Position{f64, f64} and Comp1{f64, f64} ... Comp9{f64, f64}
//   - 9000 entities with Position{f64, f64}
//
// Benchmark:
//   - Iterate all entities with Position and Comp1 ... Comp9
//   - Sum up X and Y values
package manycomponents

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	componentCount = 10
	queryCount     = 3
	sampleCount    = 1 << 12
	flipRate       = 0.5
	minEntityCount = 1 << 10
	maxEntityCount = 1 << 18
)

func report(queryName string, sum, min, max time.Duration, entityCount int) {
	sum /= sampleCount

	fmt.Printf("total:\t\tmin:\t%s,\tavg: %s,\tmax: %s (%s)\n", min, sum, max, queryName)

	ecf := float64(entityCount)
	minf := float64(min) / ecf
	sumf := float64(sum) / ecf
	maxf := float64(max) / ecf

	fmt.Printf("per entity:\tmin: %0.2fns,\tavg: %0.2fns,\tmax: %0.2fns\n\n", minf, sumf, maxf)
}

func flipCoin() bool {
	return rand.Float64() < flipRate
}
