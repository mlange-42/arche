# Changelog

## [[v0.6.2]](https://github.com/mlange-42/arche/compare/v0.6.1...v0.6.2)

### Performance

* Speed up generating world stats by factor 10 be re-using stats object (#210)

## [[v0.6.1]](https://github.com/mlange-42/arche/compare/v0.6.0...v0.6.1)

### Documentation

* Extend documentation and benchmarks for `Entity` (#201)
* Add a section with links to the Arche tools [arche-model](https://github.com/mlange-42/arche-model) and [arche-pixel](https://github.com/mlange-42/arche-pixel) (#202)

## [[v0.6.0]](https://github.com/mlange-42/arche/compare/v0.5.1...v0.6.0)

Arche v0.6.0 features fast batch entity creation and deletion, cached filters, and many internal optimizations.

### Highlights

* Batch creation and deletion of entities, with up to 4x and 10x speedup, respectively. Even more when combined with `World.Reset()`.
* Cached filters for handling many archetypes and complex queries without slowdown.
* A lot of internal performance optimizations.

### Breaking changes

* Generic mappers do no longer return all components when creating entities or components (#145)
* Resources API moved out of the world, to a helper to get by `World.Resources()` (#150)
* `World.Reset()` does no longer remove the component change listener (#157)
* Removes methods `filter.ALL.Not()` and `filter.ANY.Not()`, use `NoneOf()` and `AnyNot()` instead (#160)
* World listener function takes a pointer to the `EntityEvent` instead of a copy as argument (#162)

### Features

* Adds method `World.Reset()`, to allow for more efficient systematic simulations (#138)
* Adds `World.Batch()` helper for performing optimized batch-creation and batch-removal of entities (#149)
* Adds method `Mask.Exclusive()` to create a filter matching an exact component composition (#149, #188)
* Generic mappers (`Map1`, ...) have methods `NewEntities`, `NewEntitiesWith` and `RemoveEntities` for batch operations (#151)
* Batch-creation methods (ID-based and generic) have variants like `NewEntitiesQuery` that return a query over the created entities (#152)
* Notification during batch-creation is delayed until the resp. query is closed (#157)
* Batch-remove methods (`RemoveEntities()`) return the number of removed entities (#173)
* Filters can be cached and tracked by the `World` to speed up queries when there are many archetypes (#178)
* Function `AddResource[T](*World)` returns the ID of the resource (#183)

### Performance

* Speedup of archetype mask checks by 10% by checking mask before empty archetype (#139)
* Speedup of generic queries and mappers to come closer to ID-based access (#144)
* Speedup of archetype mask checks by casting filter interface to concrete type when possible (#148)
* Optimized batch creation of entities (#159)
* More efficiently clear the memory of removed components, with 2-3x speedup (#165)
* Do not clear memory when adding entities to archetypes, not required anymore as of #147 (#165)
* Speed up copying entity to archetype by getting entity pointer without reflection (#166)
* Avoid slice allocations in generic mapper methods (#170)
* Avoid type checks in query when iterating archetypes (#179)
* Speed up counting entities in queries with a cached filter (#182)
* Implements a fast and memory-efficient lookup data structure for components ID keys, to reduce the memory footprint of archetypes and the archetype graph (#192)
* Speedup of archetype creation by 40% by using a `const` for archetype storage page sizes (#197)

### Bugfixes

* Archetype storage buffers are "zeroed" when removing entities, to allow GC on pointers and slices in components (#147)
* Use slices instead of arrays inside paged archetype list to ensure pointer persistence (#184)

### Documentation

* Adds an example for batch-creation and batch-removal of entities (#173)
* Adds code examples to most public types, methods and functions (#183, #189)

### Other

* Restructure and extend benchmarks (#146, #153, #155, #156)
* Add an ECS competition benchmark for adding and removing components (#170)
* Add benchmarks for different ways to implement parent-child relations between entities (#194, #195)

## [[v0.5.1]](https://github.com/mlange-42/arche/compare/v0.5.0...v0.5.1)

### Performance

* Speedup of archetype access by 5-10% by merging storages into archetypes (#137)

### Documentation

* Document all private functions, types and methods  (#136)
* Adds a section  and plot on benchmarks against other Go ECS implementations to the README (#138)

### Other

* Internal code refactoring (#136)
  * Move method `nextArchetype` from `World` to `Query`.
  * Remove internal type `queryIter`.
  * Move repetitive pointer copying code in `storage` into a private method.
  * Move repetitive entity creation code in  `World` into a private method.

## [[v0.5.0]](https://github.com/mlange-42/arche/compare/v0.4.6...v0.5.0)

Feature release. Does not break v0.4.x projects.

### Features

* The World handles ECS resources, i.e. component-like global data (#132)
* Generic access to world resources (#132)

### Documentation

* Adds an example for the use of resources (#132)

## [[v0.4.6]](https://github.com/mlange-42/arche/compare/v0.4.5...v0.4.6)

### Performance

* Speedup archetype access by 10%, by elimination of bounds checks (#126)
* Speedup entity access from queries by 50% by using a component storage for them (#131)
* Minor optimizations of component storage (#128)

### Documentation

* Adds an example to demonstrate how to implement classical ECS systems (#129)

## [[v0.4.5]](https://github.com/mlange-42/arche/compare/v0.4.4...v0.4.5)

### Features

* Adds memory per entity to archetype stats (#124)

### Other

* Adds benchmarks of Arche vs. Array of Structs (AoS) and Array of Pointers (AoP), for different memory per entity and number of entities (#123)

## [[v0.4.4]](https://github.com/mlange-42/arche/compare/v0.4.3...v0.4.4)

### Features

* `Query` has methods `Count()` and `Step(int)`, primarily for effective random sampling (#119)

### Documentation

* Adds example `random_sampling` to demonstrate usage of `Query.Count()` and `Query.Step(int)` (#119)

### Bugfixes

* `Query.Next`, `Query.Get`, etc. now always panic when called on a closed query (#117)

### Other

* Update to [go-gameengine-ecs](https://github.com/marioolofo/go-gameengine-ecs) v0.9.0 in benchmarks (#116)
* Remove internal wrapper structs in generic queries and maps (#120)

## [[v0.4.3]](https://github.com/mlange-42/arche/compare/v0.4.2...v0.4.3)

### Bugfixes

* `EntityEvent` has more consistent values when an entity is removed (#115)
  * `EntityEvent.NewMask` is zero
  * `EntityEvent.Removed` is contains all former components
  * `EntityEvent.Current` is `nil`

## [[v0.4.2]](https://github.com/mlange-42/arche/compare/v0.4.1...v0.4.2)

### Performance

* Avoid creation of unused archetypes by splitting the archetype graph out of the actual archetypes (#113)
* Use slice instead of fixed-size array for type lookup in component registry (#113)
* Avoid copying `entityIndex` structs by using pointers (#114)

## [[v0.4.1]](https://github.com/mlange-42/arche/compare/v0.4.0...v0.4.1)

### Bugfixes

* Fix units symbol for bytes from `b` to `B` in string formatting of world statistics (#111)

### Other

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

### Performance

* Generic queries are compiled to masks and cached on first build (#62)
* Optimization of adding/removing components, with 2-3x speedup and vast reduction of (number of) allocations (#93)
* Speed up component world access by use of nil pointer check instead of bitmask (#96)

### Other

* Overhaul of the module structure, with generics and filters in separate packages (#55, #57, #61, #64)
* Boilerplate code for generic filters and queries is auto-generated with `go generate` (#64)
* Ensure 100% test coverage by adding a CI check for it (#68)
* `World.RemEntity(Entity)` is now `World.RemoveEntity(Entity)` (#87)
* More examples as user documentation (#83, #95)
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

### Performance

* Use of an archetype graph to speed up finding the target archetype for component addition/removal (#42)
* Minor optimization of component access by queries (#50)

### Other

* Reduced dependencies by moving profiling and benchmarking to sub-modules (#46)
* Smaller integer type for component identifiers (#47)

## [[v0.2.0]](https://github.com/mlange-42/arche/compare/v0.1.4...v0.2.0)

### Features

* `World` has method `Exchange` to add and remove components in one go (#38)
* `World` has method `Assign` add and assign components in one go (#38)
* `World` has method `AssignN` add and assign multiple components in one go (#38)

### Performance

* Optimization of `Query` iteration, avoids allocations and makes it approx. 30% faster (#35)
* Much smaller archetype data structure at the cost of one more index lookup (#37)

### Other

* Removed method `Query.Count()`, as it was a by-product of the allocations in the above point (#35)
* Archetypes are stored in a paged collection to use more efficient access by pointers (#36)

## [[v0.1.4]](https://github.com/mlange-42/arche/compare/v0.1.3...v0.1.4)

### Documentation

* Extended and improved documentation (#34)

## [[v0.1.3]](https://github.com/mlange-42/arche/compare/v0.1.2...v0.1.3)

### Features

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
