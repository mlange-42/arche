package ecs

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func all(ids ...ID) *Mask {
	mask := All(ids...)
	return &mask
}

func TestBitMask(t *testing.T) {
	big := uint8(MaskTotalBits - 2)
	mask := All(id(1), id(2), id(13), id(27), id(big))

	assert.Equal(t, 5, mask.TotalBitsSet())

	assert.True(t, mask.Get(id(1)))
	assert.True(t, mask.Get(id(2)))
	assert.True(t, mask.Get(id(13)))
	assert.True(t, mask.Get(id(27)))
	assert.True(t, mask.Get(id(big)))

	assert.False(t, mask.Get(id(0)))
	assert.False(t, mask.Get(id(3)))
	assert.False(t, mask.Get(id(big-1)))
	assert.False(t, mask.Get(id(big+1)))

	mask.Set(id(0), true)
	mask.Set(id(1), false)

	assert.True(t, mask.Get(id(0)))
	assert.False(t, mask.Get(id(1)))

	other1 := All(id(1), id(2), id(32))
	other2 := All(id(0), id(2))

	assert.False(t, mask.Contains(&other1))
	assert.True(t, mask.Contains(&other2))

	mask.Reset()
	assert.Equal(t, 0, mask.TotalBitsSet())

	mask = All(id(1), id(2), id(13), id(27))
	other1 = All(id(1), id(32))
	other2 = All(id(0), id(32))

	assert.True(t, mask.ContainsAny(&other1))
	assert.False(t, mask.ContainsAny(&other2))
}

func TestBitMaskLogic(t *testing.T) {
	big := uint8(MaskTotalBits - 2)

	assert.Equal(t, All(id(5)), all(id(0), id(5)).And(all(id(5), id(big))))
	assert.Equal(t, All(id(0), id(5), id(big)), all(id(0), id(5)).Or(all(id(5), id(big))))
	assert.Equal(t, All(id(0), id(big)), all(id(0), id(5)).Xor(all(id(5), id(big))))
}

func TestBitMaskCopy(t *testing.T) {
	big := uint8(MaskTotalBits - 2)

	mask := All(id(1), id(2), id(13), id(27), id(big))
	mask2 := mask
	mask3 := &mask

	mask2.Set(id(1), false)
	mask3.Set(id(2), false)

	assert.True(t, mask.Get(id(1)))
	assert.False(t, mask2.Get(id(1)))

	assert.True(t, mask2.Get(id(2)))
	assert.False(t, mask.Get(id(2)))
	assert.False(t, mask3.Get(id(2)))
}

func TestBitMaskWithoutExclusive(t *testing.T) {
	mask := All(id(1), id(2), id(13))

	assert.True(t, mask.Matches(all(id(1), id(2), id(13))))
	assert.True(t, mask.Matches(all(id(1), id(2), id(13), id(27))))

	assert.False(t, mask.Matches(all(id(1), id(2))))

	without := mask.Without(id(3))

	assert.True(t, without.Matches(all(id(1), id(2), id(13))))
	assert.True(t, without.Matches(all(id(1), id(2), id(13), id(27))))

	assert.False(t, without.Matches(all(id(1), id(2), id(3), id(13))))
	assert.False(t, without.Matches(all(id(1), id(2))))

	excl := mask.Exclusive()

	assert.True(t, excl.Matches(all(id(1), id(2), id(13))))
	assert.False(t, excl.Matches(all(id(1), id(2), id(13), id(27))))
	assert.False(t, excl.Matches(all(id(1), id(2), id(3), id(13))))
}

func TestBitMask256(t *testing.T) {
	for i := 0; i < MaskTotalBits; i++ {
		mask := All(id(uint8(i)))
		assert.Equal(t, 1, mask.TotalBitsSet())
		assert.True(t, mask.Get(id(uint8(i))))
	}
	mask := Mask{}
	assert.Equal(t, 0, mask.TotalBitsSet())

	for i := 0; i < MaskTotalBits; i++ {
		mask.Set(id(uint8(i)), true)
		assert.Equal(t, i+1, mask.TotalBitsSet())
		assert.True(t, mask.Get(id(uint8(i))))
	}

	big := uint8(MaskTotalBits - 10)

	mask = All(id(1), id(2), id(13), id(27), id(big), id(big+1), id(big+2))

	assert.True(t, mask.Contains(all(id(1), id(2), id(big), id(big+1))))
	assert.False(t, mask.Contains(all(id(1), id(2), id(big), id(big+5))))

	assert.True(t, mask.ContainsAny(all(id(6), id(big+2), id(big+6))))
	assert.False(t, mask.ContainsAny(all(id(6), id(big+3), id(big+5))))
}

func TestBitMask64(t *testing.T) {
	mask := newBitMask64(id(1))
	assert.True(t, mask.Get(1))
	for i := 0; i < 64; i++ {
		mask.Set(uint8(i), true)
		assert.True(t, mask.Get(uint8(i)))
		mask.Set(uint8(i), false)
		assert.False(t, mask.Get(uint8(i)))
	}
}

func BenchmarkBitmask64Get(b *testing.B) {
	b.StopTimer()
	mask := newBitMask64()
	for i := 0; i < MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(uint8(i), true)
		}
	}
	idx := id(uint8(rand.Intn(MaskTotalBits / 4)))
	b.StartTimer()

	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Get(idx.id)
	}
	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkBitmask256Get(b *testing.B) {
	b.StopTimer()
	mask := All()
	for i := 0; i < MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(id(uint8(i)), true)
		}
	}
	idx := id(uint8(rand.Intn(MaskTotalBits)))
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
			mask.Set(id(uint8(i)), true)
		}
	}
	filter := All(id(uint8(rand.Intn(MaskTotalBits))))
	b.StartTimer()

	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Contains(&filter)
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
			mask.Set(id(uint8(i)), true)
		}
	}
	filter := All(id(uint8(rand.Intn(MaskTotalBits))))
	b.StartTimer()

	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.ContainsAny(&filter)
	}

	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkMaskFilter(b *testing.B) {
	b.StopTimer()
	mask := All(id(0), id(1), id(2)).Without(id(3))
	bits := All(id(0), id(1), id(2))
	b.StartTimer()
	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Matches(&bits)
	}
	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkMaskFilterNoPointer(b *testing.B) {
	b.StopTimer()
	mask := maskFilterPointer{All(id(0), id(1), id(2)), All(id(3))}
	bits := All(id(0), id(1), id(2))
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
	mask := maskPointer(All(id(0), id(1), id(2)))
	bits := All(id(0), id(1), id(2))
	b.StartTimer()
	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Matches(bits)
	}
	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkMaskMatch(b *testing.B) {
	b.StopTimer()
	mask := All(id(0), id(1), id(2))
	bits := All(id(0), id(1), id(2))
	b.StartTimer()
	var v bool
	for i := 0; i < b.N; i++ {
		v = mask.Matches(&bits)
	}
	b.StopTimer()
	v = !v
	_ = v
}

func BenchmarkMaskCopy(b *testing.B) {
	b.StopTimer()
	mask := All(id(0), id(1), id(2))
	var tempMask Mask
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tempMask = mask
	}
	b.StopTimer()
	_ = tempMask
}

// bitMask64 is there just for performance comparison with the new 256 bit Mask.
type bitMask64 uint64

func newBitMask64(ids ...ID) bitMask64 {
	var mask bitMask64
	for _, id := range ids {
		mask.Set(id.id, true)
	}
	return mask
}
func (e bitMask64) Get(bit uint8) bool {
	mask := bitMask64(1 << bit)
	return e&mask == mask
}

func (e *bitMask64) Set(bit uint8, value bool) {
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

// Matches a filter against a mask.
func (f maskFilterPointer) Matches(bits Mask) bool {
	return bits.Contains(&f.Mask) &&
		(f.Exclude.IsZero() || !bits.ContainsAny(&f.Exclude))
}

type maskPointer Mask

// Matches a filter against a mask.
func (f *maskPointer) Matches(bits Mask) bool {
	m := Mask(*f)
	return bits.Contains(&m)
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
