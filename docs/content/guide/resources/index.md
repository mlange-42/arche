+++
title = 'ECS Resources'
type = "docs"
weight = 100
description = 'ECS resources in Arche.'
+++
So far, we stored all data in components, associated to entities.
However, sometimes this is not optimal.
This particularly applies to non-ECS data structures,
and "things" that appear only once in a model or game.

For these cases, Arche provides so-called {{< api ecs Resources >}}.
A resource can be thought of as a component that only appears once, globally,
in an {{< api ecs World >}}.

## Resource types

Just like with components, any Go struct (or other Go type) can be a resource.
An example:

```go
type Grid struct {
    Data [][]ecs.Entity
    Width int
    Height int
}
```

## Resource IDs

Also analogous to components, resources are identified by a unique ID that can be obtained or registered using {{< api ecs ResourceID >}}:

{{< code-func resources_test.go TestResourceID >}}

As with components, a resource is registered and assigned an ID on first use automatically.

## Adding resources

The most simple way to add a resource to the world is the function {{< api ecs AddResource >}}:

{{< code-func resources_test.go TestResourceAdd >}}

An ID is automatically assigned to type `Grid` here, if it was not registered before with {{< api ecs ResourceID >}}.

{{< api ecs AddResource >}}, however, is not particularly efficient.
If a resource needs to be added (and removed) repeatedly, use {{< api ecs World.Resources >}}:

{{< code-func resources_test.go TestResourceAdd2 >}}

## Accessing resources

Access to resources is obtained via {{< api ecs World.Resources >}} in the ID-based API,
and via {{< api generic Resource >}} in the generic API:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func resources_test.go TestResourceGetGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func resources_test.go TestResourceGet >}}
{{< /tab >}}
{{< /tabs >}}

Note that in the ID-based example, we need to cast the pointer retrieved from
{{< api ecs Resources.Get >}} to `*Grid`, similar to the cast in ID-based component access.
However, the syntax is a bit different here as we cast an `interface{}`,
rather than an `unsafe.Pointer` for components.

As with components, the generic API is the recommended way for normal usage.
