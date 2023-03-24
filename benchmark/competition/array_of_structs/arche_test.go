package arrayofstructs

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

type Struct16B0 struct{ ecs.Mask }
type Struct16B1 struct{ ecs.Mask }
type Struct16B2 struct{ ecs.Mask }
type Struct16B3 struct{ ecs.Mask }
type Struct16B4 struct{ ecs.Mask }
type Struct16B5 struct{ ecs.Mask }
type Struct16B6 struct{ ecs.Mask }
type Struct16B7 struct{ ecs.Mask }
type Struct16B8 struct{ ecs.Mask }
type Struct16B9 struct{ ecs.Mask }
type Struct16B10 struct{ ecs.Mask }
type Struct16B11 struct{ ecs.Mask }
type Struct16B12 struct{ ecs.Mask }
type Struct16B13 struct{ ecs.Mask }
type Struct16B14 struct{ ecs.Mask }
type Struct16B15 struct{ ecs.Mask }

func runArche16B(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	id0 := ecs.ComponentID[Struct16B0](&world)

	world.Batch().NewEntities(count, id0)
	var filter ecs.Filter = ecs.All(id0)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		for query.Next() {
			s0 := (*Struct16B0)(query.Get(id0))
			s0.Hi++
			s0.Lo++
		}
	}
}

func runArche32B(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	id0 := ecs.ComponentID[Struct16B0](&world)
	id1 := ecs.ComponentID[Struct16B1](&world)

	world.Batch().NewEntities(count, id0, id1)
	var filter ecs.Filter = ecs.All(id0)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		for query.Next() {
			s0 := (*Struct16B0)(query.Get(id0))
			s0.Hi++
			s0.Lo++
		}
	}
}

func runArche64B(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	id0 := ecs.ComponentID[Struct16B0](&world)
	id1 := ecs.ComponentID[Struct16B1](&world)
	id2 := ecs.ComponentID[Struct16B2](&world)
	id3 := ecs.ComponentID[Struct16B3](&world)

	world.Batch().NewEntities(count, id0, id1, id2, id3)
	var filter ecs.Filter = ecs.All(id0)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		for query.Next() {
			s0 := (*Struct16B0)(query.Get(id0))
			s0.Hi++
			s0.Lo++
		}
	}
}

func runArche128B(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	id0 := ecs.ComponentID[Struct16B0](&world)
	id1 := ecs.ComponentID[Struct16B1](&world)
	id2 := ecs.ComponentID[Struct16B2](&world)
	id3 := ecs.ComponentID[Struct16B3](&world)
	id4 := ecs.ComponentID[Struct16B4](&world)
	id5 := ecs.ComponentID[Struct16B5](&world)
	id6 := ecs.ComponentID[Struct16B6](&world)
	id7 := ecs.ComponentID[Struct16B7](&world)

	world.Batch().NewEntities(count, id0, id1, id2, id3, id4, id5, id6, id7)
	var filter ecs.Filter = ecs.All(id0)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		for query.Next() {
			s0 := (*Struct16B0)(query.Get(id0))
			s0.Hi++
			s0.Lo++
		}
	}
}

func runArche256B(b *testing.B, count int) {
	b.StopTimer()
	world := ecs.NewWorld()

	id0 := ecs.ComponentID[Struct16B0](&world)
	id1 := ecs.ComponentID[Struct16B1](&world)
	id2 := ecs.ComponentID[Struct16B2](&world)
	id3 := ecs.ComponentID[Struct16B3](&world)
	id4 := ecs.ComponentID[Struct16B4](&world)
	id5 := ecs.ComponentID[Struct16B5](&world)
	id6 := ecs.ComponentID[Struct16B6](&world)
	id7 := ecs.ComponentID[Struct16B7](&world)
	id8 := ecs.ComponentID[Struct16B8](&world)
	id9 := ecs.ComponentID[Struct16B9](&world)
	id10 := ecs.ComponentID[Struct16B10](&world)
	id11 := ecs.ComponentID[Struct16B11](&world)
	id12 := ecs.ComponentID[Struct16B12](&world)
	id13 := ecs.ComponentID[Struct16B13](&world)
	id14 := ecs.ComponentID[Struct16B14](&world)
	id15 := ecs.ComponentID[Struct16B15](&world)

	world.Batch().NewEntities(count,
		id0, id1, id2, id3, id4, id5, id6, id7,
		id8, id9, id10, id11, id12, id13, id14, id15,
	)
	var filter ecs.Filter = ecs.All(id0)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		for query.Next() {
			s0 := (*Struct16B0)(query.Get(id0))
			s0.Hi++
			s0.Lo++
		}
	}
}

func BenchmarkArche_16B_1_000(b *testing.B) {
	runArche16B(b, 1000)
}

func BenchmarkArche_16B_10_000(b *testing.B) {
	runArche16B(b, 10000)
}

func BenchmarkArche_16B_100_000(b *testing.B) {
	runArche16B(b, 100000)
}

func BenchmarkArche_32B_1_000(b *testing.B) {
	runArche32B(b, 1000)
}

func BenchmarkArche_32B_10_000(b *testing.B) {
	runArche32B(b, 10000)
}

func BenchmarkArche_32B_100_000(b *testing.B) {
	runArche32B(b, 100000)
}

func BenchmarkArche_64B_1_000(b *testing.B) {
	runArche64B(b, 1000)
}

func BenchmarkArche_64B_10_000(b *testing.B) {
	runArche64B(b, 10000)
}

func BenchmarkArche_64B_100_000(b *testing.B) {
	runArche64B(b, 100000)
}

func BenchmarkArche_128B_1_000(b *testing.B) {
	runArche128B(b, 1000)
}

func BenchmarkArche_128B_10_000(b *testing.B) {
	runArche128B(b, 10000)
}

func BenchmarkArche_128B_100_000(b *testing.B) {
	runArche128B(b, 100000)
}

func BenchmarkArche_256B_1_000(b *testing.B) {
	runArche256B(b, 1000)
}

func BenchmarkArche_256B_10_000(b *testing.B) {
	runArche256B(b, 10000)
}

func BenchmarkArche_256B_100_000(b *testing.B) {
	runArche256B(b, 100000)
}
