# Arche

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/arche/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/arche/actions/workflows/tests.yml)
[![100% Coverage](https://img.shields.io/github/actions/workflow/status/mlange-42/arche/coverage.yml?branch=main&label=100%25%20Coverage&logo=github)](https://github.com/mlange-42/arche/actions/workflows/coverage.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/arche.svg)](https://pkg.go.dev/github.com/mlange-42/arche)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/arche)
[![MIT license](https://img.shields.io/github/license/mlange-42/arche)](https://github.com/mlange-42/arche/blob/main/LICENSE)

*Arche* is an archetype-based Entity Component System for Go.

*Arche* is designed for the use in simulation models of the
[Department for Ecological Modelling](https://www.ufz.de/index.php?en=34213) at the
[Helmholtz Centre for Environmental Research](https://www.ufz.de).

## Installations

```shell
go get github.com/mlange-42/arche
```

## Features

* Minimal core API. See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche).
* Optional rich filtering and generic query API.
* Fast iteration and component access via queries (≈2.5ns iterate + get).
* Fast random access for components of arbitrary entities. Useful for hierarchies.
* No systems. Use your own structure.
* Not thread-safe. On purpose.
* No dependencies. Except for unit tests ([100% coverage](https://github.com/mlange-42/arche/actions/workflows/coverage.yml)).

For details on Arche's architecture, see section [Architecture](#architecture).

## Usage example

Here is a minimal usage example.
You will likely create systems with a method that takes a pointer to the `World` as argument.

See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche) and
[examples](https://github.com/mlange-42/arche/tree/main/examples) for details.

```go
package main

import (
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
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

	// Create entities.
	for i := 0; i < 1000; i++ {
		// Create a new Entity with components.
		_, pos, vel := generic.NewEntity2[Position, Velocity](&world)

		// Initialize component fields.
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100
		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Create a generic filter.
	filter := generic.NewFilter2[Position, Velocity]()

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query.
		query := filter.Query(&world)
		// Iterate it
		for query.Next() {
			// Component access through the Query.
			_, pos, vel := query.GetAll()
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
```

## Design decisions

Unlike most other ECS implementations, *Arche* is designed for the development of scientific,
individual-based models rather than for game development.
This motivates some design decisions, with a focus on simplicity, safety and performance.

### Minimal core API

The `ecs.World` object is a pure and minimal ECS implementation in the sense of a data store
for entities and components, with query and iteration capabilities.
There is neither an update loop nor systems.
These should be implemented by the user.

The packages `filter` and `generic` provide a layer around the core for richer and/or safer queries and operations. They are built on top of the `ecs` package, so they could also be implemented by users.

### Determinism

Iteration order in *Arche* is deterministic and reproducible.
This does not mean that entities are iterated in their order of insertion, nor in the same order in successive iterations.
However, given the same operations on the `ecs.World`, iteration order will always be the same.

### Strict and panic

*Arche* puts an emphasis on safety and on avoiding undefined behavior.
It panics on unexpected operations, like removing a dead entity,
adding a component that is already present, or attempting to change a locked world.
This may seem not idiomatic for Go.
However, explicit error handling in performance hotspots is not an option.
Neither is silent failure, given the scientific background.

### Other limitations

* The number of component types per `World` is limited to 128. This is mainly a performance decision.
* The number of entities alive at any one time is limited to just under 5 billion (`uint32` ID).

## Architecture

*Arche* uses an archetype-based architecture.

The ASCII graph below illustrates the architecture.
Components for entities are stored in so-called archetypes, which represent unique combinations of components.
In the illustration, the first archetype holds all components for all entities with (only/exactly) the components A, B and C.

```text
 Entities   Archetypes   Bitmasks   Queries

   E         E Comps
  |0|       |2|A|B|C|    111...<-.<--------.
  |1|---.   |8|A|B|C|            |         |
  |2|   '-->|1|A|B|C|            |         |
  |3|       |3|A|B|C|            |--(A, C) |
  |4|                            |  101... |
  |6|   .-->|7|A|C|      101...<-'         |--(B)
  |7|---'   |6|A|C|                        |  010...
  |8|       |4|A|C|                        |
  |9|---.                                  |
  |.|   |   |5|B|C|      010...   <--------'
  |.|   '-->|9|B|C|
  |.|
  |.| <===> [Entity pool]
```

The exact composition of each archetype is encoded in a bitmask for fast comparison.
Thus, queries can easily identify their relevant archetypes (i.e. query bitmask contained in archetype bitmask), and then iterate entities linearly, and very fast. Components can be accessed through the query in a very efficient way.

For getting components by entity ID, e.g. for hierarchies, the world contains a list that is indexed by the entity ID, and references the entity's archetype and index in the archetype. This way, getting components for entity IDs (i.e. random access) is fast, although not as fast as in queries.

Obviously, archetypes are an optimization for iteration speed.
But they also come with a downside. Adding or removing components to/from an entity requires moving all the components of the entity to another archetype.
It is therefore recommended to add/remove/exchange multiple components at the same time rather than one after the other.

## Generic vs. ID access

*Arche* provides generic functions and types for accessing and modifying components etc., as shown in the [Usage example](#usage-example).

Generic access is built on top of ID-based access used by the `ecs.World`.
Generic functions and types provide type-safety and are more user-friendly than ID-based access.

From Go 1.20 onwards, generics incur a runtime overhead of 0-25%, depending on the machine(?).
On earlier versions of Go, the worst to expect is a doubling of iteration + access time (from 2.5ns/op to 4ns/op).
For performance-critical code, the use of the ID-based methods of `ecs.World` may be worth testing.

For more details, see the [API docs](https://pkg.go.dev/github.com/mlange-42/arche) and
[examples](https://github.com/mlange-42/arche/tree/main/examples).

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
