// Demonstrates how to implement systems.
package main

import (
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func main() {
	// Create a new Scheduler
	scheduler := Scheduler{}

	// Parametrize and add Systems
	scheduler.AddSystem(
		&InitializerSystem{Count: 100},
	)
	scheduler.AddSystem(
		&PosUpdaterSystem{},
	)

	// Run the model
	scheduler.Run(100)
}

// System interface
type System interface {
	Initialize(w *Scheduler)
	Update(w *Scheduler)
}

// Scheduler for updating systems
type Scheduler struct {
	ecs.World
	systems []System
}

// AddSystem adds a System to the Scheduler
func (s *Scheduler) AddSystem(sys System) {
	s.systems = append(s.systems, sys)
}

// Run initializes and updates all Systems
func (s *Scheduler) Run(steps int) {
	s.initialize()
	s.update()
}

func (s *Scheduler) initialize() {
	s.World = ecs.NewWorld()

	for _, sys := range s.systems {
		sys.Initialize(s)
	}
}

func (s *Scheduler) update() {
	for _, sys := range s.systems {
		sys.Update(s)
	}
}

// Position component
type Position struct {
	X float64
	Y float64
}

// Velocity component
type Velocity struct {
	X float64
	Y float64
}

// InitializerSystem to create entities
type InitializerSystem struct {
	Count int
}

// Initialize the system
func (s *InitializerSystem) Initialize(w *Scheduler) {
	mapper := generic.NewMap2[Position, Velocity](w)
	for i := 0; i < s.Count; i++ {
		_, pos, vel := mapper.NewEntity()

		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100
		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	}
}

// Update the system
func (s *InitializerSystem) Update(w *Scheduler) {}

// PosUpdaterSystem updates entity positions
type PosUpdaterSystem struct {
	filter generic.Filter2[Position, Velocity]
}

// Initialize the system
func (s *PosUpdaterSystem) Initialize(w *Scheduler) {
	s.filter = *generic.NewFilter2[Position, Velocity]()
}

// Update the system
func (s *PosUpdaterSystem) Update(w *Scheduler) {
	query := s.filter.Query(w)
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}
