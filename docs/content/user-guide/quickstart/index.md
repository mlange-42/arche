+++
title = 'Quickstart'
weight = 10
description = 'Quickstart guide to install and use Arche.'
+++
This page shows how to install Arche, and gives a minimal usage example.

## Installation

To use *Arche* in a Go project, run:

```bash
go get github.com/mlange-42/arche
```

## Usage example

Here is the classical Position/Velocity example that every ECS shows in the docs.
It uses the type-safe [generic](https://pkg.go.dev/github.com/mlange-42/arche/generic) API.

{{< code example_test.go >}}
