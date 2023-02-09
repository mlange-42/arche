# Changelog

## [[unpublished]](https://github.com/mlange-42/arche/compare/v0.2.0...main)

### Other

* Use of an archetype graph to speed up finding the target archetype for component addition/removal (#42)
* Reduced dependencies by moving profiling and benchmarking to sub-modules (#46)

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
