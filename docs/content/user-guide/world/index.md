+++
title = 'The World'
weight = 30
description = "The World as Arche's central data storage."
+++

The *World* ({{< api ecs World >}}) is the central data storage in Arche.
It manages and stores entities ({{< api ecs Entity >}}), their components, as well as [Resources](./resources).
For the internal structure of the world, see chapter [Architecture](/background/architecture).

Here, we only deal with world creation.
Most world functionality is covered in chapters [Entities & Components](./entities) and [World entity access](./world-access).

## World creation

To create a world with default settings, use {{< api ecs NewWorld >}}:

{{< code-func world_test.go TestWorldSimple >}}

A world can also be configured with a capacity increment, using an {{< api ecs Config >}}:

{{< code-func world_test.go TestWorldConfig >}}

The capacity increment determines by how many entities an archetype grows when it reaches it's capacity.

For archetypes with an [Entity Relation](./relations), a separate capacity increment can be specified:

{{< code-func world_test.go TestWorldConfigRelations >}}

## Reset the world

For systematic simulations, it is possible to reset a populated world for reuse:

{{< code-func world_test.go TestWorldReset >}}
