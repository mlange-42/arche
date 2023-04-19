package posvel

import (
	"github.com/sedyh/mizu/pkg/engine"
	"testing"
)

func BenchmarkIterMizu(b *testing.B) {
	b.StopTimer()

	_ = engine.NewGame(func(w engine.World) {
		for j := 0; j < nPos; j++ {
			w.AddEntity(&Position{})
		}
		for j := 0; j < nPosVel; j++ {
			w.AddEntity(&Position{}, &Velocity{})
		}

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Each(func(e engine.Entity) {
				pos := engine.Get[Position](e)
				vel := engine.Get[Velocity](e)

				pos.X += vel.X
				pos.Y += vel.Y
			}, engine.And[Position](), engine.And[Velocity]())
		}
	}).Update()
}

func BenchmarkBuildMizu(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = engine.NewGame(func(w engine.World) {
			for j := 0; j < nPos; j++ {
				w.AddEntity(&Position{})
			}
			for j := 0; j < nPosVel; j++ {
				w.AddEntity(&Position{}, &Velocity{})
			}
		}).Update()

	}
}
