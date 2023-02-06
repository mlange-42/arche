# Arche

[![Test status](https://github.com/mlange-42/arche/actions/workflows/tests.yml/badge.svg)](https://github.com/mlange-42/arche/actions/workflows/tests.yml)
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

## Usage example

```go
package main

import (
	"math/rand"

	"github.com/mlange-42/arche/ecs"
)

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

func main() {
	world := ecs.NewWorld()

	positionID := ecs.RegisterComponent[Position](&world)
	velocityID := ecs.RegisterComponent[Velocity](&world)

	// Create entities
	for i := 0; i < 1000; i++ {
		entity := world.NewEntity()
		world.Add(entity, positionID, velocityID)
		pos := (*Position)(world.Get(entity, positionID))
		vel := (*Velocity)(world.Get(entity, velocityID))

		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100

		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}

	// Time loop
	for t := 0; t < 1000; t++ {
		// Get a fresh query
		query := world.Query(positionID, velocityID)
		// Iterate it
		for query.Next() {
			pos := (*Position)(query.Get(positionID))
			vel := (*Velocity)(query.Get(velocityID))

			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
```

## Benchmarks

[TODO]

## References

Information and inspiration taken from:

* [Sander Mertens](https://ajmmertens.medium.com/)' blog
* Michele "skypjack" Caini's blog [skypjack on software](https://skypjack.github.io/)
* [marioolofo/go-gameengine-ecs](https://github.com/marioolofo/go-gameengine-ecs) "Fast Entity Component System in Golang"

## License

This project is distributed under the [MIT licence](./LICENSE).
