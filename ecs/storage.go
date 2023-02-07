package ecs

import (
	"math"
	"reflect"
	"unsafe"
)

type storage struct {
	typeOf            reflect.Type
	itemSize          uintptr
	capacityExponent  int
	capacityIncrement uint32
	len               uint32
	pages             []storagePage
}

// newStorage creates a new storage
func newStorage(tp reflect.Type, increment int) storage {
	capacityExponent := math.Ilogb(float64(increment))
	if 1<<capacityExponent != increment {
		panic("capacity increment must be a power of 2 value")
	}

	size := tp.Size()
	align := uintptr(tp.Align())
	size = (size + (align - 1)) / align * align

	return storage{
		typeOf:            tp,
		itemSize:          size,
		capacityExponent:  capacityExponent,
		capacityIncrement: uint32(increment),
		pages:             make([]storagePage, 0),
	}
}

// Get retrieves an unsafe pointer to an element
func (s *storage) Get(index uint32) unsafe.Pointer {
	if index >= s.len {
		return nil
	}
	ind := index >> s.capacityExponent
	offset := index & (s.capacityIncrement - 1)
	ptr := unsafe.Add(s.pages[ind].bufferAddress, uintptr(offset)*s.itemSize)
	return unsafe.Pointer(ptr)
}

func (s *storage) page(index uint32) uint32 {
	return index >> s.capacityExponent
}

// Alloc adds an empty element to the end of the storage
func (s *storage) Alloc() (index uint32) {
	s.extend()
	s.incrLen()
	s.Zero(s.len - 1)
	return s.len - 1
}

// Add adds an element to the end of the storage
func (s *storage) Add(value interface{}) (index uint32) {
	s.extend()
	s.incrLen()
	s.set(s.len-1, value)
	return s.len - 1
}

// AddPointer adds an element to the end of the storage, based on a pointer
func (s *storage) AddPointer(value unsafe.Pointer) (index uint32) {
	s.extend()
	s.incrLen()
	s.setPointer(s.len-1, value)
	return s.len - 1
}

// Zero resets a block of storage
func (s *storage) Zero(index uint32) {
	dst := s.Get(index)

	for i := uintptr(0); i < s.itemSize; i++ {
		*(*byte)(dst) = 0
		dst = unsafe.Add(dst, 1)
	}
}

// Remove swap-removes an element
func (s *storage) Remove(index uint32) bool {
	o := s.len - 1
	n := index

	// TODO shrink the underlying data arrays
	if n < o {
		size := s.itemSize

		src := s.Get(o)
		dst := s.Get(n)

		dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
		srcSlice := (*[math.MaxInt32]byte)(src)[:size:size]

		copy(dstSlice, srcSlice)

		s.decrLen()
		return true
	}

	s.decrLen()
	return false
}

func (s *storage) incrLen() {
	s.pages[s.page(s.len)].len++
	s.len++
}

func (s *storage) decrLen() {
	s.len--
	s.pages[s.page(s.len)].len--
}

// Len returns the number of items in the storage
func (s *storage) Len() int {
	return int(s.len)
}

// Len returns the current capacity of the storage
func (s *storage) Cap() int {
	return len(s.pages) * int(s.capacityIncrement)
}

func (s *storage) set(index uint32, value interface{}) {
	rValue := reflect.ValueOf(value)
	dst := s.Get(index)
	var src unsafe.Pointer
	size := s.itemSize

	src = rValue.UnsafePointer()

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(src)[:size:size]

	copy(dstSlice, srcSlice)
}

func (s *storage) setPointer(index uint32, value unsafe.Pointer) {
	dst := s.Get(index)
	size := s.itemSize

	dstSlice := (*[math.MaxInt32]byte)(dst)[:size:size]
	srcSlice := (*[math.MaxInt32]byte)(value)[:size:size]

	copy(dstSlice, srcSlice)
}

func (s *storage) extend() {
	if len(s.pages) == 0 || s.pages[len(s.pages)-1].len == s.capacityIncrement {
		s.newStoragePage()
	}
}

func (s *storage) newStoragePage() {
	buffer := reflect.New(reflect.ArrayOf(int(s.capacityIncrement), s.typeOf)).Elem()
	s.pages = append(s.pages, storagePage{
		buffer:        buffer,
		bufferAddress: buffer.Addr().UnsafePointer(),
		len:           0,
	})
}

// storage is a storage implementation that works with reflection
type storagePage struct {
	buffer        reflect.Value
	bufferAddress unsafe.Pointer
	len           uint32
}

// toSlice converts the content of a storage to a slice of structs
func toSlice[T any](s storage) []T {
	res := make([]T, s.Len())
	for i := 0; i < int(s.Len()); i++ {
		ptr := (*T)(s.Get(uint32(i)))
		res[i] = *ptr
	}
	return res
}
