+++
title = 'Entities & Components'
weight = 40
description = 'Entities, components, creation and removal in Arche.'
+++

## Entities

An *Entity* ({{< api ecs Entity >}}) in Arche is merely an ID and contains no data itself.

The only method of an entity is {{< api ecs Entity.IsZero >}}.
The only entity that can be directly created by the user is the zero entity, in two possible ways:

{{< code-func entities_test.go TestZeroEntity >}}

All other entities must be created through the {{< api ecs World >}} (see section [Creating entities](#creating-entities) below)

## Components

With each entity, an arbitrary number of *Components* can be associated.
Components are simple, user-defined Go `struct`s (or other go types):

```go
// Position component
type Position struct {
    X float64
    Y float64
}

// Heading component
type Heading struct {
    Angle float64
}
```

Components are stored in the [World](./world) and accessed through [Queries](./queries) or
through the world itself (see [World Entity Access](./world-access)).

### Component IDs

Each component type has a unique ID, which is used to access it in the [ID-based API](./apis).
Component IDs can be registered as well as obtained through {{< api ecs ComponentID >}}.

{{< code-func entities_test.go TestComponentID >}}

When a type is used as a component the first time, it is automatically registered.
Thus, it is not necessary to register all required components during initialization.

## Creating entities

The most basic way to create an entity is {{< api ecs World.NewEntity >}}:

{{< code-func entities_test.go TestEntitiesCreate >}}

Here, we get an entity without any components.
However, {{< api ecs World.NewEntity NewEntity >}} takes an arbitrary number of components IDs for the components that should be associated with the entity:

{{< code-func entities_test.go TestEntitiesCreateComponents >}}

We get an entity with `Position`, and another one with `Position` and `Heading`.
In this case, the components are initialized with their zero value.
Alternatively, entities can be created with initialized components through {{< api ecs World.NewEntityWith >}}:

{{< code-func entities_test.go TestEntitiesCreateWithComponents >}}

We get an entity with `Position` and `Heading`, initialized according to the passed pointers.

### Generic API

Creating entities using the [generic API](./apis) requires a generic *MapX*, like {{< api generic Map2 >}}:

{{< code-func entities_test.go TestEntitiesCreateGeneric >}}

We get an entity with `Position` and `Heading`, initialized to their zero values.

Equivalent to {{< api ecs World.NewEntityWith >}}, generic MapX's have {{< api generic Map2.NewWith NewWith >}}:

{{< code-func entities_test.go TestEntitiesCreateWithComponentsGeneric >}}

{{% notice style="blue" icon="lightbulb" title="Note" %}}
The `2` in `Map2` stands for the number of components.
In the generic API, there are also `FilterX` and `QueryX`.
All these types are available for `X` in range 0 (or 1) to 12.
{{% /notice %}}

### Batch Creation

For faster batch creation of many entities, see chapter [Batch Operations](./batch-ops).

## Adding and removing components

Components are added to and removed from entities through with world,
with {{< api ecs World.Add >}} and {{< api ecs World.Remove >}}.
With generics, we use a {{< api generic Map2 >}} again:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func entities_test.go TestEntitiesAddRemoveGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func entities_test.go TestEntitiesAddRemove >}}
{{< /tab >}}
{{< /tabs >}}

First, we add `Position` and `Heading` to the entity, then we remove both.

It is also possible to assign initialized components with {{< api ecs World.Assign >}}/{{< api generic Map2.Assign >}},
similar to {{< api ecs World.NewWith >}}:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func entities_test.go TestEntitiesAssignGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func entities_test.go TestEntitiesAssign >}}
{{< /tab >}}
{{< /tabs >}}

## Removing entities

Entities can be removed from the world with {{< api ecs World.RemoveEntity >}}:

{{< code-func entities_test.go TestEntitiesRemove >}}

After removal, the entity will be recycled.
For that sake, each entity has a generation variable which allows to distinguish recycled entities.
With {{< api ecs World.Alive >}}, it can be tested whether an entity is still alive:

{{< code-func entities_test.go TestEntitiesAlive >}}
