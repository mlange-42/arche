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

const nPos = 9000
const nPosAll = 1000

// Position component
type Position struct {
	X float64
	Y float64
}

// Comp1 component
type Comp1 struct {
	X float64
	Y float64
}

// Comp2 component
type Comp2 struct {
	X float64
	Y float64
}

// Comp3 component
type Comp3 struct {
	X float64
	Y float64
}

// Comp4 component
type Comp4 struct {
	X float64
	Y float64
}

// Comp5 component
type Comp5 struct {
	X float64
	Y float64
}

// Comp6 component
type Comp6 struct {
	X float64
	Y float64
}

// Comp7 component
type Comp7 struct {
	X float64
	Y float64
}

// Comp8 component
type Comp8 struct {
	X float64
	Y float64
}

// Comp9 component
type Comp9 struct {
	X float64
	Y float64
}
