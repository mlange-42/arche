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
* No dependencies (except for testing and benchmarks).

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

	// Get component IDs.
	// Registers component type if not already registered.
	positionID := ecs.ComponentID[Position](&world)
	velocityID := ecs.ComponentID[Velocity](&world)

	// Create entities
	for i := 0; i < 1000; i++ {
		// Create a new Entity.
		entity := world.NewEntity()
		// Add components to it.
		world.Add(entity, positionID, velocityID)

		// Component access through the World.
		// See below for faster access in queries.
		pos := (*Position)(world.Get(entity, positionID))
		vel := (*Velocity)(world.Get(entity, velocityID))

		// Initialize component fields.
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100

		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query.
		query := world.Query(positionID, velocityID)
		// Iterate it
		for query.Next() {
			// Component access through a Query.
			// About 20-30% faster than access through the World.
			// Can also fetch components not in the query.
			pos := (*Position)(query.Get(positionID))
			vel := (*Velocity)(query.Get(velocityID))

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
