package base

const fixedPageSize = 32

// PagedArr32 is a paged collection working with pages of length 32 arrays
type PagedArr32[T any] struct {
	pages   [][fixedPageSize]T
	len     int
	lenLast int
}

func (p *PagedArr32[T]) Add(value T) {
	if len(p.pages) == 0 || p.lenLast == fixedPageSize {
		p.pages = append(p.pages, [fixedPageSize]T{})
		p.lenLast = 0
	}
	p.pages[len(p.pages)-1][p.lenLast] = value
	p.len++
	p.lenLast++
}

func (p *PagedArr32[T]) Get(index int) *T {
	return &p.pages[index/fixedPageSize][index%fixedPageSize]
}

func (p *PagedArr32[T]) Len() int {
	return p.len
}
