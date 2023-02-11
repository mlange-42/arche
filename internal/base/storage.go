package base

import (
	"math"
	"reflect"
	"unsafe"
)

// Storage is a Storage implementation that works with reflection
type Storage struct {
	buffer            reflect.Value
	bufferAddress     unsafe.Pointer
	typeOf            reflect.Type
	itemSize          uintptr
	len               uint32
	cap               uint32
	capacityIncrement uint32
}

// newStorage creates a new ReflectStorage
func (s *Storage) init(tp reflect.Type, increment int) {
	size := tp.Size()
	align := uintptr(tp.Align())
	size = (size + (align - 1)) / align * align

	s.buffer = reflect.New(reflect.ArrayOf(increment, tp)).Elem()
	s.bufferAddress = s.buffer.Addr().UnsafePointer()
	s.typeOf = tp
	s.itemSize = size
	s.len = 0
	s.cap = uint32(increment)
	s.capacityIncrement = uint32(increment)
}

// Get retrieves an unsafe pointer to an element
func (s *Storage) Get(index uint32) unsafe.Pointer {
	ptr := unsafe.Add(s.bufferAddress, uintptr(index)*s.itemSize)
	return unsafe.Pointer(ptr)
}

// Add adds an element to the end of the storage
func (s *Storage) Add(value interface{}) (index uint32) {
	s.extend()
	s.len++
	s.set(s.len-1, value)
	return s.len - 1
}

// AddPointer adds an element to the end of the storage, based on a pointer
func (s *Storage) AddPointer(value unsafe.Pointer) (index uint32) {
	s.extend()
	s.len++
	s.setPointer(s.len-1, value)
	return s.len - 1
}

// Alloc adds an empty element to the end of the storage
func (s *Storage) Alloc() (index uint32) {
	s.extend()
	s.len++
	s.Zero(s.len - 1)
	return s.len - 1
}

func (s *Storage) extend() {
	if s.cap < s.len+1 {
		old := s.buffer
		s.cap = s.cap + s.capacityIncrement
		s.buffer = reflect.New(reflect.ArrayOf(int(s.cap), s.typeOf)).Elem()
		s.bufferAddress = s.buffer.Addr().UnsafePointer()
		reflect.Copy(s.buffer, old)
	}
}

// Remove swap-removes an element
func (s *Storage) Remove(index uint32) bool {
	o := s.len - 1
	n := index

	// TODO shrink the underlying data arrays
	if n < o {
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

func (s *Storage) set(index uint32, value interface{}) unsafe.Pointer {
	rValue := reflect.ValueOf(value)
	dst := s.Get(index)
	var src unsafe.Pointer
	size := s.itemSize

	src = rValue.UnsafePointer()

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(src)[:size:size]

	copy(dstSlice, srcSlice)
	return dst
}

func (s *Storage) setPointer(index uint32, value unsafe.Pointer) {
	dst := s.Get(index)
	size := s.itemSize

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(value)[:size:size]

	copy(dstSlice, srcSlice)
}

// Zero resets a block of storage
func (s *Storage) Zero(index uint32) {
	dst := s.Get(index)

	for i := uintptr(0); i < s.itemSize; i++ {
		*(*byte)(dst) = 0
		dst = unsafe.Add(dst, 1)
	}
}

// Len returns the number of items in the storage
func (s *Storage) Len() uint32 {
	return s.len
}

// toSlice converts the content of a storage to a slice of structs
func toSlice[T any](s Storage) []T {
	res := make([]T, s.Len())
	for i := 0; i < int(s.Len()); i++ {
		ptr := (*T)(s.Get(uint32(i)))
		res[i] = *ptr
	}
	return res
}
