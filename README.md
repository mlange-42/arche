# Arche

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/arche/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/arche/actions/workflows/tests.yml)
[![100% Coverage](https://img.shields.io/github/actions/workflow/status/mlange-42/arche/coverage.yml?branch=main&label=100%25%20Coverage&logo=github)](https://github.com/mlange-42/arche/actions/workflows/coverage.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/arche.svg)](https://pkg.go.dev/github.com/mlange-42/arche)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/arche)
[![MIT license](https://img.shields.io/github/license/mlange-42/arche)](https://github.com/mlange-42/arche/blob/main/LICENSE)

*Arche* is an archetype-based [Entity Component System](https://en.wikipedia.org/wiki/Entity_component_system) for [Go](https://go.dev/).

*Arche* is designed for the use in simulation models of the
[Department for Ecological Modelling](https://www.ufz.de/index.php?en=34213) at the
[Helmholtz Centre for Environmental Research](https://www.ufz.de).

## Features

* Simple [core API](https://pkg.go.dev/github.com/mlange-42/arche/ecs). See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche).
* Optional rich [filtering](https://pkg.go.dev/github.com/mlange-42/arche/filter) and [generic](https://pkg.go.dev/github.com/mlange-42/arche/generic) query API.
* No systems. Just queries. Use your own structure.
* Not thread-safe. On purpose.
* No dependencies. Except for unit tests ([100% coverage](https://github.com/mlange-42/arche/actions/workflows/coverage.yml)).
* Probably the fastest Go ECS out there. See the [Benchmarks](#benchmarks).

For details, see sections [Architecture](#architecture) and [Design decisions](#design-decisions).

## Installation

To use *Arche* in a Go project, run:

```shell
go get github.com/mlange-42/arche
```

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

	// Create a component mapper.
	mapper := generic.NewMap2[Position, Velocity](&world)

	// Create entities.
	for i := 0; i < 1000; i++ {
		// Create a new Entity with components.
		entity := mapper.NewEntity()
		// Get the components
		pos, vel := mapper.Get(entity)
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
			pos, vel := query.Get()
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
This motivates some design decisions, with an emphasis on simplicity, safety and performance.
Nevertheless, *Arche* can also be used for game development.

### Minimal core API

The `ecs.World` object is a pure and minimal ECS implementation in the sense of a data store
for entities and components, with query and iteration capabilities.
The core package `ecs` consists of less than 1500 lines of easy-to-read, clean and well-documented Go code.

There is neither an update loop nor systems.
These should be implemented by the user.

The packages `filter` and `generic` provide a layer around the core for richer and/or safer queries and manipulation. They are built on top of the `ecs` package, so they could also be implemented by a user.
See also [Generic vs. ID access](#generic-vs-id-access).

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
  |0|       |2|A|B|C|    111...<-.<--match-.
  |1|---.   |8|A|B|C|            |         |
  |2|   '-->|1|A|B|C|            |         |
  |3|       |3|A|B|C|            |--(A, C) |
  |4|                            |  101... |
  |6|   .-->|7|A|C|      101...<-'         |--(B)
  |7|---'   |6|A|C|                        |  010...
  |8|       |4|A|C|                        |
  |9|---.                                  |
  |.|   |   |5|B|C|      011...   <--------'
  |.|   '-->|9|B|C|
  |.|
  |.| <===> [Entity pool]
```

The exact composition of each archetype is encoded in a bitmask for fast comparison.
Thus, queries can easily identify their relevant archetypes, and then simply iterate entities linearly, which is very fast. Components can be accessed through the query in a very efficient way (&approx;1ns).

For getting components by entity ID, e.g. for hierarchies, the world contains a list that is indexed by the entity ID. For each entity, it references it's current archetype and the index in the archetype. This way, getting components for entity IDs (i.e. random access) is fast, although not as fast as in queries (≈1.5ns vs. 1ns).

Obviously, archetypes are an optimization for iteration speed.
But they also come with a downside. Adding or removing components to/from an entity requires moving all the components of the entity to another archetype.
This takes around 20ns per involved component.
It is therefore recommended to add/remove/exchange multiple components at the same time rather than one after the other.

## Generic vs. ID access

*Arche* provides generic functions and types for accessing and modifying components etc., as shown in the [Usage example](#usage-example).

Generic access is built on top of the ID-based access that is used by the `ecs.World`.
Generic functions and types provide type-safety and are more user-friendly than ID-based access.
However, when querying many components, generic queries have a runtime overhead of around 20-30%.
For performance-critical code, the use of the ID-based methods of `ecs.World` may be worth testing.

For more details, see the [API docs](https://pkg.go.dev/github.com/mlange-42/arche) and
[examples](https://github.com/mlange-42/arche/tree/main/examples).

## Benchmarks

See also the latest [Benchmarks CI run](https://github.com/mlange-42/arche/actions/workflows/benchmarks.yml).

### Arche vs. other Go ECS implementations

To the best of the author's knowledge, there are only a handful of ECS implementations in Go that are serious and somewhat maintained:

* [go-gameengine-ecs](https://github.com/marioolofo/go-gameengine-ecs)
* [Donburi](https://github.com/yohamta/donburi)
* [Ento](https://github.com/wwfranczyk/ento)
* [Entitas-Go](https://github.com/Falldot/Entitas-Go) (archived repository)

Here, *Arche* is benchmarked against these implementations.
Feel free to open an issue if you have suggestions for improvements on the benchmarking code or other engines to include.

#### Position/Velocity

Build:
* Create 1000 entities with `Pos{float64, float64}` and `Vel{float64, float64}`.
* Create 9000 entities with only `Pos{float64, float64}`.

Iterate:
* Iterate all entities with `Pos` and `Vel`, and add `Vel` to `Pos`.

<div align="center" width="100%">

![Benchmark vs. Go ECSs - Pos/Vel](https://user-images.githubusercontent.com/44003176/227806196-857b9c24-2acd-4803-9993-27d7c0460c79.svg)  
*Position/Velocity benchmarks of Arche (left-most) vs. other Go ECS implementations.
Left panel: query iteration (log scale), right panel: world setup and entity creation.*
</div>

#### Add/remove component

Build:
* Create 1000 entities with `Pos{float64, float64}`.

Iterate:
* Get all entities with `Pos`, and add `Vel{float64, float64}` component.
* Get all entities with `Pos` and `Vel`, and remove `Vel` component.

> Note: The iteration is performed once before benchmarking,
> to avoid biasing slower implementations through one-time allocations.

<div align="center" width="100%">

![Benchmark vs. Go ECSs - Add/remove](https://user-images.githubusercontent.com/44003176/227808722-45dd7e09-5a0d-4f66-83de-5b9ce1ffc250.svg)  
*Add/remove component benchmarks of Arche (left-most) vs. other Go ECS implementations.
Left panel: iteration, right panel: world setup and entity creation.*
</div>

### Arche vs. Array of Structs

The plot below shows CPU time benchmarks of Arche (black) vs. Array of Structs (AoS, red) and Array of Pointers (AoP, blue) (with structs escaped to the heap).

Arche takes a constant time of just over 2ns per entity, regardless of the memory per entity (x-axis) and the number of entities (line styles).
For AoS and AoP, time per access increases with memory per entity as well as number of entities, due to cache misses.

In the given example with components of 16 bytes each, from 64 bytes per entity onwards (i.e. 4 components or 8 `float64` values),
Arche outperforms AoS and AoP, particularly with a large number of entities.

<div align="center" width="100%">

![Benchmark vs. AoS and AoP](https://user-images.githubusercontent.com/44003176/227033985-587b5d02-c159-4e00-b2a5-dc8b2f65b124.svg)  
*CPU benchmarks of Arche (black) vs. Array of Structs (AoS, red) and Array of Pointers (AoP, blue).*
</div>

## Cite as

Lange, M. (2023): Arche &ndash; An archetype-based Entity Component System for Go. GitHub repository: https://github.com/mlange-42/arche

## License

This project is distributed under the [MIT licence](./LICENSE).
