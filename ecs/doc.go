// Package ecs contains Arche's core API.
//
// See the top-level module [github.com/mlange-42/arche] for an overview.
//
// 🕮 Also read Arche's [User Guide]!
//
// # Outline
//
//   - [World] provides most of the basic functionality,
//     like [World.Query], [World.NewEntity], [World.Add], [World.Remove], [World.RemoveEntity], etc.
//   - [Relations] provide access to and manipulation of entity relations,
//     like [Relations.Get] and [Relations.Set].
//   - [Builder] provides advanced entity creation and batched creation with
//     [Builder.NewBatch] and [Builder.NewBatchQ].
//   - [Batch] provides batch-manipulation of entities,
//     like [Batch.Add], [Batch.Remove] and [Batch.SetRelation].
//   - [Cache] serves for registering and un-registering cached filters
//     with [Cache.Register] and [Cache.Unregister].
//   - [Resources] provide a storage for global resources, with functionality like
//     [Resources.Get], [Resources.Add] and [Resources.Remove].
//   - [Listener] provides [EntityEvent] notifications for ECS operations.
//   - Useful functions: [All], [ComponentID], [ResourceID], [GetResource], [AddResource].
//
// # Sub-packages
//   - [github.com/mlange-42/arche/ecs/event] provides event subscription masks.
//   - [github.com/mlange-42/arche/ecs/stats] provides world statistics for monitoring purposes.
//
// # ECS Manipulations
//
// This section gives an overview on how to achieve typical ECS manipulation operations with Arche.
//
// Simple manipulations of a single entity:
//   - Create an entity: [World.NewEntity], [World.NewEntityWith]
//   - Remove an entity: [World.RemoveEntity]
//   - Add components: [World.Add]
//   - Remove components: [World.Remove]
//   - Exchange components: [World.Exchange]
//   - Change entity relation target: [Relations.Set]
//
// Manipulations of a single entity, with a relation target:
//   - Create an entity: [Builder.New]
//   - Add components: [Builder.Add], [Relations.Exchange]
//   - Remove components: [Relations.Exchange]
//   - Exchange components: [Relations.Exchange]
//
// Batch-manipulations of many entities:
//   - Create entities: [Builder.NewBatch], [Builder.NewBatchQ]
//   - Remove entities: [Batch.RemoveEntities]
//   - Add components: [Batch.Add], [Batch.AddQ]
//   - Remove components: [Batch.Remove], [Batch.RemoveQ]
//   - Exchange components: [Batch.Exchange]
//   - Change entity relation target: [Batch.SetRelation], [Batch.SetRelationQ]
//
// Batch-manipulations of many entities, with a relation target:
//   - Create an entity: [Builder.NewBatch], [Builder.NewBatchQ]
//   - Add components: [Relations.ExchangeBatch], [Relations.ExchangeBatchQ]
//   - Remove components: [Relations.ExchangeBatch], [Relations.ExchangeBatchQ]
//   - Exchange components: [Relations.ExchangeBatch], [Relations.ExchangeBatchQ]
//
// # Build tags
//
// Arche provides two build tags:
//   - tiny -- Reduces the maximum number of components to 64, giving a performance boost for mask-related operations.
//   - debug -- Improves error messages on [Query] misuse, at the cost of performance. Use this if you get panics from queries.
//
// When building your application, use them like this:
//
//	go build -tags tiny .
//	go build -tags debug .
//	go build -tags tiny,debug .
//
// [User Guide]: https://mlange-42.github.io/arche/
package ecs
