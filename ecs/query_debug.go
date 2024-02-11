//go:build debug

package ecs

const isDebug = true

func (q *Query) checkNext() {
	if q.nodeIndex < -1 {
		panic("query iteration already finished")
	}
}

func (q *Query) checkGet() {
	if q.access == nil {
		panic("query already iterated or iteration not started yet")
	}
}
