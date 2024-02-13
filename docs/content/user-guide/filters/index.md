+++
title = 'Filters'
weight = 60
description = "Arche's filter APIs."
+++
*Filters* provide the logic for filtering entities for [Queries](./queries).

Due to the [archetype](/background/architecture/archetypes)-based architecture of Arche :wink:, filters are very efficient.
Instead of against every single entity, they are only matched against archetypes.

The following sections present the filtering options available in Arche.

## Core filters

### Mask

The most common filter is a simple {{< api ecs Mask >}}, which is usually generated with the function {{< api ecs All >}}:

{{< code-func filters_test.go TestMask >}}

Simple {{< api ecs Mask >}} filters match all entities that have at least all the specified components.
The generic equivalent is a simple *FilterX*, e.g. {{< api generic Filter2 >}}:

{{< code-func filters_test.go TestMaskGeneric >}}

In both examples, we filter for all entities that have `Position` and `Heading`,
and anything else we are not interested in.

### Without

Particular components can be excluded with {{< api ecs Mask.Without >}} and {{< api generic Filter2.Without >}}:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func filters_test.go TestMaskWithoutGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func filters_test.go TestMaskWithout >}}
{{< /tab >}}
{{< /tabs >}}

Here, we filter for all entities that have a `Position`, but no `Heading`. Other components are allowed on the entities.

### Exclusive

With {{< api ecs Mask.Exclusive >}} and {{< api generic Filter2.Exclusive >}},
we can exclude all other components:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func filters_test.go TestMaskExclusiveGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func filters_test.go TestMaskExclusive >}}
{{< /tab >}}
{{< /tabs >}}

### Optional

With the ID-based API, queries allow access to any component, irrespective of whether it was included in the query.
Generic queries, however, can access only the queried components.
Therefore, generic filters can have optional components through {{< api generic Filter2.Optional >}}, which takes an arbitrary number of arguments:

{{< code-func filters_test.go TestGenericOptional >}}

Note that the now optional `Heading` must be specified also in the original filter.
In case an optional component is not present, `Get` returns `nil` for it.

## Relation filters

Filters for [Entity Relations](./relations) are covered in the respective chapter.

## Logic filters

## Filter caching
