package aos

import (
	"testing"
)

func runAoP16B(b *testing.B, count int) {
	b.StopTimer()
	entities := make([]*Str16, count)

	for i := 0; i < count; i++ {
		entities[i] = &Str16{}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.S0.Hi++
			e.S0.Lo++
		}
	}
}

func runAoP32B(b *testing.B, count int) {
	b.StopTimer()
	entities := make([]*Str32, count)

	for i := 0; i < count; i++ {
		entities[i] = &Str32{}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.S0.Hi++
			e.S0.Lo++
		}
	}
}

func runAoP64B(b *testing.B, count int) {
	b.StopTimer()
	entities := make([]*Str64, count)

	for i := 0; i < count; i++ {
		entities[i] = &Str64{}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.S0.Hi++
			e.S0.Lo++
		}
	}
}

func runAoP128B(b *testing.B, count int) {
	b.StopTimer()
	entities := make([]*Str128, count)

	for i := 0; i < count; i++ {
		entities[i] = &Str128{}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.S0.Hi++
			e.S0.Lo++
		}
	}
}

func runAoP256B(b *testing.B, count int) {
	b.StopTimer()
	entities := make([]*Str256, count)

	for i := 0; i < count; i++ {
		entities[i] = &Str256{}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.S0.Hi++
			e.S0.Lo++
		}
	}
}

func BenchmarkArrOfPointers_16B_1_000(b *testing.B) {
	runAoP16B(b, 1000)
}

func BenchmarkArrOfPointers_16B_10_000(b *testing.B) {
	runAoP16B(b, 10000)
}

func BenchmarkArrOfPointers_16B_100_000(b *testing.B) {
	runAoP16B(b, 100000)
}

func BenchmarkArrOfPointers_32B_1_000(b *testing.B) {
	runAoP32B(b, 1000)
}

func BenchmarkArrOfPointers_32B_10_000(b *testing.B) {
	runAoP32B(b, 10000)
}

func BenchmarkArrOfPointers_32B_100_000(b *testing.B) {
	runAoP32B(b, 100000)
}

func BenchmarkArrOfPointers_64B_1_000(b *testing.B) {
	runAoP64B(b, 1000)
}

func BenchmarkArrOfPointers_64B_10_000(b *testing.B) {
	runAoP64B(b, 10000)
}

func BenchmarkArrOfPointers_64B_100_000(b *testing.B) {
	runAoP64B(b, 100000)
}

func BenchmarkArrOfPointers_128B_1_000(b *testing.B) {
	runAoP128B(b, 1000)
}

func BenchmarkArrOfPointers_128B_10_000(b *testing.B) {
	runAoP128B(b, 10000)
}

func BenchmarkArrOfPointers_128B_100_000(b *testing.B) {
	runAoP128B(b, 100000)
}

func BenchmarkArrOfPointers_256B_1_000(b *testing.B) {
	runAoP256B(b, 1000)
}

func BenchmarkArrOfPointers_256B_10_000(b *testing.B) {
	runAoP256B(b, 10000)
}

func BenchmarkArrOfPointers_256B_100_000(b *testing.B) {
	runAoP256B(b, 100000)
}
