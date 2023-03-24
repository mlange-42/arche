package ecs

const fixedPageSize = 32

// pagedArr32 is a paged collection working with pages of length 32 arrays.
//
// Implements [archetypes].
type pagedArr32[T any] struct {
	pages   [][fixedPageSize]T
	len     int
	lenLast int
}

// Add adds a value to the paged array.
func (p *pagedArr32[T]) Add(value T) {
	if p.len == 0 || p.lenLast == fixedPageSize {
		p.pages = append(p.pages, [fixedPageSize]T{})
		p.lenLast = 0
	}
	p.pages[len(p.pages)-1][p.lenLast] = value
	p.len++
	p.lenLast++
}

// Get returns the value at the given index.
func (p *pagedArr32[T]) Get(index int) *T {
	return &p.pages[index/fixedPageSize][index%fixedPageSize]
}

// Len returns the current number of items in the paged array.
func (p *pagedArr32[T]) Len() int {
	return p.len
}
