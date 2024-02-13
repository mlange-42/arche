+++
title = 'Benchmarks'
weight = 50
+++

This document gives an overview of the runtime cost of typical Arche operations.
All time information is per entity.
Batch operations are performed in batches of 1000 entities.

Absolute numbers are not  really meaningful, as they heavily depend on the hardware.
However, all benchmarks run in the CI in the same job and hence on the same machine, and can be compared.

Benchmark code: [`benchmark/table`](https://github.com/mlange-42/arche/tree/main/benchmark/table).

## Query

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Query.Next                       |       1.0 ns |                              |
| Query.Next + 1x Query.Get        |       1.6 ns |                              |
| Query.Next + 2x Query.Get        |       2.2 ns |                              |
| Query.Next + 5x Query.Get        |       4.4 ns |                              |
| Query.Next + Query.Relation      |       2.3 ns |                              |
| Query.EntityAt, 1 arch           |      12.0 ns |                              |
| Query.EntityAt, 1 arch           |       2.8 ns | registered filter            |
| Query.EntityAt, 5 arch           |      31.6 ns |                              |
| Query.EntityAt, 5 arch           |       4.9 ns | registered filter            |
| World.Query                      |      45.1 ns |                              |
| World.Query                      |      33.4 ns | registered filter            |

## World access

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.Get                        |       2.0 ns | random, 1000 entities        |
| World.Has                        |       1.2 ns | random, 1000 entities        |
| World.Alive                      |       0.6 ns | random, 1000 entities        |
| World.Relations.Get              |       3.5 ns | random, 1000 entities        |
| World.Relations.GetUnchecked     |       0.8 ns | random, 1000 entities        |

## Entities

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.NewEntity                  |      16.2 ns | memory already allocated     |
| World.NewEntity w/ 1 Comp        |      34.0 ns | memory already allocated     |
| World.NewEntity w/ 5 Comps       |      45.1 ns | memory already allocated     |
| World.RemoveEntity               |      14.7 ns |                              |
| World.RemoveEntity w/ 1 Comp     |      27.0 ns |                              |
| World.RemoveEntity w/ 5 Comps    |      53.5 ns |                              |

## Entities, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Builder.NewBatch                 |       9.7 ns | 1000, memory already allocated |
| Builder.NewBatch w/ 1 Comp       |      10.1 ns | 1000, memory already allocated |
| Builder.NewBatch w/ 5 Comps      |      10.2 ns | 1000, memory already allocated |
| Batch.RemoveEntities             |       7.0 ns | 1000                         |
| Batch.RemoveEntities w/ 1 Comp   |       7.2 ns | 1000                         |
| Batch.RemoveEntities w/ 5 Comps  |       7.6 ns | 1000                         |

## Components

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.Add 1 Comp                 |      48.2 ns | memory already allocated     |
| World.Add 5 Comps                |      66.9 ns | memory already allocated     |
| World.Remove 1 Comp              |      58.1 ns |                              |
| World.Remove 5 Comps             |     103.6 ns |                              |
| World.Exchange 1 Comp            |      55.5 ns | memory already allocated     |

## Components, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Batch.Add 1 Comp                 |       8.4 ns | 1000, memory already allocated |
| Batch.Add 5 Comps                |       8.9 ns | 1000, memory already allocated |
| Batch.Remove 1 Comp              |      10.3 ns | 1000                         |
| Batch.Remove 5 Comps             |      15.3 ns | 1000                         |
| Batch.Exchange 1 Comp            |      10.0 ns | 1000, memory already allocated |
