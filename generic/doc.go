// Package generic contains Arche's generic API.
//
// See the top level module [github.com/mlange-42/arche] for an overview.
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
package generic
