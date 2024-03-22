package geckecs

type Comp5 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasComp5() bool {
	return e.w.comp5SStore.Has(e)
}

func (e Entity) ReadComp5() (Comp5, bool) {
	return e.w.comp5SStore.Read(e)
}

func (e Entity) RemoveComp5() Entity {
	e.w.comp5SStore.Remove(e)

	return e
}

func (e Entity) WritableComp5() (*Comp5, bool) {
	return e.w.comp5SStore.Writeable(e)
}

func (e Entity) SetComp5(other Comp5) Entity {
	e.w.comp5SStore.Set(other, e)

	e.w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(e)
	return e
}

func (w *World) SetComp5S(c Comp5, entities ...Entity) {
	w.comp5SStore.Set(c, entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

func (w *World) RemoveComp5S(entities ...Entity) {
	w.comp5SStore.Remove(entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

//#region Resources

// HasComp5 checks if the world has a Comp5}}
func (w *World) HasComp5Resource() bool {
	return w.resourceEntity.HasComp5()
}

// Retrieve the Comp5 resource from the world
func (w *World) Comp5Resource() (Comp5, bool) {
	return w.resourceEntity.ReadComp5()
}

// Set the Comp5 resource in the world
func (w *World) SetComp5Resource(c Comp5) Entity {
	w.resourceEntity.SetComp5(c)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Comp5 resource from the world
func (w *World) RemoveComp5Resource() Entity {
	w.resourceEntity.RemoveComp5()

	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type Comp5ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp5]
}

func (iter *Comp5ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp5ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp5ReadIterator) NextComp5() (Entity, Comp5) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *Comp5ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp5ReadIter() *Comp5ReadIterator {
	iter := &Comp5ReadIterator{
		w:     w,
		store: w.comp5SStore,
	}
	iter.Reset()
	return iter
}

type Comp5WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp5]
}

func (iter *Comp5WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp5WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp5WriteIterator) NextComp5() (Entity, *Comp5) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *Comp5WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp5WriteIter() *Comp5WriteIterator {
	iter := &Comp5WriteIterator{
		w:     w,
		store: w.comp5SStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp5Entities() []Entity {
	return w.comp5SStore.entities()
}

func (w *World) SetComp5SortFn(lessThan func(a, b Entity) bool) {
	w.comp5SStore.LessThan = lessThan
}

func (w *World) SortComp5S() {
	w.comp5SStore.Sort()
}
