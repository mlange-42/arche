package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDMap(t *testing.T) {
	big1 := uint8(MaskTotalBits - 20)
	big2 := uint8(MaskTotalBits - 3)

	m := newIDMap[*Entity]()

	e0 := Entity{0, 0}
	e1 := Entity{1, 0}
	e121 := Entity{eid(big1), 0}
	e200 := Entity{eid(big2), 0}

	m.Set(0, &e0)
	m.Set(1, &e1)
	m.Set(big1, &e121)
	m.Set(big2, &e200)

	e, ok := m.Get(0)
	assert.True(t, ok)
	assert.Equal(t, e0, *e)

	e, ok = m.Get(1)
	assert.True(t, ok)
	assert.Equal(t, e1, *e)

	e, ok = m.Get(big1)
	assert.True(t, ok)
	assert.Equal(t, e121, *e)

	e, ok = m.Get(big2)
	assert.True(t, ok)
	assert.Equal(t, e200, *e)

	e, ok = m.Get(15)
	assert.False(t, ok)
	assert.Nil(t, e)

	m.Remove(0)
	m.Remove(1)

	e, ok = m.Get(0)
	assert.False(t, ok)
	assert.Nil(t, e)

	assert.Nil(t, m.chunks[0])
}

func TestIDMapPointers(t *testing.T) {
	big1 := uint8(MaskTotalBits - 20)
	big2 := uint8(MaskTotalBits - 3)

	m := newIDMap[Entity]()

	e0 := Entity{0, 0}
	e1 := Entity{1, 0}
	e121 := Entity{eid(big1), 0}
	e200 := Entity{eid(big2), 0}

	m.Set(0, e0)
	m.Set(1, e1)
	m.Set(big1, e121)
	m.Set(big2, e200)

	e, ok := m.GetPointer(0)
	assert.True(t, ok)
	assert.Equal(t, e0, *e)

	e, ok = m.GetPointer(1)
	assert.True(t, ok)
	assert.Equal(t, e1, *e)

	e, ok = m.GetPointer(big1)
	assert.True(t, ok)
	assert.Equal(t, e121, *e)

	e, ok = m.GetPointer(big2)
	assert.True(t, ok)
	assert.Equal(t, e200, *e)

	e, ok = m.GetPointer(15)
	assert.False(t, ok)
	assert.Nil(t, e)

	m.Remove(0)
	m.Remove(1)

	e, ok = m.GetPointer(0)
	assert.False(t, ok)
	assert.Nil(t, e)

	assert.Nil(t, m.chunks[0])
}

func BenchmarkIdMapping_IDMap(b *testing.B) {
	b.StopTimer()

	entities := [MaskTotalBits]Entity{}
	m := newIDMap[*Entity]()

	for i := 0; i < MaskTotalBits; i++ {
		entities[i] = Entity{eid(i), 0}
		m.Set(uint8(i), &entities[i])
	}

	b.StartTimer()

	var ptr *Entity = nil
	for i := 0; i < b.N; i++ {
		ptr, _ = m.Get(uint8(i % MaskTotalBits))
	}
	_ = ptr
}

func BenchmarkIdMapping_Array(b *testing.B) {
	b.StopTimer()

	entities := [MaskTotalBits]Entity{}
	m := [MaskTotalBits]*Entity{}

	for i := 0; i < MaskTotalBits; i++ {
		entities[i] = Entity{eid(i), 0}
		m[i] = &entities[i]
	}

	b.StartTimer()

	var ptr *Entity = nil
	for i := 0; i < b.N; i++ {
		ptr = m[i%MaskTotalBits]
	}
	_ = ptr
}

func BenchmarkIdMapping_HashMap(b *testing.B) {
	b.StopTimer()

	entities := [MaskTotalBits]Entity{}
	m := make(map[uint8]*Entity, MaskTotalBits)

	for i := 0; i < MaskTotalBits; i++ {
		entities[i] = Entity{eid(i), 0}
		m[uint8(i)] = &entities[i]
	}

	b.StartTimer()

	var ptr *Entity = nil
	for i := 0; i < b.N; i++ {
		ptr = m[uint8(i%MaskTotalBits)]
	}
	_ = ptr
}

func BenchmarkIdMapping_HashMapEntity(b *testing.B) {
	b.StopTimer()

	entities := [MaskTotalBits]Entity{}
	m := make(map[Entity]*Entity, MaskTotalBits)

	for i := 0; i < MaskTotalBits; i++ {
		entities[i] = Entity{eid(i), 0}
		m[Entity{eid(i), 0}] = &entities[i]
	}

	b.StartTimer()

	var ptr *Entity = nil
	for i := 0; i < b.N; i++ {
		ptr = m[Entity{eid(i % MaskTotalBits), 0}]
	}
	_ = ptr
}
