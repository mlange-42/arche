package ecs

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityAsIndex(t *testing.T) {
	entity := Entity{id: 1}
	arr := []int{0, 1, 2}

	val := arr[entity.id]
	_ = val
}

func TestEntityID(t *testing.T) {
	e := newEntityGen(1, 2)

	assert.Equal(t, e.ID(), uint32(1))
}

func TestEntityGeneration(t *testing.T) {
	e := newEntityGen(1, 2)

	assert.Equal(t, e.Generation(), uint32(2))
}

func TestZeroEntity(t *testing.T) {
	assert.True(t, Entity{}.IsZero())
	assert.False(t, Entity{1, 0}.IsZero())
}

func TestEntityMarshal(t *testing.T) {
	e := newEntityGen(2, 3)

	jsonData, err := json.Marshal(&e)
	if err != nil {
		t.Fatal(err)
	}

	e2 := Entity{}
	err = json.Unmarshal(jsonData, &e2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, e2, e)

	err = e2.UnmarshalJSON([]byte("pft"))
	assert.NotNil(t, err)
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
