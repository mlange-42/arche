# Architecture

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

## Entity relations

*Arche* supports entity relations as first class feature.
Relations are used to represent graphs of entities, e.g. hierarchies.
Relations are added to and removed from entities just like components.

In Arche, queries can specify a target entity for a relation.
These relation queries are as fast as usual queries for component compositions.
This is achieved by subdividing archetypes with a relation component by their relation target. I.e. entities that reference a different target entity are stored in separate archetypes.

The feature is inspired by [Flecs](https://github.com/SanderMertens/flecs).
However, the implementation in *Arche* is currently limited in that it only supports a single relation per entity, and no nested relation queries.

### Benchmarks

The figure below compares the iteration time per entity for different ways of representing entity relations.
The task is to sum up a value over the children of each parent.

The following ways to represent entity relations are shown in the figure:

* **ParentList** (purple): Children form an implicit linked list. The parent references the first child.
  * Query over parents, inner loop implicit linked list of children, using world access for next child and value component.
* **ParentSlice** (red): The parent holds a slice of all it's children.
  * Query over parents, inner loop over slice of children using world access for value component.
* **Child** (green): Each child references it's parent.
  * Query over all child entities and retrieval of the parent sum component using world access.
* **Default** (blue): Using Arche's relations feature without filter caching.
  * Outer query over parents, inner loop over children using relation queries.
* **Cached** (black): Using Arche's relations feature with filter caching.
  * Same as above, using an additional component per parent to store cached filters.

The first three representations are possible in any ECS, while the last two use Arche's relations feature.

<div align="center" width="100%">

![Benchmarks Entity relations](https://user-images.githubusercontent.com/44003176/238461931-7824bfeb-4a03-49e8-9de8-0650032259c0.svg)  
*Iteration time per entity for different ways of representing entity relations. Color: ways to represent entity relations; Line style: total number of child entities; Markers: number of children per parent entity*
</div>

The benchmarks show that Arche's relations feature outperforms the other representations, except when there are very few children per parent.
Only when there is a huge number of parents and significantly fewer than 100 children per parent,
the *Child* representation should perform better.
