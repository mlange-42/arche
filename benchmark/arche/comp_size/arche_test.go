package compsize

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
)

func benchmarkCompSizeSimple1x08B(bench *testing.B, count int) {
	bench.StopTimer()
	world := ecs.NewWorld(count)

	aID := ecs.ComponentID[A](&world)

	ecs.NewBuilder(&world, aID).NewBatch(count)
	filter := ecs.All(aID)
	bench.StartTimer()

	for i := 0; i < bench.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			a := (*A)(query.Get(aID))
			a.V++
		}
	}
}

func benchmarkCompSizeSimple1x16B(bench *testing.B, count int) {
	bench.StopTimer()
	world := ecs.NewWorld(count)

	abID := ecs.ComponentID[AB](&world)

	ecs.NewBuilder(&world, abID).NewBatch(count)
	filter := ecs.All(abID)
	bench.StartTimer()

	for i := 0; i < bench.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			ab := (*AB)(query.Get(abID))

			ab.V1++
		}
	}
}

func benchmarkCompSizeSimple1x32B(bench *testing.B, count int) {
	bench.StopTimer()
	world := ecs.NewWorld(count)

	abcdID := ecs.ComponentID[ABCD](&world)

	ecs.NewBuilder(&world, abcdID).NewBatch(count)
	filter := ecs.All(abcdID)
	bench.StartTimer()

	for i := 0; i < bench.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			abcd := (*ABCD)(query.Get(abcdID))

			abcd.V1++
		}
	}
}

func benchmarkCompSizeSimple1x64B(bench *testing.B, count int) {
	bench.StopTimer()
	world := ecs.NewWorld(count)

	allID := ecs.ComponentID[All](&world)

	ecs.NewBuilder(&world, allID).NewBatch(count)
	filter := ecs.All(allID)
	bench.StartTimer()

	for i := 0; i < bench.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			all := (*All)(query.Get(allID))

			all.V1++
		}
	}
}

func benchmarkCompSizeSimple1x128B(bench *testing.B, count int) {
	bench.StopTimer()
	world := ecs.NewWorld(count)

	allID := ecs.ComponentID[All128B](&world)

	ecs.NewBuilder(&world, allID).NewBatch(count)
	filter := ecs.All(allID)
	bench.StartTimer()

	for i := 0; i < bench.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			all := (*All128B)(query.Get(allID))

			all.V1++
		}
	}
}

func benchmarkCompSize8x08B(bench *testing.B, count int) {
	bench.StopTimer()
	world := ecs.NewWorld(count)

	aID := ecs.ComponentID[A](&world)
	bID := ecs.ComponentID[B](&world)
	cID := ecs.ComponentID[C](&world)
	dID := ecs.ComponentID[D](&world)
	eID := ecs.ComponentID[E](&world)
	fID := ecs.ComponentID[F](&world)
	gID := ecs.ComponentID[G](&world)
	hID := ecs.ComponentID[H](&world)

	ecs.NewBuilder(&world, aID, bID, cID, dID, eID, fID, gID, hID).NewBatch(count)
	filter := ecs.All(aID, bID, cID, dID, eID, fID, gID, hID)
	bench.StartTimer()

	for i := 0; i < bench.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			a := (*A)(query.Get(aID))
			b := (*B)(query.Get(bID))
			c := (*C)(query.Get(cID))
			d := (*D)(query.Get(dID))
			e := (*E)(query.Get(eID))
			f := (*F)(query.Get(fID))
			g := (*G)(query.Get(gID))
			h := (*H)(query.Get(hID))

			a.V = b.V + c.V + d.V + e.V + f.V + g.V + h.V
		}
	}
}

func benchmarkCompSize4x16B(bench *testing.B, count int) {
	bench.StopTimer()
	world := ecs.NewWorld(count)

	abID := ecs.ComponentID[AB](&world)
	cdID := ecs.ComponentID[CD](&world)
	efID := ecs.ComponentID[EF](&world)
	ghID := ecs.ComponentID[GH](&world)

	ecs.NewBuilder(&world, abID, cdID, efID, ghID).NewBatch(count)
	filter := ecs.All(abID, cdID, efID, ghID)
	bench.StartTimer()

	for i := 0; i < bench.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			ab := (*AB)(query.Get(abID))
			cd := (*CD)(query.Get(cdID))
			ef := (*EF)(query.Get(efID))
			gh := (*GH)(query.Get(ghID))

			ab.V1 = ab.V2 + cd.V1 + cd.V2 + ef.V1 + ef.V2 + gh.V1 + gh.V2
		}
	}
}

func benchmarkCompSize2x32B(bench *testing.B, count int) {
	bench.StopTimer()
	world := ecs.NewWorld(count)

	abcdID := ecs.ComponentID[ABCD](&world)
	efghID := ecs.ComponentID[EFGH](&world)

	ecs.NewBuilder(&world, abcdID, efghID).NewBatch(count)
	filter := ecs.All(abcdID, efghID)
	bench.StartTimer()

	for i := 0; i < bench.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			abcd := (*ABCD)(query.Get(abcdID))
			efgh := (*EFGH)(query.Get(efghID))

			abcd.V1 = abcd.V2 + abcd.V3 + abcd.V4 + efgh.V1 + efgh.V2 + efgh.V3 + efgh.V4
		}
	}
}

func benchmarkCompSize1x64B(bench *testing.B, count int) {
	bench.StopTimer()
	world := ecs.NewWorld(count)

	allID := ecs.ComponentID[All](&world)

	ecs.NewBuilder(&world, allID).NewBatch(count)
	filter := ecs.All(allID)
	bench.StartTimer()

	for i := 0; i < bench.N; i++ {
		query := world.Query(&filter)
		for query.Next() {
			all := (*All)(query.Get(allID))

			all.V1 = all.V2 + all.V3 + all.V4 + all.V5 + all.V6 + all.V7 + all.V8
		}
	}
}

func BenchmarkCompSizeSimple_1_x_08B_10_000(bench *testing.B) {
	benchmarkCompSizeSimple1x08B(bench, 10000)
}

func BenchmarkCompSizeSimple_1_x_08B_100_000(bench *testing.B) {
	benchmarkCompSizeSimple1x08B(bench, 100000)
}

func BenchmarkCompSizeSimple_1_x_08B_1_000_000(bench *testing.B) {
	benchmarkCompSizeSimple1x08B(bench, 1000000)
}

func BenchmarkCompSizeSimple_1_x_16B_10_000(bench *testing.B) {
	benchmarkCompSizeSimple1x16B(bench, 10000)
}

func BenchmarkCompSizeSimple_1_x_16B_100_000(bench *testing.B) {
	benchmarkCompSizeSimple1x16B(bench, 100000)
}

func BenchmarkCompSizeSimple_1_x_16B_1_000_000(bench *testing.B) {
	benchmarkCompSizeSimple1x16B(bench, 1000000)
}

func BenchmarkCompSizeSimple_1_x_32B_10_000(bench *testing.B) {
	benchmarkCompSizeSimple1x32B(bench, 10000)
}

func BenchmarkCompSizeSimple_1_x_32B_100_000(bench *testing.B) {
	benchmarkCompSizeSimple1x32B(bench, 100000)
}

func BenchmarkCompSizeSimple_1_x_32B_1_000_000(bench *testing.B) {
	benchmarkCompSizeSimple1x32B(bench, 1000000)
}

func BenchmarkCompSizeSimple_1_x_64B_10_000(bench *testing.B) {
	benchmarkCompSizeSimple1x64B(bench, 10000)
}

func BenchmarkCompSizeSimple_1_x_64B_100_000(bench *testing.B) {
	benchmarkCompSizeSimple1x64B(bench, 100000)
}

func BenchmarkCompSizeSimple_1_x_64B_1_000_000(bench *testing.B) {
	benchmarkCompSizeSimple1x64B(bench, 1000000)
}

func BenchmarkCompSizeSimple_1_x_128B_10_000(bench *testing.B) {
	benchmarkCompSizeSimple1x128B(bench, 10000)
}

func BenchmarkCompSizeSimple_1_x_128B_100_000(bench *testing.B) {
	benchmarkCompSizeSimple1x128B(bench, 100000)
}

func BenchmarkCompSizeSimple_1_x_128B_1_000_000(bench *testing.B) {
	benchmarkCompSizeSimple1x128B(bench, 1000000)
}

func BenchmarkCompSize_8_x_08B_10_000(bench *testing.B) {
	benchmarkCompSize8x08B(bench, 10000)
}

func BenchmarkCompSize_8_x_08B_100_000(bench *testing.B) {
	benchmarkCompSize8x08B(bench, 100000)
}

func BenchmarkCompSize_4_x_16B_10_000(bench *testing.B) {
	benchmarkCompSize4x16B(bench, 10000)
}

func BenchmarkCompSize_4_x_16B_100_000(bench *testing.B) {
	benchmarkCompSize4x16B(bench, 100000)
}

func BenchmarkCompSize_2_x_32B_10_000(bench *testing.B) {
	benchmarkCompSize2x32B(bench, 10000)
}

func BenchmarkCompSize_2_x_32B_100_000(bench *testing.B) {
	benchmarkCompSize2x32B(bench, 100000)
}

func BenchmarkCompSize_1_x_64B_10_000(bench *testing.B) {
	benchmarkCompSize1x64B(bench, 10000)
}

func BenchmarkCompSize_1_x_64B_100_000(bench *testing.B) {
	benchmarkCompSize1x64B(bench, 100000)
}

type A struct {
	V int64
}

type B struct {
	V int64
}

type C struct {
	V int64
}

type D struct {
	V int64
}

type E struct {
	V int64
}

type F struct {
	V int64
}

type G struct {
	V int64
}

type H struct {
	V int64
}

type AB struct {
	V1 int64
	V2 int64
}

type CD struct {
	V1 int64
	V2 int64
}

type EF struct {
	V1 int64
	V2 int64
}

type GH struct {
	V1 int64
	V2 int64
}

type ABCD struct {
	V1 int64
	V2 int64
	V3 int64
	V4 int64
}

type EFGH struct {
	V1 int64
	V2 int64
	V3 int64
	V4 int64
}

type All struct {
	V1 int64
	V2 int64
	V3 int64
	V4 int64
	V5 int64
	V6 int64
	V7 int64
	V8 int64
}

type All128B struct {
	V1  int64
	V2  int64
	V3  int64
	V4  int64
	V5  int64
	V6  int64
	V7  int64
	V8  int64
	V9  int64
	V10 int64
	V11 int64
	V12 int64
	V13 int64
	V14 int64
	V15 int64
	V16 int64
}
