+++
title = 'Queries'
weight = 20
+++
## Code tests

Code inclusion test for a complete file: `{{</* code code_test.go */>}}`

{{< code code_test.go >}}

Code inclusion test for a function `TestCode`: `{{</* code-func code_test.go TestCode */>}}`

{{< code-func code_test.go TestCode >}}

Code inclusion of lines 1-6: `{{</* code-lines code_test.go 1 6 */>}}`

{{< code-lines code_test.go 1 6 >}}

A [Link](https://example.com)

An API docs link: {{< api ecs Entity.IsZero >}}, {{< api ecs Entity.IsZero "Entity.IsZero()" >}}
