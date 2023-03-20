package generic

import (
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

type filter struct {
	include  []Comp
	optional []Comp
	exclude  []Comp
	compiled compiledQuery
}

func typeOf[T any]() Comp {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func toIds(w ecs.IWorld, types []Comp) []ecs.ID {
	ids := make([]ecs.ID, len(types))
	for i, t := range types {
		ids[i] = ecs.TypeID(w, t)
	}
	return ids
}

func toMask(w ecs.IWorld, types []Comp) ecs.Mask {
	mask := ecs.Mask{}
	for _, t := range types {
		mask.Set(ecs.TypeID(w, t), true)
	}
	return mask
}

func toMaskOptional(w ecs.IWorld, include []ecs.ID, optional []Comp) ecs.Mask {
	mask := ecs.All(include...)
	for _, t := range optional {
		mask.Set(ecs.TypeID(w, t), false)
	}
	return mask
}
