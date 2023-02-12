package generic

import (
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func toIds(w *ecs.World, types []reflect.Type) []ecs.ID {
	ids := make([]ecs.ID, len(types))
	for i, t := range types {
		ids[i] = ecs.TypeID(w, t)
	}
	return ids
}

func toMask(w *ecs.World, types []reflect.Type) ecs.Mask {
	mask := ecs.BitMask(0)
	for _, t := range types {
		mask.Set(ecs.TypeID(w, t), true)
	}
	return ecs.Mask{BitMask: mask}
}

func toMaskOptional(w *ecs.World, include []ecs.ID, optional []reflect.Type) ecs.Mask {
	mask := ecs.NewBitMask(include...)
	for _, t := range optional {
		mask.Set(ecs.TypeID(w, t), false)
	}
	return ecs.Mask{BitMask: mask}
}
