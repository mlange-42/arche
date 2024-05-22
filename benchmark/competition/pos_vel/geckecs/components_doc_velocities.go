package geckecs

type Velocity struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasVelocity() bool {
	return e.w.velocitiesStore.Has(e)
}

func (e Entity) ReadVelocity() (Velocity, bool) {
	return e.w.velocitiesStore.Read(e)
}

func (e Entity) RemoveVelocity() Entity {
	e.w.velocitiesStore.Remove(e)

	return e
}

func (e Entity) WritableVelocity() (*Velocity, bool) {
	return e.w.velocitiesStore.Writeable(e)
}

func (e Entity) SetVelocity(other Velocity) Entity {
	e.w.velocitiesStore.Set(other, e)

	e.w.PositionVelocitySet.PossibleUpdate(e)
	return e
}

func (w *World) SetVelocities(c Velocity, entities ...Entity) {
	w.velocitiesStore.Set(c, entities...)
	w.PositionVelocitySet.PossibleUpdate(entities...)
}

func (w *World) RemoveVelocities(entities ...Entity) {
	w.velocitiesStore.Remove(entities...)
	w.PositionVelocitySet.PossibleUpdate(entities...)
}

//#region Resources

// HasVelocity checks if the world has a Velocity}}
func (w *World) HasVelocityResource() bool {
	return w.resourceEntity.HasVelocity()
}

// Retrieve the Velocity resource from the world
func (w *World) VelocityResource() (Velocity, bool) {
	return w.resourceEntity.ReadVelocity()
}

// Set the Velocity resource in the world
func (w *World) SetVelocityResource(c Velocity) Entity {
	w.resourceEntity.SetVelocity(c)
	w.PositionVelocitySet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Velocity resource from the world
func (w *World) RemoveVelocityResource() Entity {
	w.resourceEntity.RemoveVelocity()

	w.PositionVelocitySet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type VelocityReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Velocity]
}

func (iter *VelocityReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *VelocityReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *VelocityReadIterator) NextVelocity() (Entity, Velocity) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *VelocityReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) VelocityReadIter() *VelocityReadIterator {
	iter := &VelocityReadIterator{
		w:     w,
		store: w.velocitiesStore,
	}
	iter.Reset()
	return iter
}

type VelocityWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Velocity]
}

func (iter *VelocityWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *VelocityWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *VelocityWriteIterator) NextVelocity() (Entity, *Velocity) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *VelocityWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) VelocityWriteIter() *VelocityWriteIterator {
	iter := &VelocityWriteIterator{
		w:     w,
		store: w.velocitiesStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) VelocityEntities() []Entity {
	return w.velocitiesStore.entities()
}

func (w *World) SetVelocitySortFn(lessThan func(a, b Entity) bool) {
	w.velocitiesStore.LessThan = lessThan
}

func (w *World) SortVelocities() {
	w.velocitiesStore.Sort()
}
