+++
title = 'Filters'
weight = 60
description = "Arche's filter APIs."
+++
*Filters* provide the logic for filtering entities in [Queries](./queries).

Due to the [archetype](/background/architecture#archetypes)-based architecture of Arche :wink:, filters are very efficient.
Instead of against every single entity, they are only matched against archetypes.

The following sections present the filtering options available in Arche.

{{% notice style="blue" icon="exclamation" title="Important" %}}
Filters should be stored and re-used where possible, particularly over time steps.
Contrary, [Queries](./queries) are for one-time utilization and must be created
from a filter before every iteration loop.
{{% /notice %}}

{{< toc >}}

## Core filters

### Mask

The most common filter is a simple {{< api ecs Mask >}}, which is usually generated with the function {{< api ecs All >}}:

{{< code-func filters_test.go TestMask >}}

Simple {{< api ecs Mask >}} filters match all entities that have at least all the specified components.
The generic equivalent is a simple *FilterX*, e.g. {{< api generic Filter2 >}}:

{{< code-func filters_test.go TestMaskGeneric >}}

In both examples, we filter for all entities that have `Position` and `Heading`,
and anything else that we are not interested in.

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

With {{< api ecs Mask.Exclusive >}} and {{< api generic Filter2.Exclusive >}} {{% dev "v0.11" %}},
we can exclude all components that are not in the filter:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func filters_test.go TestMaskExclusiveGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func filters_test.go TestMaskExclusive >}}
{{< /tab >}}
{{< /tabs >}}

I.e., we get only entities with exactly the given components, and no more.

### With & Optional

With the ID-based API, queries allow access to any component, irrespective of whether it was included in the query.
Generic queries, however, can access only the queried components.
Therefore, generic filters can have optional components through {{< api generic Filter2.Optional >}}:

{{< code-func filters_test.go TestGenericOptional >}}

Note that the now optional `Heading` must be specified also in the original filter.
In case an optional component is not present, `Get` returns `nil` for it.

Further, generic filters have {{< api generic Filter2.With >}}.
This requires the respective component(s) to be present, but they are not obtained through `Get`:

{{< code-func filters_test.go TestGenericWith >}}

### Relation filters

Filters for [Entity Relations](./relations) are covered in the respective chapter.

## Logic filters

Package {{< api filter >}} provides logic combinations of filters.
Logic filters can only be used with the ID-based API.
Here are some examples:

{{< code-func filters_test.go TestLogicFilters >}}

## Filter caching

Normally, when iterating a [Query](./queries), the underlying filter is evaluated on each [archetype](/background/architecture#archetypes).
With a high number of archetypes in the world, this can slow down query iteration and other query functions.

To prevent this slowdown, filters can be registered to the {{< api ecs World.Cache >}} via
{{< api ecs Cache.Register >}}. For generic filters, there is {{< api generic Filter2.Register >}}:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func filters_test.go TestRegisterGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func filters_test.go TestRegister >}}
{{< /tab >}}
{{< /tabs >}}

For registered filters, the list of matching archetypes is cached internally.
Thus, no filter evaluations are required during iteration.
Instead, filters are only evaluated when a new archetype is created.

When a registered filter is not required anymore, it can be unregistered with
{{< api ecs Cache.Unregister >}} or {{< api generic Filter2.Unregister >}}, respectively.
However, this is rarely required as (registered) filters are usually used over an entire simulation run.
