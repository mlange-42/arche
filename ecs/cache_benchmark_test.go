package ecs

import "testing"

func benchmarkCachedFilters(b *testing.B, arches int, count int, cached bool) {
	b.StopTimer()
	world := NewWorld()

	ids := []ID{
		ComponentID[testStruct0](&world),
		ComponentID[testStruct1](&world),
		ComponentID[testStruct2](&world),
		ComponentID[testStruct3](&world),
		ComponentID[testStruct4](&world),
		ComponentID[testStruct5](&world),
		ComponentID[testStruct6](&world),
		ComponentID[testStruct7](&world),
		ComponentID[testStruct8](&world),
		ComponentID[testStruct9](&world),
		ComponentID[testStruct10](&world),
	}
	id := ids[10]

	perArch := 25

	for i := 0; i < arches; i++ {
		mask := i
		add := make([]ID, 0, 10)
		for j := 0; j < 10; j++ {
			id := ids[j]
			m := 1 << j
			if mask&m == m {
				add = append(add, id)
			}
		}
		for j := 0; j < perArch; j++ {
			world.NewEntity(add...)
		}
	}

	NewBuilder(&world, id).NewBatch(count)

	var filter Filter = All(id)
	if cached {
		f := world.Cache().Register(filter)
		filter = &f
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(filter)
		for query.Next() {
			st := (*testStruct0)(query.Get(id))
			st.Val++
		}
	}
}

// 1 of 1
func BenchmarkFilterUncached_1of1_100(b *testing.B) {
	benchmarkCachedFilters(b, 1, 100, false)
}

func BenchmarkFilterUncached_1of1_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 1, 1000, false)
}

func BenchmarkFilterUncached_1of1_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 1, 10000, false)
}

func BenchmarkFilterCached_1of1_100(b *testing.B) {
	benchmarkCachedFilters(b, 1, 100, true)
}

func BenchmarkFilterCached_1of1_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 1, 1000, true)
}

func BenchmarkFilterCached_1of1_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 1, 10000, true)
}

// 1 of 4
func BenchmarkFilterUncached_1of4_100(b *testing.B) {
	benchmarkCachedFilters(b, 4, 100, false)
}

func BenchmarkFilterUncached_1of4_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 4, 1000, false)
}

func BenchmarkFilterUncached_1of4_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 4, 10000, false)
}

func BenchmarkFilterCached_1of4_100(b *testing.B) {
	benchmarkCachedFilters(b, 4, 100, true)
}

func BenchmarkFilterCached_1of4_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 4, 1000, true)
}

func BenchmarkFilterCached_1of4_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 4, 10000, true)
}

// 1 of 16
func BenchmarkFilterUncached_1of16_100(b *testing.B) {
	benchmarkCachedFilters(b, 16, 100, false)
}

func BenchmarkFilterUncached_1of16_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 16, 1000, false)
}

func BenchmarkFilterUncached_1of16_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 16, 10000, false)
}

func BenchmarkFilterCached_1of16_100(b *testing.B) {
	benchmarkCachedFilters(b, 16, 100, true)
}

func BenchmarkFilterCached_1of16_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 16, 1000, true)
}

func BenchmarkFilterCached_1of16_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 16, 10000, true)
}

// 1 of 64
func BenchmarkFilterUncached_1of64_100(b *testing.B) {
	benchmarkCachedFilters(b, 64, 100, false)
}

func BenchmarkFilterUncached_1of64_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 64, 1000, false)
}

func BenchmarkFilterUncached_1of64_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 64, 10000, false)
}

func BenchmarkFilterCached_1of64_100(b *testing.B) {
	benchmarkCachedFilters(b, 64, 100, true)
}

func BenchmarkFilterCached_1of64_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 64, 1000, true)
}

func BenchmarkFilterCached_1of64_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 64, 10000, true)
}

// 1 of 256
func BenchmarkFilterUncached_1of256_100(b *testing.B) {
	benchmarkCachedFilters(b, 256, 100, false)
}

func BenchmarkFilterUncached_1of256_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 256, 1000, false)
}

func BenchmarkFilterUncached_1of256_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 256, 10000, false)
}

func BenchmarkFilterCached_1of256_100(b *testing.B) {
	benchmarkCachedFilters(b, 256, 100, true)
}

func BenchmarkFilterCached_1of256_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 256, 1000, true)
}

func BenchmarkFilterCached_1of256_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 256, 10000, true)
}

// 1 of 1024
func BenchmarkFilterUncached_1of1024_100(b *testing.B) {
	benchmarkCachedFilters(b, 1024, 100, false)
}

func BenchmarkFilterUncached_1of1024_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 1024, 1000, false)
}

func BenchmarkFilterUncached_1of1024_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 1024, 10000, false)
}

func BenchmarkFilterCached_1of1024_100(b *testing.B) {
	benchmarkCachedFilters(b, 1024, 100, true)
}

func BenchmarkFilterCached_1of1024_1_000(b *testing.B) {
	benchmarkCachedFilters(b, 1024, 1000, true)
}

func BenchmarkFilterCached_1of1024_10_000(b *testing.B) {
	benchmarkCachedFilters(b, 1024, 10000, true)
}
