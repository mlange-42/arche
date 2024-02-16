+++
title = 'Arche'
no_heading = true
+++
{{< html >}}
<img src="images/logo-light.svg" alt="Arche" class="light" style="width: 100%; max-width: 680px; margin:24px auto 24px auto;"/>
<img src="images/logo-dark.svg" alt="Arche" class="dark" style="width: 100%; max-width: 680px; margin:24px auto 24px auto;"/>

<div style="width 100%; text-align: center;">
<a href="https://github.com/mlange-42/arche/actions/workflows/tests.yml" style="display:inline-block">
<img alt="Test status" src="https://img.shields.io/github/actions/workflow/status/mlange-42/arche/tests.yml?branch=main&label=Tests&logo=github"></img></a>

<a href="https://badge.coveralls.io/github/mlange-42/arche?branch=main" style="display:inline-block">
<img alt="Coverage Status" src="https://img.shields.io/coverallsCoverage/github/mlange-42/arche?logo=coveralls"></img></a>

<a href="https://goreportcard.com/report/github.com/mlange-42/arche" style="display:inline-block">
<img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/mlange-42/arche"></img></a>

<a href="https://mlange-42.github.io/arche/" style="display:inline-block">
<img alt="User Guide" src="https://img.shields.io/badge/user_guide-%23007D9C?logo=go&logoColor=white&labelColor=gray"></img></a>

<a href="https://pkg.go.dev/github.com/mlange-42/arche" style="display:inline-block">
<img alt="Go Reference" src="https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray"></img></a>

<a href="https://github.com/mlange-42/arche" style="display:inline-block">
<img alt="GitHub" src="https://img.shields.io/badge/github-repo-blue?logo=github"></img></a>

<a href="https://doi.org/10.5281/zenodo.7656484" style="display:inline-block">
<img alt="DOI:10.5281/zenodo.7656484" src="https://img.shields.io/badge/10.5281%2Fzenodo.7656484-blue?label=doi"></img></a>

<a href="https://github.com/mlange-42/arche/blob/main/LICENSE" style="display:inline-block">
<img alt="MIT license" src="https://img.shields.io/badge/MIT-brightgreen?label=license"></img></a>
</div>
{{< /html >}}

*Arche* is an [archetype](/background/architecture)-based [Entity Component System](https://en.wikipedia.org/wiki/Entity_component_system) for [Go](https://go.dev/).

*Arche* is designed for the use in simulation models of the
[Department of Ecological Modelling](https://www.ufz.de/index.php?en=34213) at the
[Helmholtz Centre for Environmental Research](https://www.ufz.de).

## Arche's Features

- Simple core API. See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche).
- Optional logic [filter](https://pkg.go.dev/github.com/mlange-42/arche/filter) and type-safe [generic](https://pkg.go.dev/github.com/mlange-42/arche/generic) API.
- Entity relations as first-class feature. See the [User Guide](https://mlange-42.github.io/arche/guide/relations/).
- World serialization and deserialization with [arche-serde](https://github.com/mlange-42/arche-serde).
- No systems. Just queries. Use your own structure (or the [Tools](https://github.com/mlange-42/arche#tools)).
- No dependencies. Except for unit tests ([100% coverage](https://coveralls.io/github/mlange-42/arche)).
- Probably the fastest Go ECS out there. See the [Benchmarks](https://github.com/mlange-42/arche#benchmarks).

For more information, see the [GitHub repository](https://github.com/mlange-42/arche) and [API docs](https://pkg.go.dev/github.com/mlange-42/arche).

## Cite as

Lange, M. (2023): Arche &ndash; An archetype-based Entity Component System for Go. DOI [10.5281/zenodo.7656484](https://doi.org/10.5281/zenodo.7656484),  GitHub repository: https://github.com/mlange-42/arche

## Contributing

[Open an issue](https://github.com/mlange-42/arche/issues/new) in the {{< repo "" "GitHub repository" >}} if you have questions, feedback, feature ideas or want to report a bug.
Pull requests are welcome.

## License

*Arche* and all its sources are released under the [MIT License](https://github.com/mlange-42/arche/blob/main/LICENSE).
