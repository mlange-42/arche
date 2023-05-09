package ecs_test

import "testing"

func BenchmarkSimpleIterUInt32_1000(b *testing.B) {
	var iMax uint32 = 1000
	var sum uint32 = 0
	for i := 0; i < b.N; i++ {
		var i uint32 = 0
		for i < iMax {
			i++
			sum += i
		}
	}
}

func BenchmarkSimpleIterUInt32_Convert_1000(b *testing.B) {
	var iMax uint32 = 1000
	var sum uintptr = 0
	for i := 0; i < b.N; i++ {
		var i uint32 = 0
		for i < iMax {
			i++
			sum += uintptr(i)
		}
	}
}

func BenchmarkSimpleIterUintptr_1000(b *testing.B) {
	var iMax uintptr = 1000
	var sum uintptr = 0
	for i := 0; i < b.N; i++ {
		var i uintptr = 0
		for i < iMax {
			i++
			sum += i
		}
	}
}

func BenchmarkMultiplyUInt32_1000(b *testing.B) {
	var iMax uint32 = 1000
	var sz uint32 = 4
	var sum uintptr = 0
	for i := 0; i < b.N; i++ {
		var i uint32 = 0
		for i < iMax {
			i++
			sum += uintptr(i * sz)
		}
	}
}

func BenchmarkMultiplyUIntptr_1000(b *testing.B) {
	var iMax uintptr = 1000
	var sz uintptr = 4
	var sum uintptr = 0
	for i := 0; i < b.N; i++ {
		var i uintptr = 0
		for i < iMax {
			i++
			sum += i * sz
		}
	}
}
