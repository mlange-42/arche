package ecs

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitMask(t *testing.T) {
	mask := All(ID(1), ID(2), ID(13), ID(27))

	assert.Equal(t, 4, mask.TotalBitsSet())

	assert.True(t, mask.Get(1))
	assert.True(t, mask.Get(2))
	assert.True(t, mask.Get(13))
	assert.True(t, mask.Get(27))

	assert.False(t, mask.Get(0))
	assert.False(t, mask.Get(3))

	mask.Set(ID(0), true)
	mask.Set(ID(1), false)

	assert.True(t, mask.Get(0))
	assert.False(t, mask.Get(1))

	other1 := All(ID(1), ID(2), ID(32))
	other2 := All(ID(0), ID(2))

	assert.False(t, mask.Contains(other1))
	assert.True(t, mask.Contains(other2))

	mask.Reset()
	assert.Equal(t, 0, mask.TotalBitsSet())

	mask = All(ID(1), ID(2), ID(13), ID(27))
	other1 = All(ID(1), ID(32))
	other2 = All(ID(0), ID(32))

	assert.True(t, mask.ContainsAny(other1))
	assert.False(t, mask.ContainsAny(other2))
}

func TestBitMaskWithoutExclusive(t *testing.T) {
	mask := All(ID(1), ID(2), ID(13))
	assert.True(t, mask.Matches(All(ID(1), ID(2), ID(13))))
	assert.True(t, mask.Matches(All(ID(1), ID(2), ID(13), ID(27))))

	assert.False(t, mask.Matches(All(ID(1), ID(2))))

	without := mask.Without(ID(3))

	assert.True(t, without.Matches(All(ID(1), ID(2), ID(13))))
	assert.True(t, without.Matches(All(ID(1), ID(2), ID(13), ID(27))))

	assert.False(t, without.Matches(All(ID(1), ID(2), ID(3), ID(13))))
	assert.False(t, without.Matches(All(ID(1), ID(2))))

	excl := mask.Exclusive()

	assert.True(t, excl.Matches(All(ID(1), ID(2), ID(13))))
	assert.False(t, excl.Matches(All(ID(1), ID(2), ID(13), ID(27))))
	assert.False(t, excl.Matches(All(ID(1), ID(2), ID(3), ID(13))))
}

func TestBitMask128(t *testing.T) {
	for i := 0; i < MaskTotalBits; i++ {
		mask := All(ID(i))
		assert.Equal(t, 1, mask.TotalBitsSet())
		assert.True(t, mask.Get(ID(i)))
	}
	mask := Mask{}
	assert.Equal(t, 0, mask.TotalBitsSet())

	for i := 0; i < MaskTotalBits; i++ {
		mask.Set(ID(i), true)
		assert.Equal(t, i+1, mask.TotalBitsSet())
		assert.True(t, mask.Get(ID(i)))
	}

	mask = All(ID(1), ID(2), ID(13), ID(27), ID(63), ID(64), ID(65))

	assert.True(t, mask.Contains(All(ID(1), ID(2), ID(63), ID(64))))
	assert.False(t, mask.Contains(All(ID(1), ID(2), ID(63), ID(90))))

	assert.True(t, mask.ContainsAny(All(ID(6), ID(65), ID(111))))
	assert.False(t, mask.ContainsAny(All(ID(6), ID(66), ID(90))))
}

func TestBitMask64(t *testing.T) {
	mask := newBitMask64(ID(1))
	assert.True(t, mask.Get(ID(1)))
	for i := 0; i < wordSize; i++ {
		mask.Set(ID(i), true)
		assert.True(t, mask.Get(ID(i)))
		mask.Set(ID(i), false)
		assert.False(t, mask.Get(ID(i)))
	}
}

func BenchmarkBitmask64Get(b *testing.B) {
	b.StopTimer()
	mask := newBitMask64()
	for i := 0; i < MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(ID(i), true)
		}
	}
	idx := ID(rand.Intn(wordSize))
	b.StartTimer()

	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Get(idx)
	}
	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkBitmask128Get(b *testing.B) {
	b.StopTimer()
	mask := All()
	for i := 0; i < MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(ID(i), true)
		}
	}
	idx := ID(rand.Intn(MaskTotalBits))
	b.StartTimer()

	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Get(idx)
	}

	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkBitmaskContains(b *testing.B) {
	b.StopTimer()
	mask := All()
	for i := 0; i < MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(ID(i), true)
		}
	}
	filter := All(ID(rand.Intn(MaskTotalBits)))
	b.StartTimer()

	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Contains(filter)
	}

	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkBitmaskContainsAny(b *testing.B) {
	b.StopTimer()
	mask := All()
	for i := 0; i < MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(ID(i), true)
		}
	}
	filter := All(ID(rand.Intn(MaskTotalBits)))
	b.StartTimer()

	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.ContainsAny(filter)
	}

	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkMaskFilter(b *testing.B) {
	b.StopTimer()
	mask := All(0, 1, 2).Without()
	bits := All(0, 1, 2)
	b.StartTimer()
	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Matches(bits)
	}
	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkMaskFilterNoPointer(b *testing.B) {
	b.StopTimer()
	mask := maskFilterPointer{All(0, 1, 2), All()}
	bits := All(0, 1, 2)
	b.StartTimer()
	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Matches(bits)
	}
	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkMaskPointer(b *testing.B) {
	b.StopTimer()
	mask := maskPointer(All(0, 1, 2))
	bits := All(0, 1, 2)
	b.StartTimer()
	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Matches(bits)
	}
	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkMask(b *testing.B) {
	b.StopTimer()
	mask := All(0, 1, 2)
	bits := All(0, 1, 2)
	b.StartTimer()
	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Matches(bits)
	}
	b.StopTimer()
	v = !v
	_ = v
}

// bitMask64 is there just for performance comparison with the new 128 bit Mask.
type bitMask64 uint64

func newBitMask64(ids ...ID) bitMask64 {
	var mask bitMask64
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}
func (e bitMask64) Get(bit ID) bool {
	mask := bitMask64(1 << bit)
	return e&mask == mask
}

func (e *bitMask64) Set(bit ID, value bool) {
	if value {
		*e |= bitMask64(1 << bit)
	} else {
		*e &= bitMask64(^(1 << bit))
	}
}

type maskFilterPointer struct {
	Mask    Mask
	Exclude Mask
}

// Matches matches a filter against a mask.
func (f maskFilterPointer) Matches(bits Mask) bool {
	return bits.Contains(f.Mask) &&
		(f.Exclude.IsZero() || !bits.ContainsAny(f.Exclude))
}

type maskPointer Mask

// Matches matches a filter against a mask.
func (f *maskPointer) Matches(bits Mask) bool {
	return bits.Contains(Mask(*f))
}

func ExampleMask() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	filter := All(posID, velID)
	query := world.Query(filter)

	for query.Next() {
		// ...
	}
	// Output:
}

func ExampleMask_Without() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	filter := All(posID).Without(velID)
	query := world.Query(&filter)

	for query.Next() {
		// ...
	}
	// Output:
}

func ExampleMask_Exclusive() {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	filter := All(posID, velID).Exclusive()
	query := world.Query(&filter)

	for query.Next() {
		// ...
	}
	// Output:
}
