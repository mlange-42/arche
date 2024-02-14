+++
title = 'Architecture'
weight = 30
description = "Arche's internal ECS architecture."
+++
Arche uses an archetype-based architecture. Therefore the name :wink:.

## Archetypes

The ASCII graph below illustrates the approach.
Components for entities are stored in so-called archetypes. Archetypes represent unique combinations of components. This means that component data for all entities with exactly the same component is stored in the same archetype.

In the illustration below, the first archetype holds entities with (only/exactly) the components A, B and C,
as well as their components.
Similarly, the second archetype contains all entities with A and C, and their components.

```text
 Entities   Archetypes   Bitmasks   Queries

   E         E Comps
  |0|       |2|A|B|C|    111...<-.<--match-.
  |1|---.   |8|A|B|C|            |         |
  |2|   '-->|1|A|B|C|            |         |
  |3|       |3|A|B|C|            |--(A, C) |
  |4|                            |  101... |
  |6|   .-->|7|A|C|      101...<-'         |--(B)
  |7|---'   |6|A|C|                        |  010...
  |8|       |4|A|C|                        |
  |9|---.                                  |
  |.|   |   |5|B|C|      011...   <--------'
  |.|   '-->|9|B|C|
  |.|
  |.| <===> [Entity pool]
```
*Illustration of Arche's archetype-based architecture.*

The exact component composition of each archetype is encoded in a bitmask for fast comparison.
Thus, queries can easily identify their relevant archetypes, and then simply iterate entities linearly, which is very fast. Components can be accessed through a query in a very efficient way (&approx;1ns).

## World entity access

For getting components by entity, e.g. for hierarchies, the world contains a list that is indexed by the entity ID (left-most in the figure). For each entity, it references it's current archetype and the position of the entity in the archetype. This way, getting components for entities (i.e. random access) is fast, although not as fast as in queries (≈2ns vs. 1ns).

Note that the entities list also contains entities that are currently not alive,
because they were removed.
These entities are recycled when new entities are requested from the {{< api ecs World >}}.
Therefore, besides the ID shown in the illustration, each entity also has a generation
variable. It is incremented on each "reincarnation" of an entity.
Thus, it allows to distinguish recycled from dead entities, as well as from previous or later "incarnations".

## Performance

Obviously, archetypes are an optimization for iteration speed.
But they also come with a downside. Adding or removing components to/from an entity requires moving all the components of the entity to another archetype.
This takes roughly 10-20ns per involved component.
It is therefore recommended to add/remove/exchange multiple components at the same time rather than one after the other.

However, as the benchmarks on the {{< repo "#benchmarks" "repo README" >}} illustrate,
Arche seems to be the fastest Go ECS available.
Not only in terms of iteration speed, which is particularly tailored for.
Even when it comes to entity manipulation (adding, removing components etc.),
where sparse-set ECS implementations should shine, Arche leads the field.

For more numbers on performance, see chapter [Benchmarks](./benchmarks). 

## ToDos

 - Archetype graph
 - Relation archetypes
