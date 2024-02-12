# Arche -- Introduction

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/arche/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/arche/actions/workflows/tests.yml)
[![Coverage Status](https://badge.coveralls.io/repos/github/mlange-42/arche/badge.svg?branch=main)](https://badge.coveralls.io/github/mlange-42/arche?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/arche)](https://goreportcard.com/report/github.com/mlange-42/arche)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/arche.svg)](https://pkg.go.dev/github.com/mlange-42/arche)
[![Manual](https://img.shields.io/static/v1?label=manual&message=mdBook&color=blue&logo=mdBook)](https://mlange-42.github.io/arche/)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/arche)
[![MIT license](https://img.shields.io/github/license/mlange-42/arche)](https://github.com/mlange-42/arche/blob/main/LICENSE)

*Arche* is an [archetype](https://github.com/mlange-42/arche/blob/main/ARCHITECTURE.md)-based [Entity Component System](https://en.wikipedia.org/wiki/Entity_component_system) for [Go](https://go.dev/).

*Arche* is designed for the use in simulation models of the
[Department of Ecological Modelling](https://www.ufz.de/index.php?en=34213) at the
[Helmholtz Centre for Environmental Research](https://www.ufz.de).

## Features

* Simple [core API](https://pkg.go.dev/github.com/mlange-42/arche/ecs). See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche).
* Optional logic [filter](https://pkg.go.dev/github.com/mlange-42/arche/filter) and type-safe [generic](https://pkg.go.dev/github.com/mlange-42/arche/generic) API.
* Entity relations as first-class feature. See [Architecture](https://github.com/mlange-42/arche/blob/main/ARCHITECTURE.md#entity-relations).
* No systems. Just queries. Use your own structure (or the [Tools](#tools)).
* No dependencies. Except for unit tests ([100% coverage](https://coveralls.io/github/mlange-42/arche)).
* Probably the fastest Go ECS out there. See the [Benchmarks](#benchmarks).

## GitHub project

*Arche* is Open Source and available [on GitHub](https://github.com/mlange-42/arche).

## Contributing

For questions, feature requests and issues, please use the [issue tracker](https://github.com/mlange-42/arche/issues). Merge requests are welcome.

## License

*Arche* and all its sources are released under the [MIT License](https://github.com/mlange-42/track/blob/main/LICENSE).
