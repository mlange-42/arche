+++
title = 'Statistics'
weight = 120
description = "Arche's world statistics feature for engine insights."
+++
Arche only exposes the API required for actual use.
Therefore, internals like the number of archetypes, memory used to store components etc. are not directly accessible.

However, it might sometimes be useful to have access to such metrics,
for example in order to judge effects of different ways of implementing something.
Otherwise, users would have to rely on logic reasoning and sufficient understanding of Arche to derive these numbers.

For that sake, Arche provides statistics about its internals, prepared in a compact and digestible way.

## Accessing statistics

All internal statistics can be accessed via {{< api ecs World.Stats >}},
which returns a {{< api "ecs/stats" World stats.World >}} object.
This, in turn, contains the other stats types described below.
All these types have a method `String()` to bring them into a compact, human-readable form. 

## World stats

## Entity stats

## Node stats

## Archetype stats

