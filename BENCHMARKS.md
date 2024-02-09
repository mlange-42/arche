# Benchmarks

This document gives an overview of the runtime cost of typical Arche operations.
All time information is per entity.
Batch operations are performed in batches of 1000 entities.

Absolute numbers are not  really meaningful, as they heavily depend on the hardware.
However, all benchmarks run in the CI in the same job and hence on the same machine, and can be compared.

## Query

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Query.Next                       |       1.0 ns |                              |
| Query.Next + 1x Get              |       1.6 ns |                              |
| Query.Next + 2x Get              |       2.3 ns |                              |
| Query.Next + 5x Get              |       4.4 ns |                              |
| Query.EntityAt, 1 arch           |      12.0 ns |                              |
| Query.EntityAt, 1 arch           |       3.1 ns | registered filter            |
| Query.EntityAt, 5 arch           |      14.9 ns |                              |
| Query.EntityAt, 5 arch           |       2.8 ns | registered filter            |
| World.Query                      |      46.9 ns |                              |
| World.Query                      |      34.0 ns | registered filter            |

## Entities

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.NewEntity                  |      16.3 ns | memory already allocated     |
| World.NewEntity w/ 1 Comp        |      34.6 ns | memory already allocated     |
| World.NewEntity w/ 5 Comps       |      46.7 ns | memory already allocated     |
| World.RemoveEntity               |      15.1 ns |                              |
| World.RemoveEntity w/ 1 Comp     |      25.3 ns |                              |
| World.RemoveEntity w/ 5 Comps    |      54.8 ns |                              |

## Entities, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Builder.NewBatch                 |      10.3 ns | 1000, memory already allocated |
| Builder.NewBatch w/ 1 Comp       |       9.8 ns | 1000, memory already allocated |
| Builder.NewBatch w/ 5 Comps      |      10.6 ns | 1000, memory already allocated |
| Batch.RemoveEntities             |       7.0 ns | 1000                         |
| Batch.RemoveEntities w/ 1 Comp   |       8.1 ns | 1000                         |
| Batch.RemoveEntities w/ 5 Comps  |       8.1 ns | 1000                         |

## Components

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.Add 1 Comp                 |      47.3 ns | memory already allocated     |
| World.Add 5 Comps                |      65.0 ns | memory already allocated     |
| World.Remove 1 Comp              |      58.9 ns |                              |
| World.Remove 5 Comps             |     101.2 ns |                              |
| World.Exchange 1 Comp            |      54.7 ns | memory already allocated     |

## Components, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Batch.Add 1 Comp                 |       9.3 ns | 1000, memory already allocated |
| Batch.Add 5 Comps                |       8.9 ns | 1000, memory already allocated |
| Batch.Remove 1 Comp              |      11.3 ns | 1000                         |
| Batch.Remove 5 Comps             |      17.3 ns | 1000                         |
| Batch.Exchange 1 Comp            |      10.5 ns | 1000, memory already allocated |
