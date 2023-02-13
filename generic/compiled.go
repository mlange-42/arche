package generic

import (
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

type compiledQuery struct {
	mask     ecs.Mask
	exclude  ecs.Mask
	Ids      []ecs.ID
	compiled bool
}

func (q *compiledQuery) Compile(w *ecs.World, include, optional, exclude []reflect.Type) {
	if q.compiled {
		return
	}
	q.Ids = toIds(w, include)
	q.mask = toMaskOptional(w, q.Ids, optional)
	q.exclude = toMask(w, exclude)
	q.compiled = true
}

func (q *compiledQuery) Filter() ecs.MaskPair {
	return ecs.MaskPair{
		Mask:    q.mask,
		Exclude: q.exclude,
	}
}

func (q *compiledQuery) Reset() {
	q.compiled = false
}
