# Benchmarks

This document gives an overview of the runtime cost of typical Arche operations.
All time information is per entity.
Batch oerations are performed in batches of 1000 entities.

Absolute numbers are not  really meaningful, as they heavily depend on the hardware.
However, all benchmarks run in the CI in the same job and hence on the same machine, and can be compared.

## Query

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Query.Next                       |       1.0 ns |                              |
| Query.Next + 1x Get              |       1.6 ns |                              |
| Query.Next + 2x Get              |       2.3 ns |                              |
| Query.Next + 5x Get              |       4.4 ns |                              |
| Query.EntityAt, 1 arch           |      11.7 ns |                              |
| Query.EntityAt, 1 arch           |       2.8 ns | registered filter            |
| Query.EntityAt, 5 arch           |      14.9 ns |                              |
| Query.EntityAt, 5 arch           |       3.1 ns | registered filter            |
| World.Query                      |      46.3 ns |                              |
| World.Query                      |      34.1 ns | registered filter            |

## Entities

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.NewEntity                  |      16.1 ns | memory already allocated     |
| World.NewEntity w/ 1 Comp        |      34.2 ns | memory already allocated     |
| World.NewEntity w/ 5 Comps       |      59.4 ns | memory already allocated     |
| World.RemoveEntity               |      15.6 ns |                              |
| World.RemoveEntity w/ 1 Comp     |      26.1 ns |                              |
| World.RemoveEntity w/ 5 Comps    |      54.8 ns |                              |

## Entities, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Builder.NewBatch                 |       9.9 ns | 1000, memory already allocated |
| Builder.NewBatch w/ 1 Comp       |      10.0 ns | 1000, memory already allocated |
| Builder.NewBatch w/ 5 Comps      |      10.0 ns | 1000, memory already allocated |
| Batch.RemoveEntities             |       7.5 ns | 1000                         |
| Batch.RemoveEntities w/ 1 Comp   |       7.2 ns | 1000                         |
| Batch.RemoveEntities w/ 5 Comps  |       8.5 ns | 1000                         |

## Components

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.Add 1 Comp                 |      49.4 ns | memory already allocated     |
| World.Add 5 Comps                |      79.0 ns | memory already allocated     |
| World.Remove 1 Comp              |      59.8 ns |                              |
| World.Remove 5 Comps             |     121.8 ns |                              |
| World.Exchange 1 Comp            |      56.4 ns | memory already allocated     |

## Components, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Batch.Add 1 Comp                 |       8.8 ns | 1000, memory already allocated |
| Batch.Add 5 Comps                |       8.9 ns | 1000, memory already allocated |
| Batch.Remove 1 Comp              |      10.3 ns | 1000                         |
| Batch.Remove 5 Comps             |      16.1 ns | 1000                         |
| Batch.Exchange 1 Comp            |      10.3 ns | 1000, memory already allocated |