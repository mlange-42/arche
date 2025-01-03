// Package generic provides the generic API of Arche, an Entity Component System (ECS) for Go.
//
// See the top level module [github.com/mlange-42/arche] for an overview.
//
// 🕮 Also read Arche's [User Guide]!
//
// # Outline
//
//   - [Filter0], [Filter1], etc. provide generic filters and query generation using [Filter0.Query] and friends.
//   - [Query0], [Query1], provide the usual [ecs.Query] functionality,
//     as well as generic [Query1.Get], [Query1.Relation], etc.
//   - [Map] provides generic access to a single component using world access, like [Map.Get] and [Map.Set].
//   - [Map1], [Map2], etc. provide generic access to multiple components using world access,
//     like [Map1.Get], [Map1.Add], [Map1.Remove], etc.
//   - [Exchange] allows to add, remove and exchange components, incl. as batch operations.
//   - [Resource] provides generic access to a resource from [ecs.Resources].
//
// # ECS Manipulations
//
// This section gives an overview on how to achieve typical ECS manipulation operations using Arche's generic API.
// MapX and QueryX variants are shown as [Map2] and [Query2] here.
//
// Manipulations of a single entity, with or without a relation target:
//   - Create an entity: [Map2.New], [Map2.NewWith]
//   - Remove an entity: use ecs.World.RemoveEntity
//   - Add components: [Map2.Add], [Exchange.Add]
//   - Remove components: [Map2.Remove], [Exchange.Remove]
//   - Exchange components: [Exchange.Exchange]
//   - Change entity relation target: [Map.SetRelation]
//
// Batch-manipulations of many entities, with or without a relation target:
//   - Create entities: [Map2.NewBatch], [Map2.NewBatchQ]
//   - Remove entities: [Map2.RemoveEntities]
//   - Add components: [Map2.AddBatch], [Map2.AddBatchQ]
//   - Remove components: [Map2.RemoveBatch], [Map2.RemoveBatchQ]
//   - Exchange components: [Exchange.ExchangeBatch]
//   - Change entity relation target: [Map.SetRelationBatch], [Map.SetRelationBatchQ]
//
// [User Guide]: https://mlange-42.github.io/arche/
package generic
