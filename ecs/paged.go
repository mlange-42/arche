package ecs

const fixedPageSize = 32

// pagedSlice is a paged collection working with pages of length 32 slices.
// It's primary purpose is pointer persistence, which is not given using simple slices.
//
// Implements [archetypes].
type pagedSlice[T any] struct {
	pages    [][]T
	len      int
	lenLast  int
	pageSize int
}

func newPagedSlice[T any](pageSize int) pagedSlice[T] {
	return pagedSlice[T]{
		pageSize: pageSize,
	}
}

// Add adds a value to the paged array.
func (p *pagedSlice[T]) Add(value T) {
	if p.len == 0 || p.lenLast == fixedPageSize {
		p.pages = append(p.pages, make([]T, fixedPageSize))
		p.lenLast = 0
	}
	p.pages[len(p.pages)-1][p.lenLast] = value
	p.len++
	p.lenLast++
}

// Get returns the value at the given index.
func (p *pagedSlice[T]) Get(index int) *T {
	return &p.pages[index/fixedPageSize][index%fixedPageSize]
}

// Len returns the current number of items in the paged array.
func (p *pagedSlice[T]) Len() int {
	return p.len
}

// pagedPointerArr32 is a paged collection working with pages of length 32 arrays.
//
// Implements [archetypes].
type pagedPointerArr32[T any] struct {
	pages   [][]*T
	len     int
	lenLast int
}

// Add adds a value to the paged array.
func (p *pagedPointerArr32[T]) Add(value *T) {
	if p.len == 0 || p.lenLast == fixedPageSize {
		p.pages = append(p.pages, make([]*T, fixedPageSize))
		p.lenLast = 0
	}
	p.pages[len(p.pages)-1][p.lenLast] = value
	p.len++
	p.lenLast++
}

// Get returns the value at the given index.
func (p *pagedPointerArr32[T]) Get(index int) *T {
	return p.pages[index/fixedPageSize][index%fixedPageSize]
}

// Len returns the current number of items in the paged array.
func (p *pagedPointerArr32[T]) Len() int {
	return p.len
}
