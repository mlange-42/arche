// Package filter contains Arche's advanced logic filtering API.
//
// See the top level module [github.com/mlange-42/arche] for an overview.
//
// 🕮 Also read the Arche's [User Guide]!
//
// # Outline
//
//   - [All] creates a basic filter for components.
//   - [ANY] filters for one of multiple possible components.
//   - [NoneOF] excludes components.
//   - [AnyNOT] matches missing components.
//   - [AND], [OR], [XOR] logically combine two filters.
//   - [NOT] inverts any other filter.
//
// All filters that wrap other filters ([AND], [OR], [XOR], [NOT]) ignore potential relation targets
// of any wrapped ecs.RelationFilter (see [github.com/mlange-42/arche/ecs.RelationFilter]).
//
// [User Guide]: https://mlange-42.github.io/arche/
package filter
