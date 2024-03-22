package manycomponents

import (
	"testing"

	"github.com/mlange-42/arche/benchmark/competition/many_components/geckecs"
)

func BenchmarkIterGECK(b *testing.B) {
	b.StopTimer()
	world := geckecs.NewWorld()

	posEntities := world.Entities(nPos)
	world.SetPositions(geckecs.Position{}, posEntities...,
	)
	posAllEntities := world.Entities(nPosAll)
	world.SetPositions(geckecs.Position{}, posAllEntities...)
	world.SetComp1S(geckecs.Comp1{}, posAllEntities...)
	world.SetComp2S(geckecs.Comp2{}, posAllEntities...)
	world.SetComp3S(geckecs.Comp3{}, posAllEntities...)
	world.SetComp4S(geckecs.Comp4{}, posAllEntities...)
	world.SetComp5S(geckecs.Comp5{}, posAllEntities...)
	world.SetComp6S(geckecs.Comp6{}, posAllEntities...)
	world.SetComp7S(geckecs.Comp7{}, posAllEntities...)
	world.SetComp8S(geckecs.Comp8{}, posAllEntities...)
	world.SetComp9S(geckecs.Comp9{}, posAllEntities...)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		iter := world.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.NewIterator()
		for iter.HasNext() {
			_, pos, c1, c2, c3, c4, c5, c6, c7, c8, c9 := iter.Next()
			pos.X += c1.X + c2.X + c3.X + c4.X + c5.X + c6.X + c7.X + c8.X + c9.X
			pos.Y += c1.Y + c2.Y + c3.Y + c4.Y + c5.Y + c6.Y + c7.Y + c8.Y + c9.Y
		}
	}
}
