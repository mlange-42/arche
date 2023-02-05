package ecs

import (
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

// NewByteStorage creates a new ByteStorage
func NewByteStorage(obj interface{}) *ByteStorage {
	tp := reflect.TypeOf(obj)
	size := tp.Size()

	return &ByteStorage{
		data:              []byte{},
		itemSize:          size,
		len:               0,
		capacityIncrement: 32,
	}
}

// ByteStorage stores components in a byte slice
// It does not work with components that contain references or pointers,
// because the referenced memory will be garbage-collected.
type ByteStorage struct {
	data              []byte
	itemSize          uintptr
	len               uint32
	capacityIncrement uint32
}

// Get retrieves an unsafe pointer to an element
func (s *ByteStorage) Get(index uint32) unsafe.Pointer {
	base := unsafe.Pointer(&s.data[0])

	offset := uintptr(index) * s.itemSize
	return unsafe.Add(base, offset)
}

// Add adds an element to the end of the storage
func (s *ByteStorage) Add(value interface{}) (index uint32) {
	// TODO this allocates a new slice and should be improved
	if uint32(len(s.data)) < (s.len+1)*uint32(s.itemSize) {
		old := s.data
		s.data = make([]byte, (s.len+s.capacityIncrement)*uint32(s.itemSize))
		copy(s.data, old)
	}
	s.len++
	s.set(s.len-1, value)
	return s.len - 1
}

// Remove swap-removes an element
func (s *ByteStorage) Remove(index uint32) {
	o := s.len - 1
	n := index
	s.copy(o, n)
	s.len--
}

// Len returns the number of items in the storage
func (s *ByteStorage) Len() uint32 {
	return s.len
}

func (s *ByteStorage) set(index uint32, value interface{}) {
	dst := s.Get(index)
	src := unsafe.Pointer(reflect.ValueOf(value).Pointer())
	for i := uintptr(0); i < s.itemSize; i++ {
		*(*byte)(dst) = *(*byte)(src)
		dst = unsafe.Add(dst, 1)
		src = unsafe.Add(src, 1)
	}
}

func (s *ByteStorage) copy(from, to uint32) {
	if from == to {
		return
	}
	f := uintptr(from) * s.itemSize
	t := uintptr(to) * s.itemSize
	copy(s.data[t:t+s.itemSize], s.data[f:f+s.itemSize])
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
