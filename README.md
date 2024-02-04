[![Arche (logo)](https://user-images.githubusercontent.com/44003176/236701164-28178d13-7e52-4449-baa4-41b764183cbd.png)](https://github.com/mlange-42/arche)
[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/arche/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/arche/actions/workflows/tests.yml)
[![Coverage Status](https://badge.coveralls.io/repos/github/mlange-42/arche/badge.svg?branch=main)](https://badge.coveralls.io/github/mlange-42/arche?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/arche)](https://goreportcard.com/report/github.com/mlange-42/arche)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/arche.svg)](https://pkg.go.dev/github.com/mlange-42/arche)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/arche)
[![DOI:10.5281/zenodo.7656484](https://zenodo.org/badge/DOI/10.5281/zenodo.7656484.svg)](https://doi.org/10.5281/zenodo.7656484)
[![MIT license](https://img.shields.io/github/license/mlange-42/arche)](https://github.com/mlange-42/arche/blob/main/LICENSE)

*Arche* is an [archetype](https://github.com/mlange-42/arche/blob/main/ARCHITECTURE.md)-based [Entity Component System](https://en.wikipedia.org/wiki/Entity_component_system) for [Go](https://go.dev/).

*Arche* is designed for the use in simulation models of the
[Department of Ecological Modelling](https://www.ufz.de/index.php?en=34213) at the
[Helmholtz Centre for Environmental Research](https://www.ufz.de).

<div align="center" width="100%">

&mdash;&mdash;

[Features](#features) &nbsp; &bull; &nbsp; [Installation](#installation) &nbsp; &bull; &nbsp; [Usage](#usage) &nbsp; &bull; &nbsp; [Tools](#tools) &nbsp; &bull; &nbsp; [Design](#design) &nbsp; &bull; &nbsp; [Benchmarks](#benchmarks)
</div>

## Features

* Simple [core API](https://pkg.go.dev/github.com/mlange-42/arche/ecs). See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche).
* Optional logic [filter](https://pkg.go.dev/github.com/mlange-42/arche/filter) and type-safe [generic](https://pkg.go.dev/github.com/mlange-42/arche/generic) API.
* Entity relations as first-class feature. See [Architecture](https://github.com/mlange-42/arche/blob/main/ARCHITECTURE.md#entity-relations).
* World serialization and deserialization with [arche-serde](https://github.com/mlange-42/arche-serde).
* No systems. Just queries. Use your own structure (or the [Tools](#tools)).
* No dependencies. Except for unit tests ([100% coverage](https://coveralls.io/github/mlange-42/arche)).
* Probably the fastest Go ECS out there. See the [Benchmarks](#benchmarks).

## Installation

To use *Arche* in a Go project, run:

```shell
go get github.com/mlange-42/arche
```

## Usage

Here is a minimal usage example.
It uses the type-safe [generic](https://pkg.go.dev/github.com/mlange-42/arche/generic) API.
For a full-featured wrapper with systems, scheduling and more, see [arche-model](https://github.com/mlange-42/arche-model).

See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche) and
[examples](https://github.com/mlange-42/arche/tree/main/examples) for details.
For more complex examples, see [arche-demo](https://github.com/mlange-42/arche-demo).

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
		entity := mapper.New()
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

## Tools

Several tools for *Arche* are provided in separate modules:

* [arche-serde](https://github.com/mlange-42/arche-serde) provides JSON serialization and deserialization for *Arche*'s World.
* [arche-model](https://github.com/mlange-42/arche-model) provides a wrapper around *Arche*, and some common systems and resources.
It's purpose is to get started with prototyping and developing simulation models immediately, focussing on the model logic.
* [arche-pixel](https://github.com/mlange-42/arche-pixel) provides OpenGL graphics and live plots for *Arche* using the [Pixel](https://github.com/gopxl/pixel) game engine.
* [arche-demo](https://github.com/mlange-42/arche-demo) provides examples of *Arche* models, which can be viewed in a [live demo](https://mlange-42.github.io/arche-demo/).

## Design

Unlike most other ECS implementations, *Arche* is designed for the development of scientific,
individual-based models rather than for game development.
This motivates some design decisions, with an emphasis on simplicity, safety and performance.
Nevertheless, *Arche* can also be used for game development.

### Simple core API

The `ecs.World` object is a pure and simple ECS implementation in the sense of a data store
for entities and components, with query and iteration capabilities.
More advanced features like batch operations or entity relations are provided through separate objects.

There is neither an update loop nor systems.
These should be implemented by the user.
For a batteries-included implementation, see module [arche-model](https://github.com/mlange-42/arche-model).

The packages [filter](https://pkg.go.dev/github.com/mlange-42/arche/filter)
and [generic](https://pkg.go.dev/github.com/mlange-42/arche/generic)
provide a layer around the core for richer resp. generic queries and manipulation.
They are built on top of the `ecs` package, so they could also be implemented by a user.

### Determinism

Iteration order in *Arche* is deterministic and reproducible.
This does not mean that entities are iterated in their order of insertion, nor in the same order in successive iterations.
However, given the same operations on the `ecs.World`, iteration order will always be the same.

### Strict and panic

*Arche* puts an emphasis on safety and on avoiding undefined behavior.
It panics on unexpected operations, like removing a dead entity,
adding a component that is already present, or attempting to change a locked world.
This may not seem idiomatic for Go.
However, explicit error handling in performance hotspots is not an option.
Neither is silent failure, given the scientific background.

### Other limitations

* The number of component types per `World` is limited to 256. This is mainly a performance decision.
* The number of entities alive at any one time is limited to just under 5 billion (`uint32` ID).

## Benchmarks

See also the latest [Benchmarks CI run](https://github.com/mlange-42/arche/actions/workflows/benchmarks.yml).

### Arche vs. other Go ECS implementations

To the best of the author's knowledge, there are only a handful of ECS implementations in Go that are serious and somewhat maintained:

* [go-gameengine-ecs](https://github.com/marioolofo/go-gameengine-ecs)
* [Donburi](https://github.com/yohamta/donburi)
* [Ento](https://github.com/wwfranczyk/ento)
* [unitoftime/ecs](https://github.com/unitoftime/ecs)

Here, *Arche* is benchmarked against these implementations.
Feel free to open an issue if you have suggestions for improvements on the benchmarking code or other engines to include.

#### Position/Velocity

Build:
* Create 1000 entities with `Pos{float64, float64}` and `Vel{float64, float64}`.
* Create 9000 entities with only `Pos{float64, float64}`.

Iterate:
* Iterate all entities with `Pos` and `Vel`, and add `Vel` to `Pos`.

<div align="center" width="100%">

![Benchmark vs. Go ECSs - Pos/Vel](https://github.com/mlange-42/arche/assets/44003176/7b73f9d8-238c-4d7a-98a1-267ad0b5e4a8)  
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

![Benchmark vs. Go ECSs - Add/remove](https://github.com/mlange-42/arche/assets/44003176/7a127568-e71a-441f-91b0-6e626b3fcf19)  
*Add/remove component benchmarks of Arche (left-most) vs. other Go ECS implementations.
Left panel: iteration, right panel: world setup and entity creation.*
</div>

### Arche vs. Array of Structs

The plot below shows CPU time benchmarks of Arche (black) vs. Array of Structs (AoS, red) and Array of Pointers (AoP, blue) (with structs escaped to the heap).

Arche takes a constant time of just over 2ns per entity, regardless of the memory per entity (x-axis) and the number of entities (line styles).
For AoS and AoP, time per access increases with memory per entity as well as number of entities, due to cache misses.

In the given example with components of 16 bytes each, from 64 bytes per entity onwards (i.e. 4 components or 8 `float64` values),
Arche outperforms AoS and AoP, particularly with a large number of entities.
Note that the maximum shown here corresponds to only 25 MB of entity data!

<div align="center" width="100%">

![Benchmark vs. AoS and AoP](https://user-images.githubusercontent.com/44003176/237245154-0070bba0-c8fe-447e-a710-e370af1dcdab.svg)  
*CPU benchmarks of Arche (black) vs. Array of Structs (AoS, red) and Array of Pointers (AoP, blue).*
</div>

## Cite as

Lange, M. (2023): Arche &ndash; An archetype-based Entity Component System for Go. DOI [10.5281/zenodo.7656484](https://doi.org/10.5281/zenodo.7656484),  GitHub repository: https://github.com/mlange-42/arche

## License

This project is distributed under the [MIT license](./LICENSE).
