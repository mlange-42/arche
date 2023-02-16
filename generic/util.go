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

type query struct {
	ecs.Query
	ids []ecs.ID
}

type mapper struct {
	ids   []ecs.ID
	world *ecs.World
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
