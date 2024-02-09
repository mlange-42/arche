# Benchmarks

This document gives an overview over the runtime cost of typical Arche operations.
If not stated differently, time is per entity or per single operation.

## Query

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Query.Next                       |       1.6 ns |                              |
| Query.Next + 1x Get              |       1.9 ns |                              |
| Query.Next + 2x Get              |       2.6 ns |                              |
| Query.Next + 5x Get              |       4.8 ns |                              |
| Query.EntityAt, 1 arch           |      11.3 ns |                              |
| Query.EntityAt, 1 arch           |       3.0 ns | registered filter            |
| Query.EntityAt, 5 arch           |      15.9 ns |                              |
| Query.EntityAt, 5 arch           |       3.1 ns | registered filter            |
| World.Query                      |      50.7 ns |                              |
| World.Query                      |      30.6 ns | registered filter            |

## Entities

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.NewEntity                  |      16.5 ns | memory already allocated     |
| World.NewEntity w/ 1 Comp        |      33.1 ns | memory already allocated     |
| World.NewEntity w/ 5 Comps       |      58.5 ns | memory already allocated     |
| World.RemoveEntity               |      13.2 ns |                              |
| World.RemoveEntity w/ 1 Comp     |      20.0 ns |                              |
| World.RemoveEntity w/ 5 Comps    |      44.0 ns |                              |

## Entities, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Builder.NewBatch                 |       8.7 ns | 1/1000, memory already allocated |
| Builder.NewBatch w/ 1 Comp       |       9.3 ns | 1000, memory already allocated |
| Builder.NewBatch w/ 5 Comps      |      10.4 ns | 1000, memory already allocated |
| Batch.RemoveEntities             |       6.6 ns | 1/1000                       |
| Batch.RemoveEntities w/ 1 Comp   |       5.8 ns | 1000                         |
| Batch.RemoveEntities w/ 5 Comps  |       6.6 ns | 1000                         |
