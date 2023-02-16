package ecs

import (
	"math"
	"reflect"
	"unsafe"
)

// storage is a storage implementation that works with reflection
type storage struct {
	buffer            reflect.Value
	bufferAddress     unsafe.Pointer
	typeOf            reflect.Type
	itemSize          uintptr
	len               uint32
	cap               uint32
	capacityIncrement uint32
}

// Init initializes a storage
func (s *storage) Init(tp reflect.Type, increment int) {
	size := tp.Size()
	align := uintptr(tp.Align())
	size = (size + (align - 1)) / align * align

	s.buffer = reflect.New(reflect.ArrayOf(1, tp)).Elem()
	s.bufferAddress = s.buffer.Addr().UnsafePointer()
	s.typeOf = tp
	s.itemSize = size
	s.len = 0
	s.cap = 1
	s.capacityIncrement = uint32(increment)
}

// Get retrieves an unsafe pointer to an element
func (s *storage) Get(index uint32) unsafe.Pointer {
	ptr := unsafe.Add(s.bufferAddress, uintptr(index)*s.itemSize)
	return unsafe.Pointer(ptr)
}

// Add adds an element to the end of the storage
func (s *storage) Add(value interface{}) (index uint32) {
	s.extend()
	s.len++
	s.Set(s.len-1, value)
	return s.len - 1
}

// AddPointer adds an element to the end of the storage, based on a pointer
func (s *storage) AddPointer(value unsafe.Pointer) (index uint32) {
	s.extend()
	s.len++
	s.SetPointer(s.len-1, value)
	return s.len - 1
}

// Alloc adds an empty element to the end of the storage
func (s *storage) Alloc() (index uint32) {
	s.extend()
	s.len++
	s.Zero(s.len - 1)
	return s.len - 1
}

func (s *storage) extend() {
	if s.itemSize > 0 && s.cap < s.len+1 {
		old := s.buffer
		s.cap = s.capacityIncrement * ((s.cap + s.capacityIncrement) / s.capacityIncrement)
		s.buffer = reflect.New(reflect.ArrayOf(int(s.cap), s.typeOf)).Elem()
		s.bufferAddress = s.buffer.Addr().UnsafePointer()
		reflect.Copy(s.buffer, old)
	}
}

// Remove swap-removes an element
func (s *storage) Remove(index uint32) bool {
	o := s.len - 1
	n := index

	// TODO shrink the underlying data arrays
	if s.itemSize > 0 && n < o {
		size := s.itemSize

		src := unsafe.Add(s.bufferAddress, uintptr(o)*s.itemSize)
		dst := unsafe.Add(s.bufferAddress, uintptr(n)*s.itemSize)

		dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
		srcSlice := (*[math.MaxInt32]byte)(src)[:size:size]

		copy(dstSlice, srcSlice)

		s.len--
		return true
	}

	s.len--
	return false
}

// Set sets the storage at the given index
func (s *storage) Set(index uint32, value interface{}) unsafe.Pointer {
	rValue := reflect.ValueOf(value)
	dst := s.Get(index)

	if s.itemSize == 0 {
		return dst
	}

	var src unsafe.Pointer
	size := s.itemSize

	src = rValue.UnsafePointer()

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(src)[:size:size]

	copy(dstSlice, srcSlice)
	return dst
}

func (s *storage) SetPointer(index uint32, value unsafe.Pointer) unsafe.Pointer {
	if s.itemSize == 0 {
		return s.Get(index)
	}

	dst := s.Get(index)
	size := s.itemSize

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(value)[:size:size]

	copy(dstSlice, srcSlice)

	return dst
}

// Zero resets a block of storage
func (s *storage) Zero(index uint32) {
	if s.itemSize == 0 {
		return
	}

	dst := s.Get(index)

	for i := uintptr(0); i < s.itemSize; i++ {
		*(*byte)(dst) = 0
		dst = unsafe.Add(dst, 1)
	}
}

// Len returns the number of items in the storage
func (s *storage) Len() uint32 {
	return s.len
}

// Cap returns the capacity of the storage
func (s *storage) Cap() uint32 {
	return s.cap
}

// toSlice converts the content of a storage to a slice of structs
func toSlice[T any](s storage) []T {
	res := make([]T, s.Len())
	var i uint32
	for i = 0; i < s.Len(); i++ {
		ptr := (*T)(s.Get(i))
		res[i] = *ptr
	}
	return res
}
