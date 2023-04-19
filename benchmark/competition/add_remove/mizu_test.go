package addremove

import (
	"github.com/sedyh/mizu/pkg/engine"
	"testing"
)

func BenchmarkIterMizu(b *testing.B) {
	b.StopTimer()

	_ = engine.NewGame(func(w engine.World) {
		position := &Position{}
		velocity := &Velocity{}

		for i := 0; i < nEntities; i++ {
			w.AddEntity(&position)
		}

		// Iterate once for more fairness
		w.Each(func(e engine.Entity) {
			engine.Set(e, velocity)
		}, engine.And[Position]())
		w.Each(func(e engine.Entity) {
			engine.Rem[Velocity](e)
		}, engine.And[Position]())

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Each(func(e engine.Entity) {
				engine.Set(e, velocity)
			}, engine.And[Position]())
			w.Each(func(e engine.Entity) {
				engine.Rem[Velocity](e)
			}, engine.And[Position]())
		}
	}).Update()
}

func BenchmarkBuildMizu(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = engine.NewGame(func(w engine.World) {
			for i := 0; i < nEntities; i++ {
				w.AddEntity(&Position{})
			}
		}).Update()
	}
}
