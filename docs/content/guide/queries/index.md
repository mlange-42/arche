+++
title = 'Queries'
weight = 50
description = "Usage of Arche's queries."
+++

Queries ({{< api ecs Query >}}) are the heart of Arche's query engine.
They allow for very fast retrieval and iteration of entities with certain components.

{{% notice style="blue" icon="exclamation" title="Important" %}}
Queries are for one-time utilization.
A query can be iterated through only once, and a new one must be created before every loop.
Contrary, the underlying [Filters](./filters) should be stored and re-used.
{{% /notice %}}

{{< toc >}}

## Query creation & iteration

Queries are created through the {{< api ecs World >}} using a *Filter* (interface {{< api ecs Filter >}}).
The most basic type of filter is {{< api ecs Mask >}}. For more advanced filters, see chapter [Filters](./filters).

Here, we create a filter that gives us all entities with all the given components, and potentially further components. Then, we create an {{< api ecs Query >}} (or generic *QueryX*, e.g. {{< api generic Query2 >}}) and iterate it.

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func queries_test.go TestQueryIterateGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func queries_test.go TestQueryIterate >}}
{{< /tab >}}
{{< /tabs >}}

Where {{< api ecs Query.Get >}} (resp. {{< api generic Query2.Get >}}) return components of the entity at the
current query iterator position.

Comparing the two versions of the code above, one can clearly observe the advantages of the generic API
over the ID-based API (see chapter on [APIs](./apis)).
Firstly, the generic code is shorter and more readable.
But even more importantly, it much safer.
A little mistake in line 9 or 10 of the ID-based version could result in silently casting a component
to the wrong type, which would lead to bugs that are hard to track down.

{{% notice style="blue" icon="bug" title="Tip" %}}
If you get error messages like "index out of range [-1]" or "invalid memory address or nil pointer dereference" from queries, you are probably using them in the wrong way. Try running with build tag `debug` for more helpful error messages:

```
go run -tags debug .
```
{{% /notice %}}

## World lock

When a query gets created, the {{< api ecs World >}} gets locked for modifications.
When locked, no entities can be created or removed, and no components can be added to
or removed from entities.

When a query is fully iterated, the world gets unlocked again. When a query is not fully iterated
for some reason (see next section for examples), it must be closed with {{< api ecs Query.Close >}}.

Due to the world lock, denied operations like entity creation or removal must be deferred:

{{< code-func queries_test.go TestQueryRemoveEntities >}}

Where {{< api ecs Query.Entity >}} returns the entity at the current query iterator position.

## Other functionality

Besides {{< api ecs Query.Next >}}, {{< api ecs Query.Get >}} and {{< api ecs Query.Entity >}}
that we used above, queries have a few more useful methods.

### Query.Count

{{< api ecs Query.Count >}} allows for counting the entities in a query, very fast:

{{< code-func queries_test.go TestQueryCount >}}

Note that we need to call {{< api ecs Query.Close >}} here, as the query was not (fully) iterated!
After {{< api ecs Query.Count >}}, the query could also be iterated as usual.

### Query.EntityAt

{{% dev "v0.11" %}}

With {{< api ecs Query.EntityAt >}}, queries also support access by index.
This is particularly useful to select random entities from a query,
like in this example:

{{< code-func queries_test.go TestQueryEntityAt >}}

Note that we need to close the query manually, again!
To access components of the retrieved entities, see chapter [World Entity Access](./world-access).

Note that {{< api ecs Query.EntityAt >}} may be slow when working with a large number of archetypes.
Often, it is useful to register the underlying filter for speedup.
See chapter [Filter](./filters), section [Filter caching](./filters#filter-caching) for details.
See the [query benchmarks](/background/benchmarks#query) for some numbers on performance.
