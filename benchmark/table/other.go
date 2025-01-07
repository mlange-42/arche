package main

import (
	"testing"

	"github.com/mlange-42/arche/benchmark"
	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func benchesOther() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "ecs.NewWorld", Desc: "", F: newWorld, N: 1, Factor: 0.001, Units: "μs"},
		{Name: "ecs.ComponentID", Desc: "registered component", F: componentID, N: 1},
	}
}

func newWorld(b *testing.B) {
	var w ecs.World

	for i := 0; i < b.N; i++ {
		w = ecs.NewWorld()
	}
	b.StopTimer()

	assert.False(b, w.IsLocked())
}

func componentID(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	origID := ecs.ComponentID[comp1](&w)

	var id ecs.ID

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		id = ecs.ComponentID[comp1](&w)
	}
	b.StopTimer()

	assert.Equal(b, origID, id)
}
