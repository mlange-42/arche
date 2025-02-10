+++
title = 'Generic & ID-based API'
type = "docs"
weight = 20
description = "Overview of Arche's generic and ID-based APIs."

[params]
prev = "/guide"
+++

Arche provides two different APIs:

A **generic API** that is often the most convenient. But perhaps more importantly, it is type safe.
It is the recommended way of usage for most users.

An **ID-based API** that is slightly faster than the generic one in some places.
Further, it is more flexible and may be more appropriate for tasks like automated serialization.

Both APIs can be mixed as needed.

> [!TIP]
> In this user guide, most code examples will be presented with two tabs, one for each API:

{{< tabs items="generic,ID-based" >}}
{{< tab >}}
{{< code-func api_test.go TestGeneric >}}
{{< /tab >}}
{{< tab >}}
{{< code-func api_test.go TestIDs >}}
{{< /tab >}}
{{< /tabs >}}
