# Changelog

## [[unpublished]](https://github.com/mlange-42/arche/compare/v0.1.1...main)

### Other

* Use aligned item size in component storage, for approx. 20% faster query iteration (#25)

## [[v0.1.1]](https://github.com/mlange-42/arche/compare/v0.1.0...v0.1.1)

### Other

* Avoid allocation in `World.Has(entity, compID)` (#16)
* `World.RemEntity(entity)` panics on dead entity, like all other `World` methods (#18)
* Reserve zero value `Entity` to serve as nil/undefined value (#23)

## [[v0.1.0]](https://github.com/mlange-42/arche/tree/v0.1.0)

Initial release.

Basic ECS implementation.
