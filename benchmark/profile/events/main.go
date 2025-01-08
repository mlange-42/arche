package main

// Profiling:
// go build ./benchmark/profile/events
// events
// go tool pprof -http=":8000" -nodefraction=0.001 events cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 events mem.pprof

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/event"
	"github.com/pkg/profile"
)

type position struct {
	X int
	Y int
}

type velocity struct {
	X int
	Y int
}

func main() {
	count := 100
	iters := 1000
	entities := 1000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()

	stop = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()
}

func run(rounds, iters, numEntities int) {
	for i := 0; i < rounds; i++ {
		world := ecs.NewWorld(1024)

		posID := ecs.ComponentID[position](&world)
		velID := ecs.ComponentID[velocity](&world)

		builder := ecs.NewBuilder(&world, posID)
		builder.NewBatch(1000)

		pos := []ecs.ID{posID}
		vel := []ecs.ID{velID}

		filterPos := ecs.All(posID)
		filterVel := ecs.All(velID)

		var temp event.Subscription
		listener := dummyListener{}
		world.SetListener(&listener)

		hasPos := true
		for j := 0; j < iters; j++ {
			if hasPos {
				world.Batch().Exchange(filterPos, vel, pos)
			} else {
				world.Batch().Exchange(filterVel, pos, vel)
			}
			hasPos = !hasPos
		}
		_ = temp
	}
}

type dummyListener struct {
	temp event.Subscription
}

func (l *dummyListener) Notify(w *ecs.World, evt ecs.EntityEvent) {
	l.temp = evt.EventTypes
}

func (l *dummyListener) Subscriptions() event.Subscription {
	return event.Components
}

func (l *dummyListener) Components() *ecs.Mask {
	return nil
}
