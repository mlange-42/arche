package generic

import (
	"github.com/mlange-42/arche/ecs"
)

// compiledQuery is a helper for compiling a generic filter into a [ecs.Filter].
type compiledQuery struct {
	filter   ecs.MaskFilter
	Ids      []ecs.ID
	compiled bool
}

// Compile compiles a generic filter.
func (q *compiledQuery) Compile(w *ecs.World, include, optional, exclude []Comp) {
	if q.compiled {
		return
	}
	q.Ids = toIds(w, include)
	q.filter = ecs.MaskFilter{
		Include: toMaskOptional(w, q.Ids, optional),
		Exclude: toMask(w, exclude),
	}
	q.compiled = true
}

// Reset sets the compiledQuery to not compiled.
func (q *compiledQuery) Reset() {
	q.compiled = false
}
