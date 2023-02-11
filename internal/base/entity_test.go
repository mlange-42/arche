package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityAsIndex(t *testing.T) {
	entity := Entity{1, 0}
	arr := []int{0, 1, 2}

	val := arr[entity.ID]
	_ = val
}

func TestZeroEntity(t *testing.T) {
	assert.True(t, Entity{}.IsZero())
	assert.False(t, Entity{1, 0}.IsZero())
}
