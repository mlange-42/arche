# Changelog

## [[unpublished]](https://github.com/mlange-42/arche/compare/v0.4.2...main)

Nothing

## [[v0.4.2]](https://github.com/mlange-42/arche/compare/v0.4.1...v0.4.2)

## Other

* Avoid creation of unused archetypes by splitting the archetype graph out of the actual archetypes (#113)
* Use slice instead of fixed-size array for type lookup in component registry (#113)
* Avoid copying `entityIndex` structs by using pointers (#114)

## [[v0.4.1]](https://github.com/mlange-42/arche/compare/v0.4.0...v0.4.1)

## Bugfixes

* Fix units symbol for bytes from `b` to `B` in string formatting of world statistics (#111)

## Other

* Adds [github.com/wfranczyk/ento](https://github.com/wfranczyk/ento) to benchmarks (#110)

## [[v0.4.0]](https://github.com/mlange-42/arche/compare/v0.3.1...v0.4.0)

API revision, split out generics and filters into separate packages.

### Features

* Generic queries support optional, additional and excluded components (#53, #56, #58, #59, #60, #63)
* Logic filters for complex queries (#54, #58, #61)
* `Query` and `World` have a method `Mask(Entity)` to access archetype masks (#65)
* Generic query `Get` method returns all queried components (#83)
* Added method `World.Stats()` for inspecting otherwise inaccessible world statistics (#67)
* Entities can be initialized with components, via ID as well as using generics (#76)
* A listener function can be registered to the world, for notification on entity changes (#77)
* Support for up to 128 distinct component types per world (was limited to 64 before) (#78)
* Generic entity manipulation through types `Map1`, `Map2`, ... and `Exchange` (#79, #84, #87)

### Other

* Overhaul of the module structure, with generics and filters in separate packages (#55, #57, #61, #64)
* Generic queries are compiled to masks and cached on first build (#62)
* Boilerplate code for generic filters and queries is auto-generated with `go generate` (#64)
* Ensure 100% test coverage by adding a CI check for it (#68)
* `World.RemEntity(Entity)` is now `World.RemoveEntity(Entity)` (#87)
* Optimization of adding/removing components, with 2-3x speedup and vast reduction of (number of) allocations (#93)
* More examples as user documentation (#83, #95)
* Speed up component world access by use of nil pointer check instead of bitmask (#96)
* General API cleanup with renaming of several types and methods (#100)

## [[v0.3.1]](https://github.com/mlange-42/arche/compare/v0.3.0...v0.3.1)

### Other

* Fix failing https://pkg.go.dev to fetch Arche version v0.3.0

## [[v0.3.0]](https://github.com/mlange-42/arche/compare/v0.2.0...v0.3.0)

### Features

* Added a layer of generic access as alternative for using component IDs, for type safety and ergonomics (#47, #48)
  * Generic queries like `Query1[T]`, `Query2[T, U]`, ... (#47)
  * Generic add, assign and remove (`Add[T]()`, `Add2[T, U](), ...`) (#47)
  * Generic get, has, and set through component mapper `Map[T]` (#47)

### Other

* Use of an archetype graph to speed up finding the target archetype for component addition/removal (#42)
* Reduced dependencies by moving profiling and benchmarking to sub-modules (#46)
* Smaller integer type for component identifiers (#47)
* Minor optimization of component access by queries (#50)

## [[v0.2.0]](https://github.com/mlange-42/arche/compare/v0.1.4...v0.2.0)

### Features

* `World` has method `Exchange` to add and remove components in one go (#38)
* `World` has method `Assign` add and assign components in one go (#38)
* `World` has method `AssignN` add and assign multiple components in one go (#38)

### Other

* Optimization of `Query` iteration, avoids allocations and makes it approx. 30% faster (#35)
* Removed method `Query.Count()`, as it was a by-product of the allocations in the above point (#35)
* Archetypes are stored in a paged collection to use more efficient access by pointers (#36)
* Much smaller archetype data structure at the cost of one more index lookup (#37)

## [[v0.1.4]](https://github.com/mlange-42/arche/compare/v0.1.3...v0.1.4)

### Other

* Extended and improved documentation (#34)

## [[v0.1.3]](https://github.com/mlange-42/arche/compare/v0.1.2...v0.1.3)

## Features

* Add `Config` to allow for configuration of the world (currently only storage capacity increment) (#28)
* `Query` has a method `Count()`, reporting the total number of matching entities (#30)

## [[v0.1.2]](https://github.com/mlange-42/arche/compare/v0.1.1...v0.1.2)

### Other

* Use aligned item size in component storage for faster query iteration (#25)
* Queries lock the World, and automatically unlock it after iteration (#26)

## [[v0.1.1]](https://github.com/mlange-42/arche/compare/v0.1.0...v0.1.1)

### Other

* Avoid allocation in `World.Has(entity, compID)` (#16)
* `World.RemEntity(entity)` panics on dead entity, like all other `World` methods (#18)
* Reserve zero value `Entity` to serve as nil/undefined value (#23)

## [[v0.1.0]](https://github.com/mlange-42/arche/tree/v0.1.0)

Initial release.

Basic ECS implementation.
