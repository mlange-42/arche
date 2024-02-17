+++
title = 'Benchmarks'
weight = 60
description = 'An overview of the runtime cost of typical Arche operations.'
+++

This chapter gives an overview of the runtime cost of typical Arche operations.
All time information is per entity.
Batch operations are performed in batches of 1000 entities.

Absolute numbers are not  really meaningful, as they heavily depend on the hardware.
However, all benchmarks run in the CI in the same job and hence on the same machine, and can be compared.

Benchmark code: {{< repo "/tree/main/benchmark/table" "benchmark/table" >}} in the {{< repo "" "GitHub repository" >}}.

Benchmarks are run automatically in the GitHub CI, and are updated on this page on every merge into the `main` branch.
They always reflect the latest development state of Arche, but may change depending on the hardware the CI runs on.
For a benchmark comparison with other ECS implementations, see the benchmarks on the [repo README](https://github.com/mlange-42/arche#benchmarks).

{{< toc >}}

{{% include file="/background/_benchmarks.md" %}}
