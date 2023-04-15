## Archetype-based architecture

*Arche* uses an archetype-based architecture.

The ASCII graph below illustrates the approach.
Components for entities are stored in so-called archetypes, which represent unique combinations of components.
In the illustration, the first archetype holds all components for all entities with (only/exactly) the components A, B and C.

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

The exact composition of each archetype is encoded in a bitmask for fast comparison.
Thus, queries can easily identify their relevant archetypes, and then simply iterate entities linearly, which is very fast. Components can be accessed through the query in a very efficient way (&approx;1ns).

For getting components by entity ID, e.g. for hierarchies, the world contains a list that is indexed by the entity ID. For each entity, it references it's current archetype and the index in the archetype. This way, getting components for entity IDs (i.e. random access) is fast, although not as fast as in queries (≈1.5ns vs. 1ns).

Obviously, archetypes are an optimization for iteration speed.
But they also come with a downside. Adding or removing components to/from an entity requires moving all the components of the entity to another archetype.
This takes around 20ns per involved component.
It is therefore recommended to add/remove/exchange multiple components at the same time rather than one after the other.
