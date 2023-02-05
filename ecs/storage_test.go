package ecs

import (
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

func TestByteStorageItemSize(t *testing.T) {
	obj1 := struct{}{}
	obj2 := struct{ bool }{true}
	obj3 := struct{ int8 }{0}
	obj4 := struct{ int16 }{0}
	obj5 := struct{ int32 }{0}
	obj6 := struct{ int64 }{0}
	obj7 := struct {
		int64
		bool
	}{0, true}

	s1 := NewByteStorage(obj1)
	s2 := NewByteStorage(obj2)
	s3 := NewByteStorage(obj3)
	s4 := NewByteStorage(obj4)
	s5 := NewByteStorage(obj5)
	s6 := NewByteStorage(obj6)
	s7 := NewByteStorage(obj7)

	assert.Equal(
		t,
		[]uintptr{0, 1, 1, 2, 4, 8, 16},
		[]uintptr{s1.itemSize, s2.itemSize, s3.itemSize, s4.itemSize, s5.itemSize, s6.itemSize, s7.itemSize},
		"Unexpected struct size",
	)
}

func TestByteStorageAddGet(t *testing.T) {
	obj1 := testStruct{}
	obj2 := testStruct{1, 2, true, false}
	s := NewByteStorage(obj1)

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

	assert.Equal(t, []testStruct{{}, {1, 1001, true, false}}, ToSlice[testStruct](s), "Wrong extracted struct slice")
}

func TestByteStorageRemove(t *testing.T) {
	ref := simpleStruct{}
	s := NewByteStorage(ref)

	for i := 0; i < 5; i++ {
		obj := simpleStruct{i}
		s.Add(&obj)
	}

	assert.Equal(t, uint32(5), s.Len(), "Wrong storage length")

	s.Remove(4)
	assert.Equal(t, uint32(4), s.Len(), "Wrong storage length")
	assert.Equal(t, []simpleStruct{{0}, {1}, {2}, {3}}, ToSlice[simpleStruct](s), "Wrong slice after remove")

	s.Remove(1)
	assert.Equal(t, uint32(3), s.Len(), "Wrong storage length")
	assert.Equal(t, []simpleStruct{{0}, {3}, {2}}, ToSlice[simpleStruct](s), "Wrong slice after remove")
}

func TestByteStorageDataSize(t *testing.T) {
	ref := simpleStruct{}
	s := NewByteStorage(ref)
	s.capacityIncrement = 1

	size := int(s.itemSize)

	for i := 0; i < 5; i++ {
		obj := simpleStruct{i}
		s.Add(&obj)
	}

	assert.Equal(t, 5*size, len(s.data))

	s.Remove(0)
	assert.Equal(t, 5*size, len(s.data))
	s.Remove(0)
	assert.Equal(t, 5*size, len(s.data))

	s.Add(&simpleStruct{})
	assert.Equal(t, 5*size, len(s.data))
	s.Add(&simpleStruct{})
	assert.Equal(t, 5*size, len(s.data))

	s.Add(&simpleStruct{})
	assert.Equal(t, 6*size, len(s.data))
}

func TestNewByteStorage(t *testing.T) {
	_ = NewByteStorage(simpleStruct{})
}

func BenchmarkIterByteStorage(b *testing.B) {
	ref := testStruct{}
	s := NewByteStorage(ref)
	for i := 0; i < 1000; i++ {
		s.Add(&testStruct{})
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < int(s.Len()); j++ {
			a := (*testStruct)(s.Get(uint32(j)))
			_ = a
		}
	}
}

func BenchmarkIterSlice(b *testing.B) {
	s := []testStruct{}
	for i := 0; i < 1000; i++ {
		s = append(s, testStruct{})
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			a := s[j]
			_ = a
		}
	}
}

func BenchmarkIterSliceInterface(b *testing.B) {
	s := []interface{}{}
	for i := 0; i < 1000; i++ {
		s = append(s, testStruct{})
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			a := s[j].(testStruct)
			_ = a
		}
	}
}
