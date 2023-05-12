// Package filter contains Arche's advanced logic filtering API.
//
// See the top level module [github.com/mlange-42/arche] for an overview.
//
// # Outline
//
//   - [All] creates a basic filter for components.
//   - [ANY] filters for one of multiple possible components.
//   - [NoneOF] excludes components.
//   - [AnyNOT] matches missing components.
//   - [AND], [OR], [XOR] logically combine two filters.
//   - [NOT] inverts any other filter.
package filter
