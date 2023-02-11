package base

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Int32 int32
	Int64 int64
	Bool1 bool
	Bool2 bool
}

type simpleStruct struct {
	Index int
}

func TestReflectStorageAddGet(t *testing.T) {
	obj1 := testStruct{}
	s := Storage{}
	s.Init(reflect.TypeOf(obj1), 1)
	storageAddGet(t, s)
}

func storageAddGet(t *testing.T, s Storage) {
	obj1 := testStruct{}
	obj2 := testStruct{1, 2, true, false}

	idx := s.Add(&obj1)
	assert.Equal(t, idx, uint32(0), "Index of first insertion should be 0")

	idx = s.Add(&obj2)
	assert.Equal(t, idx, uint32(1), "Index of second insertion should be 1")

	ret1 := (*testStruct)(s.Get(0))
	assert.Equal(t, obj1, *ret1, "First element not as passed in")

	ret2 := (*testStruct)(s.Get(1))
	assert.Equal(t, obj2, *ret2, "Second element not as passed in")

	ret2.Int64 = 1001
	ret2 = (*testStruct)(s.Get(1))
	assert.Equal(t, testStruct{1, 1001, true, false}, *ret2, "Manipulating element does not change data")

	assert.Equal(t, []testStruct{{}, {1, 1001, true, false}}, toSlice[testStruct](s), "Wrong extracted struct slice")

	s.Alloc()
	assert.Equal(t, []testStruct{{}, {1, 1001, true, false}, {}}, toSlice[testStruct](s), "Wrong extracted struct slice")
}

func TestReflectStorageRemove(t *testing.T) {
	ref := simpleStruct{}
	s := Storage{}
	s.Init(reflect.TypeOf(ref), 32)

	storageRemove(t, s)
}

func storageRemove(t *testing.T, s Storage) {
	for i := 0; i < 5; i++ {
		obj := simpleStruct{i}
		s.Add(&obj)
	}

	assert.Equal(t, uint32(5), s.Len(), "Wrong storage length")

	s.Remove(4)
	assert.Equal(t, uint32(4), s.Len(), "Wrong storage length")
	assert.Equal(t, []simpleStruct{{0}, {1}, {2}, {3}}, toSlice[simpleStruct](s), "Wrong slice after remove")

	s.Remove(1)
	assert.Equal(t, uint32(3), s.Len(), "Wrong storage length")
	assert.Equal(t, []simpleStruct{{0}, {3}, {2}}, toSlice[simpleStruct](s), "Wrong slice after remove")
}

func TestReflectStorageDataSize(t *testing.T) {
	ref := simpleStruct{}
	s := Storage{}
	s.Init(reflect.TypeOf(ref), 1)

	for i := 0; i < 5; i++ {
		obj := simpleStruct{i}
		s.Add(&obj)
	}

	assert.Equal(t, 5, int(s.cap))

	s.Remove(0)
	assert.Equal(t, 5, int(s.cap))
	s.Remove(0)
	assert.Equal(t, 5, int(s.cap))

	s.Add(&simpleStruct{})
	assert.Equal(t, 5, int(s.cap))
	s.Add(&simpleStruct{})
	assert.Equal(t, 5, int(s.cap))

	s.Add(&simpleStruct{})
	assert.Equal(t, 6, int(s.cap))
}

func TestNewReflectStorage(t *testing.T) {
	s := Storage{}
	s.Init(reflect.TypeOf(simpleStruct{}), 32)
}

func BenchmarkIterReflectStorage_1000(b *testing.B) {
	b.StopTimer()
	ref := testStruct{}
	s := Storage{}
	s.Init(reflect.TypeOf(ref), 128)
	for i := 0; i < 1000; i++ {
		s.Add(&testStruct{})
	}
	assert.Equal(b, 1000, int(s.Len()))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		len := int(s.Len())
		for j := 0; j < len; j++ {
			a := (*testStruct)(s.Get(uint32(j)))
			_ = a
		}
	}
}

func BenchmarkIterSlice_1000(b *testing.B) {
	b.StopTimer()
	s := []testStruct{}
	for i := 0; i < 1000; i++ {
		s = append(s, testStruct{})
	}
	assert.Equal(b, 1000, len(s))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			a := s[j]
			_ = a
		}
	}
}

func BenchmarkIterSliceInterface_1000(b *testing.B) {
	b.StopTimer()
	s := []interface{}{}
	for i := 0; i < 1000; i++ {
		s = append(s, testStruct{})
	}
	assert.Equal(b, 1000, len(s))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			a := s[j].(testStruct)
			_ = a
		}
	}
}

func BenchmarkAddReflectStorage_1000(b *testing.B) {
	ref := testStruct{}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := Storage{}
		s.Init(reflect.TypeOf(ref), 1024)
		b.StartTimer()

		for i := 0; i < 1000; i++ {
			s.Add(&ref)
		}
	}
}

func BenchmarkAddSlice_1000(b *testing.B) {
	ref := testStruct{}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := make([]testStruct, 0)
		b.StartTimer()

		for i := 0; i < 1000; i++ {
			s = append(s, ref)
		}

		b.StopTimer()
		_ = s
		b.StartTimer()
	}
}

func BenchmarkAddSliceInterface_1000(b *testing.B) {
	ref := testStruct{}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := make([]interface{}, 0)
		b.StartTimer()

		for i := 0; i < 1000; i++ {
			s = append(s, ref)
		}

		b.StopTimer()
		_ = s
		b.StartTimer()
	}
}

func BenchmarkRemoveReflectStorage_1000(b *testing.B) {
	ref := testStruct{}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := Storage{}
		s.Init(reflect.TypeOf(ref), 1024)
		for i := 0; i < 1000; i++ {
			s.Add(&ref)
		}
		b.StartTimer()

		for i := 0; i < 1000; i++ {
			s.Remove(0)
		}
	}
}

func BenchmarkRemoveSlice_1000(b *testing.B) {
	b.StopTimer()
	ref := testStruct{}
	template := make([]testStruct, 0)
	for i := 0; i < 1000; i++ {
		template = append(template, ref)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := append([]testStruct{}, template...)
		b.StartTimer()

		l := len(s) - 1
		for i := 0; i < 1000; i++ {
			s[0], s[l] = s[l], s[0]
			s = s[:l]
			l--
		}

		b.StopTimer()
		_ = s
		b.StartTimer()
	}
}

func BenchmarkRemoveSliceInterface_1000(b *testing.B) {
	b.StopTimer()
	ref := testStruct{}
	template := make([]interface{}, 0)
	for i := 0; i < 1000; i++ {
		template = append(template, ref)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := append([]interface{}{}, template...)
		b.StartTimer()

		l := len(s) - 1
		for i := 0; i < 1000; i++ {
			s[0], s[l] = s[l], s[0]
			s = s[:l]
			l--
		}

		b.StopTimer()
		_ = s
		b.StartTimer()
	}
}
