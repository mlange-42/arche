package arrayofstructs

import (
	"testing"
)

type Item16 struct {
	Next *Item16
	S0   Struct16B0
}

type Item32 struct {
	Next *Item32
	S0   Struct16B0
	S1   Struct16B1
}

type Item64 struct {
	Next *Item64
	S0   Struct16B0
	S1   Struct16B1
	S2   Struct16B2
	S3   Struct16B3
}

type Item128 struct {
	Next *Item128
	S0   Struct16B0
	S1   Struct16B1
	S2   Struct16B2
	S3   Struct16B3
	S4   Struct16B4
	S5   Struct16B5
	S6   Struct16B6
	S7   Struct16B7
}

type Item256 struct {
	Next *Item256
	S0   Struct16B0
	S1   Struct16B1
	S2   Struct16B2
	S3   Struct16B3
	S4   Struct16B4
	S5   Struct16B5
	S6   Struct16B6
	S7   Struct16B7
	S8   Struct16B8
	S9   Struct16B9
	S10  Struct16B10
	S11  Struct16B11
	S12  Struct16B12
	S13  Struct16B13
	S14  Struct16B14
	S15  Struct16B15
}

func runLL16B(b *testing.B, count int) {
	b.StopTimer()
	first := &Item16{}
	curr := first

	for i := 0; i < count-1; i++ {
		curr.Next = &Item16{}
		curr = curr.Next
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		curr = first
		for curr != nil {
			curr.S0.Hi++
			curr.S0.Lo++
			curr = curr.Next
		}
	}
}

func runLL32B(b *testing.B, count int) {
	b.StopTimer()
	first := &Item32{}
	curr := first

	for i := 0; i < count-1; i++ {
		curr.Next = &Item32{}
		curr = curr.Next
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		curr = first
		for curr != nil {
			curr.S0.Hi++
			curr.S0.Lo++
			curr = curr.Next
		}
	}
}

func runLL64B(b *testing.B, count int) {
	b.StopTimer()
	first := &Item64{}
	curr := first

	for i := 0; i < count-1; i++ {
		curr.Next = &Item64{}
		curr = curr.Next
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		curr = first
		for curr != nil {
			curr.S0.Hi++
			curr.S0.Lo++
			curr = curr.Next
		}
	}
}

func runLL128B(b *testing.B, count int) {
	b.StopTimer()
	first := &Item128{}
	curr := first

	for i := 0; i < count-1; i++ {
		curr.Next = &Item128{}
		curr = curr.Next
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		curr = first
		for curr != nil {
			curr.S0.Hi++
			curr.S0.Lo++
			curr = curr.Next
		}
	}
}

func runLL256B(b *testing.B, count int) {
	b.StopTimer()
	first := &Item256{}
	curr := first

	for i := 0; i < count-1; i++ {
		curr.Next = &Item256{}
		curr = curr.Next
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		curr = first
		for curr != nil {
			curr.S0.Hi++
			curr.S0.Lo++
			curr = curr.Next
		}
	}
}

func BenchmarkLinkedList_16B_1_000(b *testing.B) {
	runLL16B(b, 1000)
}

func BenchmarkLinkedList_16B_10_000(b *testing.B) {
	runLL16B(b, 10000)
}

func BenchmarkLinkedList_16B_100_000(b *testing.B) {
	runLL16B(b, 100000)
}

func BenchmarkLinkedList_32B_1_000(b *testing.B) {
	runLL32B(b, 1000)
}

func BenchmarkLinkedList_32B_10_000(b *testing.B) {
	runLL32B(b, 10000)
}

func BenchmarkLinkedList_32B_100_000(b *testing.B) {
	runLL32B(b, 100000)
}

func BenchmarkLinkedList_64B_1_000(b *testing.B) {
	runLL64B(b, 1000)
}

func BenchmarkLinkedList_64B_10_000(b *testing.B) {
	runLL64B(b, 10000)
}

func BenchmarkLinkedList_64B_100_000(b *testing.B) {
	runLL64B(b, 100000)
}

func BenchmarkLinkedList_128B_1_000(b *testing.B) {
	runLL128B(b, 1000)
}

func BenchmarkLinkedList_128B_10_000(b *testing.B) {
	runLL128B(b, 10000)
}

func BenchmarkLinkedList_128B_100_000(b *testing.B) {
	runLL128B(b, 100000)
}

func BenchmarkLinkedList_256B_1_000(b *testing.B) {
	runLL256B(b, 1000)
}

func BenchmarkLinkedList_256B_10_000(b *testing.B) {
	runLL256B(b, 10000)
}

func BenchmarkLinkedList_256B_100_000(b *testing.B) {
	runLL256B(b, 100000)
}
