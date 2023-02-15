package generic

import (
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

// Comp is an alias for component types.
type Comp reflect.Type

// T provides a component type from a generic type argument.
func T[A any]() Comp {
	return Comp(typeOf[A]())
}

func typeOf[T any]() Comp {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func toIds(w *ecs.World, types []Comp) []ecs.ID {
	ids := make([]ecs.ID, len(types))
	for i, t := range types {
		ids[i] = ecs.TypeID(w, t)
	}
	return ids
}

func toMask(w *ecs.World, types []Comp) ecs.Mask {
	mask := ecs.BitMask{}
	for _, t := range types {
		mask.Set(ecs.TypeID(w, t), true)
	}
	return ecs.Mask{BitMask: mask}
}

func toMaskOptional(w *ecs.World, include []ecs.ID, optional []Comp) ecs.Mask {
	mask := ecs.NewBitMask(include...)
	for _, t := range optional {
		mask.Set(ecs.TypeID(w, t), false)
	}
	return ecs.Mask{BitMask: mask}
}
