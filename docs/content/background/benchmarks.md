+++
title = 'Benchmarks'
weight = 60
description = 'An overview of the runtime cost of typical Arche operations.'
+++

This document gives an overview of the runtime cost of typical Arche operations.
All time information is per entity.
Batch operations are performed in batches of 1000 entities.

Absolute numbers are not  really meaningful, as they heavily depend on the hardware.
However, all benchmarks run in the CI in the same job and hence on the same machine, and can be compared.

Benchmark code: {{< repo "/tree/main/benchmark/table" "benchmark/table" >}} in the {{< repo "" "GitHub repository" >}}.

{{% include file="background/_benchmarks.md" %}}
