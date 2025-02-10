Last run: Mon, 06 Jan 2025 23:23:47 CET  
Version: Arche v0.14.6-dev  
CPU: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz

## Query

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Query.Next                       |       4.0 ns |                              |
| Query.Next + 1x Query.Get        |       6.2 ns |                              |
| Query.Next + 2x Query.Get        |       8.0 ns |                              |
| Query.Next + 5x Query.Get        |      16.4 ns |                              |
| Query.Next + Query.Entity        |       4.4 ns |                              |
| Query.Next + Query.Relation      |       6.2 ns |                              |
| Query.EntityAt, 1 arch           |      36.5 ns |                              |
| Query.EntityAt, 1 arch           |      12.0 ns | registered filter            |
| Query.EntityAt, 5 arch           |     105.9 ns |                              |
| Query.EntityAt, 5 arch           |      21.1 ns | registered filter            |
| World.Query                      |      91.3 ns |                              |
| World.Query                      |      87.0 ns | registered filter            |

## World access

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.Get                        |       9.5 ns | random, 1000 entities        |
| World.GetUnchecked               |       4.7 ns | random, 1000 entities        |
| World.Has                        |       2.4 ns | random, 1000 entities        |
| World.HasUnchecked               |       1.4 ns | random, 1000 entities        |
| World.Alive                      |       1.1 ns | random, 1000 entities        |
| World.Relations.Get              |       9.1 ns | random, 1000 entities        |
| World.Relations.GetUnchecked     |       2.6 ns | random, 1000 entities        |

## Entities

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Entity.IsZero                    |       0.3 ns |                              |
| World.NewEntity                  |      66.8 ns | memory already alloc.     |
| World.NewEntity w/ 1 Comp        |     121.6 ns | memory already alloc.     |
| World.NewEntity w/ 5 Comps       |     162.5 ns | memory already alloc.     |
| World.RemoveEntity               |      59.1 ns |                              |
| World.RemoveEntity w/ 1 Comp     |      79.6 ns |                              |
| World.RemoveEntity w/ 5 Comps    |     141.1 ns |                              |
| Map1.NewWith w/ 1 Comp           |     102.9 ns | memory already alloc.     |
| Map5.NewWith w/ 5 Comps          |     243.4 ns | memory already alloc.     |

## Entities, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Builder.NewBatch                 |      33.5 ns | 1000, memory already alloc. |
| Builder.NewBatch w/ 1 Comp       |      36.4 ns | 1000, memory already alloc. |
| Builder.NewBatch w/ 5 Comps      |      36.2 ns | 1000, memory already alloc. |
| Batch.RemoveEntities             |      23.4 ns | 1000                         |
| Batch.RemoveEntities w/ 1 Comp   |      23.4 ns | 1000                         |
| Batch.RemoveEntities w/ 5 Comps  |      27.7 ns | 1000                         |

## Components

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| World.Add 1 Comp                 |     174.4 ns | memory already alloc.     |
| World.Add 5 Comps                |     301.4 ns | memory already alloc.     |
| World.Remove 1 Comp              |     178.7 ns |                              |
| World.Remove 5 Comps             |     340.1 ns |                              |
| World.Exchange 1 Comp            |     167.8 ns | memory already alloc.     |
| Map1.Assign 1 Comps              |     178.7 ns | memory already alloc.     |
| Map5.Assign 5 Comps              |     441.2 ns | memory already alloc.     |

## Components, batched

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| Batch.Add 1 Comp                 |      29.9 ns | 1000, memory already alloc. |
| Batch.Add 5 Comps                |      28.5 ns | 1000, memory already alloc. |
| Batch.Remove 1 Comp              |      39.6 ns | 1000                         |
| Batch.Remove 5 Comps             |      65.4 ns | 1000                         |
| Batch.Exchange 1 Comp            |      34.8 ns | 1000, memory already alloc. |

## Other

| Operation                        | Time         | Remark                       |
|----------------------------------|-------------:|------------------------------|
| ecs.NewWorld                     |      40.1 μs |                              |
| World.Reset                      |     251.6 ns |                              |
| ecs.ComponentID                  |      65.6 ns | registered component         |

