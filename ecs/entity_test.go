package ecs

import (
	"testing"
)

func TestEntityAsIndex(t *testing.T) {
	entity := Entity{1, 0}
	arr := []int{0, 1, 2}

	val := arr[entity.id]
	_ = val
}
