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

		q.filter = &ecs.RelationFilter{
			Filter: &q.maskFilter,
			Target: target,
		}
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
