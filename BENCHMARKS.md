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
| Query.Next + 1x Query.Get        |       1.6 ns |                              |
| Query.Next + 2x Query.Get        |       2.3 ns |                              |
| Query.Next + 5x Query.Get        |       4.4 ns |                              |
| Query.Next + Query.Relation      |       2.2 ns |                              |
| Query.EntityAt, 1 arch           |      12.0 ns |                              |
| Query.EntityAt, 1 arch           |       2.8 ns | registered filter            |
| Query.EntityAt, 5 arch           |      16.1 ns |                              |
| Query.EntityAt, 5 arch           |       2.8 ns | registered filter            |
| World.Query                      |      46.0 ns |                              |
| World.Query                      |      33.6 ns | registered filter            |

## Entities

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.NewEntity                  |      16.6 ns | memory already allocated     |
| World.NewEntity w/ 1 Comp        |      33.6 ns | memory already allocated     |
| World.NewEntity w/ 5 Comps       |      47.0 ns | memory already allocated     |
| World.RemoveEntity               |      15.4 ns |                              |
| World.RemoveEntity w/ 1 Comp     |      25.7 ns |                              |
| World.RemoveEntity w/ 5 Comps    |      54.4 ns |                              |

## Entities, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Builder.NewBatch                 |       9.6 ns | 1000, memory already allocated |
| Builder.NewBatch w/ 1 Comp       |       9.8 ns | 1000, memory already allocated |
| Builder.NewBatch w/ 5 Comps      |       9.6 ns | 1000, memory already allocated |
| Batch.RemoveEntities             |       6.9 ns | 1000                         |
| Batch.RemoveEntities w/ 1 Comp   |       7.1 ns | 1000                         |
| Batch.RemoveEntities w/ 5 Comps  |       7.9 ns | 1000                         |

## Components

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.Add 1 Comp                 |      50.7 ns | memory already allocated     |
| World.Add 5 Comps                |      67.6 ns | memory already allocated     |
| World.Remove 1 Comp              |      59.0 ns |                              |
| World.Remove 5 Comps             |     101.0 ns |                              |
| World.Exchange 1 Comp            |      60.5 ns | memory already allocated     |

## Components, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Batch.Add 1 Comp                 |       9.0 ns | 1000, memory already allocated |
| Batch.Add 5 Comps                |       9.1 ns | 1000, memory already allocated |
| Batch.Remove 1 Comp              |       9.8 ns | 1000                         |
| Batch.Remove 5 Comps             |      14.9 ns | 1000                         |
| Batch.Exchange 1 Comp            |      10.0 ns | 1000, memory already allocated |

## World access

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.Get                        |       2.0 ns | random, 1000 entities        |
| World.Has                        |       1.2 ns | random, 1000 entities        |
| World.Alive                      |       0.6 ns | random, 1000 entities        |
| World.Relations.Get              |       3.5 ns | random, 1000 entities        |
| World.Relations.GetUnchecked     |       0.8 ns | random, 1000 entities        |
