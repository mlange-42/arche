+++
title = 'World Statistics'
type = "docs"
weight = 120
description = "Arche's world statistics feature for engine insights."
+++
Arche only exposes the API required for actual use.
Therefore, internals like the number of archetypes, memory used to store components etc. are not directly accessible.

However, it might sometimes be useful to have access to such metrics,
for example in order to judge effects of different ways of implementing something.
Otherwise, users would have to rely on logic reasoning and sufficient understanding of Arche to derive these numbers.

For that sake, Arche provides statistics about its internals, prepared in a compact and digestible form.

## Accessing statistics

All internal statistics can be accessed via {{< api ecs World.Stats >}},
which returns a {{< api "ecs/stats" World "*stats.World" >}}.
This, in turn, contains the other stats types described below.
All these types have a method `String()` to bring them into a compact, human-readable form. 

{{< code-func stats_test.go TestWorldStats >}}

Which prints:

```text
World -- Components: 2, Nodes: 3, Filters: 0, Memory: 7.0 kB, Locked: false
  Components: Position, Heading
Entities -- Used: 100, Recycled: 0, Total: 100, Capacity: 128
Node -- Components:  0, Entities:      0, Capacity:      1, Memory:     0.0 kB, Per entity:    0 B
  Components:
Node -- Components:  2, Entities:    100, Capacity:    128, Memory:     4.0 kB, Per entity:   24 B
  Components: Position, Heading
```

## World stats

{{< api "ecs/stats" World stats.World >}} provides world information like a list of all component types
and the total memory reserved for entities and components.
Further, it contains {{< api "ecs/stats" Entities stats.Entities >}} and
a {{< api "ecs/stats" Node stats.Node >}} for each active [archetype node](/background/architecture#archetype-graph).

## Entity stats

{{< api "ecs/stats" Entities stats.Entities >}} contains information about the entity pool,
live capacity, alive entities and available entities for recycling.

## Node stats

{{< api "ecs/stats" Node stats.Node >}} provides information about an [archetype node](/background/architecture#archetype-graph), like its components, memory in total and per entity,
and more state information.

Further, it contains a {{< api "ecs/stats" Node stats.Archetype >}} for each archetype.

## Archetype stats

{{< api "ecs/stats" Archetype stats.Archetype >}} contains size, capacity and memory information for an archetype.
