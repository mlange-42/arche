package event_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/arche/ecs/event"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptions(t *testing.T) {
	m1 := event.EntityCreated | event.TargetChanged

	assert.True(t, m1.Contains(event.EntityCreated))
	assert.False(t, m1.Contains(event.EntityRemoved))

	assert.True(t, m1.ContainsAny(event.ComponentAdded|event.TargetChanged))
	assert.False(t, m1.Contains(event.ComponentAdded|event.RelationChanged))
}

func ExampleSubscription() {
	mask := event.EntityCreated | event.EntityRemoved

	fmt.Printf("%08b contains\n%08b -> %t\n\n", mask, event.EntityRemoved, mask.Contains(event.EntityRemoved))
	fmt.Printf("%08b contains\n%08b -> %t\n\n", mask, event.ComponentAdded, mask.Contains(event.ComponentAdded))

	fmt.Printf("%08b contains any\n%08b -> %t\n\n", mask, event.EntityRemoved|event.ComponentAdded, mask.ContainsAny(event.EntityRemoved|event.ComponentAdded))
	fmt.Printf("%08b contains any\n%08b -> %t\n\n", mask, event.ComponentAdded|event.ComponentRemoved, mask.ContainsAny(event.ComponentAdded|event.ComponentRemoved))
	// Output: 00000011 contains
	// 00000010 -> true
	//
	// 00000011 contains
	// 00000100 -> false
	//
	// 00000011 contains any
	// 00000110 -> true
	//
	// 00000011 contains any
	// 00001100 -> false
}
