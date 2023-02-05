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
func NewReflectStorage(obj interface{}) *ReflectStorage {
	incr := 32

	tp := reflect.TypeOf(obj)
	size := tp.Size()

	buffer := reflect.New(reflect.ArrayOf(0, tp)).Elem()
	return &ReflectStorage{
		buffer:            buffer,
		bufferAddress:     buffer.Addr().UnsafePointer(),
		typeOf:            tp,
		itemSize:          size,
		len:               0,
		cap:               0,
		capacityIncrement: uint32(incr),
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
		newCap := s.cap + s.capacityIncrement

		old := s.buffer
		oldAddr := s.bufferAddress

		buffer := reflect.New(reflect.ArrayOf(int(newCap), s.typeOf)).Elem()
		newAddr := buffer.Addr().UnsafePointer()

		if s.len > 0 {
			dstSlice := (*[math.MaxInt32]byte)(newAddr)[: s.itemSize : s.itemSize*uintptr(s.len)]
			srcSlice := (*[math.MaxInt32]byte)(oldAddr)[: s.itemSize : s.itemSize*uintptr(s.len)]

			copy(dstSlice, srcSlice)
		}

		s.cap = newCap
		s.buffer = buffer
		s.bufferAddress = newAddr
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
