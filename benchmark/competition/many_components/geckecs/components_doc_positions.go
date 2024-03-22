package geckecs

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasPosition() bool {
	return e.w.positionsStore.Has(e)
}

func (e Entity) ReadPosition() (Position, bool) {
	return e.w.positionsStore.Read(e)
}

func (e Entity) RemovePosition() Entity {
	e.w.positionsStore.Remove(e)

	return e
}

func (e Entity) WritablePosition() (*Position, bool) {
	return e.w.positionsStore.Writeable(e)
}

func (e Entity) SetPosition(other Position) Entity {
	e.w.positionsStore.Set(other, e)

	e.w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(e)
	return e
}

func (w *World) SetPositions(c Position, entities ...Entity) {
	w.positionsStore.Set(c, entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

func (w *World) RemovePositions(entities ...Entity) {
	w.positionsStore.Remove(entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

//#region Resources

// HasPosition checks if the world has a Position}}
func (w *World) HasPositionResource() bool {
	return w.resourceEntity.HasPosition()
}

// Retrieve the Position resource from the world
func (w *World) PositionResource() (Position, bool) {
	return w.resourceEntity.ReadPosition()
}

// Set the Position resource in the world
func (w *World) SetPositionResource(c Position) Entity {
	w.resourceEntity.SetPosition(c)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Position resource from the world
func (w *World) RemovePositionResource() Entity {
	w.resourceEntity.RemovePosition()

	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type PositionReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Position]
}

func (iter *PositionReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *PositionReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *PositionReadIterator) NextPosition() (Entity, Position) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *PositionReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) PositionReadIter() *PositionReadIterator {
	iter := &PositionReadIterator{
		w:     w,
		store: w.positionsStore,
	}
	iter.Reset()
	return iter
}

type PositionWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Position]
}

func (iter *PositionWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *PositionWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *PositionWriteIterator) NextPosition() (Entity, *Position) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *PositionWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) PositionWriteIter() *PositionWriteIterator {
	iter := &PositionWriteIterator{
		w:     w,
		store: w.positionsStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) PositionEntities() []Entity {
	return w.positionsStore.entities()
}

func (w *World) SetPositionSortFn(lessThan func(a, b Entity) bool) {
	w.positionsStore.LessThan = lessThan
}

func (w *World) SortPositions() {
	w.positionsStore.Sort()
}
