package ecs

import (
	"reflect"
	"unsafe"
)

// Storage is the interface for component storage
type Storage interface {
	Get(index uint32) unsafe.Pointer
	Add(value interface{}) (index uint32)
	Remove(index uint32) (old, new uint32)
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

func (s *storage) Get(index uint32) unsafe.Pointer {
	if index >= s.len {
		return nil
	}
	base := unsafe.Pointer(&s.data[0])

	offset := uintptr(index) * s.itemSize
	return unsafe.Add(base, offset)
}
func (s *storage) Add(value interface{}) (index uint32) {
	// TODO this allocates a new slice and should be improved
	s.data = append(s.data, make([]byte, s.itemSize)...)
	s.len++
	s.Set(s.len-1, value)
	return s.len - 1
}
func (s *storage) Remove(index uint32) (old, new uint32) {
	return 0, 0
}

func (s *storage) Set(index uint32, value interface{}) {
	dst := s.Get(index)
	src := unsafe.Pointer(reflect.ValueOf(value).Pointer())
	for i := uintptr(0); i < s.itemSize; i++ {
		*(*byte)(dst) = *(*byte)(src)
		dst = unsafe.Add(dst, 1)
		src = unsafe.Add(src, 1)
	}
}

func (s *storage) Zero(index uint32) {
	dst := s.Get(index)
	for i := uintptr(0); i < s.itemSize; i++ {
		*(*byte)(dst) = 0
		dst = unsafe.Add(dst, 1)
	}
}

func (s *storage) Swap(i, j uint32) {
	ii := uintptr(i) * s.itemSize
	jj := uintptr(j) * s.itemSize
	for k := uintptr(0); k < s.itemSize; k++ {
		s.data[ii+k], s.data[jj+k] = s.data[jj+k], s.data[ii+k]
	}
}

// Len returns the number of items in the storage
func (s *storage) Len() uint32 {
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
