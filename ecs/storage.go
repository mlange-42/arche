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
	itemSize          uintptr
	len               uint32
	cap               uint32
	capacityIncrement uint32
}

// Init initializes a storage
func (s *storage) Init(tp reflect.Type, increment int, forStorage bool) {
	size, align := tp.Size(), uintptr(tp.Align())
	size = (size + (align - 1)) / align * align

	cap := 1
	if forStorage {
		cap = increment
	}
	s.buffer = reflect.New(reflect.ArrayOf(cap, tp)).Elem()
	s.bufferAddress = s.buffer.Addr().UnsafePointer()
	s.itemSize = size
	s.len = 0
	s.cap = uint32(cap)
	s.capacityIncrement = uint32(increment)
}

// Get retrieves an unsafe pointer to an element
func (s *storage) Get(index uintptr) unsafe.Pointer {
	if s == nil {
		return nil
	}
	return unsafe.Add(s.bufferAddress, index*s.itemSize)
}

// Add adds an element to the end of the storage
func (s *storage) Add(value interface{}) (index uint32) {
	s.extend()
	s.len++
	s.Set(uintptr(s.len-1), value)
	return s.len - 1
}

// AddPointer adds an element to the end of the storage, based on a pointer
func (s *storage) AddPointer(value unsafe.Pointer) (index uint32) {
	s.extend()
	s.len++
	s.SetPointer(uintptr(s.len-1), value)
	return s.len - 1
}

// Alloc adds an empty element to the end of the storage.
// It does not zero the storage!
func (s *storage) Alloc() (index uintptr) {
	s.extend()
	s.len++
	return uintptr(s.len - 1)
}

// extend the storage's capacity by capacityIncrement.
//
// Extends to a multiple of capacityIncrement.
func (s *storage) extend() {
	if s.cap > s.len || s.itemSize == 0 {
		return
	}

	old := s.buffer
	s.cap = s.capacityIncrement * ((s.cap + s.capacityIncrement) / s.capacityIncrement)
	s.buffer = reflect.New(reflect.ArrayOf(int(s.cap), s.buffer.Type().Elem())).Elem()
	s.bufferAddress = s.buffer.Addr().UnsafePointer()
	reflect.Copy(s.buffer, old)
}

// Remove swap-removes an element
func (s *storage) Remove(index uintptr) bool {
	o := uintptr(s.len - 1)
	n := uintptr(index)

	if n == o || s.itemSize == 0 {
		s.len--
		return false
	}

	src := unsafe.Add(s.bufferAddress, o*s.itemSize)
	dst := unsafe.Add(s.bufferAddress, n*s.itemSize)
	s.copy(src, dst, s.itemSize)

	s.len--
	return true
}

// Set sets the storage at the given index.
func (s *storage) Set(index uintptr, value interface{}) unsafe.Pointer {
	dst := s.Get(index)

	if s.itemSize == 0 {
		return dst
	}
	rValue := reflect.ValueOf(value)

	src := rValue.UnsafePointer()
	s.copy(src, dst, s.itemSize)

	return dst
}

// SetPointer sets the storage at the given index from the data behind an unsafe pointer.
func (s *storage) SetPointer(index uintptr, value unsafe.Pointer) unsafe.Pointer {
	dst := s.Get(index)
	if s.itemSize == 0 {
		return dst
	}

	s.copy(value, dst, s.itemSize)

	return dst
}

// Zero resets a block of storage
func (s *storage) Zero(index uintptr) {
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

// copy from one pointer to another.
func (s *storage) copy(src, dst unsafe.Pointer, itemSize uintptr) {
	dstSlice := (*[math.MaxInt32]byte)(dst)[:itemSize:itemSize]
	srcSlice := (*[math.MaxInt32]byte)(src)[:itemSize:itemSize]
	copy(dstSlice, srcSlice)
}

// toSlice converts the content of a storage to a slice of structs
func toSlice[T any](s storage) []T {
	res := make([]T, s.Len())
	var i uintptr
	len := uintptr(s.Len())
	for i = 0; i < len; i++ {
		ptr := (*T)(s.Get(i))
		res[i] = *ptr
	}
	return res
}
