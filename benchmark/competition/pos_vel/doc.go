// Package posvel benchmarks the basic position/velocity example
//
// Setup:
//   - 1000 entities with Position{f64, f64} and Velocity{f64, f64}
//   - 9000 entities with Position{f64, f64}
//
// Benchmark:
//   - Iterate all entities with Position and Velocity
//   - Add Velocity to Position
package posvel

const nPos = 9000
const nPosVel = 1000

// Position component
type Position struct {
	X float64
	Y float64
}

// Velocity component
type Velocity struct {
	X float64
	Y float64
}
