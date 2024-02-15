package ecs

import (
	"testing"

	"github.com/mlange-42/arche/ecs/event"
	"github.com/mlange-42/arche/ecs/stats"
)

func BenchmarkEntityAlive_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld(NewConfig().WithCapacityIncrement(1024))
	posID := ComponentID[Position](&world)

	entities := make([]Entity, 0, 1000)
	q := world.newEntitiesQuery(1000, ID{}, false, Entity{}, posID)
	for q.Next() {
		entities = append(entities, q.Entity())
	}

	b.StartTimer()

	var alive bool
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			alive = world.Alive(e)
		}
	}

	_ = alive
}

func BenchmarkGetResource(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	AddResource(&w, &Position{1, 2})
	posID := ResourceID[Position](&w)

	b.StartTimer()

	var res *Position
	for i := 0; i < b.N; i++ {
		res = w.Resources().Get(posID).(*Position)
	}

	_ = res
}

func BenchmarkGetResourceShortcut(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	AddResource(&w, &Position{1, 2})

	b.StartTimer()

	var res *Position
	for i := 0; i < b.N; i++ {
		res = GetResource[Position](&w)
	}

	_ = res
}

func BenchmarkNewEntities_10_000_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := NewWorld(NewConfig().WithCapacityIncrement(1024))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		for i := 0; i < 10000; i++ {
			_ = world.NewEntity(posID, velID)
		}
	}
}

func BenchmarkNewEntitiesBatch_10_000_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := NewWorld(NewConfig().WithCapacityIncrement(1024))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		world.newEntities(10000, ID{}, false, Entity{}, posID, velID)
	}
}

func BenchmarkNewEntities_10_000_Reset(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(1024))

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	for i := 0; i < 10000; i++ {
		_ = world.NewEntity(posID, velID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		world.Reset()
		for i := 0; i < 10000; i++ {
			_ = world.NewEntity(posID, velID)
		}
	}
}

func BenchmarkNewEntitiesBatch_10_000_Reset(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(1024))

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	for i := 0; i < 10000; i++ {
		_ = world.NewEntity(posID, velID)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		world.Reset()
		world.newEntities(10000, ID{}, false, Entity{}, posID, velID)
	}
}

func BenchmarkRemoveEntities_10_000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := NewWorld(NewConfig().WithCapacityIncrement(10000))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		entities := make([]Entity, 10000)
		q := world.newEntitiesQuery(10000, ID{}, false, Entity{}, posID, velID)

		cnt := 0
		for q.Next() {
			entities[cnt] = q.Entity()
			cnt++
		}

		b.StartTimer()

		for _, e := range entities {
			world.RemoveEntity(e)
		}
	}
}

func BenchmarkWorldNewQuery(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(10000))
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	NewBuilder(&world, posID, velID).NewBatch(25)

	filter := All(posID, velID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q := world.Query(filter)
		q.Close()
	}
}

func BenchmarkWorldNewQueryNext(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(10000))

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	NewBuilder(&world, posID, velID).NewBatch(25)

	filter := All(posID, velID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q := world.Query(filter)
		q.Next()
		q.Close()
	}
}

func BenchmarkWorldNewQueryCached(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(10000))
	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	NewBuilder(&world, posID, velID).NewBatch(25)

	filter := All(posID, velID)
	cf := world.Cache().Register(filter)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q := world.Query(&cf)
		q.Close()
	}
}

func BenchmarkWorldNewQueryNextCached(b *testing.B) {
	b.StopTimer()
	world := NewWorld(NewConfig().WithCapacityIncrement(10000))

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	NewBuilder(&world, posID, velID).NewBatch(25)

	filter := All(posID, velID)
	cf := world.Cache().Register(filter)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q := world.Query(&cf)
		q.Next()
		q.Close()
	}
}

func BenchmarkRemoveEntitiesBatch_10_000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		world := NewWorld(NewConfig().WithCapacityIncrement(10000))

		posID := ComponentID[Position](&world)
		velID := ComponentID[Velocity](&world)

		q := world.newEntitiesQuery(10000, ID{}, false, Entity{}, posID, velID)
		q.Close()
		b.StartTimer()
		world.Batch().RemoveEntities(All(posID, velID))
	}
}

func BenchmarkWorldStats_1Arch(b *testing.B) {
	b.StopTimer()

	w := NewWorld()
	w.NewEntity()

	b.StartTimer()

	var st *stats.World
	for i := 0; i < b.N; i++ {
		st = w.Stats()
	}
	_ = st
}

func BenchmarkWorldStats_10Arch(b *testing.B) {
	b.StopTimer()

	w := NewWorld()

	ids := []ID{
		ComponentID[testStruct0](&w),
		ComponentID[testStruct1](&w),
		ComponentID[testStruct2](&w),
		ComponentID[testStruct3](&w),
		ComponentID[testStruct4](&w),
		ComponentID[testStruct5](&w),
		ComponentID[testStruct6](&w),
		ComponentID[testStruct7](&w),
		ComponentID[testStruct8](&w),
		ComponentID[testStruct9](&w),
	}

	for _, id := range ids {
		w.NewEntity(id)
	}

	b.StartTimer()

	var st *stats.World
	for i := 0; i < b.N; i++ {
		st = w.Stats()
	}
	_ = st
}

func BenchmarkWorldNewEntityNoListener_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	filterPos := All(posID)

	builder := NewBuilder(&world, posID)
	builder.NewBatch(1000)
	world.Batch().RemoveEntities(filterPos)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			world.NewEntity(posID)
		}
		b.StopTimer()
		world.Batch().RemoveEntities(filterPos)
	}
}

func BenchmarkWorldNewEntityNoEvent_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	filterPos := All(posID)

	builder := NewBuilder(&world, posID)
	builder.NewBatch(1000)
	world.Batch().RemoveEntities(filterPos)

	var temp event.Subscription
	listener := newTestListener(func(world *World, e EntityEvent) { temp = e.EventTypes })
	listener.Subscribe = 0

	world.SetListener(&listener)
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			world.NewEntity(posID)
		}
		b.StopTimer()
		world.Batch().RemoveEntities(filterPos)
	}
	_ = temp
}

func BenchmarkWorldNewEntityEvent_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	filterPos := All(posID)

	builder := NewBuilder(&world, posID)
	builder.NewBatch(1000)
	world.Batch().RemoveEntities(filterPos)

	listener := dummyListener{}
	world.SetListener(&listener)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			world.NewEntity(posID)
		}
		b.StopTimer()
		world.Batch().RemoveEntities(filterPos)
	}
}

func BenchmarkWorldNewEntityEventCallback_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	filterPos := All(posID)

	builder := NewBuilder(&world, posID)
	builder.NewBatch(1000)
	world.Batch().RemoveEntities(filterPos)

	var temp event.Subscription
	listener := newTestListener(func(world *World, e EntityEvent) { temp = e.EventTypes })
	world.SetListener(&listener)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			world.NewEntity(posID)
		}
		b.StopTimer()
		world.Batch().RemoveEntities(filterPos)
	}
	_ = temp
}

func BenchmarkWorldExchangeNoListener_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	builder := NewBuilder(&world, posID)
	entities := make([]Entity, 0, 1000)
	query := builder.NewBatchQ(1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	filterPos := All(posID)
	filterVel := All(velID)

	pos := []ID{posID}
	vel := []ID{velID}

	world.Batch().Exchange(filterPos, vel, pos)
	world.Batch().Exchange(filterVel, pos, vel)

	b.StartTimer()
	hasPos := true
	for i := 0; i < b.N; i++ {
		if hasPos {
			for _, e := range entities {
				world.Exchange(e, vel, pos)
			}
		} else {
			for _, e := range entities {
				world.Exchange(e, pos, vel)
			}
		}
		hasPos = !hasPos
	}
}

func BenchmarkWorldExchangeNoEvent_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	builder := NewBuilder(&world, posID)
	entities := make([]Entity, 0, 1000)
	query := builder.NewBatchQ(1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	filterPos := All(posID)
	filterVel := All(velID)

	pos := []ID{posID}
	vel := []ID{velID}

	world.Batch().Exchange(filterPos, vel, pos)
	world.Batch().Exchange(filterVel, pos, vel)

	var temp event.Subscription
	listener := newTestListener(func(world *World, e EntityEvent) { temp = e.EventTypes })
	listener.Subscribe = 0
	world.SetListener(&listener)

	b.StartTimer()
	hasPos := true
	for i := 0; i < b.N; i++ {
		if hasPos {
			for _, e := range entities {
				world.Exchange(e, vel, pos)
			}
		} else {
			for _, e := range entities {
				world.Exchange(e, pos, vel)
			}
		}
		hasPos = !hasPos
	}
	_ = temp
}

func BenchmarkWorldExchangeEvent_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	builder := NewBuilder(&world, posID)
	entities := make([]Entity, 0, 1000)
	query := builder.NewBatchQ(1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	filterPos := All(posID)
	filterVel := All(velID)

	pos := []ID{posID}
	vel := []ID{velID}

	world.Batch().Exchange(filterPos, vel, pos)
	world.Batch().Exchange(filterVel, pos, vel)

	listener := dummyListener{}
	world.SetListener(&listener)

	b.StartTimer()
	hasPos := true
	for i := 0; i < b.N; i++ {
		if hasPos {
			for _, e := range entities {
				world.Exchange(e, vel, pos)
			}
		} else {
			for _, e := range entities {
				world.Exchange(e, pos, vel)
			}
		}
		hasPos = !hasPos
	}
}

func BenchmarkWorldExchangeEventCallback_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	builder := NewBuilder(&world, posID)
	entities := make([]Entity, 0, 1000)
	query := builder.NewBatchQ(1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	filterPos := All(posID)
	filterVel := All(velID)

	pos := []ID{posID}
	vel := []ID{velID}

	world.Batch().Exchange(filterPos, vel, pos)
	world.Batch().Exchange(filterVel, pos, vel)

	var temp event.Subscription
	listener := newTestListener(func(world *World, e EntityEvent) { temp = e.EventTypes })
	world.SetListener(&listener)

	b.StartTimer()
	hasPos := true
	for i := 0; i < b.N; i++ {
		if hasPos {
			for _, e := range entities {
				world.Exchange(e, vel, pos)
			}
		} else {
			for _, e := range entities {
				world.Exchange(e, pos, vel)
			}
		}
		hasPos = !hasPos
	}
	_ = temp
}

func BenchmarkWorldExchangeBatchNoListener_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	builder := NewBuilder(&world, posID)
	builder.NewBatch(1000)

	filterPos := All(posID)
	filterVel := All(velID)

	pos := []ID{posID}
	vel := []ID{velID}

	world.Batch().Exchange(filterPos, vel, pos)
	world.Batch().Exchange(filterVel, pos, vel)

	b.StartTimer()
	hasPos := true
	for i := 0; i < b.N; i++ {
		if hasPos {
			world.Batch().Exchange(filterPos, vel, pos)
		} else {
			world.Batch().Exchange(filterVel, pos, vel)
		}
		hasPos = !hasPos
	}
}

func BenchmarkWorldExchangeBatchNoEvent_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	builder := NewBuilder(&world, posID)
	builder.NewBatch(1000)

	filterPos := All(posID)
	filterVel := All(velID)

	pos := []ID{posID}
	vel := []ID{velID}

	world.Batch().Exchange(filterPos, vel, pos)
	world.Batch().Exchange(filterVel, pos, vel)

	var temp event.Subscription
	listener := newTestListener(func(world *World, e EntityEvent) { temp = e.EventTypes })
	listener.Subscribe = 0
	world.SetListener(&listener)

	b.StartTimer()
	hasPos := true
	for i := 0; i < b.N; i++ {
		if hasPos {
			world.Batch().Exchange(filterPos, vel, pos)
		} else {
			world.Batch().Exchange(filterVel, pos, vel)
		}
		hasPos = !hasPos
	}
	_ = temp
}

func BenchmarkWorldExchangeBatchEvent_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	builder := NewBuilder(&world, posID)
	builder.NewBatch(1000)

	filterPos := All(posID)
	filterVel := All(velID)

	pos := []ID{posID}
	vel := []ID{velID}

	world.Batch().Exchange(filterPos, vel, pos)
	world.Batch().Exchange(filterVel, pos, vel)

	listener := dummyListener{}
	world.SetListener(&listener)

	b.StartTimer()
	hasPos := true
	for i := 0; i < b.N; i++ {
		if hasPos {
			world.Batch().Exchange(filterPos, vel, pos)
		} else {
			world.Batch().Exchange(filterVel, pos, vel)
		}
		hasPos = !hasPos
	}
}

func BenchmarkWorldExchangeBatchEventCallback_1000(b *testing.B) {
	b.StopTimer()

	world := NewWorld()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	builder := NewBuilder(&world, posID)
	builder.NewBatch(1000)

	filterPos := All(posID)
	filterVel := All(velID)

	pos := []ID{posID}
	vel := []ID{velID}

	world.Batch().Exchange(filterPos, vel, pos)
	world.Batch().Exchange(filterVel, pos, vel)

	var temp event.Subscription
	listener := newTestListener(func(world *World, e EntityEvent) { temp = e.EventTypes })
	world.SetListener(&listener)

	b.StartTimer()
	hasPos := true
	for i := 0; i < b.N; i++ {
		if hasPos {
			world.Batch().Exchange(filterPos, vel, pos)
		} else {
			world.Batch().Exchange(filterVel, pos, vel)
		}
		hasPos = !hasPos
	}
	_ = temp
}

type dummyListener struct {
	temp event.Subscription
}

func (l *dummyListener) Notify(w *World, evt EntityEvent) {
	l.temp = evt.EventTypes
}

func (l *dummyListener) Subscriptions() event.Subscription {
	return event.All
}

func (l *dummyListener) Components() *Mask {
	return nil
}
