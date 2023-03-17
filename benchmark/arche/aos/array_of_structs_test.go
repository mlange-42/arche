package aos

import (
	"testing"
)

type Str16 struct {
	S0 Struct16B0
}

type Str32 struct {
	S0 Struct16B0
	S1 Struct16B1
}

type Str64 struct {
	S0 Struct16B0
	S1 Struct16B1
	S2 Struct16B2
	S3 Struct16B3
}

type Str128 struct {
	S0 Struct16B0
	S1 Struct16B1
	S2 Struct16B2
	S3 Struct16B3
	S4 Struct16B4
	S5 Struct16B5
	S6 Struct16B6
	S7 Struct16B6
}

func runAoS16B(b *testing.B, count int) {
	b.StopTimer()
	entities := make([]Str16, count)

	for i := 0; i < count; i++ {
		entities[i] = Str16{}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.S0.Hi++
			e.S0.Lo++
		}
	}
}

func runAoS32B(b *testing.B, count int) {
	b.StopTimer()
	entities := make([]Str32, count)

	for i := 0; i < count; i++ {
		entities[i] = Str32{}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.S0.Hi++
			e.S0.Lo++
		}
	}
}

func runAoS64B(b *testing.B, count int) {
	b.StopTimer()
	entities := make([]Str64, count)

	for i := 0; i < count; i++ {
		entities[i] = Str64{}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.S0.Hi++
			e.S0.Lo++
		}
	}
}

func runAoS128B(b *testing.B, count int) {
	b.StopTimer()
	entities := make([]Str128, count)

	for i := 0; i < count; i++ {
		entities[i] = Str128{}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.S0.Hi++
			e.S0.Lo++
		}
	}
}

func BenchmarkArrOfStructs_16B_1_000(b *testing.B) {
	runAoS16B(b, 1000)
}

func BenchmarkArrOfStructs_16B_10_000(b *testing.B) {
	runAoS16B(b, 10000)
}

func BenchmarkArrOfStructs_16B_100_000(b *testing.B) {
	runAoS16B(b, 100000)
}

func BenchmarkArrOfStructs_32B_1_000(b *testing.B) {
	runAoS32B(b, 1000)
}

func BenchmarkArrOfStructs_32B_10_000(b *testing.B) {
	runAoS32B(b, 10000)
}

func BenchmarkArrOfStructs_32B_100_000(b *testing.B) {
	runAoS32B(b, 100000)
}

func BenchmarkArrOfStructs_64B_1_000(b *testing.B) {
	runAoS64B(b, 1000)
}

func BenchmarkArrOfStructs_64B_10_000(b *testing.B) {
	runAoS64B(b, 10000)
}

func BenchmarkArrOfStructs_64B_100_000(b *testing.B) {
	runAoS64B(b, 100000)
}

func BenchmarkArrOfStructs_128B_1_000(b *testing.B) {
	runAoS128B(b, 1000)
}

func BenchmarkArrOfStructs_128B_10_000(b *testing.B) {
	runAoS128B(b, 10000)
}

func BenchmarkArrOfStructs_128B_100_000(b *testing.B) {
	runAoS128B(b, 100000)
}
