package generic

import (
	"github.com/mlange-42/arche/ecs"
)

type compiledQuery struct {
	filter   ecs.MaskFilter
	Ids      []ecs.ID
	compiled bool
}

func (q *compiledQuery) Compile(world ecs.IWorld, include, optional, exclude []Comp) {
	if q.compiled {
		return
	}
	q.Ids = toIds(world, include)
	q.filter = ecs.MaskFilter{
		Include: toMaskOptional(world, q.Ids, optional),
		Exclude: toMask(world, exclude),
	}
	q.compiled = true
}

func (q *compiledQuery) Reset() {
	q.compiled = false
}
