package ecs

const (
	numChunks = 8
	chunkSize = 16
)

// idMap maps component IDs.
type idMap[T any] struct {
	chunks    [][]*T
	used      Mask
	chunkUsed []uint8
}

func newIDMap[T any]() idMap[T] {
	return idMap[T]{
		chunks:    make([][]*T, numChunks),
		used:      Mask{},
		chunkUsed: make([]uint8, numChunks),
	}
}

func (m *idMap[T]) Get(index uint8) (*T, bool) {
	if !m.used.Get(index) {
		return nil, false
	}
	chunk := index / chunkSize
	subIndex := index % chunkSize
	return m.chunks[chunk][subIndex], true
}

func (m *idMap[T]) Set(index uint8, value *T) {
	chunk := index / chunkSize
	subIndex := index % chunkSize
	if m.chunks[chunk] == nil {
		m.chunks[chunk] = make([]*T, chunkSize)
	}
	m.chunks[chunk][subIndex] = value
	m.used.Set(index, true)
	m.chunkUsed[chunk]++
}

func (m *idMap[T]) Remove(index uint8) {
	chunk := index / chunkSize
	m.used.Set(index, false)
	m.chunkUsed[chunk]--
	if m.chunkUsed[chunk] == 0 {
		m.chunks[chunk] = nil
	}
}
