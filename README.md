# Arche

[![Test status](https://github.com/mlange-42/arche/actions/workflows/tests.yml/badge.svg)](https://github.com/mlange-42/arche/actions/workflows/tests.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/arche.svg)](https://pkg.go.dev/github.com/mlange-42/arche)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/arche)
[![MIT license](https://img.shields.io/github/license/mlange-42/arche)](https://github.com/mlange-42/arche/blob/main/LICENSE)

*Arche* is an archetype-based Entity Component System for Go.

*Arche* is designed for the use in simulation models of the
[Department for Ecological Modelling](https://www.ufz.de/index.php?en=34213) at the
[Helmholtz Centre for Environmental Research](https://www.ufz.de).

:warning: This project is in an early stage of development! :warning:

## Installations

```shell
go get github.com/mlange-42/arche
```

## Features

* Minimal API. See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche).
* Very fast iteration and component access via `Query` (benchmarks for comparison in progress).
* Fast random access for components of arbitrary entities. Useful for hierarchies.
* No systems. Use your own structure.
* Not thread-safe. On purpose.
* No dependencies (except for unit tests).

## Usage example

Here is a minimal usage example.
You will likely create systems with a method that takes a pointer to the `World` as argument.

See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche) for details.

```go
package main

import (
	"math/rand"

	"github.com/mlange-42/arche/ecs"
)

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

func main() {
	// Create a World.
	world := ecs.NewWorld()

	// Create entities
	for i := 0; i < 1000; i++ {
		// Create a new Entity.
		entity := world.NewEntity()
		// Add components to it.
		pos, vel := ecs.Add2[Position, Velocity](&world, entity)

		// Initialize component fields.
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100

		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query.
		// Generic queries support up to 5 components.
		// For more components, use World.Query()
		query := ecs.Query2[Position, Velocity](&world)
		// Iterate it
		for query.Next() {
			// Component access through a Query.
			pos, vel := query.GetAll()
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
```

## Benchmarks

[TODO]

For now, see the latest [Benchmarks CI run](https://github.com/mlange-42/arche/actions/workflows/benchmarks.yml).

## References

Information and inspiration taken from:

* [Sander Mertens](https://ajmmertens.medium.com/)' blog
* Michele "skypjack" Caini's blog [skypjack on software](https://skypjack.github.io/)
* [marioolofo/go-gameengine-ecs](https://github.com/marioolofo/go-gameengine-ecs) "Fast Entity Component System in Golang"

## License

This project is distributed under the [MIT licence](./LICENSE).
