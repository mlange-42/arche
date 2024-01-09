package ecs_test

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestBitMask(t *testing.T) {
	mask := ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(13), ecs.ID(27), ecs.ID(200))

	assert.Equal(t, 5, mask.TotalBitsSet())

	assert.True(t, mask.Get(1))
	assert.True(t, mask.Get(2))
	assert.True(t, mask.Get(13))
	assert.True(t, mask.Get(27))
	assert.True(t, mask.Get(200))

	assert.False(t, mask.Get(0))
	assert.False(t, mask.Get(3))
	assert.False(t, mask.Get(199))
	assert.False(t, mask.Get(201))

	mask.Set(ecs.ID(0), true)
	mask.Set(ecs.ID(1), false)

	assert.True(t, mask.Get(0))
	assert.False(t, mask.Get(1))

	other1 := ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(32))
	other2 := ecs.All(ecs.ID(0), ecs.ID(2))

	assert.False(t, mask.Contains(other1))
	assert.True(t, mask.Contains(other2))

	mask.Reset()
	assert.Equal(t, 0, mask.TotalBitsSet())

	mask = ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(13), ecs.ID(27))
	other1 = ecs.All(ecs.ID(1), ecs.ID(32))
	other2 = ecs.All(ecs.ID(0), ecs.ID(32))

	assert.True(t, mask.ContainsAny(other1))
	assert.False(t, mask.ContainsAny(other2))
}

func TestBitMaskWithoutExclusive(t *testing.T) {
	mask := ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(13))
	assert.True(t, mask.Matches(ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(13))))
	assert.True(t, mask.Matches(ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(13), ecs.ID(27))))

	assert.False(t, mask.Matches(ecs.All(ecs.ID(1), ecs.ID(2))))

	without := mask.Without(ecs.ID(3))

	assert.True(t, without.Matches(ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(13))))
	assert.True(t, without.Matches(ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(13), ecs.ID(27))))

	assert.False(t, without.Matches(ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(3), ecs.ID(13))))
	assert.False(t, without.Matches(ecs.All(ecs.ID(1), ecs.ID(2))))

	excl := mask.Exclusive()

	assert.True(t, excl.Matches(ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(13))))
	assert.False(t, excl.Matches(ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(13), ecs.ID(27))))
	assert.False(t, excl.Matches(ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(3), ecs.ID(13))))
}

func TestBitMask256(t *testing.T) {
	for i := 0; i < ecs.MaskTotalBits; i++ {
		mask := ecs.All(ecs.ID(i))
		assert.Equal(t, 1, mask.TotalBitsSet())
		assert.True(t, mask.Get(ecs.ID(i)))
	}
	mask := ecs.Mask{}
	assert.Equal(t, 0, mask.TotalBitsSet())

	for i := 0; i < ecs.MaskTotalBits; i++ {
		mask.Set(ecs.ID(i), true)
		assert.Equal(t, i+1, mask.TotalBitsSet())
		assert.True(t, mask.Get(ecs.ID(i)))
	}

	mask = *ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(13), ecs.ID(27), ecs.ID(63), ecs.ID(64), ecs.ID(65))

	assert.True(t, mask.Contains(ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(63), ecs.ID(64))))
	assert.False(t, mask.Contains(ecs.All(ecs.ID(1), ecs.ID(2), ecs.ID(63), ecs.ID(90))))

	assert.True(t, mask.ContainsAny(ecs.All(ecs.ID(6), ecs.ID(65), ecs.ID(111))))
	assert.False(t, mask.ContainsAny(ecs.All(ecs.ID(6), ecs.ID(66), ecs.ID(90))))
}

func TestBitMask64(t *testing.T) {
	mask := newBitMask64(ecs.ID(1))
	assert.True(t, mask.Get(ecs.ID(1)))
	for i := 0; i < 64; i++ {
		mask.Set(ecs.ID(i), true)
		assert.True(t, mask.Get(ecs.ID(i)))
		mask.Set(ecs.ID(i), false)
		assert.False(t, mask.Get(ecs.ID(i)))
	}
}

func BenchmarkBitmask64Get(b *testing.B) {
	b.StopTimer()
	mask := newBitMask64()
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(ecs.ID(i), true)
		}
	}
	idx := ecs.ID(rand.Intn(ecs.MaskTotalBits / 2))
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
	mask := ecs.All()
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(ecs.ID(i), true)
		}
	}
	idx := ecs.ID(rand.Intn(ecs.MaskTotalBits))
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
	mask := ecs.All()
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(ecs.ID(i), true)
		}
	}
	filter := ecs.All(ecs.ID(rand.Intn(ecs.MaskTotalBits)))
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
	mask := ecs.All()
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(ecs.ID(i), true)
		}
	}
	filter := ecs.All(ecs.ID(rand.Intn(ecs.MaskTotalBits)))
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
	mask := ecs.All(0, 1, 2).Without()
	bits := ecs.All(0, 1, 2)
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
	mask := maskFilterPointer{*ecs.All(0, 1, 2), *ecs.All()}
	bits := ecs.All(0, 1, 2)
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
	mask := maskPointer(*ecs.All(0, 1, 2))
	bits := ecs.All(0, 1, 2)
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
	mask := ecs.All(0, 1, 2)
	bits := ecs.All(0, 1, 2)
	b.StartTimer()
	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Matches(bits)
	}
	b.StopTimer()
	v = !v
	_ = v
}

// bitMask64 is there just for performance comparison with the new 256 bit Mask.
type bitMask64 uint64

func newBitMask64(ids ...ecs.ID) bitMask64 {
	var mask bitMask64
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}
func (e bitMask64) Get(bit ecs.ID) bool {
	mask := bitMask64(1 << bit)
	return e&mask == mask
}

func (e *bitMask64) Set(bit ecs.ID, value bool) {
	if value {
		*e |= bitMask64(1 << bit)
	} else {
		*e &= bitMask64(^(1 << bit))
	}
}

type maskFilterPointer struct {
	Mask    ecs.Mask
	Exclude ecs.Mask
}

// Matches a filter against a mask.
func (f maskFilterPointer) Matches(bits *ecs.Mask) bool {
	return bits.Contains(&f.Mask) &&
		(f.Exclude.IsZero() || !bits.ContainsAny(&f.Exclude))
}

type maskPointer ecs.Mask

// Matches a filter against a mask.
func (f *maskPointer) Matches(bits *ecs.Mask) bool {
	m := ecs.Mask(*f)
	return bits.Contains(&m)
}

func ExampleMask() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	filter := ecs.All(posID, velID)
	query := world.Query(filter)

	for query.Next() {
		// ...
	}
	// Output:
}

func ExampleMask_Without() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	filter := ecs.All(posID).Without(velID)
	query := world.Query(&filter)

	for query.Next() {
		// ...
	}
	// Output:
}

func ExampleMask_Exclusive() {
	world := ecs.NewWorld()
	posID := ecs.ComponentID[Position](&world)
	velID := ecs.ComponentID[Velocity](&world)

	filter := ecs.All(posID, velID).Exclusive()
	query := world.Query(&filter)

	for query.Next() {
		// ...
	}
	// Output:
}
