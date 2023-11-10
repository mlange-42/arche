// Package ecs contains Arche's core API.
//
// See the top-level module [github.com/mlange-42/arche] for an overview.
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
//   - Useful functions: [All], [ComponentID], [ResourceID], [GetResource], [AddResource].
package ecs
