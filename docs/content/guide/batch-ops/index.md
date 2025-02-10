+++
title = 'Batch Operations'
type = "docs"
weight = 90
description = 'Batch operations in Arche.'
+++
Compared to [Queries](./queries) and [World Entity Access](./world-access),
creation and removal of entities or components are relatively costly operations.
See the [Benchmarks](/background/benchmarks) for some numbers.

For these operations, Arche provides batched versions.
This allows to create or manipulate a large number of entities much faster than one by one.
Most batch methods come in two flavors. A "normal" one, and one suffixed with `Q` that returns a query over the affected entities.

## Creating entities

Entity creation is probably the most common use case for batching.
When the number of similar entities that are to be created are known,
creation can be batched with {{< api ecs Batch.New >}}.
In the generic API, *MapX* provide e.g. {{< api generic Map2.NewBatch >}}.

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func batch_test.go TestBatchCreateGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func batch_test.go TestBatchCreate >}}
{{< /tab >}}
{{< /tabs >}}

However, this is only sometimes useful, as we can't initialize component fields here.

With the query variant of the methods, suffixed with `Q`, we can fix this:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func batch_test.go TestBatchCreateQueryGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func batch_test.go TestBatchCreateQuery >}}
{{< /tab >}}
{{< /tabs >}}

Here, we obtain a query over exactly the entities we just created, and can initialize their components.

## Components

Components can be added, removed or exchanged in batch operations.
For these operations, Arche provides {{< api ecs World.Batch >}}.
Component batch operations take an {{< api ecs Filter >}} as an argument to determine the affected entities.

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func batch_test.go TestBatchAddQueryGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func batch_test.go TestBatchAddQuery >}}
{{< /tab >}}
{{< /tabs >}}

Methods of interest for the ID-based API are:
 - {{< api ecs Batch.Add >}}, {{< api ecs Batch.AddQ >}}
 - {{< api ecs Batch.Remove >}}, {{< api ecs Batch.RemoveQ >}}
 - {{< api ecs Batch.Exchange >}}, {{< api ecs Batch.ExchangeQ >}}

Methods of interest for the generic API are:
 - {{< api generic Map2.AddBatch >}}, {{< api generic Map2.AddBatchQ >}}
 - {{< api generic Map2.RemoveBatch >}}, {{< api generic Map2.RemoveBatchQ >}}
 - {{< api generic Exchange.ExchangeBatch >}}

## Relations

Entity relations can be changed in batches, too.
In the ID-based API, both {{< api ecs Batch.SetRelation >}}/{{< api ecs Batch.SetRelationQ >}}
and {{< api ecs Relations.SetBatch >}}/{{< api ecs Relations.SetBatchQ >}} can be used.
In the generic API, use {{< api generic Map.SetRelation >}}/{{< api generic Map.SetRelationQ >}}:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func batch_test.go TestBatchRelationsGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func batch_test.go TestBatchRelations >}}
{{< /tab >}}
{{< /tabs >}}

## Removing entities

Entities can be removed in batches using {{< api ecs Batch.RemoveEntities >}}:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func batch_test.go TestBatchRemoveEntitiesGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func batch_test.go TestBatchRemoveEntities >}}
{{< /tab >}}
{{< /tabs >}}
