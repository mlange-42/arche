[![Arche (logo)](https://user-images.githubusercontent.com/44003176/236701164-28178d13-7e52-4449-baa4-41b764183cbd.png)](https://github.com/mlange-42/arche)
[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/arche/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/arche/actions/workflows/tests.yml)
[![Coverage Status](https://img.shields.io/coverallsCoverage/github/mlange-42/arche?logo=coveralls)](https://badge.coveralls.io/github/mlange-42/arche?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/arche)](https://goreportcard.com/report/github.com/mlange-42/arche)
[![User Guide](https://img.shields.io/badge/user_guide-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://mlange-42.github.io/arche/)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/arche)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/arche)
[![DOI:10.5281/zenodo.7656484](https://img.shields.io/badge/10.5281%2Fzenodo.7656484-blue?label=doi)](https://doi.org/10.5281/zenodo.7656484)
[![MIT license](https://img.shields.io/badge/MIT-brightgreen?label=license)](https://github.com/mlange-42/arche/blob/main/LICENSE)

*Arche* is an [archetype](https://mlange-42.github.io/arche/background/architecture/)-based [Entity Component System](https://en.wikipedia.org/wiki/Entity_component_system) for [Go](https://go.dev/).

<div align="center" width="100%">

&mdash;&mdash;

[Features](#features) &nbsp; &bull; &nbsp; [Installation](#installation) &nbsp; &bull; &nbsp; [Usage](#usage) &nbsp; &bull; &nbsp; [Tools](#tools) &nbsp; &bull; &nbsp; [Design](#design) &nbsp; &bull; &nbsp; [Benchmarks](#benchmarks)
</div>

## Features

* Designed for performance and highly optimized. See the [Benchmarks](#benchmarks).
* Well-documented [API](https://pkg.go.dev/github.com/mlange-42/arche) and comprehensive [User Guide](https://mlange-42.github.io/arche/).
* No systems. Just queries. Use your own structure (or the [Tools](#tools)).
* No dependencies. Except for unit tests (100% [test coverage](https://coveralls.io/github/mlange-42/arche)).
* World serialization and deserialization with [arche-serde](https://github.com/mlange-42/arche-serde).

> [!IMPORTANT]
> Arche has a successor: [Ark](https://github.com/mlange-42/ark).
>  - If you are new here, use [Ark](https://github.com/mlange-42/ark).
>  - If you use Arche's generic API, consider migrating to [Ark](https://github.com/mlange-42/ark). It will be easy.
>  - If you use Arche's ID-based API, stay with Arche.
> 
> Arche will still be maintained further.

## Installation

To use *Arche* in a Go project, run:

```shell
go get github.com/mlange-42/arche
```

## Usage

Here is the classical Position/Velocity example that every ECS shows in the docs.
It uses the type-safe [generic](https://pkg.go.dev/github.com/mlange-42/arche/generic) API.

See the [User Guide](https://mlange-42.github.io/arche/), [API docs](https://pkg.go.dev/github.com/mlange-42/arche) and
[examples](https://github.com/mlange-42/arche/tree/main/_examples) for details.
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
its purpose is to get started with prototyping and developing simulation models immediately, focussing on the model logic.
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

The type-safe generic API and advanced logic filters are provided in the packages
[generic](https://pkg.go.dev/github.com/mlange-42/arche/generic) and
[filter](https://pkg.go.dev/github.com/mlange-42/arche/filter), respectively.
Both packages are built on top of the core [ecs](https://pkg.go.dev/github.com/mlange-42/arche/ecs) package, so they could also be implemented by a user.

### Determinism

Iteration order in *Arche* is deterministic and reproducible.
This does not mean that entities are iterated in their order of insertion, nor in the same order in successive iterations.
However, given the same operations on the `ecs.World`, iteration order will always be the same.

### Strict and panic

*Arche* puts an emphasis on safety and on avoiding undefined behavior.
It panics on unexpected operations, like removing a dead entity,
adding a component that is already present, or attempting to change a locked world.
This may not seem idiomatic for Go.
However, explicit error handling in performance hot spots is not an option.
Neither is silent failure, given the scientific background.

### Limitations

* The number of component types per `World` is limited to 256. This is mainly a performance decision.
* The number of entities alive at any one time is limited to just under 5 billion (`uint32` ID).

## Benchmarks

A tabular overview of the runtime cost of typical *Arche* ECS operations is provided under [benchmarks](https://mlange-42.github.io/arche/background/benchmarks/) in the Arche's [User Guide](https://mlange-42.github.io/arche/).

For a benchmark comparison with other ECS implementations,
see the [go-ecs-benchmarks](https://github.com/mlange-42/go-ecs-benchmarks) repository.

## Cite as

Lange, M. (2023): Arche &ndash; An archetype-based Entity Component System for Go. DOI [10.5281/zenodo.7656484](https://doi.org/10.5281/zenodo.7656484),  GitHub repository: https://github.com/mlange-42/arche

## License

This project is distributed under the [MIT license](./LICENSE).
