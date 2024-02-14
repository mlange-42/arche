package generic

import (
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

// filter is a helper to simplify generated generic filter code.
type filter struct {
	include    []Comp
	optional   []Comp
	exclude    []Comp
	exclusive  bool
	targetType Comp
	target     ecs.Entity
	hasTarget  bool
	compiled   compiledQuery
}

func newFilter(include ...Comp) filter {
	return filter{
		include:  include,
		compiled: newCompiledQuery(),
	}
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

func newEntity(w *ecs.World, ids []ecs.ID, relation ecs.ID, hasRelation bool, target ...ecs.Entity) ecs.Entity {
	if len(target) == 0 {
		return w.NewEntity(ids...)
	}
	if !hasRelation {
		panic("map has no relation defined, can't set a target")
	}
	return ecs.NewBuilder(w, ids...).WithRelation(relation).New(target[0])
}

func newBatch(w *ecs.World, count int, ids []ecs.ID, relation ecs.ID, hasRelation bool, target ...ecs.Entity) {
	if len(target) == 0 {
		ecs.NewBuilder(w, ids...).NewBatch(count)
		return
	}
	if !hasRelation {
		panic("map has no relation defined, can't set a target")
	}
	ecs.NewBuilder(w, ids...).WithRelation(relation).NewBatch(count, target[0])
}

func newQuery(w *ecs.World, count int, ids []ecs.ID, relation ecs.ID, hasRelation bool, target ...ecs.Entity) ecs.Query {
	var query ecs.Query

	if len(target) == 0 {
		query = ecs.NewBuilder(w, ids...).NewBatchQ(count)
	} else {
		if !hasRelation {
			panic("map has no relation defined, can't set a target")
		}
		query = ecs.NewBuilder(w, ids...).WithRelation(relation).NewBatchQ(count, target[0])
	}

	return query
}
