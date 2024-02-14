+++
title = 'Benchmarks'
weight = 50
description = 'An overview of the runtime cost of typical Arche operations.'
+++

This document gives an overview of the runtime cost of typical Arche operations.
All time information is per entity.
Batch operations are performed in batches of 1000 entities.

Absolute numbers are not  really meaningful, as they heavily depend on the hardware.
However, all benchmarks run in the CI in the same job and hence on the same machine, and can be compared.

Benchmark code: [`benchmark/table`](https://github.com/mlange-42/arche/tree/main/benchmark/table).

{{% include file="background/_benchmarks.md" %}}
