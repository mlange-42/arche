package generic

import (
	"fmt"
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

var relationType = reflect.TypeOf((*ecs.Relation)(nil)).Elem()

// compiledQuery is a helper for compiling a generic filter into a [ecs.Filter].
type compiledQuery struct {
	maskFilter     ecs.MaskFilter
	cachedFilter   ecs.CachedFilter
	filter         ecs.Filter
	Ids            []ecs.ID
	TargetID       int8
	compiled       bool
	targetCompiled bool
	locked         bool
}

func newCompiledQuery() compiledQuery {
	return compiledQuery{
		TargetID: -1,
	}
}

// Compile compiles a generic filter.
func (q *compiledQuery) Compile(w *ecs.World, include, optional, exclude []Comp, targetType Comp, target ecs.Entity) {
	if q.targetCompiled {
		return
	}

	if !q.compiled {
		q.Ids = toIds(w, include)
		q.maskFilter = ecs.MaskFilter{
			Include: toMaskOptional(w, q.Ids, optional),
			Exclude: toMask(w, exclude),
		}
	}

	if targetType == nil {
		q.filter = &q.maskFilter
		q.TargetID = -1
	} else {
		targetID := ecs.TypeID(w, targetType)

		if targetID != uint8(q.TargetID) {
			q.TargetID = int8(targetID)

			if !q.maskFilter.Include.Get(targetID) {
				panic(fmt.Sprintf("relation component %v not in filter", targetType))
			}
			isRelation := false
			if targetType.NumField() > 0 {
				field := targetType.Field(0)
				isRelation = field.Type == relationType && field.Name == relationType.Name()
			}
			if !isRelation {
				panic(fmt.Sprintf("component type %v is not a relation", targetType))
			}
		}

		q.filter = ecs.RelationFilter(&q.maskFilter, target)
	}
	q.targetCompiled = true
	q.compiled = true
}

// Reset sets the compiledQuery to not compiled.
func (q *compiledQuery) Reset(targetOnly bool) {
	q.targetCompiled = false
	if !targetOnly {
		q.compiled = false
	}
}

// Register the compiledQuery for caching.
func (q *compiledQuery) Register(w *ecs.World) {
	q.cachedFilter = w.Cache().Register(q.filter)
	q.filter = &q.cachedFilter
	q.locked = true
}

// Unregister the compiledQuery from caching.
func (q *compiledQuery) Unregister(w *ecs.World) {
	if cf, ok := q.filter.(*ecs.CachedFilter); ok {
		q.filter = w.Cache().Unregister(cf)
	} else {
		panic("can't unregister a filter that is not cached")
	}
	q.locked = false
}
