+++
title = 'World Entity Access'
type = "docs"
weight = 70
description = 'Access to components through the world, by entity ID.'
+++
So far, we only used {{< api ecs Query >}} to access components.
Component access in queries is highly efficient, but it does not provide
access to the components of a specific entity.
This is possible through {{< api ecs World >}} methods, or using a generic *MapX* (like {{< api generic Map2 >}}) or {{< api generic Map >}}.

## Getting components

For a given entity, components can be accessed using {{< api ecs World.Get >}}
or {{< api generic Map2.Get >}}, respectively:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func access_test.go TestGetGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func access_test.go TestGet >}}
{{< /tab >}}
{{< /tabs >}}

Similarly, it is also possible to check if an entity has a given component with
{{< api ecs World.Has >}} or {{< api generic Map.Has >}}, respectively:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func access_test.go TestHasGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func access_test.go TestHas >}}
{{< /tab >}}
{{< /tabs >}}

Note that we have to use {{< api generic Map >}} here, which is similar to
{{< api generic Map1 >}} for a single component, but offers more functionality.

> [!IMPORTANT]
> Note that the component pointers obtained here should never be stored persistently.

## Unchecked access

The `Get` and `Has` methods shown above all have a companion `GetUnchecked` and `HasUnchecked`,
which is faster, but should be used with care.
Particularly, they do not check whether the entity is still alive.
Like `Get`, they panic when called on a removed entity.
However, for a removed and subsequently recycled entity, they lead to undefined behavior.

It is safe to use methods like {{< api ecs World.GetUnchecked >}} after a usual `Get` was used on the same entity:

{{< code-func access_test.go TestGetUnchecked >}}

Note that, following this use case, generic *MapX* internally use
{{< api ecs World.GetUnchecked >}} for all but the first component.
