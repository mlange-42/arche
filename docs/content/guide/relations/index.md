+++
title = 'Entity Relations'
type = "docs"
weight = 80
description = "Arche's entity relationships feature."
+++
In a basic ECS, relations between entities, like hierarchies, can be represented
by storing entities in components.
E.g., we could have a child component like this:

```go
type ChildOf struct {
    Parent ecs.Entity
}
```

Or, alternatively, a parent component with many children:

```go
type Parent struct {
    Children []ecs.Entity
}
```

In conjunction with [World Entity Access](../world-access), this is often sufficient.
However, we are not able to leverage the power of queries to e.g. get all children of a particular parent.

To make entity relations even more useful and efficient, Arche supports them as first class feature.
Relations are added to and removed from entities just like components,
and hence can be queried like components, with the usual efficiency.
This is achieved by creating separate [archetypes](../../background/architecture#archetypes)
for relations with different target entities.

## Relation components

To use entity relations, create components that have *embedded* an {{< api ecs Relation >}} as their first member:

```go
type ChildOf struct {
    ecs.Relation
}
```

That's all to make a component be treated as an entity relation by Arche.
Thus, we have created a relation type. When added to an entity, a target entity for the relation can be defined.

> [!TIP]
> Note that each entity can only have one relation component. See section [Limitations](#limitations).

## Creating relations

### On new entities

When creating entities, we can use an {{< api ecs Builder >}} to set a relation target.
In the generic API, we use a *MapX* (e.g. {{< api generic Map2 >}}).

{{< tabs items="generic,ID-based" >}}
{{< tab >}}
{{< code-func relations_test.go TestCreateEntityGeneric >}}
{{< /tab >}}
{{< tab >}}
{{< code-func relations_test.go TestCreateEntity >}}
{{< /tab >}}
{{< /tabs >}}

### When adding components

A relation target can also be given when adding a relation component.
With the ID-based API, we use the helper {{< api ecs World.Relations >}} for this,
like for most operations on entity relations.
In the generic API, we use a *MapX* (e.g. {{< api generic Map2 >}}) again.

{{< tabs items="generic,ID-based" >}}
{{< tab >}}
{{< code-func relations_test.go TestAddRelationGeneric >}}
{{< /tab >}}
{{< tab >}}
{{< code-func relations_test.go TestAddRelation >}}
{{< /tab >}}
{{< /tabs >}}

Alternatively, we can use a {{< api generic.Exchange >}}:

{{< code-func relations_test.go TestAddRelationGenericExchange >}}

## Set and get relations

We can also change the target of an already assigned relation component.
This is done via {{< api ecs Relations.Set >}} or {{< api generic Map.SetRelation >}}:

{{< tabs items="generic,ID-based" >}}
{{< tab >}}
{{< code-func relations_test.go TestSetRelationGeneric >}}
{{< /tab >}}
{{< tab >}}
{{< code-func relations_test.go TestSetRelation >}}
{{< /tab >}}
{{< /tabs >}}

Similarly, relation targets can be obtained with {{< api ecs Relations.Get >}} or {{< api generic Map.GetRelation >}}:

{{< tabs items="generic,ID-based" >}}
{{< tab >}}
{{< code-func relations_test.go TestGetRelationGeneric >}}
{{< /tab >}}
{{< tab >}}
{{< code-func relations_test.go TestGetRelation >}}
{{< /tab >}}
{{< /tabs >}}

## Querying relations

And now for the best: querying for entities that have a certain relation and target.

In the ID-based API, relation targets can be queries with {{< api ecs RelationFilter >}}.
In the generic API, it is supported by all *FilterX* via e.g. {{< api generic Filter2.WithRelation >}}.

{{< tabs items="generic,ID-based" >}}
{{< tab >}}
{{< code-func relations_test.go TestRelationQueryGeneric >}}
{{< /tab >}}
{{< tab >}}
{{< code-func relations_test.go TestRelationQuery >}}
{{< /tab >}}
{{< /tabs >}}

## Limitations

Entity relations in Arche are inspired by [Flecs](https://github.com/SanderMertens/flecs).
However, the implementation in *Arche* is currently limited in that it only supports a single relation per entity, and no chained (or nested) relation queries.

## When to use, and when not

When using Arche's entity relations, an archetype is created for each target entity of a relation.
Thus, entity relations are not efficient if the number of target entities is high (tens of thousands),
while only a low number of entities has a relation to each particular target (less than a few dozens).
Particularly in the extreme case of 1:1 relations, storing entities in components
as explained in the introduction of this chapter is more efficient.

However, with a moderate number of relation targets, particularly with many entities per target,
entity relations are very efficient. See section [Benchmarks](#benchmarks) below, for a comparison of different ways to represent entity relations.

Beyond use cases where the relation target is a "physical" entity that appears
in a simulation or game, targets can also be more abstract, like categories.
Examples:

 - Different tree species in a forest model
 - Behavioral states in a finite state machine
 - The opposing factions in a strategy game
 - Render layers in a game or other graphical application

This concept is particularly useful for things that would best be expressed by components,
but the possible components (or categories) are only known at runtime.
Thus, it is not possible to create ordinary components for them.
However, these categories can be represented by entities, which are used as relation targets.

See the last section of this chapter ([Longer example](#longer-example))
for an implementation of the tree species example above.

## Benchmarks

The figure below compares the iteration time per entity for different ways of representing entity relations.
The task is to sum up a value over the children of each parent.

The following ways to represent entity relations are shown in the figure:

* *ParentList* (purple): Children form an implicit linked list. The parent references the first child.
  * Query over parents, inner loop implicit linked list of children, using world access for next child and value component.
* *ParentSlice* (red): The parent holds a slice of all its children.
  * Query over parents, inner loop over slice of children using world access for value component.
* *Child* (green): Each child references its parent.
  * Query over all child entities and retrieval of the parent sum component using world access.
* *Default* (blue): Using Arche's relations feature without filter caching.
  * Outer query over parents, inner loop over children using relation queries.
* *Cached* (black): Using Arche's relations feature with filter caching.
  * Same as above, using an additional component per parent to store cached filters.

The first three representations are possible in any ECS, while the last two use Arche's entity relations feature.

![Benchmarks Entity relations](https://user-images.githubusercontent.com/44003176/238461931-7824bfeb-4a03-49e8-9de8-0650032259c0.svg)  
*Iteration time per entity for different ways of representing entity relations. Color: ways to represent entity relations; Line style: total number of child entities; Markers: number of children per parent entity*

The benchmarks show that Arche's relations feature outperforms the other representations, except when there are very few children per parent.
Only when there is a huge number of parents and significantly fewer than 100 children per parent,
the *Child* representation should perform better.

The {{< repo "tree/main/benchmark/competition" "benchmark code" >}}
can be found in the {{< repo "" "GitHub repository" >}}.

## Longer example

To conclude this chapter, here is a longer example that uses Arche' entity relations feature
to represent tree species in a forest model.

{{< code relations_example_test.go >}}
