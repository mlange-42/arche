package geckecs

type Comp9 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//#region Events
//#endregion

func (e Entity) HasComp9() bool {
	return e.w.comp9SStore.Has(e)
}

func (e Entity) ReadComp9() (Comp9, bool) {
	return e.w.comp9SStore.Read(e)
}

func (e Entity) RemoveComp9() Entity {
	e.w.comp9SStore.Remove(e)

	return e
}

func (e Entity) WritableComp9() (*Comp9, bool) {
	return e.w.comp9SStore.Writeable(e)
}

func (e Entity) SetComp9(other Comp9) Entity {
	e.w.comp9SStore.Set(other, e)

	e.w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(e)
	return e
}

func (w *World) SetComp9S(c Comp9, entities ...Entity) {
	w.comp9SStore.Set(c, entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

func (w *World) RemoveComp9S(entities ...Entity) {
	w.comp9SStore.Remove(entities...)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(entities...)
}

//#region Resources

// HasComp9 checks if the world has a Comp9}}
func (w *World) HasComp9Resource() bool {
	return w.resourceEntity.HasComp9()
}

// Retrieve the Comp9 resource from the world
func (w *World) Comp9Resource() (Comp9, bool) {
	return w.resourceEntity.ReadComp9()
}

// Set the Comp9 resource in the world
func (w *World) SetComp9Resource(c Comp9) Entity {
	w.resourceEntity.SetComp9(c)
	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

// Remove the Comp9 resource from the world
func (w *World) RemoveComp9Resource() Entity {
	w.resourceEntity.RemoveComp9()

	w.Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet.PossibleUpdate(w.resourceEntity)
	return w.resourceEntity
}

//#endregion

//#region Iterators

type Comp9ReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp9]
}

func (iter *Comp9ReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *Comp9ReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *Comp9ReadIterator) NextComp9() (Entity, Comp9) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *Comp9ReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) Comp9ReadIter() *Comp9ReadIterator {
	iter := &Comp9ReadIterator{
		w:     w,
		store: w.comp9SStore,
	}
	iter.Reset()
	return iter
}

type Comp9WriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Comp9]
}

func (iter *Comp9WriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *Comp9WriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *Comp9WriteIterator) NextComp9() (Entity, *Comp9) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *Comp9WriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) Comp9WriteIter() *Comp9WriteIterator {
	iter := &Comp9WriteIterator{
		w:     w,
		store: w.comp9SStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) Comp9Entities() []Entity {
	return w.comp9SStore.entities()
}

func (w *World) SetComp9SortFn(lessThan func(a, b Entity) bool) {
	w.comp9SStore.LessThan = lessThan
}

func (w *World) SortComp9S() {
	w.comp9SStore.Sort()
}
