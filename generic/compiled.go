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
	relationFilter ecs.RelationFilter
	cachedFilter   ecs.CachedFilter
	filter         ecs.Filter
	Ids            []ecs.ID
	TargetComp     ecs.ID
	Target         ecs.Entity
	HasTarget      bool
	compiled       bool
	targetCompiled bool
	locked         bool
}

func newCompiledQuery() compiledQuery {
	return compiledQuery{}
}

// Compile compiles a generic filter.
func (q *compiledQuery) Compile(w *ecs.World, include, optional, exclude []Comp, targetType Comp, target ecs.Entity, hasTarget bool) {
	if q.compiled {
		return
	}

	q.Ids = toIds(w, include)
	q.maskFilter = ecs.MaskFilter{
		Include: toMaskOptional(w, q.Ids, optional),
		Exclude: toMask(w, exclude),
	}

	if targetType == nil {
		q.filter = &q.maskFilter
		q.TargetComp = ecs.ID{}
		q.HasTarget = false
	} else {

		targetID := ecs.TypeID(w, targetType)

		q.TargetComp = targetID
		q.HasTarget = true

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

		if hasTarget {
			q.HasTarget = true
			q.Target = target
			q.relationFilter = ecs.NewRelationFilter(&q.maskFilter, target)
			q.filter = &q.relationFilter
		} else {
			q.filter = &q.maskFilter
		}
	}
	q.targetCompiled = true
	q.compiled = true
}

// Reset sets the compiledQuery to not compiled.
func (q *compiledQuery) Reset() {
	q.compiled = false
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
