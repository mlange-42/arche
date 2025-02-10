+++
title = 'Performance Tips'
type = "docs"
weight = 130
description = "Performance tips on how to get the most out of Arche."
+++
Arche is highly optimized and can compete with mature ECS implementations in terms of performance.
It will probably not get into your way in this regard.
Experience shows that in simulation built with Arche, ECS code like queries, entity creation etc.
takes between 5% and 15% of the total CPU time.
Keep in mind that this is not "on top" of the simulation, but replaces the overhead any other implementation for storing and iterating entities would have.
Due to its cache-friendliness, [archetype-based ECS](../../background/architecture) can outperform e.g. an [Array of Structs](https://en.wikipedia.org/wiki/AoS_and_SoA) implementation, particularly for simulations with many entities and/or many variables per entity.

Nevertheless, each ECS has its strengths and weaknesses.
This chapter provides tips on what you should pay attention to in order to get the most out of Arche.

## Optimized for Iteration

Being an [archetype-based ECS](../../background/architecture), Arche is optimized for queries and iteration.
Adding and removing components is comparatively costly with this architecture,
because components must be moved around between archetypes.
The runtime difference between accessing a component and adding/removing a component is at least one order of magnitude.
For some numbers for comparison, see the [Benchmarks](../../background/benchmarks).

## Queries & Components

The largest potential for optimizing queries is the clever design of components.
The aim should be that queries access only data that is actually required,
while at the same time reducing the amount of accessed components.
Accessing fewer data means fewer cache misses, while accessing fewer components avoids lookups.
To access only data that it actually required primarily means that the accessed components should contain only data that is used by the query.

A component should only contain closely related data that is mostly or always accessed together.
A `Position` component with `X` and `Y` is a good example.
Vice versa, closely related data should be in the same component.
What should be avoided are all-in-one components that mimic OOP classes to represent entities.
A good (or rather, bad) example is a `Tree` component with `X`, `Y`, `Biomass`, `Height`, `StemDiameter` and `LeaveAreaIndex` (or more).

For fast memory access, the use of slices in components should be avoided. Use fixed-size arrays where possible.

## Filter caching

When working with many archetypes, queries can be sped up by caching the underlying filter.
This way, the filter is not checked against archetypes during query iteration.
Instead, the archetypes relevant for the filter are cached,
and checks are only required when new archetypes are created.

For details, see the section on caching in chapter [Filters](../filters#filter-caching).

## World access

World access to components with {{< api ecs World.Get >}} is per se slower than access through a query,
as there is one more indirection and the alive status of the entity is checked for safety.
Queries should be preferred over world access where possible.

Further, world access can't benefit from the cache-friendly linearity of query iterations.
This becomes more severe when the length of "jumps" between entities increases.
Thus, is it more efficient to randomly access among e.g. 1000 entities compared to 100k entities.

As an example, say we have 1000 parent entities, 100k child entities, and don't use [Entity Relations](../relations).
Here, it would be better to use a query over the children and access the parent of each child by world access. We jump around between 1000 entities.
Alternatively, we could query the parents and access the children of each parent by world access.
The number of accesses through the world would be the same, but we would jump between 100k entities,
which would be slower.

See also the [benchmarks](../relations/#benchmarks) for [Entity Relations](../relations).

## World operations

As explained above, operations like adding and removing components or creating entities are relatively
costly in an archetype-based ECS.
However, Arche provides some optimizations here,
and following a few principles can help keeping the performance cost at a minimum.

### Avoiding

Different components are a great way to represent different states of otherwise similar entities.
For example, it is completely valid to build a [finite state machine](https://en.wikipedia.org/wiki/Finite-state_machine)
to model behavior, using components to represent states.
However, each state transition results in moving an entity and its components between archetypes.
Thus, when transitions occur frequently (say faster than every 20 ticks),
different components are not the most efficient way to represent states.
Alternatively, states could be represented by a variable in a single component,
avoiding the overhead of moving entities between archetypes,
at the cost of overhead in the queries.

It is a matter of weighting, and potentially benchmarking,
to decide on what is represented by components in a query-able way,
and what is left to be managed inside query loops.

### Multiple at once, Exchange

As explained above, moving entities between archetypes is relatively costly.
It is necessary when adding or removing components,
but multiple components can be added or removed with a single transition between archetypes.

For that sake, methods like {{<api ecs World.Add>}} and {{<api ecs World.Remove>}}
take an arbitrary number of component IDs as arguments.
The generic API provides all types and functions for up to 12 components.
Add or remove components together instead of one after another!

Further, {{<api ecs World.Exchange>}} and {{<api generic Exchange>}}
allow to add some components and remove others at the same time.
This also requires only a single transition between archetypes.

### Batching

Arche provides batched variants of all operations like creating entities, adding and removing components, etc.
Batching can speed up operations by up to an order of magnitude.
It allows for bulk allocation of component memory and entities,
and cuts off the overhead that is otherwise required for each entity, repeatedly. 
Entity creation is the most common use case for batching.
For details, see the chapter on [Batch Operations](../batch-ops).

See also the [Benchmarks](../../background/benchmarks#entities) for batched vs. un-batched operations.
