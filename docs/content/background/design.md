+++
title = 'Design Philosophy'
weight = 20
description = 'Specific design considerations behind Arche.'
+++
Unlike most other ECS implementations, Arche is designed for the development of scientific,
individual-based models rather than for game development.
This motivates some design decisions, with an emphasis on simplicity, safety and performance.
Nevertheless, Arche can also be used for game development.

### Simple core API

The {{< api ecs World >}} object is a pure and simple ECS implementation in the sense of a data store
for entities and components, with query and iteration capabilities.
More advanced features like batch operations or entity relations are provided through separate objects.

There is neither an update loop nor systems.
These should be implemented by the user.
For a batteries-included implementation, see module [arche-model](https://github.com/mlange-42/arche-model).

The packages {{< api filter >}} and {{< api generic >}}
provide a layer around the core for richer resp. generic queries and manipulation.
They are built on top of the {{< api ecs >}} package, so they could also be implemented by a user.

### Determinism

Iteration order in Arche is deterministic and reproducible.
This does not mean that entities are iterated in their order of insertion, nor in the same order in successive iterations.
However, given the same operations on the {{< api ecs World >}}, iteration order will always be the same.

### Strict and panic

Arche puts an emphasis on safety and on avoiding undefined behavior.
It panics on unexpected operations, like removing a dead entity,
adding a component that is already present, or attempting to change a locked world.
This may not seem idiomatic for Go.
However, explicit error handling in performance hotspots is not an option.
Neither is silent failure, given the scientific background.

### Other limitations

* The number of component types per `World` is limited to 256. This is mainly a performance decision.
* The number of entities alive at any one time is limited to just under 5 billion (`uint32` ID).
