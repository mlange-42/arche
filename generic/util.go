package generic

import (
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

// filter is a helper to simplify generated generic filter code.
type filter struct {
	include  []Comp
	optional []Comp
	exclude  []Comp
	compiled compiledQuery
}

// typeOf is a shortcut for getting the reflection type of a generic type argument.
func typeOf[T any]() Comp {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// toIds extracts [ecs.ID]s from a sequence of [Comp]s.
func toIds(w *ecs.World, types []Comp) []ecs.ID {
	ids := make([]ecs.ID, len(types))
	for i, t := range types {
		ids[i] = ecs.TypeID(w, t)
	}
	return ids
}

// toMask extracts an [ecs.Mask] from a sequence of [Comp]s.
func toMask(w *ecs.World, types []Comp) ecs.Mask {
	mask := ecs.Mask{}
	for _, t := range types {
		mask.Set(ecs.TypeID(w, t), true)
	}
	return mask
}

// toMaskOptional extracts an [ecs.Mask] from a sequence of [Comp]s, ignoring the given optional [Comp]s.
func toMaskOptional(w *ecs.World, include []ecs.ID, optional []Comp) ecs.Mask {
	mask := ecs.All(include...)
	for _, t := range optional {
		mask.Set(ecs.TypeID(w, t), false)
	}
	return mask
}
