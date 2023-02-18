package generic

import (
	"github.com/mlange-42/arche/ecs"
)

type compiledQuery struct {
	filter   ecs.MaskPair
	Ids      []ecs.ID
	compiled bool
}

func (q *compiledQuery) Compile(w *ecs.World, include, optional, exclude []Comp) {
	if q.compiled {
		return
	}
	q.Ids = toIds(w, include)
	q.filter = ecs.MaskPair{
		Include: toMaskOptional(w, q.Ids, optional),
		Exclude: toMask(w, exclude),
	}
	q.compiled = true
}

func (q *compiledQuery) Reset() {
	q.compiled = false
}
