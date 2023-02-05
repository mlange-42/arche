package ecs

import (
	"math"
	"reflect"
	"unsafe"
)

// Storage is the interface for component storage
type Storage interface {
	Get(index uint32) unsafe.Pointer
	Add(value interface{}) (index uint32)
	Remove(index uint32)
	Len() uint32
}

// ReflectStorage is a storage implementation that works with reflection
type ReflectStorage struct {
	buffer            reflect.Value
	bufferAddress     unsafe.Pointer
	typeOf            reflect.Type
	itemSize          uintptr
	len               uint32
	cap               uint32
	capacityIncrement uint32
}

// NewReflectStorage creates a new ReflectStorage
func NewReflectStorage(obj interface{}, increment int) *ReflectStorage {
	tp := reflect.TypeOf(obj)
	size := tp.Size()

	buffer := reflect.New(reflect.ArrayOf(increment, tp)).Elem()
	return &ReflectStorage{
		buffer:            buffer,
		bufferAddress:     buffer.Addr().UnsafePointer(),
		typeOf:            tp,
		itemSize:          size,
		len:               0,
		cap:               uint32(increment),
		capacityIncrement: uint32(increment),
	}
}

// Get retrieves an unsafe pointer to an element
func (s *ReflectStorage) Get(index uint32) unsafe.Pointer {
	ptr := unsafe.Add(s.bufferAddress, uintptr(index)*s.itemSize)
	return unsafe.Pointer(ptr)
}

// Add adds an element to the end of the storage
func (s *ReflectStorage) Add(value interface{}) (index uint32) {
	if s.cap < s.len+1 {
		old := s.buffer
		s.cap = s.cap + s.capacityIncrement
		s.buffer = reflect.New(reflect.ArrayOf(int(s.cap), s.typeOf)).Elem()
		s.bufferAddress = s.buffer.Addr().UnsafePointer()
		reflect.Copy(s.buffer, old)
	}
	s.len++
	s.set(s.len-1, value)
	return s.len - 1
}

// Remove swap-removes an element
func (s *ReflectStorage) Remove(index uint32) {
	o := s.len - 1
	n := index
	size := s.itemSize

	src := unsafe.Add(s.bufferAddress, uintptr(o)*s.itemSize)
	dst := unsafe.Add(s.bufferAddress, uintptr(n)*s.itemSize)

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(src)[:size:size]

	copy(dstSlice, srcSlice)

	s.len--
	// TODO shrink the underlying data arrays
}

func (s *ReflectStorage) set(index uint32, value interface{}) {
	rValue := reflect.ValueOf(value)
	dst := s.Get(index)
	var src unsafe.Pointer
	size := s.itemSize

	src = rValue.UnsafePointer()

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(src)[:size:size]

	copy(dstSlice, srcSlice)
}

// Len returns the number of items in the storage
func (s *ReflectStorage) Len() uint32 {
	return s.len
}

// ToSlice converts the content of a storage to a slice of structs
func ToSlice[T any](s Storage) []T {
	res := make([]T, s.Len())
	for i := 0; i < int(s.Len()); i++ {
		ptr := (*T)(s.Get(uint32(i)))
		res[i] = *ptr
	}
	return res
}
