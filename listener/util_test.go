package listener

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/stretchr/testify/assert"
)

type comp1 struct{}
type comp2 struct{}
type comp3 struct{}

func TestSubscribes(t *testing.T) {
	world := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&world)
	id2 := ecs.ComponentID[comp2](&world)
	id3 := ecs.ComponentID[comp3](&world)

	m1 := ecs.All(id1)
	m12 := ecs.All(id1, id2)
	m2 := ecs.All(id2)
	m3 := ecs.All(id3)

	assert.True(t,
		subscribes(event.ComponentAdded, &m1, &m12, nil, nil),
	)
	assert.True(t,
		subscribes(event.ComponentAdded, &m12, &m2, nil, nil),
	)
	assert.False(t,
		subscribes(event.ComponentAdded, &m12, &m3, nil, nil),
	)

	assert.True(t,
		subscribes(event.RelationChanged, &ecs.Mask{}, &m12, nil, &id1),
	)
	assert.True(t,
		subscribes(event.RelationChanged, &ecs.Mask{}, &m12, &id1, &id3),
	)
	assert.False(t,
		subscribes(event.RelationChanged, &ecs.Mask{}, &m1, &id2, &id3),
	)

	assert.True(t,
		subscribes(event.TargetChanged, &ecs.Mask{}, &m12, &id1, &id1),
	)
	assert.False(t,
		subscribes(event.TargetChanged, &ecs.Mask{}, &m12, &id3, &id3),
	)

	assert.True(t,
		subscribes(event.ComponentAdded|event.ComponentRemoved|event.TargetChanged, &m12, &m3, &id3, &id3),
	)
	assert.False(t,
		subscribes(event.ComponentAdded|event.ComponentRemoved|event.TargetChanged, &m1, &m3, &id2, &id2),
	)
}
