package event_test

import (
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
