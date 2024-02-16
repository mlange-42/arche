+++
title = 'Quickstart'
weight = 1
description = 'Quickstart guide to install and use Arche.'
+++
This page shows how to install Arche, and gives a minimal usage example.

Finally, it points into possible directions to continue.

## Installation

To use Arche in a Go project, run:

```bash
go get github.com/mlange-42/arche
```

## Usage example

Here is the classical Position/Velocity example that every ECS shows in the docs.
It uses the type-safe  {{< api generic >}} API.

{{< code example_test.go >}}

## What's next?

If you ask **"What is ECS?"**, take a look at the great [**ECS FAQ**](https://github.com/SanderMertens/ecs-faq) by Sander Mertens, the author of the [Flecs](http://flecs.dev) ECS.

To learn how to use Arche, read the [User Guide](/guide),
browse the [API documentation](https://pkg.go.dev/github.com/mlange-42/arche),
or take a look at the {{< repo "tree/main/_examples" examples >}} in the {{< repo "" "GitHub repository" >}}.

You can also read about Arche's [Design Philosophy](/background/design)
and [Architecture](/background/architecture) for more background information.

See the [Benchmarks](/background/benchmarks) if you are interested in some numbers on Arche's performance.
