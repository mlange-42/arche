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
	Relation       ecs.ID
	Target         ecs.Entity
	HasRelation    bool
	compiled       bool
	targetCompiled bool
	locked         bool
}

func newCompiledQuery() compiledQuery {
	return compiledQuery{}
}

// Compile compiles a generic filter.
func (q *compiledQuery) Compile(w *ecs.World, include, optional, exclude []Comp, exclusive bool, targetType Comp, target ecs.Entity, hasTarget bool) {
	if q.compiled {
		return
	}

	q.Ids = toIds(w, include)

	incl := toMaskOptional(w, q.Ids, optional)
	var excl ecs.Mask
	if exclusive {
		excl = incl.Not()
	} else {
		excl = toMask(w, exclude)
	}
	q.maskFilter = ecs.MaskFilter{
		Include: incl,
		Exclude: excl,
	}
	noExclude := !exclusive && len(exclude) == 0

	if targetType == nil {
		if noExclude {
			q.filter = q.maskFilter.Include
		} else {
			q.filter = &q.maskFilter
		}
		q.Relation = ecs.ID{}
		q.HasRelation = false
	} else {

		targetID := ecs.TypeID(w, targetType)

		q.Relation = targetID
		q.HasRelation = true

		if !q.maskFilter.Include.Get(targetID) {
			panic(fmt.Sprintf("relation component %v not in filter", targetType))
		}
		isRelation := false
		if targetType.Kind() == reflect.Struct && targetType.NumField() > 0 {
			field := targetType.Field(0)
			isRelation = field.Type == relationType && field.Name == relationType.Name()
		}
		if !isRelation {
			panic(fmt.Sprintf("component type %v is not a relation", targetType))
		}

		if hasTarget {
			q.Target = target
			q.relationFilter = ecs.NewRelationFilter(&q.maskFilter, target)
			q.filter = &q.relationFilter
		} else {
			if noExclude {
				q.filter = q.maskFilter.Include
			} else {
				q.filter = &q.maskFilter
			}
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
