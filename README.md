# Arche

[![Test status](https://github.com/mlange-42/arche/actions/workflows/tests.yml/badge.svg)](https://github.com/mlange-42/arche/actions/workflows/tests.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/arche.svg)](https://pkg.go.dev/github.com/mlange-42/arche)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/arche)
[![MIT license](https://img.shields.io/github/license/mlange-42/arche)](https://github.com/mlange-42/arche/blob/main/LICENSE)

*Arche* is an archetype-based Entity Component System for Go.

*Arche* is designed for the use in simulation models of the
[Department for Ecological Modelling](https://www.ufz.de/index.php?en=34213) at the
[Helmholtz Centre for Environmental Research](https://www.ufz.de).

**:warning: Arche is still under rapid development! Be prepared for frequent API changes. :warning:**  


## Installations

```shell
go get github.com/mlange-42/arche
```

## Features

* Minimal API. See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche).
* Fast iteration and component access via queries (≈2.5ns iterate + get).
* Fast random access for components of arbitrary entities. Useful for hierarchies.
* No systems. Use your own structure.
* Not thread-safe. On purpose.
* No dependencies. Except for unit tests (100% coverage).

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

	// Create entities
	for i := 0; i < 1000; i++ {
		// Create a new Entity.
		entity := world.NewEntity()
		// Add components to it.
		pos, vel := generic.Add2[Position, Velocity](&world, entity)

		// Initialize component fields.
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100

		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query.
		// Generic queries support up to 8 components.
		// For more components, use World.Query()
		query := generic.Query2[Position, Velocity](&world)
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

Depending on the machine the code is running on, generics may or may not incur an overhead.
The worst to expect is a doubling of iteration + access time (from 2.5ns/op to 4ns/op),
while on other machines both approaches are equally fast.

For performance-critical code, the use of the ID-based methods of `ecs.World` may be worth testing.
Component IDs are retrieved like this:

```go
posID := ComponentID[Position](&world)
```

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
