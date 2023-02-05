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

func TestStorageItemSize(t *testing.T) {
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

	s1 := newStorage(obj1)
	s2 := newStorage(obj2)
	s3 := newStorage(obj3)
	s4 := newStorage(obj4)
	s5 := newStorage(obj5)
	s6 := newStorage(obj6)
	s7 := newStorage(obj7)

	assert.Equal(
		t,
		[]uintptr{0, 1, 1, 2, 4, 8, 16},
		[]uintptr{s1.itemSize, s2.itemSize, s3.itemSize, s4.itemSize, s5.itemSize, s6.itemSize, s7.itemSize},
		"Unexpected struct size",
	)
}

func TestStorageAddGet(t *testing.T) {
	obj1 := testStruct{}
	obj2 := testStruct{1, 2, true, false}
	s := newStorage(obj1)

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
