package ecs

import (
	"fmt"
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

// NewStorage creates a new component storage
func NewStorage(obj interface{}) Storage {
	return newStorage(obj)
}

func newStorage(obj interface{}) *storage {
	tp := reflect.TypeOf(obj)
	size := tp.Size()

	return &storage{
		data:     []byte{},
		itemSize: size,
		len:      0,
	}
}

type storage struct {
	data     []byte
	itemSize uintptr
	len      uint32
}

// Get retrieves an unsafe pointer to an element
func (s *storage) Get(index uint32) unsafe.Pointer {
	if index >= s.len {
		panic(fmt.Sprintf("Index %d out of range %d", index, s.len))
	}
	base := unsafe.Pointer(&s.data[0])

	offset := uintptr(index) * s.itemSize
	return unsafe.Add(base, offset)
}

// Add adds an element to the end of the storage
func (s *storage) Add(value interface{}) (index uint32) {
	// TODO this allocates a new slice and should be improved
	s.data = append(s.data, make([]byte, s.itemSize)...)
	s.len++
	s.set(s.len-1, value)
	return s.len - 1
}

// Remove swap-removes an element
func (s *storage) Remove(index uint32) {
	if index >= s.len {
		panic(fmt.Sprintf("Index %d out of range %d", index, s.len))
	}
	o := s.len - 1
	n := index
	s.swap(o, n)
	s.len--
}

// Len returns the number of items in the storage
func (s *storage) Len() uint32 {
	return s.len
}

func (s *storage) set(index uint32, value interface{}) {
	dst := s.Get(index)
	src := unsafe.Pointer(reflect.ValueOf(value).Pointer())
	for i := uintptr(0); i < s.itemSize; i++ {
		*(*byte)(dst) = *(*byte)(src)
		dst = unsafe.Add(dst, 1)
		src = unsafe.Add(src, 1)
	}
}

func (s *storage) swap(i, j uint32) {
	if i == j {
		return
	}
	ii := uintptr(i) * s.itemSize
	jj := uintptr(j) * s.itemSize
	for k := uintptr(0); k < s.itemSize; k++ {
		s.data[ii+k], s.data[jj+k] = s.data[jj+k], s.data[ii+k]
	}
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
