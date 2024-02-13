+++
title = 'Generic vs. ID-based API'
weight = 20
description = "Overview of Arche's generic and ID-based APIs."
+++

Arche provides two different APIs:

A **generic API** that is often the most convenient, and it is type safe.
It is the recommended way of usage for most users.

An **ID-based API** that is slightly faster than the generic one in some places.
Further, it is more flexible and may be more appropriate for tasks like automated serialization.

Both APIs can be mixed as needed.




{{< notice style="blue" icon="lightbulb" title="Tip" >}}
In the this user guide, most code examples will be presented with two tabs, one for each API:

{{< tabs >}}
{{< tab title="generic" >}}
{{< code-func api_test.go TestGeneric >}}
{{< /tab >}}
{{< tab title="ID-based" >}}
{{< code-func api_test.go TestIDs >}}
{{< /tab >}}
{{< /tabs >}}

{{< /notice >}}

