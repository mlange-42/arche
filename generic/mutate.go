package generic

import "github.com/mlange-42/arche/ecs"

// Mutate is a helper for mutating components.
type Mutate struct {
	add    []ecs.ID
	remove []ecs.ID
	world  *ecs.World
}

// NewMutate creates a new Mutate object.
func NewMutate(w *ecs.World) *Mutate {
	return &Mutate{
		world: w,
	}
}

// Adds sets components to add in calls to [Mutate.Add] and [Mutate.Exchange].
//
// Create the required mask items with [T].
func (m *Mutate) Adds(add ...Comp) *Mutate {
	m.add = toIds(m.world, add)
	return m
}

// Removes sets components to remove in calls to [Mutate.Remove] and [Mutate.Exchange].
//
// Create the required mask items with [T].
func (m *Mutate) Removes(remove ...Comp) *Mutate {
	m.remove = toIds(m.world, remove)
	return m
}

// NewEntity creates a new [ecs.Entity] with the components set via [Mutate.ToAdd].
//
// See also [ecs.World.NewEntity].
func (m *Mutate) NewEntity() ecs.Entity {
	entity := m.world.NewEntity(m.add...)
	return entity
}

// Add the components set via [Mutate.Adds] to the given entity.
//
// See also [ecs.World.Add].
func (m *Mutate) Add(entity ecs.Entity) {
	m.world.Add(entity, m.add...)
}

// Remove the components set via [Mutate.Removes] from the given entity.
//
// See also [ecs.World.Remove].
func (m *Mutate) Remove(entity ecs.Entity) {
	m.world.Remove(entity, m.remove...)
}

// Exchange components on an entity.
//
// Removes the components set via [Mutate.Removes].
// Adds the components set via [Mutate.Adds].
//
// See also [ecs.World.Exchange].
func (m *Mutate) Exchange(entity ecs.Entity) {
	m.world.Exchange(entity, m.add, m.remove)
}
