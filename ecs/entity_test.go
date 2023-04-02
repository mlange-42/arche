package ecs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityAsIndex(t *testing.T) {
	entity := Entity{1, 0}
	arr := []int{0, 1, 2}

	val := arr[entity.id]
	_ = val
}

func TestZeroEntity(t *testing.T) {
	assert.True(t, Entity{}.IsZero())
	assert.False(t, Entity{1, 0}.IsZero())
}

func BenchmarkEntityIsZero(b *testing.B) {
	e := Entity{}

	isZero := false
	for i := 0; i < b.N; i++ {
		isZero = e.IsZero()
	}
	_ = isZero
}

func ExampleEntity() {
	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	e1 := world.NewEntity()
	e2 := world.NewEntity(posID, velID)

	fmt.Println(e1.IsZero(), e2.IsZero())
	// Output: false false
}

func ExampleEntity_IsZero() {
	world := NewWorld()

	var e1 Entity
	var e2 Entity = world.NewEntity()

	fmt.Println(e1.IsZero(), e2.IsZero())
	// Output: true false
}
