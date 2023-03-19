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
	size := tp.Size()
	align := uintptr(tp.Align())
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
	//s.Zero(s.len - 1)
	return uintptr(s.len - 1)
}

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

	// TODO shrink the underlying data arrays?
	size := s.itemSize

	src := unsafe.Add(s.bufferAddress, o*s.itemSize)
	dst := unsafe.Add(s.bufferAddress, n*s.itemSize)

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(src)[:size:size]

	copy(dstSlice, srcSlice)

	s.len--
	return true
}

// Set sets the storage at the given index
func (s *storage) Set(index uintptr, value interface{}) unsafe.Pointer {
	dst := s.Get(index)

	if s.itemSize == 0 {
		return dst
	}
	rValue := reflect.ValueOf(value)

	var src unsafe.Pointer
	size := s.itemSize

	src = rValue.UnsafePointer()

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(src)[:size:size]

	copy(dstSlice, srcSlice)
	return dst
}

func (s *storage) SetPointer(index uintptr, value unsafe.Pointer) unsafe.Pointer {
	dst := s.Get(index)
	if s.itemSize == 0 {
		return dst
	}

	size := s.itemSize

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(value)[:size:size]

	copy(dstSlice, srcSlice)

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

type genericStorage[T any] struct {
	buffer            []T
	capacityIncrement uint32
}

// Init initializes a genericStorage
func (s *genericStorage[T]) Init(increment int, forStorage bool) {
	cap := 1
	if forStorage {
		cap = increment
	}
	s.buffer = make([]T, 0, cap)
	s.capacityIncrement = uint32(increment)
}

// Get retrieves an element
func (s *genericStorage[T]) Get(index uint32) T {
	return s.buffer[index]
}

// Set sets an element
func (s *genericStorage[T]) Set(index uint32, value T) {
	s.buffer[index] = value
}

// Add adds an element to the end of the genericStorage
func (s *genericStorage[T]) Add(value T) (index uint32) {
	s.extend()
	s.buffer = append(s.buffer, value)
	return s.Len() - 1
}

func (s *genericStorage[T]) extend() {
	currLen := len(s.buffer)
	currCap := cap(s.buffer)
	if currCap > currLen {
		return
	}

	inc := int(s.capacityIncrement)
	old := s.buffer
	cap := inc * ((currCap + inc) / inc)
	s.buffer = make([]T, currLen, cap)
	copy(s.buffer, old)
}

// Remove swap-removes an element
func (s *genericStorage[T]) Remove(index uint32) bool {
	o := len(s.buffer) - 1
	n := int(index)

	if n == o {
		s.buffer = s.buffer[:o]
		return false
	}

	s.buffer[n] = s.buffer[o]
	s.buffer = s.buffer[:o]
	return true
}

// Len returns the number of items in the genericStorage
func (s *genericStorage[T]) Len() uint32 {
	return uint32(len(s.buffer))
}

// Cap returns the capacity of the genericStorage
func (s *genericStorage[T]) Cap() uint32 {
	return uint32(cap(s.buffer))
}
