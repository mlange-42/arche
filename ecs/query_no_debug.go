//go:build !debug

package ecs

const isDebug = false

func (q *Query) checkNext() {}

func (q *Query) checkGet() {}
